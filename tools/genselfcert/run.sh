#!/usr/bin/env bash
set -ex;

NS=${NAMESPACE:-default,kube-system}
SUFFIX=${SVC_SUFFIX:-cluster.local}
# \([^,]\+\), capture group that matches one or more characters that are not commas.
NSHOSTS=$(printf $NS | sed 's/\([^,]\+\)/*.\1.svc/g')
SUFFIXED_NSHOSTS=$(printf $NSHOSTS | sed "s/\([^,]\+\)/\1.$SUFFIX/g")

/gencert --host "${HOSTS:+$HOSTS,}$NSHOSTS,$SUFFIXED_NSHOSTS,127.0.0.1,::1"  --ecdsa-curve P256 --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h

cat tls.crt
cat tls.key