package cmd

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/term"
	"k8s.io/klog/v2"

	"github.com/phosae/x-kubernetes/apiserver-proxy/pkg/apiserver"
)

type Options struct {
	SecureServing *genericoptions.SecureServingOptionsWithLoopback
	Kubeconfig    string
	Features      *genericoptions.FeatureOptions

	Authentication *genericoptions.DelegatingAuthenticationOptions
	Authorization  *genericoptions.DelegatingAuthorizationOptions
}

func (o *Options) Flags() (fs cliflag.NamedFlagSets) {
	msfs := fs.FlagSet("kube-apiserver-proxy")
	msfs.StringVar(&o.Kubeconfig, "kubeconfig", o.Kubeconfig, "The path to the kubeconfig used to connect to the Kubernetes API server (defaults to in-cluster config)")

	o.SecureServing.AddFlags(fs.FlagSet("apiserver secure serving"))
	o.Features.AddFlags(fs.FlagSet("features"))

	o.Authentication.AddFlags(fs.FlagSet("apiserver authentication"))
	o.Authorization.AddFlags(fs.FlagSet("apiserver authorization"))
	return fs
}

// Complete fills in fields required to have valid data
func (o *Options) Complete() error {
	return nil
}

// Validate validates ServerOptions
func (o Options) Validate(args []string) error {
	var errs []error
	errs = append(errs, o.Authentication.Validate()...)
	errs = append(errs, o.Authorization.Validate()...)
	return utilerrors.NewAggregate(errs)
}

type ServerConfig struct {
	Apiserver *genericapiserver.Config
	Rest      *rest.Config
}

func (o Options) ServerConfig() (*apiserver.Config, error) {
	apiservercfg, err := o.ApiserverConfig()
	if err != nil {
		return nil, err
	}

	r, err := o.restConfig()
	if err != nil {
		return nil, err
	}

	return &apiserver.Config{
		GenericConfig: apiservercfg,
		ExtraConfig: apiserver.ExtraConfig{
			Rest: r,
		},
	}, nil
}

func (o Options) ApiserverConfig() (*genericapiserver.RecommendedConfig, error) {
	if err := o.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	serverConfig := genericapiserver.NewRecommendedConfig(apiserver.Codecs)
	if err := o.SecureServing.ApplyTo(&serverConfig.SecureServing, &serverConfig.LoopbackClientConfig); err != nil {
		return nil, err
	}

	if err := o.Authentication.ApplyTo(&serverConfig.Authentication, serverConfig.SecureServing, nil); err != nil {
		return nil, err
	}
	if err := o.Authorization.ApplyTo(&serverConfig.Authorization); err != nil {
		return nil, err
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

// NewProxyServerCommand provides a CLI handler for the metrics server entrypoint
func NewProxyServerCommand(stopCh <-chan struct{}) *cobra.Command {
	opts := &Options{
		SecureServing:  genericoptions.NewSecureServingOptions().WithLoopback(),
		Authentication: genericoptions.NewDelegatingAuthenticationOptions(),
		Authorization:  genericoptions.NewDelegatingAuthorizationOptions(),
	}

	cmd := &cobra.Command{
		Short: "Launch kube-apiserver-proxy",
		Long:  "Launch kube-apiserver-proxy",
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
