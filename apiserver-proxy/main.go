package main

import (
	"os"

	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/cli"

	"github.com/phosae/x-kubernetes/apiserver-proxy/pkg/cmd"
)

func main() {
	stopCh := genericapiserver.SetupSignalHandler()
	cmd := cmd.NewProxyServerCommand(stopCh)
	code := cli.Run(cmd)
	os.Exit(code)
}
