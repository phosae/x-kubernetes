#!/usr/bin/env bash
set -eux

PROJECT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

# controller-gen schemapatch:manifests=./artifacts/crd paths=./... output:dir=./artifacts/crd
# controller-gen crd:crdVersions=v1 paths=./... output:dir=./artifacts/crd
## control OpenAPI description length by maxDescLen, crd:crdVersions=v1,maxDescLen=0 disable field description 
docker run -it --rm -u root \
    -v ${PROJECT_ROOT}/..:/src\
    -e GOPROXY="https://goproxy.cn"\
    -e GO_PROJECT_ROOT="/src"\
    -e CRD_FLAG="schemapatch:manifests=./api/artifacts/crd"\
    -e CRD_TYPES_PATH="/src/api"\
    -e CRD_OUT_PATH="/src/api/artifacts/crd"\
    zengxu/kube-code-generator:v1.28.0-alpha.1\
    update-crd.sh