package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"

	apix "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/install"
	apixv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kjson "k8s.io/apimachinery/pkg/runtime/serializer/json"

	hellov1 "github.com/phosae/x-kubernetes/api/hello.zeng.dev/v1"
	hellov2 "github.com/phosae/x-kubernetes/api/hello.zeng.dev/v2"
)

var (
	scheme          = runtime.NewScheme()
	kjsonSerializer = kjson.NewSerializer(kjson.DefaultMetaFactory, scheme, scheme, false)

	tlsCertDir = ""
)

func init() {
	apix.Install(scheme)
	metav1.AddMetaToScheme(scheme)
	hellov1.AddToScheme(scheme)
	hellov2.AddToScheme(scheme)

	cert := os.Getenv("TLS_CERT")
	key := os.Getenv("TLS_KEY")

	if len(cert) > 0 && len(key) > 0 {
		certDir, err := os.MkdirTemp("", "tls")
		if err != nil {
			panic(err)
		}
		tlsCertDir = certDir

		writeTLSFileOrDie("tls.crt", cert)
		writeTLSFileOrDie("tls.key", key)
	}
}

func writeTLSFileOrDie(name string, base64content string) {
	content, err := base64.StdEncoding.DecodeString(base64content)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(tlsCertDir, name), content, 0666); err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/convert/hello.zeng.dev", func(w http.ResponseWriter, r *http.Request) {
		rr, _ := httputil.DumpRequest(r, true)
		log.Println("rx", string(rr))
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		Convert(w, r)
	})
	if len(tlsCertDir) > 0 {
		log.Println("serve https on :8443")
		log.Fatal(http.ListenAndServeTLS(":8443", filepath.Join(tlsCertDir, "tls.crt"), filepath.Join(tlsCertDir, "tls.key"), nil))
	}
	log.Println("serve http on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func Convert(w http.ResponseWriter, r *http.Request) {
	conversionReview := &apixv1.ConversionReview{}
	err := json.NewDecoder(r.Body).Decode(conversionReview)
	if err != nil {
		log.Println("failed to read conversion request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if conversionReview.Request == nil {
		log.Println("conversion request is nil")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := ConvertHello(conversionReview.Request)
	if err != nil {
		log.Printf("failed to convert, request/%s, err:%v\n", conversionReview.Request.UID, err)
		conversionReview.Response.Result = metav1.Status{
			Status:  metav1.StatusFailure,
			Message: err.Error(),
		}
	} else {
		resp.Result = metav1.Status{Status: metav1.StatusSuccess}
		conversionReview.Response = resp
	}
	conversionReview.Response.UID = conversionReview.Request.UID
	conversionReview.Request = nil

	err = json.NewEncoder(w).Encode(conversionReview)
	if err != nil {
		log.Println("failed to write response", err)
		return
	}
}

func ConvertHello(req *apixv1.ConversionRequest) (*apixv1.ConversionResponse, error) {
	resp := apixv1.ConversionResponse{}

	switch req.DesiredAPIVersion {
	default:
		return nil, fmt.Errorf("unsupported apiVersion/" + req.DesiredAPIVersion)
	case "hello.zeng.dev/v1":
		for _, o := range req.Objects {
			src, _, err := kjsonSerializer.Decode(o.Raw, nil, nil)
			if err != nil {
				return nil, err
			}

			switch in := src.(type) {
			case *hellov2.Foo:
				objv1 := &hellov1.Foo{TypeMeta: metav1.TypeMeta{Kind: "Foo", APIVersion: "hello.zeng.dev/v1"}}
				convertV2ToV1(in, objv1)
				resp.ConvertedObjects = append(resp.ConvertedObjects, runtime.RawExtension{Object: objv1})
			default:
				return nil, fmt.Errorf("unsupported type %v", in)
			}
		}
	case "hello.zeng.dev/v2":
		for _, o := range req.Objects {
			src, _, err := kjsonSerializer.Decode(o.Raw, nil, nil)
			if err != nil {
				return nil, fmt.Errorf("err decode ConversionReview.request.objects: %s", err)
			}

			switch in := src.(type) {
			case *hellov1.Foo:
				objv2 := &hellov2.Foo{TypeMeta: metav1.TypeMeta{Kind: "Foo", APIVersion: "hello.zeng.dev/v2"}}
				convertV1ToV2(in, objv2)
				resp.ConvertedObjects = append(resp.ConvertedObjects, runtime.RawExtension{Object: objv2})
			default:
				return nil, fmt.Errorf("unsupported type %v", in)
			}
		}
	}
	return &resp, nil
}

const AnnotationImage = "spec.image"

func convertV1ToV2(in *hellov1.Foo, out *hellov2.Foo) {
	out.ObjectMeta = in.ObjectMeta
	out.Spec.Image = out.Annotations[AnnotationImage]
	out.Spec.Config.Msg = in.Spec.Msg
	out.Spec.Config.Msg1 = in.Spec.Msg1
}

func convertV2ToV1(in *hellov2.Foo, out *hellov1.Foo) {
	out.ObjectMeta = in.ObjectMeta
	if out.Annotations == nil {
		out.Annotations = map[string]string{}
	}
	out.Annotations[AnnotationImage] = in.Spec.Image
	out.Spec.Msg = in.Spec.Config.Msg
	out.Spec.Msg1 = in.Spec.Config.Msg1
}
