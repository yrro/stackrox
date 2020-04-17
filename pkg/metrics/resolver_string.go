// Code generated by "stringer -type=Resolver"; DO NOT EDIT.

package metrics

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Cluster-0]
	_ = x[Compliance-1]
	_ = x[ComlianceControl-2]
	_ = x[CVEs-3]
	_ = x[Deployments-4]
	_ = x[Groups-5]
	_ = x[Images-6]
	_ = x[ImageComponents-7]
	_ = x[K8sRoles-8]
	_ = x[Namespaces-9]
	_ = x[Nodes-10]
	_ = x[Notifiers-11]
	_ = x[Policies-12]
	_ = x[Roles-13]
	_ = x[Root-14]
	_ = x[Secrets-15]
	_ = x[ServiceAccounts-16]
	_ = x[Subjects-17]
	_ = x[Tokens-18]
	_ = x[Violations-19]
	_ = x[Pods-20]
	_ = x[ContainerInstances-21]
}

const _Resolver_name = "ClusterComplianceComlianceControlCVEsDeploymentsGroupsImagesImageComponentsK8sRolesNamespacesNodesNotifiersPoliciesRolesRootSecretsServiceAccountsSubjectsTokensViolationsPodsContainerInstances"

var _Resolver_index = [...]uint8{0, 7, 17, 33, 37, 48, 54, 60, 75, 83, 93, 98, 107, 115, 120, 124, 131, 146, 154, 160, 170, 174, 192}

func (i Resolver) String() string {
	if i < 0 || i >= Resolver(len(_Resolver_index)-1) {
		return "Resolver(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Resolver_name[_Resolver_index[i]:_Resolver_index[i+1]]
}
