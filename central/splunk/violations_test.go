package splunk

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gogo/protobuf/types"
	"github.com/golang/mock/gomock"
	"github.com/stackrox/rox/central/alert/datastore/mocks"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/pkg/utils"
	"github.com/stretchr/testify/suite"
	"k8s.io/client-go/util/jsonpath"
)

var (
	// In case later we'd need to make adjustments or take different samples, the following structs were dumped from a
	// live system by simply using "github.com/mitranim/repr" module and printing them with `repr.Println(alert)` call.
	// Next, I removed some redundant `&` operators that compiler did not like, adjusted enums to use symbols (such as
	// storage.Severity_HIGH_SEVERITY) instead of integer values (e.g. 3) and made timestamps look human-friendly by
	// means of makeTimestamp() calls.

	processAlert = storage.Alert{
		Id: "f2d0efaa-2c54-402c-aeed-5b88ed5ccb8a",
		Policy: &storage.Policy{
			Id:          "f0bacecd-87be-4f51-89a5-8f86ad523620",
			Name:        "nmap Execution",
			Description: "Alerts when the nmap process launches in a container during run time",
			Rationale:   "Nmap can be used to probe a running container's network to enumerate open ports and perform other actions such as OS version detection and launching over-the-network scripted attacks",
			Remediation: "Consider removing package managers during the build process that could be used to download such software. Check that exposed ports don't allow for remote code execution",
			Categories: []string{
				"Network Tools",
			},
			LifecycleStages:    []storage.LifecycleStage{storage.LifecycleStage_RUNTIME},
			Severity:           storage.Severity_HIGH_SEVERITY,
			SORTName:           "nmap Execution",
			SORTLifecycleStage: "RUNTIME",
			PolicyVersion:      "1.1",
			PolicySections: []*storage.PolicySection{{
				PolicyGroups: []*storage.PolicyGroup{{
					FieldName: "Process Name",
					Values: []*storage.PolicyValue{{
						Value: "nmap",
					}},
				}},
			}},
		},
		LifecycleStage: storage.LifecycleStage_RUNTIME,
		Entity: &storage.Alert_Deployment_{
			Deployment: &storage.Alert_Deployment{
				Id:          "0f709d63-f2cc-4825-a984-b9cfd25b02cd",
				Name:        "debian-test",
				Type:        "Pod",
				Namespace:   "stackrox",
				NamespaceId: "dff6a17e-f246-4dc0-98b3-c70ee59c3cea",
				Labels: map[string]string{
					"run": "debian-test",
				},
				ClusterId:   "098e0e05-a96b-43ca-95af-3ef72cd32828",
				ClusterName: "remote",
				Containers: []*storage.Alert_Deployment_Container{{
					Image: &storage.ContainerImage{
						Id: "sha256:b16f66714660c4b3ea14d273ad8c35079b81b35d65d1e206072d226c7ff78299",
						Name: &storage.ImageName{
							Registry: "docker.io",
							Remote:   "library/debian",
							Tag:      "latest",
							FullName: "docker.io/library/debian:latest",
						},
					},
					Name: "debian-test",
				}},
				Annotations: map[string]string{
					"cni.projectcalico.org/podIP": "10.65.48.8/32",
				},
			},
		},
		ProcessViolation: &storage.Alert_ProcessViolation{
			Message: "Binary '/usr/bin/nmap' executed with arguments '-v -A localhost' under user ID 0",
			Processes: []*storage.ProcessIndicator{{
				Id:            "8472f6e2-53d2-4ddf-ad59-ecc43a8d98d2",
				DeploymentId:  "0f709d63-f2cc-4825-a984-b9cfd25b02cd",
				ContainerName: "debian-test",
				PodId:         "debian-test",
				PodUid:        "e20f2691-1371-588f-ae6d-bd0ef24af78b",
				Signal: &storage.ProcessSignal{
					Id:           "2569b112-64b1-11eb-9541-f65aedf20953",
					ContainerId:  "111bf6d5e461",
					Time:         makeTimestamp("2021-02-01T17:18:49.421852357Z"),
					Name:         "nmap",
					Args:         "-v -A localhost",
					ExecFilePath: "/usr/bin/nmap",
					Pid:          64307,
					LineageInfo: []*storage.ProcessSignal_LineageInfo{{
						ParentExecFilePath: "/bin/bash",
					}},
				},
				Namespace:          "stackrox",
				ContainerStartTime: makeTimestamp("2021-02-01T16:17:32Z"),
			}, {
				Id:            "cfc994d5-11bd-4471-a82e-b1735ad94e06",
				DeploymentId:  "0f709d63-f2cc-4825-a984-b9cfd25b02cd",
				ContainerName: "debian-test",
				PodId:         "debian-test",
				PodUid:        "e20f2691-1371-588f-ae6d-bd0ef24af78b",
				Signal: &storage.ProcessSignal{
					Id:           "8c84d098-64b1-11eb-9541-f65aedf20953",
					ContainerId:  "111bf6d5e461",
					Time:         makeTimestamp("2021-02-01T17:15:56.457252Z"),
					Name:         "nmap",
					Args:         "-v -A localhost",
					ExecFilePath: "/usr/bin/nmap",
					Pid:          65923,
					LineageInfo: []*storage.ProcessSignal_LineageInfo{{
						ParentExecFilePath: "/bin/bash",
					}},
				},
				Namespace:          "stackrox",
				ContainerStartTime: makeTimestamp("2021-02-01T16:17:32Z"),
			}},
		},
		Time:          makeTimestamp("2021-02-01T17:18:49.439085673Z"),
		FirstOccurred: makeTimestamp("2021-02-01T17:15:56.474524288Z"),
	}

	k8sAlert = storage.Alert{
		Id: "90e0feed-662c-4593-b414-e55d1eaff017",
		Policy: &storage.Policy{
			Id:          "8ab0f199-4904-4808-9461-3501da1d1b77",
			Name:        "Kubernetes Actions: Exec into Pod",
			Description: "Alerts when Kubernetes API receives request to execute command in container",
			Rationale:   "'pods/exec' is non-standard approach for interacting with containers. Attackers with permissions could execute malicious code and compromise resources within a cluster",
			Remediation: "Restrict RBAC access to the 'pods/exec' resource according to the Principle of Least Privilege. Limit such usage only to development, testing or debugging (non-production) activities",
			Categories: []string{
				"Kubernetes Events",
			},
			LifecycleStages:    []storage.LifecycleStage{storage.LifecycleStage_RUNTIME},
			Severity:           storage.Severity_HIGH_SEVERITY,
			SORTName:           "Kubernetes Actions: Exec into Pod",
			SORTLifecycleStage: "RUNTIME",
			PolicyVersion:      "1.1",
			PolicySections: []*storage.PolicySection{{
				PolicyGroups: []*storage.PolicyGroup{{
					FieldName: "Kubernetes Resource",
					Values: []*storage.PolicyValue{{
						Value: "PODS_EXEC",
					}},
				}},
			}},
		},
		LifecycleStage: storage.LifecycleStage_RUNTIME,
		Entity: &storage.Alert_Deployment_{
			Deployment: &storage.Alert_Deployment{
				Id:          "587556aa-5885-4a6e-8389-9d4e1c36e42a",
				Name:        "central",
				Type:        "Deployment",
				Namespace:   "stackrox",
				NamespaceId: "75868f7c-5949-4de3-bfc0-579f80148d45",
				Labels: map[string]string{
					"app":                          "central",
					"app.kubernetes.io/component":  "central",
					"app.kubernetes.io/instance":   "stackrox-central-services",
					"app.kubernetes.io/managed-by": "Helm",
					"app.kubernetes.io/name":       "stackrox",
					"app.kubernetes.io/part-of":    "stackrox-central-services",
					"app.kubernetes.io/version":    "3.0.55.x-118-gec7dc725f2-dirty",
					"helm.sh/chart":                "stackrox-central-services-55.0.118-gec7dc725f2-dirty",
				},
				ClusterId:   "943451bd-54c8-437d-98fa-820f5b9ad431",
				ClusterName: "remote",
				Containers: []*storage.Alert_Deployment_Container{{
					Image: &storage.ContainerImage{
						Id: "sha256:d4c1df40d209978307551e4b0a000067105e07578b66dfc8a4929f59dce86368",
						Name: &storage.ImageName{
							Registry: "docker.io",
							Remote:   "stackrox/main",
							Tag:      "3.0.55.x-118-gec7dc725f2-dirty",
							FullName: "docker.io/stackrox/main:3.0.55.x-118-gec7dc725f2-dirty",
						},
						NotPullable: true,
					},
					Name: "central",
				}},
				Annotations: map[string]string{
					"meta.helm.sh/release-name":      "stackrox-central-services",
					"meta.helm.sh/release-namespace": "stackrox",
					"owner":                          "stackrox",
					"email":                          "support@stackrox.com",
				},
			},
		},
		Violations: []*storage.Alert_Violation{{
			Message: "Kubernetes API received exec '/go/bin/dlv --headless --listen=:40000 --api-version=2 --accept-multiclient attach 1 --continue' request into pod 'central-6c8f4d4d8d-9hxpt' container 'central'",
			MessageAttributes: &storage.Alert_Violation_KeyValueAttrs_{
				KeyValueAttrs: &storage.Alert_Violation_KeyValueAttrs{
					Attrs: []*storage.Alert_Violation_KeyValueAttrs_KeyValueAttr{{
						Key:   "pod",
						Value: "central-6c8f4d4d8d-9hxpt",
					}, {
						Key:   "container",
						Value: "central",
					}, {
						Key:   "commands",
						Value: "/go/bin/dlv --headless --listen=:40000 --api-version=2 --accept-multiclient attach 1 --continue",
					}},
				},
			},
			Type: storage.Alert_Violation_K8S_EVENT,
			Time: makeTimestamp("2021-02-15T19:04:36.843302212Z"),
		}, {
			Message: "Kubernetes API received exec '/bin/sh -c [ -e /proc/sys/kernel/yama/ptrace_scope ] && cat /proc/sys/kernel/yama/ptrace_scope || echo 0' request into pod 'central-6c8f4d4d8d-9hxpt' container 'central'",
			MessageAttributes: &storage.Alert_Violation_KeyValueAttrs_{
				KeyValueAttrs: &storage.Alert_Violation_KeyValueAttrs{
					Attrs: []*storage.Alert_Violation_KeyValueAttrs_KeyValueAttr{{
						Key:   "pod",
						Value: "central-6c8f4d4d8d-9hxpt",
					}, {
						Key:   "container",
						Value: "central",
					}, {
						Key:   "commands",
						Value: "/bin/sh -c [ -e /proc/sys/kernel/yama/ptrace_scope ] && cat /proc/sys/kernel/yama/ptrace_scope || echo 0",
					}},
				},
			},
			Type: storage.Alert_Violation_K8S_EVENT,
			Time: makeTimestamp("2021-02-15T19:04:36.659410153Z"),
		}},
		Time:          makeTimestamp("2021-02-15T19:04:36.843516328Z"),
		FirstOccurred: makeTimestamp("2021-02-15T19:04:36.662294945Z"),
	}

	deployAlert = storage.Alert{
		Id: "f56ffae8-adf9-4983-8e56-e260f1ab3dc9",
		Policy: &storage.Policy{
			Id:          "2db9a279-2aec-4618-a85d-7f1bdf4911b1",
			Name:        "90-Day Image Age",
			Description: "Alert on deployments with images that haven't been updated in 90 days",
			Rationale:   "Base images are updated frequently with bug fixes and vulnerability patches. Image age exceeding 90 days may indicate a higher risk of vulnerabilities existing in the image.",
			Remediation: "Rebuild your image, push a new minor version (with a new immutable tag), and update your service to use it.",
			Categories: []string{
				"DevOps Best Practices",
				"Security Best Practices",
			},
			LifecycleStages: []storage.LifecycleStage{storage.LifecycleStage_BUILD, storage.LifecycleStage_DEPLOY},
			Exclusions: []*storage.Exclusion{{
				Name: "Don't alert on kube-system namespace",
				Deployment: &storage.Exclusion_Deployment{
					Scope: &storage.Scope{
						Namespace: "kube-system",
					},
				},
			}, {
				Name: "Don't alert on istio-system namespace",
				Deployment: &storage.Exclusion_Deployment{
					Scope: &storage.Scope{
						Namespace: "istio-system",
					},
				},
			}},
			Severity:           storage.Severity_LOW_SEVERITY,
			SORTName:           "90-Day Image Age",
			SORTLifecycleStage: "BUILD,DEPLOY",
			PolicyVersion:      "1.1",
			PolicySections: []*storage.PolicySection{{
				PolicyGroups: []*storage.PolicyGroup{{
					FieldName: "Image Age",
					Values: []*storage.PolicyValue{{
						Value: "90",
					}},
				}},
			}},
		},
		Entity: &storage.Alert_Deployment_{
			Deployment: &storage.Alert_Deployment{
				Id:          "565bdd7a-eb3e-4367-9b73-87a9bcb8f4e7",
				Name:        "monitoring",
				Type:        "Deployment",
				Namespace:   "stackrox",
				NamespaceId: "dff6a17e-f246-4dc0-98b3-c70ee59c3cea",
				Labels: map[string]string{
					"app":                          "monitoring",
					"app.kubernetes.io/managed-by": "Helm",
					"app.kubernetes.io/name":       "stackrox",
				},
				ClusterId:   "098e0e05-a96b-43ca-95af-3ef72cd32828",
				ClusterName: "remote",
				Containers: []*storage.Alert_Deployment_Container{{
					Image: &storage.ContainerImage{
						Id: "sha256:488ce940267b9b7e281779845d45c6aef36774ed4ca54b2aef67104bf70dee23",
						Name: &storage.ImageName{
							Registry: "docker.io",
							Remote:   "stackrox/monitoring",
							Tag:      "1.0.0",
							FullName: "docker.io/stackrox/monitoring:1.0.0",
						},
					},
					Name: "grafana",
				}, {
					Image: &storage.ContainerImage{
						Id: "sha256:488ce940267b9b7e281779845d45c6aef36774ed4ca54b2aef67104bf70dee23",
						Name: &storage.ImageName{
							Registry: "docker.io",
							Remote:   "stackrox/monitoring",
							Tag:      "1.0.0",
							FullName: "docker.io/stackrox/monitoring:1.0.0",
						},
					},
					Name: "influxdb",
				}, {
					Image: &storage.ContainerImage{
						Id: "sha256:488ce940267b9b7e281779845d45c6aef36774ed4ca54b2aef67104bf70dee23",
						Name: &storage.ImageName{
							Registry: "docker.io",
							Remote:   "stackrox/monitoring",
							Tag:      "1.0.0",
							FullName: "docker.io/stackrox/monitoring:1.0.0",
						},
					},
					Name: "telegraf-proxy",
				}},
				Annotations: map[string]string{
					"owner":                          "stackrox",
					"email":                          "support@stackrox.com",
					"meta.helm.sh/release-name":      "stackrox-monitoring",
					"meta.helm.sh/release-namespace": "stackrox",
				},
			},
		},
		Violations: []*storage.Alert_Violation{{
			Message: "Container 'grafana' has image created at 2020-09-28 17:03:00 (UTC)",
		}, {
			Message: "Container 'influxdb' has image created at 2020-09-28 17:03:00 (UTC)",
		}, {
			Message: "Container 'telegraf-proxy' has image created at 2020-09-28 17:03:00 (UTC)",
		}},
		Time:          makeTimestamp("2021-02-01T16:09:02.193352817Z"),
		FirstOccurred: makeTimestamp("2021-02-01T16:09:02.128791072Z"),
	}
)

