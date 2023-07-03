package cmd

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	openapinamer "k8s.io/apiserver/pkg/endpoints/openapi"
	"k8s.io/apiserver/pkg/features"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/options"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/apiserver/pkg/storage/storagebackend"
	"k8s.io/apiserver/pkg/util/feature"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/term"
	"k8s.io/klog/v2"

	"github.com/phosae/x-kubernetes/api-aggregation-lib/pkg/admisssion/disallow"
	myapiserver "github.com/phosae/x-kubernetes/api-aggregation-lib/pkg/apiserver"
	generatedopenapi "github.com/phosae/x-kubernetes/api/generated/openapi"
	hellov1 "github.com/phosae/x-kubernetes/api/hello.zeng.dev/v1"
)

const defaultEtcdPathPrefix = "/registry/hello.zeng.dev"

type Options struct {
	SecureServing *genericoptions.SecureServingOptionsWithLoopback
	Kubeconfig    string
	Features      *genericoptions.FeatureOptions

	EnableEtcdStorage bool
	Etcd              *genericoptions.EtcdOptions

	EnableAuth     bool
	Authentication *genericoptions.DelegatingAuthenticationOptions
	Authorization  *genericoptions.DelegatingAuthorizationOptions

	EnableAdmission bool
	Admission       *genericoptions.AdmissionOptions
}

func (o *Options) Flags() (fs cliflag.NamedFlagSets) {
	msfs := fs.FlagSet("hello.zeng.dev-server")
	msfs.StringVar(&o.Kubeconfig, "kubeconfig", o.Kubeconfig, "The path to the kubeconfig used to connect to the Kubernetes API server (defaults to in-cluster config)")

	o.SecureServing.AddFlags(fs.FlagSet("apiserver secure serving"))
	o.Features.AddFlags(fs.FlagSet("features"))

	msfs.BoolVar(&o.EnableEtcdStorage, "enable-etcd-storage", false, "If true, store objects in etcd")
	o.Etcd.AddFlags(fs.FlagSet("Etcd"))

	msfs.BoolVar(&o.EnableAuth, "enable-auth", o.EnableAuth, "If true, enable authn and authz")
	o.Authentication.AddFlags(fs.FlagSet("apiserver authentication"))
	o.Authorization.AddFlags(fs.FlagSet("apiserver authorization"))

	msfs.BoolVar(&o.EnableAdmission, "enable-admission", o.EnableAdmission, "If true, enable admission plugins")
	return fs
}

// Complete fills in fields required to have valid data
func (o *Options) Complete() error {
	disallow.Register(o.Admission.Plugins)
	o.Admission.RecommendedPluginOrder = append(o.Admission.RecommendedPluginOrder, "DisallowFoo")
	return nil
}

// Validate validates ServerOptions
func (o Options) Validate(args []string) error {
	var errs []error
	if o.EnableEtcdStorage {
		errs = o.Etcd.Validate()
	}
	if o.EnableAuth {
		errs = append(errs, o.Authentication.Validate()...)
		errs = append(errs, o.Authorization.Validate()...)
	}
	return utilerrors.NewAggregate(errs)
}

type ServerConfig struct {
	Apiserver *genericapiserver.Config
	Rest      *rest.Config
}

