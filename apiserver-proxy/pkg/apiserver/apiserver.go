package apiserver

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/version"
	genericapiserver "k8s.io/apiserver/pkg/server"
	kubescheme "k8s.io/client-go/kubernetes/scheme"
	clientrest "k8s.io/client-go/rest"

	"github.com/phosae/x-kubernetes/apiserver-proxy/pkg/proxy"
)

var (
	// Scheme defines methods for serializing and deserializing API objects.
	Scheme = kubescheme.Scheme
	// Codecs provides methods for retrieving codecs and serializers for specific
	// versions and content types.
	Codecs = serializer.NewCodecFactory(Scheme)
)

// ExtraConfig holds custom apiserver config
type ExtraConfig struct {
	Rest *clientrest.Config
}

// Config defines the config for the apiserver
type Config struct {
	GenericConfig *genericapiserver.RecommendedConfig
	ExtraConfig   ExtraConfig
}

type completedConfig struct {
	GenericConfig genericapiserver.CompletedConfig
	ExtraConfig   *ExtraConfig
}

// CompletedConfig embeds a private pointer that cannot be instantiated outside of this package.
type CompletedConfig struct {
	*completedConfig
}

// ProxyApiServer contains state for a Kubernetes cluster master/api server.
type ProxyApiServer struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() CompletedConfig {
	c := completedConfig{
		cfg.GenericConfig.Complete(),
		&cfg.ExtraConfig,
	}

	c.GenericConfig.EnableIndex = false     // disable it so as to install our '/' handler
	c.GenericConfig.EnableDiscovery = false // disable it so as forward to kube-apiserver

	c.GenericConfig.Version = &version.Info{
		Major: "1",
		Minor: "0",
	}

	return CompletedConfig{&c}
}

// New returns a new instance of WardleServer from the given config.
func (c completedConfig) New() (*ProxyApiServer, error) {
	genericServer, err := c.GenericConfig.New("kube-apiserver-proxy", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	s := &ProxyApiServer{
		GenericAPIServer: genericServer,
	}

	pxy, err := proxy.New(c.ExtraConfig.Rest, Codecs, Scheme)
	if err != nil {
		return nil, err
	}

	s.GenericAPIServer.Handler.NonGoRestfulMux.HandlePrefix("/", pxy)
	s.GenericAPIServer.AddPostStartHookOrDie("start-cache-informers", func(ctx genericapiserver.PostStartHookContext) error {
		cctx, cancel := context.WithCancel(context.Background())
		go func() {
			<-ctx.StopCh
			cancel()
		}()
		pxy.Start(cctx)
		return nil
	})

	return s, nil
}
