/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package apiserver

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	clientrest "k8s.io/client-go/rest"

	hello "github.com/phosae/x-kubernetes/api-aggregation-lib/pkg/api/hello.zeng.dev"
	intallhello "github.com/phosae/x-kubernetes/api-aggregation-lib/pkg/api/hello.zeng.dev/install"
	intalltransformation "github.com/phosae/x-kubernetes/api-aggregation-lib/pkg/api/transformation/install"
	fooregistry "github.com/phosae/x-kubernetes/api-aggregation-lib/pkg/registry/hello.zeng.dev/foo"
)

var (
	// Scheme defines methods for serializing and deserializing API objects.
	Scheme = runtime.NewScheme()
	// Codecs provides methods for retrieving codecs and serializers for specific
	// versions and content types.
	Codecs = serializer.NewCodecFactory(Scheme)
)

func init() {
	intallhello.Install(Scheme)
	intalltransformation.Install(Scheme)

	// we need to add the options to empty v1
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Group: "", Version: "v1"})

	// TODO: keep the generic API server from wanting this
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	Scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)
}

// ExtraConfig holds custom apiserver config
type ExtraConfig struct {
	Rest              *clientrest.Config
	EnableEtcdStorage bool
}

// Config defines the config for the apiserver
type Config struct {
	GenericConfig *genericapiserver.RecommendedConfig
	ExtraConfig   ExtraConfig
}

// WardleServer contains state for a Kubernetes cluster master/api server.
type HelloApiServer struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

type completedConfig struct {
	GenericConfig genericapiserver.CompletedConfig
	ExtraConfig   *ExtraConfig
}

// CompletedConfig embeds a private pointer that cannot be instantiated outside of this package.
type CompletedConfig struct {
	*completedConfig
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() CompletedConfig {
	c := completedConfig{
		cfg.GenericConfig.Complete(),
		&cfg.ExtraConfig,
	}

	c.GenericConfig.Version = &version.Info{
		Major: "1",
		Minor: "0",
	}

	return CompletedConfig{&c}
}

// New returns a new instance of WardleServer from the given config.
func (c completedConfig) New() (*HelloApiServer, error) {
	genericServer, err := c.GenericConfig.New("hello.zeng.dev-apiserver", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	s := &HelloApiServer{
		GenericAPIServer: genericServer,
	}

	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(hello.GroupName, Scheme, metav1.ParameterCodec, Codecs)

	if !c.ExtraConfig.EnableEtcdStorage {
		restStorage := fooregistry.NewMemStore()
		v1storage := map[string]rest.Storage{"foos": restStorage}
		v2storage := map[string]rest.Storage{"foos": restStorage}
		apiGroupInfo.VersionedResourcesStorageMap["v1"] = v1storage
		apiGroupInfo.VersionedResourcesStorageMap["v2"] = v2storage
	} else {
		restStorage, err := fooregistry.NewREST(Scheme, c.GenericConfig.RESTOptionsGetter)
		if err != nil {
			return nil, err
		}

		v1storage := map[string]rest.Storage{"foos": restStorage.Foo, "foos/base64": restStorage.Base64}
		v2storage := map[string]rest.Storage{"foos": restStorage.Foo, "foos/config": restStorage.Config, "foos/status": restStorage.Status, "foos/base64": restStorage.Base64}
		apiGroupInfo.VersionedResourcesStorageMap["v1"] = v1storage
		apiGroupInfo.VersionedResourcesStorageMap["v2"] = v2storage
	}

	if err := s.GenericAPIServer.InstallAPIGroup(&apiGroupInfo); err != nil {
		return nil, err
	}

	return s, nil
}
