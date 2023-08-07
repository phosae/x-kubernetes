#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

LOCAL_IP_LIST=$(ip a | grep inet |  awk '{print $2}' | cut -d/ -f1 | paste -sd "," -)
OUT_DIR=${PKIDIR:-$PWD/pki}
HACKDIR=$(dirname "${BASH_SOURCE[0]}")
ROOTDIR=$(dirname "${BASH_SOURCE[0]}")/..

[ -d $OUT_DIR ] || mkdir -p $OUT_DIR

docker run -v $PWD:/output -w /output -e HOSTS=localhost,apiserver,kube-apiserver,$LOCAL_IP_LIST -e NAMESPACE=default,kube-system zengxu/genselfcert \
	&& mv tls.crt ./pki/apiserver.crt && mv ./tls.key ./pki/apiserver.key

openssl ecparam -name prime256v1 -genkey -noout -out $OUT_DIR/sa-ecdsa.key \
	&& openssl ec -in $OUT_DIR/sa-ecdsa.key -pubout -out $OUT_DIR/sa-ecdsa.pub

printf 'admin-token,admin,,"system:masters"' > $OUT_DIR/token-users.csv

cat api.kubeconfig | sed -e 's/localhost/apiserver/g' -e 's#./pki/apiserver.crt#/etc/kubernetes/pki/apiserver.crt#g' > $OUT_DIR/api.kubeconfig