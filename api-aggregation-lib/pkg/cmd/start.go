package cmd

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
	openapinamer "k8s.io/apiserver/pkg/endpoints/openapi"
	"k8s.io/apiserver/pkg/features"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/term"
	"k8s.io/klog/v2"

	myapiserver "github.com/phosae/x-kubernetes/api-aggregation-lib/pkg/apiserver"
	generatedopenapi "github.com/phosae/x-kubernetes/api/generated/openapi"
)

type Options struct {
	// RecommendedOptions *genericoptions.RecommendedOptions // - EtcdOptions
	SecureServing *genericoptions.SecureServingOptionsWithLoopback
	Kubeconfig    string
	Features      *genericoptions.FeatureOptions
}

func (o *Options) Flags() (fs cliflag.NamedFlagSets) {
	msfs := fs.FlagSet("hello.zeng.dev-server")
	msfs.StringVar(&o.Kubeconfig, "kubeconfig", o.Kubeconfig, "The path to the kubeconfig used to connect to the Kubernetes API server and the Kubelets (defaults to in-cluster config)")

	o.SecureServing.AddFlags(fs.FlagSet("apiserver secure serving"))
	o.Features.AddFlags(fs.FlagSet("features"))

	return fs
}

// Complete fills in fields required to have valid data
func (o *Options) Complete() error { return nil }

// Validate validates ServerOptions
func (o Options) Validate(args []string) error { return nil }

type ServerConfig struct {
	Apiserver *genericapiserver.Config
	Rest      *rest.Config
}

func (o Options) ServerConfig() (*myapiserver.Config, error) {
	apiservercfg, err := o.ApiserverConfig()
	if err != nil {
		return nil, err
	}

	//apiservercfg.ClientConfig, err = o.restConfig()
	//if err != nil {
	//	return nil, err
	//}
	return &myapiserver.Config{
		GenericConfig: apiservercfg,
		ExtraConfig:   myapiserver.ExtraConfig{},
	}, nil
}

func (o Options) ApiserverConfig() (*genericapiserver.RecommendedConfig, error) {
	if err := o.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	serverConfig := genericapiserver.NewRecommendedConfig(myapiserver.Codecs)
	if err := o.SecureServing.ApplyTo(&serverConfig.SecureServing, &serverConfig.LoopbackClientConfig); err != nil {
		return nil, err
	}

	// enable OpenAPI schemas
	namer := openapinamer.NewDefinitionNamer(myapiserver.Scheme)
	serverConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(generatedopenapi.GetOpenAPIDefinitions, namer)
	serverConfig.OpenAPIConfig.Info.Title = "hello.zeng.dev-server"
	serverConfig.OpenAPIConfig.Info.Version = "0.1"

	if utilfeature.DefaultFeatureGate.Enabled(features.OpenAPIV3) {
		serverConfig.OpenAPIV3Config = genericapiserver.DefaultOpenAPIV3Config(generatedopenapi.GetOpenAPIDefinitions, namer)
		serverConfig.OpenAPIV3Config.Info.Title = "hello.zeng.dev-server"
		serverConfig.OpenAPIV3Config.Info.Version = "0.1"
	}

	return serverConfig, nil
}

func (o Options) restConfig() (*rest.Config, error) {
	var config *rest.Config
	var err error
	if len(o.Kubeconfig) > 0 {
		loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: o.Kubeconfig}
		loader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})

		config, err = loader.ClientConfig()
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, fmt.Errorf("unable to construct lister client config: %v", err)
	}
	// Use protobufs for communication with apiserver
	config.ContentType = "application/vnd.kubernetes.protobuf"
	rest.SetKubernetesDefaults(config)
	return config, err
}

// NewHelloServerCommand provides a CLI handler for the metrics server entrypoint
func NewHelloServerCommand(stopCh <-chan struct{}) *cobra.Command {
	opts := &Options{
		SecureServing: genericoptions.NewSecureServingOptions().WithLoopback(),
	}

	cmd := &cobra.Command{
		Short: "Launch hello.zeng.dev-server",
		Long:  "Launch hello.zeng.dev-server",
		RunE: func(c *cobra.Command, args []string) error {
			if err := opts.Complete(); err != nil {
				return err
			}
			if err := opts.Validate(args); err != nil {
				return err
			}
			if err := runCommand(opts, stopCh); err != nil {
				return err
			}
			return nil
		},
	}

	fs := cmd.Flags()
	nfs := opts.Flags()
	for _, f := range nfs.FlagSets {
		fs.AddFlagSet(f)
	}
	local := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	klog.InitFlags(local)
	nfs.FlagSet("logging").AddGoFlagSet(local)

	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStderr(), nfs, cols)
		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStdout(), nfs, cols)
	})
	return cmd
}

func runCommand(o *Options, stopCh <-chan struct{}) error {
	servercfg, err := o.ServerConfig()
	if err != nil {
		return err
	}

	server, err := servercfg.Complete().New()
	if err != nil {
		return err
	}

	return server.GenericAPIServer.PrepareRun().Run(stopCh)
}
