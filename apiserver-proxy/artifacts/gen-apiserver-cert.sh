#!/bin/bash
set -eux

openssl req -new -nodes -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 \
    -subj "/O=system:masters/CN=kube-apiserver" \
    -keyout apiserver-proxy.key -out apiserver-proxy.csr

openssl x509 -req -in apiserver-proxy.csr -CA ca.crt -CAkey ca.key \
  -CAcreateserial -out apiserver-proxy.crt -days 3650 -extensions v3_req \
  -extfile <(printf "[v3_req]\nbasicConstraints=critical,CA:FALSE\nextendedKeyUsage=serverAuth\nsubjectAltName=DNS:localhost,DNS:kubernetes,DNS:kubernetes.default,DNS:kubernetes.default.svc,DNS:kubernetes.default.svc.cluster.localDNS:apiserver-proxy,DNS:apiserver-proxy.hello,DNS:apiserver-proxy.hello.svc,DNS:apiserver-proxy.hello.svc.cluster.local,IP:127.0.0.1,IP:172.18.0.2")