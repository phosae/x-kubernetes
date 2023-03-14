package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	tlsKeyName  = "tls.key"
	tlsCertName = "tls.crt"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Hello)
	mux.HandleFunc("/apis/hello.zeng.dev/v1beta1", Apis)
	// LIST /apis/hello.zeng.dev/v1beta1/namespaces/default/foos
	// GET  /apis/hello.zeng.dev/v1beta1/namespaces/default/foos/myfoo
	// POST /apis/hello.zeng.dev/v1beta1/namespaces/default/foos/
	// PUT  /apis/hello.zeng.dev/v1beta1/namespaces/default/foos/myfoo
	// DEL  /apis/hello.zeng.dev/v1beta1/namespaces/default/foos/myfoo

	// /apis/hello.zeng.dev/v1beta1/namespaces/default/foos/bar
	// /apis/hello.zeng.dev/v1beta1/namespaces/default/foos
	mux.HandleFunc("/apis/hello.zeng.dev/v1beta1/namespaces/", fooHandler)

	if certDir := os.Getenv("CERT_DIR"); certDir != "" {
		certFile := filepath.Join(certDir, tlsCertName)
		keyFile := filepath.Join(certDir, tlsKeyName)
		log.Println("serving https on 0.0.0.0:8443")
		log.Fatal(http.ListenAndServeTLS(":8443", certFile, keyFile, mux))
	} else {
		log.Println("serving http on 0.0.0.0:8000")
		log.Fatal(http.ListenAndServe(":8000", mux))
	}
}

const apis = `kind: APIResourceList
apiVersion: v1beta1
groupVersion: hello.zeng.dev/v1beta1
resources:
- name: foos
  singularName: ''
  namespaced: true
  kind: Foo
  verbs:
  - create
  - delete
  - get
  - list
  - update
  shortNames:
  - foo
  categories:
  - all`

func Apis(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	w.Write([]byte(apis))
}

// const kapiTpl = `apiVersion: hello.zeng.dev/v1beta1
// kind: Foo
// metadata:
//   name: {{.name}}
//   namespace:
// spec:
//   msg: 'hello world'`

type Foo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec struct {
		Msg string `json:"msg"`
	} `json:"spec"`
}

type FooList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Foo `json:"items"`
}

var foos = map[string]Foo{}

func init() {
	foos["bar"] = Foo{
		TypeMeta:   metav1.TypeMeta{APIVersion: "hello.zeng.dev/v1beta1", Kind: "Foo"},
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "bar"},
		Spec: struct {
			Msg string "json:\"msg\""
		}{
			Msg: "hello world",
		},
	}
}

func fooHandler(w http.ResponseWriter, r *http.Request) {
	nsResource := strings.TrimLeft(r.URL.Path, "/apis/hello.zeng.dev/v1beta1/namespaces/")
	parts := strings.Split(nsResource, "/")
	if len(parts) == 2 { // GET/POST default/foos
		switch r.Method {
		case http.MethodGet:
			GetAllFoos(w, r)
		case http.MethodPost:
			PostFoo(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else if len(parts) == 3 { // GET/PUT/DELETE default/foos/myfoo
		name := parts[2]
		switch r.Method {
		case http.MethodGet:
			GetFoo(w, r, name)
		case http.MethodPut:
			PutFoo(w, r, name)
		case http.MethodDelete:
			DeleteFoo(w, r, name)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// renderJSON renders 'v' as JSON and writes it as a response into w.
func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func GetFoo(w http.ResponseWriter, _ *http.Request, name string) {
	f, ok := foos[name]
	if !ok {
		http.Error(w, fmt.Sprintf("foo/%s not found", name), http.StatusInternalServerError)
		return
	}
	renderJSON(w, f)
}

func GetAllFoos(w http.ResponseWriter, _ *http.Request) {
	flist := FooList{
		TypeMeta: metav1.TypeMeta{Kind: "FooList", APIVersion: "hello.zeng.dev/v1beta1"},
	}
	for _, f := range foos {
		flist.Items = append(flist.Items, f)
	}
	renderJSON(w, flist)
}

func PostFoo(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var f Foo
	if err := dec.Decode(&f); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := foos[f.Name]; ok { // already exists
		w.WriteHeader(http.StatusConflict)
	}

	foos[f.Name] = f
}

func PutFoo(w http.ResponseWriter, r *http.Request, name string) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var f Foo
	if err := dec.Decode(&f); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := foos[name]; !ok { // not exists
		w.WriteHeader(http.StatusNotFound)
	}
	foos[f.Name] = f
}

func DeleteFoo(w http.ResponseWriter, _ *http.Request, name string) {
	if _, ok := foos[name]; !ok { // not exists
		w.WriteHeader(http.StatusNotFound)
	}
	delete(foos, name)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	rr, _ := httputil.DumpRequest(r, true)
	log.Println("rx", string(rr))

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"msg": "hello world"}`))
}