func makeTimestamp(timeStr string) *types.Timestamp {
	ts, err := parseTimestamp(timeStr)
	utils.Must(err)
	return ts
}

func TestViolations(t *testing.T) {
	suite.Run(t, &violationsTestSuite{})
}

type violationsTestSuite struct {
	suite.Suite
	processAlert, k8sAlert, deployAlert storage.Alert
}

func (s *violationsTestSuite) SetupTest() {
	s.processAlert = *processAlert.Clone()
	s.k8sAlert = *k8sAlert.Clone()
	s.deployAlert = *deployAlert.Clone()
}

func (s *violationsTestSuite) TestProcessAlert() {
	vs := s.getViolations(s.requestAndGetBody("", &s.processAlert))
	s.Len(vs, 2)

	for _, v := range vs {
		s.Equal("PROCESS_EVENT", s.extr(v, ".violationInfo.violationType"))
		s.Equal(float64(0), s.extr(v, ".processInfo.processUid"))
		s.Equal(float64(0), s.extr(v, ".processInfo.processGid"))
		s.Equal(s.extr(v, ".processInfo.processCreationTime"), s.extr(v, ".violationInfo.violationTime"))

		s.checkViolationInfo(v)
		s.checkProcessInfo(v)
		s.checkAlertInfo(v, ".lifecycleStage")
		s.checkDeploymentInfo(v)
		s.checkPolicy(v)
	}
}

