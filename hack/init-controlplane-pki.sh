#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

LOCAL_IP_LIST=$(ip a | grep inet |  awk '{print $2}' | cut -d/ -f1 | paste -sd "," -)
OUT_DIR=${PKIDIR:-$PWD/pki}
HACKDIR=$(dirname "${BASH_SOURCE[0]}")
ROOTDIR=$(dirname "${BASH_SOURCE[0]}")/..

[ -d $OUT_DIR ] || mkdir -p $OUT_DIR

go run $HACKDIR/gen_cert.go --host "kubernetes,kubernetes.default,kubernetes.default.svc,kubernetes.default.svc.cluster.local,localhost,apiserver,kube-apiserver,172.20.0.1,$LOCAL_IP_LIST"  --ecdsa-curve P256 --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h \
	&& mv tls.crt ./pki/apiserver.crt && mv ./tls.key ./pki/apiserver.key

openssl ecparam -name prime256v1 -genkey -noout -out $OUT_DIR/sa-ecdsa.key \
	&& openssl ec -in $OUT_DIR/sa-ecdsa.key -pubout -out $OUT_DIR/sa-ecdsa.pub

printf 'admin-token,admin,,"system:masters"' > $OUT_DIR/token-users.csv

cat $ROOTDIR/api.kubeconfig | sed -e 's/localhost/apiserver/g' -e 's#./pki/apiserver.crt#/etc/kubernetes/pki/apiserver.crt#g' > $OUT_DIR/api.kubeconfig