func (o Options) ServerConfig() (*myapiserver.Config, error) {
	apiservercfg, err := o.ApiserverConfig()
	if err != nil {
		return nil, err
	}

	if o.EnableEtcdStorage {
		if err := o.Etcd.Complete(apiservercfg.Config.StorageObjectCountTracker, apiservercfg.Config.DrainedNotify(), apiservercfg.Config.AddPostStartHook); err != nil {
			return nil, err
		}

		// set apiservercfg's RESTOptionsGetter as StorageFactoryRestOptionsFactory{..., StorageFactory: DefaultStorageFactory}
		// like https://github.com/kubernetes/kubernetes/blob/e1ad9bee5bba8fbe85a6bf6201379ce8b1a611b1/cmd/kube-apiserver/app/server.go#L407-L415
		// DefaultStorageFactory#NewConfig provides a way to negotiate StorageSerializer/DeSerializer by Etcd.DefaultStorageMediaType option
		//
		// DefaultStorageFactory's NewConfig will be called by interface genericregistry.RESTOptionsGetter#GetRESTOptions (struct StorageFactoryRestOptionsFactory)
		// interface genericregistry.RESTOptionsGetter#GetRESTOptions will be called by genericregistry.Store#CompleteWithOptions
		// Finally all RESTBackend Options will be passed to genericregistry.Store implementations
		if o.Etcd.ApplyWithStorageFactoryTo(serverstorage.NewDefaultStorageFactory(
			o.Etcd.StorageConfig,
			o.Etcd.DefaultStorageMediaType,
			myapiserver.Codecs,
			serverstorage.NewDefaultResourceEncodingConfig(myapiserver.Scheme),
			apiservercfg.MergedResourceConfig,
			nil), &apiservercfg.Config); err != nil {
			return nil, err
		}
		klog.Infof("etcd cfg: %v", o.Etcd)
		o.Etcd.StorageConfig.Paging = utilfeature.DefaultFeatureGate.Enabled(features.APIListChunking)
		// apiservercfg.ClientConfig, err = o.restConfig()
		// if err != nil {
		// 	return nil, err
		// }
	}

	return &myapiserver.Config{
		GenericConfig: apiservercfg,
		ExtraConfig: myapiserver.ExtraConfig{
			EnableEtcdStorage: o.EnableEtcdStorage,
		},
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

	if o.EnableAuth {
		if err := o.Authentication.ApplyTo(&serverConfig.Authentication, serverConfig.SecureServing, nil); err != nil {
			return nil, err
		}
		if err := o.Authorization.ApplyTo(&serverConfig.Authorization); err != nil {
			return nil, err
		}
	}

	if o.EnableAdmission {
		(&options.CoreAPIOptions{}).ApplyTo(serverConfig) // init SharedInformerFactory

		// we can use LoopbackClientConfig for local resources
		// client, err := helloclientset.NewForConfig(serverConfig.LoopbackClientConfig)
		// informerFactory := helloinformers.NewSharedInformerFactory(client, serverConfig.LoopbackClientConfig.Timeout)
		// initializers := []admission.PluginInitializer{//} */

		o.Admission.ApplyTo(&serverConfig.Config, serverConfig.SharedInformerFactory, serverConfig.ClientConfig, feature.DefaultFeatureGate)
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
		// if just encode as json and store to etcd, just do this
		// Etcd: genericoptions.NewEtcdOptions(storagebackend.NewDefaultConfig(defaultEtcdPathPrefix, myapiserver.Codecs.LegacyCodec(hellov1.SchemeGroupVersion))),
		// but if we want to encode as json and pb, just assign nil to Codec here
		// like the official kube-apiserver https://github.com/kubernetes/kubernetes/blob/e1ad9bee5bba8fbe85a6bf6201379ce8b1a611b1/cmd/kube-apiserver/app/options/options.go#L96
		// when new/complete apiserver config, use EtcdOptions#ApplyWithStorageFactoryTo server.Config, which
		// finally init server.Config.RESTOptionsGetter as StorageFactoryRestOptionsFactory
		Etcd:           genericoptions.NewEtcdOptions(storagebackend.NewDefaultConfig(defaultEtcdPathPrefix, nil)),
		Authentication: genericoptions.NewDelegatingAuthenticationOptions(),
		Authorization:  genericoptions.NewDelegatingAuthorizationOptions(),
		Admission:      genericoptions.NewAdmissionOptions(),
	}
	opts.Etcd.StorageConfig.EncodeVersioner = runtime.NewMultiGroupVersioner(hellov1.SchemeGroupVersion, schema.GroupKind{Group: hellov1.GroupName})
	// opts.Etcd.DefaultStorageMediaType = "application/vnd.kubernetes.protobuf"
	opts.Etcd.DefaultStorageMediaType = "application/json"
	opts.SecureServing.BindPort = 6443

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
