PROJECT_PACKAGE=github.com/phosae/x-kubernetes/api

docker run -it --rm -u root \
    -e GOPROXY="https://goproxy.cn"\
    -e CODEGEN_PKG=/go/src/k8s.io/code-generator\
    -v ${PWD}:/go/src/${PROJECT_PACKAGE}\
    -w /go/src/${PROJECT_PACKAGE}\
    zengxu/kube-code-generator:v1.27.1\
    ./hack/update-codegen.sh

docker run -it --rm -u root \
    -e GOPROXY="https://goproxy.cn"\
    -v ${PWD}/hack/custom-boilerplate.go.txt:/tmp/fake-boilerplate.txt \
    -v ${PWD}:/go/src/${PROJECT_PACKAGE}\
    -w /go/src/${PROJECT_PACKAGE}\
    zengxu/kube-code-generator:v1.27.1-pb\
    update-generated-protobuf.sh github.com/phosae/x-kubernetes/api/hello.zeng.dev/v1