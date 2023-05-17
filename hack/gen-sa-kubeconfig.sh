#!/usr/bin/env bash
set -eu

namespace=${1:-default}
sa=${2}
cluster=$(kubectl config view --minify --raw -o jsonpath='{.clusters[].name}')
ca=$(kubectl config view --minify --raw -o jsonpath='{.clusters[].cluster.certificate-authority-data}')
version=$(kubectl version --output json | jq -r '.serverVersion.minor' | tr -d '+')

# Compare the minor version and output the appropriate message
if [ "$version" -ge 24 ]; then
    kubectl config set-cluster ${cluster} --server=$(kubectl config view --minify --raw -o jsonpath='{.clusters[].cluster.server}') --certificate-authority=${ca}  --kubeconfig=${namespace}-${sa}.kubeconfig
    kubectl config set-credentials ${namespace}-${sa} --token=$(kubectl -n ${namespace} create token ${sa} --duration=8760h) --kubeconfig=${namespace}-${sa}.kubeconfig
    kubectl config set-context default --cluster ${cluster} --user ${namespace}-${sa} --kubeconfig=${namespace}-${sa}.kubeconfig
    kubectl config use-context default --kubeconfig=${namespace}-${sa}.kubeconfig
else
    secret=$(kubectl --namespace $namespace get serviceAccount $sa -o jsonpath='{.secrets[0].name}')
    kubectl config set-cluster ${cluster} --server=$(kubectl config view --minify --raw -o jsonpath='{.clusters[].cluster.server}') --certificate-authority=$(kubectl --namespace ${namespace} get secret/${secret} -o jsonpath='{.data.ca\.crt}')  --kubeconfig=${namespace}-${sa}.kubeconfig
    kubectl config set-credentials ${namespace}-${sa} --token=$(kubectl -n ${namespace} get secret/${secret} -o jsonpath={.data.token} | base64 -d) --kubeconfig=${namespace}-${sa}.kubeconfig
    kubectl config set-context default --cluster ${cluster} --user ${namespace}-${sa} --kubeconfig=${namespace}-${sa}.kubeconfig
    kubectl config use-context default --kubeconfig=${namespace}-${sa}.kubeconfig
fi

os_name="$(uname)"
if [ "$os_name" = "Darwin" ]; then
    sed -i '' 's/certificate-authority/certificate-authority-data/g' ${namespace}-${sa}.kubeconfig
elif [ "$os_name" = "Linux" ]; then
    sed -i 's/certificate-authority/certificate-authority-data/g' ${namespace}-${sa}.kubeconfig
fi