#!/usr/bin/env bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd)"

KUBE_COMMAND=${KUBE_COMMAND:-{{.K8sConfig.Command}}}

${KUBE_COMMAND} get namespace stackrox &>/dev/null || ${KUBE_COMMAND} create namespace stackrox
${KUBE_COMMAND} -n stackrox annotate namespace/stackrox --overwrite openshift.io/node-selector=""

if ! ${KUBE_COMMAND} get secret/{{.K8sConfig.ScannerSecretName}} -n stackrox > /dev/null; then
  registry_auth="$("${DIR}/docker-auth.sh" -m k8s "{{.K8sConfig.ScannerRegistry}}")"
  [[ -n "$registry_auth" ]] || { echo >&2 "Unable to get registry auth info." ; exit 1 ; }
  ${KUBE_COMMAND} create --namespace "stackrox" -f - <<EOF
apiVersion: v1
data:
  .dockerconfigjson: ${registry_auth}
kind: Secret
metadata:
  name: {{.K8sConfig.ScannerSecretName}}
  namespace: stackrox
  labels:
    app.kubernetes.io/name: stackrox
type: kubernetes.io/dockerconfigjson
EOF
fi
