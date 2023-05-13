#!/usr/bin/env bash
set -eux

PROJECT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

# controller-gen crd:crdVersions=v1 paths=./... output:dir=./artifacts/crd
# controller-gen schemapatch:manifests=./artifacts/crd paths=./... output:dir=./artifacts/crd
docker run -it --rm -u root \
    -v ${PROJECT_ROOT}/..:/src\
    -e GOPROXY="https://goproxy.cn"\
    -e GO_PROJECT_ROOT="/src"\
    -e CRD_FLAG="schemapatch:manifests=./api/artifacts/crd"\
    -e CRD_TYPES_PATH="/src/api"\
    -e CRD_OUT_PATH="/src/api/artifacts/crd"\
    zengxu/kube-code-generator:v1.27.1-pb\
    update-crd.sh