func (s *violationsTestSuite) TestK8sAlert() {
	vs := s.getViolations(s.requestAndGetBody("", &s.k8sAlert))
	s.Len(vs, 2)

	for _, v := range vs {
		s.Equal("K8S_EVENT", s.extr(v, ".violationInfo.violationType"))

		s.checkViolationInfo(v, ".violationMessageAttributes")
		s.checkAlertInfo(v, ".lifecycleStage")
		s.checkDeploymentInfo(v)
		s.checkPolicy(v)
	}
	s.Equal("2021-02-15T19:04:36.659410153Z", s.extr(vs[0], ".violationInfo.violationTime"))
	s.Equal("2021-02-15T19:04:36.843302212Z", s.extr(vs[1], ".violationInfo.violationTime"))
}

func (s *violationsTestSuite) TestDeployAlert() {
	vs := s.getViolations(s.requestAndGetBody("", &s.deployAlert))
	s.Len(vs, 3)

	for _, v := range vs {
		s.Equal("GENERIC", s.extr(v, ".violationInfo.violationType"))
		s.Equal("2021-02-01T16:09:02.128791072Z", s.extr(v, ".alertInfo.alertFirstOccurred"))
		s.Equal("2021-02-01T16:09:02.193352817Z", s.extr(v, ".violationInfo.violationTime"))

		s.checkViolationInfo(v)
		s.checkAlertInfo(v)
		s.checkDeploymentInfo(v)
		s.checkPolicy(v)
	}
}

