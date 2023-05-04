package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kruntime "k8s.io/apimachinery/pkg/runtime"
	kstrategicpatch "k8s.io/apimachinery/pkg/util/strategicpatch"
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

	// LIST /apis/hello.zeng.dev/v1/foos
	// LIST /apis/hello.zeng.dev/v1/namespaces/{namespace}/foos
	// GET  /apis/hello.zeng.dev/v1/namespaces/{namespace}/foos/{name}
	// POST /apis/hello.zeng.dev/v1/namespaces/{namespace}/foos/
	// PUT  /apis/hello.zeng.dev/v1/namespaces/{namespace}/foos/{name}
	// DEL  /apis/hello.zeng.dev/v1/namespaces/{namespace}/foos/{name}
	mux.Handle("/apis/hello.zeng.dev/v1/", logHandler(ContentTypeJSONHandler(http.HandlerFunc(fooHandler)))) // ends with '/' for prefix matching...

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
//	@Router			/openapi/v2 [get]
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
//	    "msg": "%s",
//	    "msg1": "%s"
//	  }
//	}`
type Foo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec struct {
		Msg  string `json:"msg"`
		Msg1 string `json:"msg1"`
	} `json:"spec"`
}

type FooList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Foo `json:"items"`
}

var foos = map[string]Foo{}

func init() {
	foos["default/bar"] = Foo{
		TypeMeta:   metav1.TypeMeta{APIVersion: "hello.zeng.dev/v1", Kind: "Foo"},
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "bar"},
		Spec: struct {
			Msg  string "json:\"msg\""
			Msg1 string "json:\"msg1\""
		}{
			Msg: "hello world",
		},
	}
}

type ctxkey string

