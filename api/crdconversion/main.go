package main

import (
	"bytes"
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
	"k8s.io/apimachinery/pkg/runtime/schema"
	kjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/runtime/serializer/versioning"

	helloinstall "github.com/phosae/x-kubernetes/api/crdconversion/internal/api/hello.zeng.dev/install"
)

var (
	scheme          = runtime.NewScheme()
	kjsonSerializer = kjson.NewSerializer(kjson.DefaultMetaFactory, scheme, scheme, false)

	tlsCertDir = ""
)

func init() {
	apix.Install(scheme)
	metav1.AddMetaToScheme(scheme)
	helloinstall.Install(scheme)

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

	desiredGV, err := schema.ParseGroupVersion(req.DesiredAPIVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to parse desired apiVersion: %v", err)
	}

	groupVersioner := schema.GroupVersions([]schema.GroupVersion{desiredGV})
	codec := versioning.NewCodec(
		kjsonSerializer,                       // decoder
		kjsonSerializer,                       // encoder
		runtime.UnsafeObjectConvertor(scheme), // convertor
		scheme,                                // creator
		scheme,                                // typer
		nil,                                   // defaulter
		groupVersioner,                        // encodeVersion
		runtime.InternalGroupVersioner,        // decodeVersion
		scheme.Name(),                         // originalSchemeName
	)

	convertedObjects := make([]runtime.RawExtension, len(req.Objects))
	for i, raw := range req.Objects {
		decodedObject, _, err := codec.Decode(raw.Raw, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to decode into apiVersion: %v", err)
		}
		buf := bytes.Buffer{}
		if err := codec.Encode(decodedObject, &buf); err != nil {
			return nil, fmt.Errorf("failed to convert to desired apiVersion: %v", err)
		}
		convertedObjects[i] = runtime.RawExtension{Raw: buf.Bytes()}
	}
	resp.ConvertedObjects = convertedObjects
	return &resp, nil
}