func (s *violationsTestSuite) TestViolationsAreOrdered() {
	vs := s.getViolations(s.requestAndGetBody("", &s.processAlert, &s.k8sAlert, &s.deployAlert))

	s.Greater(len(vs), 2)
	for i := range vs {
		if i == 0 {
			continue
		}
		s.LessOrEqual(s.extr(vs[i-1], ".violationInfo.violationTime"), s.extr(vs[i], ".violationInfo.violationTime"))
	}
}

func (s *violationsTestSuite) TestViolationIdsAreDistinct() {
	vs := s.getViolations(s.requestAndGetBody("", &s.processAlert, &s.k8sAlert, &s.deployAlert))

	ids := set.StringSet{}
	for _, v := range vs {
		id, ok := s.extr(v, ".violationInfo.violationId").(string)
		s.Truef(ok, "Detected violationId that is not a string: %v", id)
		s.Truef(ids.Add(id), "violationId=%q is not unique. Already seen ids: %v", id, ids)
	}
}

func (s *violationsTestSuite) TestWithDeploymentImage() {
	alert := s.processAlert.Clone()
	// Change alert's Entity from Alert_Deployment to Alert_Image. Conveniently the former Alert_Deployment has a ContainerImage we can use for testing.
	alert.Entity = &storage.Alert_Image{
		Image: alert.GetDeployment().Containers[0].GetImage(),
	}

	vs := s.getViolations(s.requestAndGetBody("", alert))

	s.assertPresent(vs[0], ".deploymentInfo",
		// deploymentImage must obviously be present coming from above
		".deploymentImage",
		// other deployment details are obtained from ProcessViolation
		".deploymentId", ".deploymentNamespace")
}

