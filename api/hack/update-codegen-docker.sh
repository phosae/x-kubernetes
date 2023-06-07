#!/usr/bin/env bash
set -eu
PROJECT_PACKAGE=github.com/phosae/x-kubernetes/api

docker run -it --rm -u root \
    -v ${PWD}/..:/go/src/github.com/phosae/x-kubernetes\
    -v ${PWD}/hack/custom-boilerplate.go.txt:/tmp/fake-boilerplate.txt \
    -v ${PWD}/api-rules/violation_exceptions.list:/tmp/validation-exceptions.list \
    -e GOPROXY="https://goproxy.cn"\
    -e PROJECT_PACKAGE=${PROJECT_PACKAGE} \
    -e CLIENT_GENERATOR_OUT=${PROJECT_PACKAGE}/generated \
    -e OPENAPI_GENERATOR_OUT=${PROJECT_PACKAGE}/generated \
    -e APIS_ROOT=${PROJECT_PACKAGE} \
    -e GENERATION_TARGETS="helpers,client,openapi" \
    -e OPENAPI_VIOLATION_EXCEPTIONS="/tmp/validation-exceptions.list" \
    -e WITH_APPLYCONFIG="true" \
    -e WITH_WATCH="" \
    zengxu/kube-code-generator:v1.28.0-alpha.1

docker run -it --rm -u root \
    -e GOPROXY="https://goproxy.cn"\
    -v ${PWD}/hack/custom-boilerplate.go.txt:/tmp/fake-boilerplate.txt \
    -v ${PWD}:/go/src/${PROJECT_PACKAGE}\
    -w /go/src/${PROJECT_PACKAGE}\
    zengxu/kube-code-generator:v1.28.0-alpha.1\
    update-generated-protobuf.sh github.com/phosae/x-kubernetes/api/hello.zeng.dev/v1,github.com/phosae/x-kubernetes/api/hello.zeng.dev/v2