func fooHandler(w http.ResponseWriter, r *http.Request) {
	nsResource := strings.TrimPrefix(r.URL.Path, "/apis/hello.zeng.dev/v1/namespaces/")

	if nsResource == r.URL.Path && r.URL.Path == "/apis/hello.zeng.dev/v1/foos" {
		GetAllFoos(w, r)
		return
	}

	parts := strings.Split(nsResource, "/")
	if len(parts) == 2 { // GET/POST {namespace}/foos
		r = r.WithContext(context.WithValue(r.Context(), ctxkey("namespace"), parts[0]))
		switch r.Method {
		case http.MethodGet:
			GetAllFoosInNamespace(w, r)
		case http.MethodPost:
			PostFoo(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else if len(parts) == 3 { // GET/PUT/DELETE {namespace}/foos/{fooname}
		r = r.WithContext(context.WithValue(r.Context(), ctxkey("namespace"), parts[0]))
		name := parts[2]
		switch r.Method {
		case http.MethodGet:
			GetFoo(w, r, name)
		case http.MethodPut:
			PutFoo(w, r, name)
		case http.MethodPatch:
			PatchFoo(w, r, name)
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
		writeErrStatus(w, "", http.StatusInternalServerError, err.Error())
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

func writeErrStatus(w http.ResponseWriter, name string, status int, msg string) {
	var errStatus string
	switch status {
	case http.StatusNotFound:
		errStatus = fmt.Sprintf(kstatusTmplate, fmt.Sprintf(`foos '%s' not found`, name), http.StatusText(http.StatusNotFound), name, http.StatusNotFound)
	case http.StatusConflict:
		errStatus = fmt.Sprintf(kstatusTmplate, fmt.Sprintf(`foos '%s' already exists`, name), http.StatusText(http.StatusConflict), name, http.StatusConflict)
	default:
		errStatus = fmt.Sprintf(kstatusTmplate, msg, http.StatusText(status), name, status)
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
func GetFoo(w http.ResponseWriter, r *http.Request, name string) {
	ns := r.Context().Value(ctxkey("namespace"))
	nsname := fmt.Sprintf("%s/%s", ns, name)

	f, ok := foos[nsname]
	if !ok {
		writeErrStatus(w, nsname, http.StatusNotFound, "")
		return
	}
	renderJSON(w, f)
}

// GetAllFoosInNamespace swag doc
// @Summary      List all Foos in some namespace
// @Description  List all Foos in some namespace
// @Tags         foos
// @Produce      json
// @Param        namespace	path	string  true  "Namepsace"
// @Success      200  {object}  FooList
// @Router       /apis/hello.zeng.dev/v1/namespaces/{namespace}/foos [get]
func GetAllFoosInNamespace(w http.ResponseWriter, r *http.Request) {
	flist := FooList{
		TypeMeta: metav1.TypeMeta{Kind: "FooList", APIVersion: "hello.zeng.dev/v1"},
	}
	for _, f := range foos {
		if f.Namespace == r.Context().Value(ctxkey("namespace")) {
			flist.Items = append(flist.Items, f)
		}
	}
	renderJSON(w, flist)
}

// GetAllFoos swag doc
// @Summary      List all Foos
// @Description  List all Foos
// @Tags         foos
// @Produce      json
// @Success      200  {object}  FooList
// @Router       /apis/hello.zeng.dev/v1/foos [get]
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
		writeErrStatus(w, "", http.StatusBadRequest, err.Error())
		return
	}

	ns := r.Context().Value(ctxkey("namespace"))
	nsname := fmt.Sprintf("%s/%s", ns, f.Name)

	if _, ok := foos[nsname]; ok { // already exists
		writeErrStatus(w, nsname, http.StatusConflict, "")
		return
	}

	foos[nsname] = f
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
	ns := r.Context().Value(ctxkey("namespace"))
	nsname := fmt.Sprintf("%s/%s", ns, name)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var f Foo
	if err := dec.Decode(&f); err != nil {
		writeErrStatus(w, nsname, http.StatusBadRequest, err.Error())
		return
	}

	if _, ok := foos[nsname]; !ok { // not exists
		w.WriteHeader(http.StatusCreated)
	}
	foos[nsname] = f
	renderJSON(w, f) // follow official API, return the replacement
}

// PatchFoo swag doc
// @Summary      partially update the specified Foo
// @Description  partially update the specified Foo
// @Tags         foos
// @Consume      json
// @Produce      json
// @Param        namespace	path	string  true  "Namepsace"
// @Param        name	path	string  true  "Resource Name"
// @Success      200  {object}  Foo "OK"
// @Router       /apis/hello.zeng.dev/v1/namespaces/{namespace}/foos/{name} [patch]
func PatchFoo(w http.ResponseWriter, r *http.Request, name string) {
	ns := r.Context().Value(ctxkey("namespace"))
	nsname := fmt.Sprintf("%s/%s", ns, name)

	patchBytes, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrStatus(w, nsname, http.StatusBadRequest, err.Error())
		return
	}

	var originalObjMap map[string]interface{}
	var originalBytes []byte
	var schema Foo
	if old, ok := foos[nsname]; !ok { // not exists
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		schema = old
		originalBytes, _ = json.Marshal(old)
		originalObjMap, _ = kruntime.DefaultUnstructuredConverter.ToUnstructured(&old)
	}

	var patchedFoo []byte
	switch r.Header.Get("Content-Type") {
	case "application/merge-patch+json":
		patchedFoo, err = jsonpatch.MergePatch(originalBytes, patchBytes)
		if err != nil {
			writeErrStatus(w, nsname, http.StatusBadRequest, err.Error())
			return
		}
	case "application/json-patch+json":
		patch, err := jsonpatch.DecodePatch(patchBytes)
		if err != nil {
			writeErrStatus(w, nsname, http.StatusBadRequest, err.Error())
			return
		}
		patchedFoo, err = patch.Apply(originalBytes)
		if err != nil {
			writeErrStatus(w, nsname, http.StatusBadRequest, err.Error())
			return
		}
	case "application/strategic-merge-patch+json":
		var patchMap map[string]interface{}
		if err = json.Unmarshal(patchBytes, &patchMap); err != nil {
			writeErrStatus(w, nsname, http.StatusBadRequest, err.Error())
			return
		}
		fmt.Println("orig map", originalObjMap)
		if patchedObjMap, err := kstrategicpatch.StrategicMergeMapPatch(originalObjMap, patchMap, schema); err != nil {
			writeErrStatus(w, nsname, http.StatusBadRequest, err.Error())
			return
		} else {
			var theFoo Foo
			fmt.Println("patched map", patchedObjMap)
			if err = kruntime.DefaultUnstructuredConverter.FromUnstructuredWithValidation(patchedObjMap, &theFoo, false); err != nil {
				writeErrStatus(w, nsname, http.StatusBadRequest, err.Error())
				return
			} else {
				patchedFoo, _ = json.Marshal(theFoo)
			}
		}
	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(bytes.NewReader(patchedFoo))
	dec.DisallowUnknownFields()
	var f Foo
	if err := dec.Decode(&f); err != nil {
		writeErrStatus(w, nsname, http.StatusBadRequest, err.Error())
		return
	}
	foos[nsname] = f
	renderJSON(w, f)
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
func DeleteFoo(w http.ResponseWriter, r *http.Request, name string) {
	ns := r.Context().Value(ctxkey("namespace"))
	nsname := fmt.Sprintf("%s/%s", ns, name)

	if f, ok := foos[nsname]; !ok { // not exists
		writeErrStatus(w, nsname, http.StatusNotFound, "")
		return
	} else {
		delete(foos, nsname)
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