func (s *violationsTestSuite) TestAlertWithoutPolicy() {
	alert := s.processAlert.Clone()
	alert.Policy = nil
	alert.ProcessViolation.Processes = alert.ProcessViolation.Processes[:1]
	vs := s.getViolations(s.requestAndGetBody("", alert))
	s.Nil(s.extr(vs[0], ".policyInfo"))
}

func (s *violationsTestSuite) TestProcessAlertWithoutProcessIndicators() {
	alert := s.processAlert.Clone()
	alert.ProcessViolation.Processes = []*storage.ProcessIndicator{}
	s.Empty(s.getViolations(s.requestAndGetBody("", alert)))
}

func (s *violationsTestSuite) TestProcessAlertWithoutProcessSignal() {
	alert := s.processAlert.Clone()
	alert.ProcessViolation.Processes = alert.ProcessViolation.Processes[:1]
	alert.ProcessViolation.Processes[0].Signal = nil
	vs := s.getViolations(s.requestAndGetBody("", alert))
	s.checkViolationInfo(vs[0])
	s.checkAlertInfo(vs[0])
	// That's all it can gather from ProcessIndicator without ProcessSignal
	s.assertPresent(vs[0], ".processInfo",
		".processViolationId", ".podId", ".podUid", ".containerName", ".containerStartTime")
	s.checkDeploymentInfo(vs[0])
	s.checkPolicy(vs[0])
}

