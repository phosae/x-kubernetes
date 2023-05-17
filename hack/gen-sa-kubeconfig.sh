#!/usr/bin/env bash
set -eu

namespace=${1:-default}
sa=${2}
cluster=$(kubectl config view --minify --raw -o jsonpath='{.clusters[].name}')
ca=$(kubectl config view --minify --raw -o jsonpath='{.clusters[].cluster.certificate-authority-data}')
version=$(kubectl version --output json | jq -r '.serverVersion.minor' | tr -d '+')

configure_cluster() {
  kubectl config set-cluster ${cluster} --server=$(kubectl config view --minify --raw -o jsonpath='{.clusters[].cluster.server}') "$@" --kubeconfig=${namespace}-${sa}.kubeconfig
}

configure_credentials() {
  kubectl config set-credentials ${namespace}-${sa} "$@" --kubeconfig=${namespace}-${sa}.kubeconfig
}

set_context_and_use() {
  kubectl config set-context default --cluster ${cluster} --user ${namespace}-${sa} --kubeconfig=${namespace}-${sa}.kubeconfig
  kubectl config use-context default --kubeconfig=${namespace}-${sa}.kubeconfig
}

if [ "$version" -ge 24 ]; then
    configure_cluster --certificate-authority=${ca}
    configure_credentials --token=$(kubectl -n ${namespace} create token ${sa} --duration=8760h)
    set_context_and_use
else
    secret=$(kubectl --namespace $namespace get serviceAccount $sa -o jsonpath='{.secrets[0].name}')
    configure_cluster --certificate-authority=$(kubectl --namespace ${namespace} get secret/${secret} -o jsonpath='{.data.ca\.crt}')
    configure_credentials --token=$(kubectl -n ${namespace} get secret/${secret} -o jsonpath={.data.token} | base64 -d)
    set_context_and_use
fi

os_name="$(uname)"
sed_opts=""
if [ "$os_name" = "Darwin" ]; then
    sed_opts="-i ''"
elif [ "$os_name" = "Linux" ]; then
    sed_opts="-i"
fi

sed ${sed_opts} 's/certificate-authority/certificate-authority-data/g' ${namespace}-${sa}.kubeconfig