#!/bin/sh
set -o errexit

IMAGE=${IMAGE:-kindest/node:v1.30.0}

# create registry container unless it already exists
reg_name='kind-registry'
reg_port='5000'
if [ "$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)" != 'true' ]; then
  docker run \
    -d --restart=always -p "127.0.0.1:${reg_port}:5000" --name "${reg_name}" \
    registry:2
fi

# create a cluster with the local registry enabled in containerd
cat << EOF | kind create cluster --config -
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: kind
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:${reg_port}"]
    endpoint = ["http://${reg_name}:5000"]
featureGates:
  "UserNamespacesSupport": true # beta 1.30
runtimeConfig:
  "api/all": true # enable all built-in APIs
nodes:
  - role: control-plane
    image: $IMAGE
    extraPortMappings:
    - containerPort: 30443 # expose apiserver-proxy on host
      hostPort: 30443
      protocol: TCP
  - role: worker
    image: $IMAGE
  - role: worker
    image: $IMAGE
networking:
  podSubnet: "10.244.0.0/16"
  serviceSubnet: "10.96.0.0/12"
EOF

# connect the registry to the cluster network if not already connected
if [ "$(docker inspect -f='{{json .NetworkSettings.Networks.kind}}' "${reg_name}")" = 'null' ]; then
  docker network connect "kind" "${reg_name}"
fi

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:${reg_port}"
    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF

kubectl taint node kind-control-plane node-role.kubernetes.io/control-plane:NoSchedule-