func (s *violationsTestSuite) TestAlertWithoutViolations() {
	alert := s.deployAlert.Clone()
	alert.Violations = []*storage.Alert_Violation{}
	s.Empty(s.getViolations(s.requestAndGetBody("", alert)))
}

func (s *violationsTestSuite) TestK8sAlertWithoutDeployment() {
	alert := s.k8sAlert.Clone()
	alert.Entity = nil
	alert.Violations = alert.Violations[:1]
	vs := s.getViolations(s.requestAndGetBody("", alert))
	s.Empty(s.extr(vs[0], ".deploymentInfo"))
}

func (s *violationsTestSuite) TestProcessAlertWithoutDeployment() {
	alert := s.processAlert.Clone()
	alert.Entity = nil
	alert.ProcessViolation.Processes = alert.ProcessViolation.Processes[:1]
	vs := s.getViolations(s.requestAndGetBody("", alert))
	// deploymentInfo still has some attributes because they came from ProcessIndicator-s
	s.assertPresent(vs[0], ".deploymentInfo", ".deploymentId", ".deploymentNamespace")
}

func (s *violationsTestSuite) TestProcessAlertNotMatchingDeploymentId() {
	alert := s.processAlert.Clone()
	alert.ProcessViolation.Processes = alert.ProcessViolation.Processes[:1]
	alert.ProcessViolation.Processes[0].DeploymentId = "blah"
	vs := s.getViolations(s.requestAndGetBody("", alert))
	// DeploymentId value from ProcessIndicator should take priority
	s.Equal("blah", s.extr(vs[0], ".deploymentInfo.deploymentId"))
	s.NotEmpty(s.extr(vs[0], ".deploymentInfo.deploymentNamespace"))
}

func (s *violationsTestSuite) TestProcessAlertNotMatchingDeploymentInfo() {
	alert := s.processAlert.Clone()
	alert.ProcessViolation.Processes = alert.ProcessViolation.Processes[:1]
	alert.ProcessViolation.Processes[0].ClusterId = "blah-cluster"
	alert.ProcessViolation.Processes[0].Namespace = "blah-namespace"
	vs := s.getViolations(s.requestAndGetBody("", alert))
	s.Equal("blah-cluster", s.extr(vs[0], ".deploymentInfo.clusterId"))
	s.Equal("blah-namespace", s.extr(vs[0], ".deploymentInfo.deploymentNamespace"))
}

func (s *violationsTestSuite) TestDefaultCheckpointAndNoViolations() {
	body := s.requestAndGetBody("")
	s.Empty(s.getViolations(body))
	// API should return some default checkpoint for subsequent querying.
	cp := makeTimestamp(s.extr(body, ".newCheckpoint").(string))
	// That timestamp should be in the past. For example, earlier than the first commit to Kubernetes repo.
	s.True(cp.Compare(makeTimestamp("2014-06-06T23:40:48Z")) < 0)
}

func (s *violationsTestSuite) TestCheckpointInTheFuture() {
	body := s.requestAndGetBody("2130-12-31T23:59:59Z", &s.processAlert, &s.k8sAlert, &s.deployAlert)
	s.Empty(s.getViolations(body))
	// Incoming checkpoint should be echoed on the output.
	s.Equal("2130-12-31T23:59:59Z", s.extr(body, ".newCheckpoint"))
}

