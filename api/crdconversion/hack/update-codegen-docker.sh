PROJECT_ROOT=github.com/phosae/x-kubernetes

docker run -it --rm -u root \
    -v ${PWD}/hack/custom-boilerplate.go.txt:/tmp/fake-boilerplate.txt \
    -v ${PWD}/../..:/go/src/${PROJECT_ROOT}\
    -e GOPROXY="https://goproxy.cn"\
    -e PROJECT_PACKAGE=${PROJECT_ROOT}/api/crdconversion \
    -e APIS_ROOT=${PROJECT_ROOT}/api/crdconversion/internal/api \
    -e GENERATION_TARGETS="helpers" \
    zengxu/kube-code-generator:v1.28.0-alpha.1