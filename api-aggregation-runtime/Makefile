
LOCALBIN=$(PWD)/bin

code-gen:
	GOBIN=$(LOCALBIN) $(LOCALBIN)/apiserver-runtime-gen

build:
	go build -o ./bin/apiserver cmd/sample-apiserver/main.go