func (s *violationsTestSuite) TestCheckpointFiltering() {
	fromStr := "2021-02-01T17:18:48Z"         // Somewhere in the middle of violation timestamps
	toStr := "2021-02-15T19:04:36.843302212Z" // Timestamp of the last violation met in the data
	fromTs := makeTimestamp(fromStr)
	toTs := makeTimestamp(toStr)
	alerts := []*storage.Alert{&s.processAlert, &s.k8sAlert, &s.deployAlert}

	body := s.requestAndGetBody(fromStr, alerts...)
	vs := s.getViolations(body)
	s.NotEmpty(vs)
	s.Less(len(vs), len(s.getViolations(s.requestAndGetBody("", alerts...))))

	for _, v := range vs {
		ts := makeTimestamp(s.extr(v, ".violationInfo.violationTime").(string))
		s.True(ts.Compare(fromTs) > 0)
		s.True(ts.Compare(toTs) <= 0)
	}

	s.Equal(toStr, s.extr(body, ".newCheckpoint"))
}

func (s *violationsTestSuite) TestCheckpointBeforeData() {
	alerts := []*storage.Alert{&s.processAlert, &s.k8sAlert, &s.deployAlert}

	// Query with checkpoint. The checkpoint's date is before violation timestamps in the data.
	vs1 := s.getViolations(s.requestAndGetBody("2021-01-01T00:00:00Z", alerts...))
	// Query without checkpoint.
	vs2 := s.getViolations(s.requestAndGetBody("", alerts...))

	s.Equal(len(vs1), len(vs2))
}

func (s *violationsTestSuite) TestInvalidCheckpoint() {
	w := httptest.NewRecorder()

	s.request(w, nil, url.QueryEscape("This isn't any good timestamp"))

	s.Equal(http.StatusBadRequest, w.Code)
	body := w.Body.String()
	s.Contains(body, "could not parse query parameter")
	s.Contains(body, "from_checkpoint")
	s.Contains(body, "This isn't any good timestamp")
}

func (s *violationsTestSuite) TestFirstCheckpointWins() {
	checkpointParam := "2130-12-31T23:59:59Z&from_checkpoint=2005-01-01T00:00:00Z"
	vs := s.getViolations(s.requestAndGetBody(checkpointParam, &s.k8sAlert))
	s.Empty(vs)
}

func (s *violationsTestSuite) assertPresent(violation interface{}, prefix string, attributes ...string) {
	for _, attr := range attributes {
		s.NotEmpty(s.extr(violation, prefix+attr))
	}
}

func (s *violationsTestSuite) checkViolationInfo(violation interface{}, extraAttrs ...string) {
	s.assertPresent(violation, ".violationInfo",
		".violationId",
		".violationMessage",
		".violationType",
		".violationTime")
	s.assertPresent(violation, ".violationInfo", extraAttrs...)
}

func (s *violationsTestSuite) checkAlertInfo(violation interface{}, extraAttrs ...string) {
	s.assertPresent(violation, ".alertInfo", ".alertId")
	s.assertPresent(violation, ".alertInfo", extraAttrs...)
}

func (s *violationsTestSuite) checkProcessInfo(violation interface{}) {
	s.assertPresent(violation, ".processInfo",
		".processViolationId",
		".podId",
		".podUid",
		".containerName",
		".containerStartTime",
		".containerId",
		".processSignalId",
		".processCreationTime",
		".processName",
		".processArgs",
		".execFilePath",
		".pid",
		".processLineageInfo")
}

func (s *violationsTestSuite) checkDeploymentInfo(violation interface{}) {
	s.assertPresent(violation, ".deploymentInfo",
		".deploymentId",
		".deploymentName",
		".deploymentType",
		".deploymentNamespace",
		".deploymentNamespaceId",
		".deploymentLabels",
		".clusterId",
		".clusterName",
		".deploymentContainers",
		".deploymentAnnotations")
}

func (s *violationsTestSuite) checkPolicy(violation interface{}) {
	s.assertPresent(violation, ".policyInfo",
		".policyId",
		".policyName",
		".policyDescription",
		".policyRationale",
		".policyCategories",
		".policyLifecycleStages",
		".policySeverity",
		".policyVersion")
}

