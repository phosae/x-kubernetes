# docker build -t zengxu/k8s-api .
FROM registry.k8s.io/etcd:3.5.7-0 as etcd-bin
FROM registry.k8s.io/kube-apiserver:v1.26.2 as kube-bin

FROM ubuntu:jammy

COPY --from=etcd-bin /usr/local/bin/etcd /usr/local/bin/etcd
COPY --from=etcd-bin /usr/local/bin/etcdctl /usr/local/bin/etcdctl
COPY --from=kube-bin /usr/local/bin/kube-apiserver /usr/local/bin/kube-apiserver

