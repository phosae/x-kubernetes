package main

import (
	"embed"
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
	mux.Handle("/", logHandler(http.NotFoundHandler()))
	mux.HandleFunc("/apis/hello.zeng.dev/v1", Apis)
	mux.HandleFunc("/openapi/v2", OpenapiV2)

	// LIST /apis/hello.zeng.dev/v1/namespaces/default/foos
	// GET  /apis/hello.zeng.dev/v1/namespaces/default/foos/myfoo
	// POST /apis/hello.zeng.dev/v1/namespaces/default/foos/
	// PUT  /apis/hello.zeng.dev/v1/namespaces/default/foos/myfoo
	// DEL  /apis/hello.zeng.dev/v1/namespaces/default/foos/myfoo
	mux.Handle("/apis/hello.zeng.dev/v1/namespaces/", logHandler(ContentTypeJSONHandler(http.HandlerFunc(fooHandler)))) // ends with '/' for prefix matching...

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

//go:embed docs/*
var embedFS embed.FS

// Get APIResourceList
//
//	@Summary		Get APIResourceList for group version 'hello.zeng.dev/v1'
//	@Description	List APIResource Info about group version 'hello.zeng.dev/v1'
//	@Produce		json
//	@Success		200	{string} apis
//	@Router			/apis/hello.zeng.dev/v1 [get]
func OpenapiV2(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json, _ := embedFS.ReadFile("docs/swagger.json")
	w.Write([]byte(json))
}

const apis = `{
  "kind": "APIResourceList",
  "apiVersion": "v1",
  "groupVersion": "hello.zeng.dev/v1",
  "resources": [
    {
      "name": "foos",
      "singularName": "",
      "namespaced": true,
      "kind": "Foo",
      "verbs": [
        "create",
        "delete",
        "get",
        "list",
        "update"
      ],
      "shortNames": [
        "foo"
      ],
      "categories": [
        "all"
      ]
    }
  ]
}`

// Get APIResourceList
//
//	@Summary		Get APIResourceList for group version 'hello.zeng.dev/v1'
//	@Description	List APIResource Info about group version 'hello.zeng.dev/v1'
//	@Produce		json
//	@Success		200	{string} apis
//	@Router			/apis/hello.zeng.dev/v1 [get]
func Apis(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(apis))
}

// Foo is some object like
//
//	`{
//	  "apiVersion": "hello.zeng.dev/v1",
//	  "kind": "Foo",
//	  "metadata": {
//	    "name": "%s",
//	    "namespace": "default"
//	  },
//	  "spec": {
//	    "msg": "%s"
//	  }
//	}`
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
		TypeMeta:   metav1.TypeMeta{APIVersion: "hello.zeng.dev/v1", Kind: "Foo"},
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "bar"},
		Spec: struct {
			Msg string "json:\"msg\""
		}{
			Msg: "hello world",
		},
	}
}

func fooHandler(w http.ResponseWriter, r *http.Request) {
	nsResource := strings.TrimLeft(r.URL.Path, "/apis/hello.zeng.dev/v1/namespaces/")
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

	w.Write(js)
}

const kstatusTmplate = `{
	"kind":"Status",
	"apiVersion":"v1",
	"metadata":{},
	"status":"Failure",
	"message":"%s",
	"reason":"%s",
	"details":{"name":"%s","kind":"foos"},
	"code":%d
}`

func writeErrStatus(w http.ResponseWriter, name string, status int) {
	var errStatus string
	switch status {
	case http.StatusNotFound:
		errStatus = fmt.Sprintf(kstatusTmplate, fmt.Sprintf(`foos '%s' not found`, name), http.StatusText(http.StatusNotFound), name, http.StatusNotFound)
	case http.StatusConflict:
		errStatus = fmt.Sprintf(kstatusTmplate, fmt.Sprintf(`foos '%s' already exists`, name), http.StatusText(http.StatusConflict), name, http.StatusConflict)
	default:
		errStatus = "{}"
	}
	w.Write([]byte(errStatus))
	w.WriteHeader(status)
}