func (s *violationsTestSuite) requestAndGetBody(checkpointParam string, data ...*storage.Alert) map[string]interface{} {
	w := httptest.NewRecorder()
	s.request(w, nil, checkpointParam, data...)
	s.Equal(200, w.Code)

	var parsed map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &parsed)
	s.NoError(err)

	return parsed
}

func (s *violationsTestSuite) request(responseWriter http.ResponseWriter, searchAlertsError error, checkpointParam string, searchAlertsData ...*storage.Alert) {
	mockCtrl := gomock.NewController(s.T())
	defer mockCtrl.Finish()

	mockDS := mocks.NewMockDataStore(mockCtrl)

	mockDS.EXPECT().SearchRawAlerts(gomock.Any(), gomock.Any()).AnyTimes().Return(searchAlertsData, searchAlertsError)

	handler := NewViolationsHandler(mockDS)

	url := "/ignored"
	if checkpointParam != "" {
		url += "?from_checkpoint=" + checkpointParam
	}

	r := httptest.NewRequest("GET", url, nil)

	handler.ServeHTTP(responseWriter, r)
}

// getViolations extracts "violations" attribute as a slice for later querying them with JSONPath.
func (s *violationsTestSuite) getViolations(body map[string]interface{}) []interface{} {
	violations := s.extr(body, ".violations")
	if violations == nil {
		return nil
	}
	return violations.([]interface{})
}

// extr extracts value from input according to provided jsonPath. Returns nil if given attribute does not exist.
func (s *violationsTestSuite) extr(input interface{}, jsonPath string) interface{} {
	jp := jsonpath.New("")

	err := jp.Parse("{" + jsonPath + "}")
	s.NoError(err)

	val, err := jp.FindResults(input)
	if err != nil && strings.HasSuffix(err.Error(), " is not found") {
		return nil
	}
	s.NoError(err)

	return val[0][0].Interface()
}

func (s *violationsTestSuite) TestViolationsHandlerError() {
	w := httptest.NewRecorder()

	s.request(w, errors.New("mock error"), "")

	s.Equal(http.StatusInternalServerError, w.Code)
	s.Contains(w.Body.String(), "mock error")
}

// failingResponseWriter is an implementation of http.ResponseWriter that returns error on attempt to write to it
// to emulate e.g. closed connection.
type failingResponseWriter struct {
	header http.Header
}

func (f failingResponseWriter) Header() http.Header {
	return f.header
}
func (f failingResponseWriter) WriteHeader(_ int) {}
func (f failingResponseWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("mock http write error")
}

func (s *violationsTestSuite) TestViolationsHandlerWriteError() {
	w := failingResponseWriter{}
	s.PanicsWithError("net/http: abort Handler", func() {
		s.request(w, nil, "", &s.processAlert)
	})
}

func (s *violationsTestSuite) TestGenerateViolationId() {
	v1Empty := storage.Alert_Violation{}
	v2Empty := storage.Alert_Violation{}

	id1, err := generateViolationID("alert1", &v1Empty)
	s.Require().NoError(err)
	id2, err := generateViolationID("alert1", &v2Empty)
	s.Require().NoError(err)
	s.Equal(id1, id2)

	id2Other, err := generateViolationID("other-alert", &v2Empty)
	s.Require().NoError(err)
	s.NotEqual(id2, id2Other)

	v3 := storage.Alert_Violation{
		Message: "mock message",
		Type:    storage.Alert_Violation_K8S_EVENT,
		Time: &types.Timestamp{
			Seconds: 123,
			Nanos:   456,
		},
	}
	id3, err := generateViolationID("alert1", &v3)
	s.Require().NoError(err)

	v4 := v3.Clone()
	v4.Message = "mock message4"
	id4, err := generateViolationID("alert1", v4)
	s.Require().NoError(err)

	s.NotEqual(id3, id4)
}

func BenchmarkGenerateViolationId(b *testing.B) {
	violations := deployAlert.Clone().Violations
	violations = append(violations, k8sAlert.Clone().Violations...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := violations[i%len(violations)]
		id, err := generateViolationID(deployAlert.GetId(), v)
		if err != nil || len(id) != 36 {
			b.FailNow()
		}
	}
}
