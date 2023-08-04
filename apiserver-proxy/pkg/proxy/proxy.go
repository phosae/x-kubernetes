package proxy

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metainternalversionscheme "k8s.io/apimachinery/pkg/apis/meta/internalversion/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	restClient "k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	runtimecache "sigs.k8s.io/controller-runtime/pkg/cache"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

func New(config *restClient.Config, codecs serializer.CodecFactory, scheme *runtime.Scheme) (*Proxy, error) {
	transport, err := restClient.TransportFor(config)
	if err != nil {
		return nil, err
	}
	kubeurl, err := url.Parse(config.Host)
	if err != nil {
		return nil, err
	}

	reverseProxy := *httputil.NewSingleHostReverseProxy(kubeurl)
	reverseProxy.Transport = transport
	reverseProxy.ModifyResponse = nil
	reverseProxy.ErrorHandler = nil

	hc, err := restClient.HTTPClientFor(config)
	if err != nil {
		return nil, err
	}
	mapper, err := apiutil.NewDynamicRESTMapper(config, hc)
	if err != nil {
		return nil, err
	}
	rc, err := runtimecache.New(config, runtimecache.Options{
		HTTPClient: hc,
		Mapper:     mapper,
		Scheme:     scheme,
	})
	if err != nil {
		return nil, err
	}

	return &Proxy{
		codecs:         codecs,
		scheme:         scheme,
		restMapper:     mapper,
		rc:             rc,
		k8sProxy:       reverseProxy,
		jsonSerializer: json.NewSerializerWithOptions(json.DefaultMetaFactory, nil, nil, json.SerializerOptions{}),
		yamlSerializer: json.NewSerializerWithOptions(json.DefaultMetaFactory, nil, nil, json.SerializerOptions{Yaml: true}),
	}, nil
}

type Proxy struct {
	codecs     serializer.CodecFactory
	scheme     *runtime.Scheme
	restMapper meta.RESTMapper
	rc         runtimecache.Cache
	k8sProxy   httputil.ReverseProxy

	jsonSerializer *json.Serializer
	yamlSerializer *json.Serializer
}

func (p *Proxy) Start(ctx context.Context) {
	p.rc.Start(ctx)
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	info, _ := apirequest.RequestInfoFrom(ctx)

	switch info.Verb {
	case "get":
		if !info.IsResourceRequest {
			break
		}

		gvr := schema.GroupVersionResource{Group: info.APIGroup, Version: info.APIVersion, Resource: info.Resource}
		gvk, err := p.restMapper.KindFor(gvr)
		if err != nil {
			klog.Errorf("err Kind for %v: %s", gvr, err)
			break
		}
		klog.V(6).Infof("map gvr: %v to gvk: %v", gvr, gvk)

		var obj runtimeclient.Object
		runtimeObj, err := p.scheme.New(gvk)
		if err == nil {
			obj = runtimeObj.(runtimeclient.Object)
		} else {
			us := &unstructured.Unstructured{}
			us.SetGroupVersionKind(gvk)
			obj = us
		}
		err = p.rc.Get(ctx, types.NamespacedName{Namespace: info.Namespace, Name: info.Name}, obj)
		if err != nil {
			klog.Error("err get object from cache", err)
			break
		}
		klog.V(6).Infof("serve %s from cache", r.URL)
		responsewriters.WriteObjectNegotiated(p.codecs, negotiation.DefaultEndpointRestrictions, gvk.GroupVersion(), w, r, http.StatusOK, obj, false)
		return
	case "list":
		gvr := schema.GroupVersionResource{Group: info.APIGroup, Version: info.APIVersion, Resource: info.Resource}
		gvk, err := p.restMapper.KindFor(gvr)
		if err != nil {
			klog.Errorf("err Kind for %v: %s", gvr, err)
			break
		}
		klog.V(6).Infof("map gvr: %v to gvk: %v", gvr, gvk)

		opts := metav1.ListOptions{}
		if err := metainternalversionscheme.ParameterCodec.DecodeParameters(r.URL.Query(), metav1.SchemeGroupVersion, &opts); err != nil {
			err = apierrors.NewBadRequest(err.Error())
			responsewriters.ErrorNegotiated(err, p.codecs, gvk.GroupVersion(), w, r)
			return
		}

		var obj unstructured.UnstructuredList
		listgvk := schema.GroupVersionKind{Group: gvk.Group, Version: gvk.Version, Kind: gvk.Kind + "List"}
		obj.SetGroupVersionKind(listgvk)

		klog.V(6).Infof("serve %s from cache with options %v", r.URL.Path, opts)
		err = p.rc.List(ctx, &obj, &runtimeclient.ListOptions{Namespace: info.Namespace, Raw: &opts})
		if err != nil {
			status := responsewriters.ErrorToAPIStatus(err)
			responsewriters.WriteRawJSON(int(status.Code), status, w)
			return
		}

		if _, err := p.scheme.New(listgvk); err != nil {
			_, serializer, _ := negotiation.NegotiateOutputMediaType(r, p.codecs, negotiation.DefaultEndpointRestrictions)
			switch serializer.MediaType {
			case "":
				fallthrough
			case "application/json":
				w.Header().Add("Content-Type", "application/json")
				p.jsonSerializer.Encode(&obj, w)
			case "application/yaml":
				w.Header().Add("Content-Type", "application/yaml")
				p.yamlSerializer.Encode(&obj, w)
			default:
				status := metav1.Status{
					Status:  metav1.StatusFailure,
					Code:    http.StatusNotAcceptable,
					Reason:  metav1.StatusReasonNotAcceptable,
					Message: "only the following media types are accepted: application/json, application/yaml",
				}
				responsewriters.WriteRawJSON(int(status.Code), status, w)
			}
			return
		}

		responsewriters.WriteObjectNegotiated(p.codecs, negotiation.DefaultEndpointRestrictions, gvk.GroupVersion(), w, r, http.StatusOK, &obj, false)
		return
	default:
	}

	p.k8sProxy.ServeHTTP(w, r)
}