// GetFoo swag doc
// @Summary      Get an Foo Object
// @Description  Get an Foo by Resource Name
// @Tags         foos
// @Produce      json
// @Param        namespace	path	string  true  "Namepsace"
// @Param        name	path	string  true  "Resource Name"
// @Success      200  {object}  Foo
// @Router       /apis/hello.zeng.dev/v1/namespaces/{namespace}/foos/{name} [get]
func GetFoo(w http.ResponseWriter, _ *http.Request, name string) {
	f, ok := foos[name]
	if !ok {
		writeErrStatus(w, name, http.StatusNotFound)
		return
	}
	renderJSON(w, f)
}

// GetAllFoos swag doc
// @Summary      List all Foos
// @Description  List all Foos
// @Tags         foos
// @Produce      json
// @Param        namespace	path	string  true  "Namepsace"
// @Success      200  {object}  Foo
// @Router       /apis/hello.zeng.dev/v1/namespaces/{namespace}/foos [get]
func GetAllFoos(w http.ResponseWriter, _ *http.Request) {
	flist := FooList{
		TypeMeta: metav1.TypeMeta{Kind: "FooList", APIVersion: "hello.zeng.dev/v1"},
	}
	for _, f := range foos {
		flist.Items = append(flist.Items, f)
	}
	renderJSON(w, flist)
}

// PostFoo swag doc
// @Summary      Create a Foo Object
// @Description  Create a Foo Object
// @Tags         foos
// @Consume      json
// @Produce      json
// @Param        namespace	path	string  true  "Namepsace"
// @Success      201  {object}  Foo
// @Router       /apis/hello.zeng.dev/v1/namespaces/{namespace}/foos [post]
func PostFoo(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var f Foo
	if err := dec.Decode(&f); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := foos[f.Name]; ok { // already exists
		writeErrStatus(w, f.Name, http.StatusConflict)
		return
	}

	foos[f.Name] = f
	w.WriteHeader(http.StatusCreated)
	renderJSON(w, f) // follow official API, return the created object
}

// PutFoo swag doc
// @Summary      Replace a Foo Object
// @Description  Replace a Foo Object by Creation or Update
// @Tags         foos
// @Consume      json
// @Produce      json
// @Param        namespace	path	string  true  "Namepsace"
// @Param        name	path	string  true  "Resource Name"
// @Success      201  {object}  Foo	"created"
// @Success      200  {object}  Foo "updated"
// @Router       /apis/hello.zeng.dev/v1/namespaces/{namespace}/foos/{name} [put]
func PutFoo(w http.ResponseWriter, r *http.Request, name string) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var f Foo
	if err := dec.Decode(&f); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := foos[name]; !ok { // not exists
		w.WriteHeader(http.StatusCreated)
	}
	foos[f.Name] = f
	renderJSON(w, f) // follow official API, return the replacement
}

// DeleteFoo swag doc
// @Summary      Delete a Foo Object
// @Description  Delete a Foo Object by name in some Namespace
// @Tags         foos
// @Produce      json
// @Param        namespace	path	string  true  "Namepsace"
// @Param        name	path	string  true  "Resource Name"
// @Success      200  {object}  Foo "deleted"
// @Router       /apis/hello.zeng.dev/v1/namespaces/{namespace}/foos/{name} [delete]
func DeleteFoo(w http.ResponseWriter, _ *http.Request, name string) {
	if f, ok := foos[name]; !ok { // not exists
		writeErrStatus(w, name, http.StatusNotFound)
		return
	} else {
		delete(foos, name)
		now := metav1.Now()
		var noWait int64 = 0
		f.DeletionTimestamp = &now
		f.DeletionGracePeriodSeconds = &noWait
		renderJSON(w, f) // follow official API, return the deleted object
	}
}

func logHandler(ha http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rr, _ := httputil.DumpRequest(r, true)
		log.Println("rx", string(rr))
		ha.ServeHTTP(w, r)
	})
}

func ContentTypeJSONHandler(ha http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ha.ServeHTTP(w, r)
	})
}
