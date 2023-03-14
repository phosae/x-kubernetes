#!/bin/bash

LOCAL_IP_LIST=$(ip a | grep inet |  awk '{print $2}' | cut -d/ -f1 | paste -sd "," -)

go run ./hack/gen_cert.go --host "kubernetes,kubernetes.default,kubernetes.default.svc,kubernetes.default.svc.cluster.local,localhost,apiserver,kube-apiserver,$LOCAL_IP_LIST"  --ecdsa-curve P256 --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h