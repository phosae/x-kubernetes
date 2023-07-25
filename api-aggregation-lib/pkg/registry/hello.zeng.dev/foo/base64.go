package foo

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"

	hello "github.com/phosae/x-kubernetes/api-aggregation-lib/pkg/api/hello.zeng.dev"
	transformationapi "github.com/phosae/x-kubernetes/api-aggregation-lib/pkg/api/transformation"
	transformationv1beta1 "github.com/phosae/x-kubernetes/api/transformation/v1beta1"
)

type getter interface {
	Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error)
}

var _ = rest.NamedCreater(&Base64REST{})
var _ = rest.GroupVersionKindProvider(&Base64REST{})

type Base64REST struct {
	foos   getter
	scheme runtime.Scheme
}

func NewBase64REST(fooGetter getter, s *runtime.Scheme) *Base64REST {
	return &Base64REST{foos: fooGetter, scheme: *s}
}

func (r *Base64REST) New() runtime.Object {
	return &transformationapi.Base64{}
}

// Destroy cleans up resources on shutdown.
func (r *Base64REST) Destroy() {
	// Given no underlying store, we don't destroy anything
	// here explicitly.
}

var gvk = schema.GroupVersionKind{
	Group:   transformationv1beta1.SchemeGroupVersion.Group,
	Version: transformationv1beta1.SchemeGroupVersion.Version,
	Kind:    "Base64",
}

func (r *Base64REST) Create(ctx context.Context, name string, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	req := obj.(*transformationapi.Base64)

	// Get the namespace from the context (populated from the URL).
	namespace, ok := genericapirequest.NamespaceFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("namespace is required")
	}

	info, ok := genericapirequest.RequestInfoFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("request info is required")
	}

	// require name/namespace in the body to match URL if specified
	if len(req.Name) > 0 && req.Name != name {
		errs := field.ErrorList{field.Invalid(field.NewPath("metadata").Child("name"), req.Name, "must match the foo name if specified")}
		return nil, errors.NewInvalid(gvk.GroupKind(), name, errs)
	}
	if len(req.Namespace) > 0 && req.Namespace != namespace {
		errs := field.ErrorList{field.Invalid(field.NewPath("metadata").Child("namespace"), req.Namespace, "must match the foo namespace if specified")}
		return nil, errors.NewInvalid(gvk.GroupKind(), name, errs)
	}

	// Lookup foo
	fooObj, err := r.foos.Get(ctx, name, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	foo := fooObj.(*hello.Foo)
	groupVersioner := schema.GroupVersions([]schema.GroupVersion{{Group: info.APIGroup, Version: info.APIVersion}})
	target := unstructured.Unstructured{}
	if err != r.scheme.Convert(foo, &target, groupVersioner) {
		return nil, errors.NewBadRequest(fmt.Sprintf("unknown request version %s, %s", groupVersioner, err))
	}

	targetBytes, err := target.MarshalJSON()
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	if len(req.Spec.FieldPath) > 0 {
		fields := strings.Split(req.Spec.FieldPath, ".")
		v, ok, err := unstructured.NestedFieldNoCopy(target.Object, fields...)
		if !ok || err != nil {
			return nil, errors.NewBadRequest(fmt.Sprintf("undefined field path %s", req.Spec.FieldPath))
		}
		targetBytes, err = json.Marshal(v)
		if err != nil {
			return nil, errors.NewBadRequest(fmt.Sprintf("unexpected nested object %s", err))
		}
	}

	out := req.DeepCopy()
	out.Status = transformationapi.Base64Status{
		Output: base64.StdEncoding.EncodeToString(targetBytes),
	}
	return out, nil
}

func (r *Base64REST) GroupVersionKind(schema.GroupVersion) schema.GroupVersionKind {
	return gvk
}
