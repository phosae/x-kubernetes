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
	"sync"
	"time"

	jsonpatch "github.com/evanphx/json-patch"
	gnosticopenapiv2 "github.com/google/gnostic/openapiv2"
	"google.golang.org/protobuf/proto"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kruntime "k8s.io/apimachinery/pkg/runtime"
	durationutil "k8s.io/apimachinery/pkg/util/duration"
	kstrategicpatch "k8s.io/apimachinery/pkg/util/strategicpatch"
)

const (
	tlsKeyName  = "apiserver.key"
	tlsCertName = "apiserver.crt"
)

// @title           hello.zeng.dev-server
// @version         0.1
// @description     K8s apiserver style http server from scratch
// @BasePath  /apis
func main() {
	mux := BuildMux()
	if certDir := os.Getenv("CERT_DIR"); certDir != "" {
		certFile := filepath.Join(certDir, tlsCertName)
		keyFile := filepath.Join(certDir, tlsKeyName)
		log.Println("serving https on 0.0.0.0:6443")
		log.Fatal(http.ListenAndServeTLS(":6443", certFile, keyFile, mux))
	} else {
		log.Println("serving http on 0.0.0.0:8000")
		log.Fatal(http.ListenAndServe(":8000", mux))
	}
}

func BuildMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", logHandler(http.NotFoundHandler()))
	mux.Handle("/apis", logHandler(http.HandlerFunc(APIs)))
	mux.Handle("/apis/hello.zeng.dev", logHandler(http.HandlerFunc(APIGroupHelloV1)))
	mux.Handle("/apis/hello.zeng.dev/v1", logHandler(http.HandlerFunc(APIGroupHelloV1Resources)))
	mux.Handle("/openapi/v2", logHandler(http.HandlerFunc(OpenapiV2)))

	// LIST /apis/hello.zeng.dev/v1/foos
	// LIST /apis/hello.zeng.dev/v1/namespaces/{namespace}/foos
	// GET  /apis/hello.zeng.dev/v1/namespaces/{namespace}/foos/{name}
	// POST /apis/hello.zeng.dev/v1/namespaces/{namespace}/foos/
	// PUT  /apis/hello.zeng.dev/v1/namespaces/{namespace}/foos/{name}
	// DEL  /apis/hello.zeng.dev/v1/namespaces/{namespace}/foos/{name}
	mux.Handle("/apis/hello.zeng.dev/v1/", logHandler(ContentTypeJSONHandler(http.HandlerFunc(fooHandler)))) // ends with '/' for prefix matching...
	return mux
}

var apis = metav1.APIGroupList{
	TypeMeta: metav1.TypeMeta{
		Kind:       "APIGroupList",
		APIVersion: "v1",
	},
	Groups: []metav1.APIGroup{
		{
			TypeMeta: metav1.TypeMeta{
				Kind:       "APIGroup",
				APIVersion: "v1",
			},
			Name: "hello.zeng.dev",
			Versions: []metav1.GroupVersionForDiscovery{
				{
					GroupVersion: "hello.zeng.dev/v1",
					Version:      "v1",
				},
			},
			PreferredVersion: metav1.GroupVersionForDiscovery{GroupVersion: "hello.zeng.dev/v1", Version: "v1"},
		},
	},
}

var apidiscoveries = `{
	"apiVersion": "apidiscovery.k8s.io/v2beta1",
	"kind": "APIGroupDiscoveryList",
	"metadata": {},
	"items": [
	  {
		"metadata": {
		  "name": "hello.zeng.dev"
		},
		"versions": [
		  {
			"version": "v1",
			"resources": [
			  {
				"resource": "foos",
				"responseKind": {
				  "group": "hello.zeng.dev",
				  "kind": "Foo",
				  "version": "v1"
				},
				"scope": "Namespaced",
				"shortNames": [
				  "fo"
				],
				"singularResource": "foo",
				"verbs": [
				  "delete",
				  "get",
				  "list",
				  "patch",
				  "create",
				  "update"
				]
			  }
			]
		  }
		]
	  }
	]
  }`

// List APIGroups
//
//	@Summary		List all APIGroups of this apiserver
//	@Description	List all APIGroups of this apiserver
//	@Produce		json
//	@Success		200	{object} metav1.APIGroupList
//	@Router			/apis [get]
func APIs(w http.ResponseWriter, r *http.Request) {
	var gvk [3]string
	// 1.27+ kubectl discovery APIGroups and APIResourceList only by /apis with Header
	//   Accept: application/json;g=apidiscovery.k8s.io;v=v2beta1;as=APIGroupDiscoveryList
	// 1.27- kubectl discovery APIGroups and APIResourceList by /apis, /apis/{group}, /apis/{group}/{version}
	for _, acceptPart := range strings.Split(r.Header.Get("Accept"), ";") {
		if g_v_k := strings.Split(acceptPart, "="); len(g_v_k) == 2 {
			switch g_v_k[0] {
			case "g":
				gvk[0] = g_v_k[1]
			case "v":
				gvk[1] = g_v_k[1]
			case "as":
				gvk[2] = g_v_k[1]
			}
		}
	}

	if gvk[0] == "apidiscovery.k8s.io" && gvk[2] == "APIGroupDiscoveryList" {
		w.Header().Set("Content-Type", "application/json;g=apidiscovery.k8s.io;v=v2beta1;as=APIGroupDiscoveryList")
		w.Write([]byte(apidiscoveries))
	} else {
		w.Header().Set("Content-Type", "application/json")
		renderJSON(w, apis)
	}
}

// Get APIGroupHelloV1
//
//	@Summary		Get APIGroupHelloV1 info of 'hello.zeng.dev'
//	@Description	Get APIGroupHelloV1 'hello.zeng.dev' detail, including version list and preferred version
//	@Produce		json
//	@Success		200	{object} metav1.APIGroup
//	@Router			/apis/hello.zeng.dev [get]
func APIGroupHelloV1(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	renderJSON(w, apis.Groups[0])
}

const hellov1Resources = `{
	"kind": "APIResourceList",
	"apiVersion": "v1",
	"groupVersion": "hello.zeng.dev/v1",
	"resources": [
	  {
		"name": "foos",
		"singularName": "foo",
		"namespaced": true,
		"kind": "Foo",
		"verbs": [
		  "create",
		  "delete",
		  "get",
		  "list",
		  "update",
		  "patch"
		],
		"shortNames": [
		  "fo"
		],
		"categories": [
		  "all"
		]
	  }
	]
  }`

// Get APIGroupHelloV1Resources
//
//	@Summary		Get APIGroupHelloV1Resources for group version 'hello.zeng.dev/v1'
//	@Description	List APIResource Info about group version 'hello.zeng.dev/v1'
//	@Produce		json
//	@Success		200	{string} hellov1Resources
//	@Router			/apis/hello.zeng.dev/v1 [get]
func APIGroupHelloV1Resources(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(hellov1Resources))
}

//go:embed docs/*
var embedFS embed.FS

// Get OpenAPI Spec v2 doc
//
//	@Summary		Get OpenAPI Spec v2 doc of this server
//	@Description	Get OpenAPI Spec v2 doc of this server
//	@Produce		json
//	@Produce		application/com.github.proto-openapi.spec.v2@v1.0+protobuf
//	@Success		200	{string} swagger.json
//	@Router			/openapi/v2 [get]
func OpenapiV2(w http.ResponseWriter, r *http.Request) {
	jsonbytes, _ := embedFS.ReadFile("docs/swagger.json")

	// 😭 kubectl (v1.26.2, v1.27.1 ...) api discovery module (which fetch /openapi/v2, /openapi/v3)
	//    only accept application/com.github.proto-openapi.spec.v2@v1.0+protobuf
	if !strings.Contains(r.Header.Get("Accept"), "application/json") && strings.Contains(r.Header.Get("Accept"), "protobuf") {
		w.Header().Set("Content-Type", "application/com.github.proto-openapi.spec.v2.v1.0+protobuf")
		if pbbytes, err := ToProtoBinary(jsonbytes); err != nil {
			w.Header().Set("Content-Type", "application/json")
			writeErrStatus(w, "", http.StatusInternalServerError, err.Error())
			return
		} else {
			w.Write(pbbytes)
			return
		}
	}

	// 😄 kube apiserver aggregation module accept application/json
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonbytes)
}

func ToProtoBinary(json []byte) ([]byte, error) {
	document, err := gnosticopenapiv2.ParseDocument(json)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(document)
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
		// Msg says hello world!
		Msg string `json:"msg"`
		// Msg1 provides verbose information
		Msg1 string `json:"msg1"`
	} `json:"spec"`
}

func (f *Foo) DeepCopyObject() kruntime.Object {
	cf := *f
	return &cf
}

type FooList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Foo `json:"items"`
}

var x sync.RWMutex
var foos = map[string]Foo{}

func init() {
	foos["default/bar"] = Foo{
		TypeMeta:   metav1.TypeMeta{APIVersion: "hello.zeng.dev/v1", Kind: "Foo"},
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "bar", CreationTimestamp: metav1.Now()},
		Spec: struct {
			Msg  string "json:\"msg\""
			Msg1 string "json:\"msg1\""
		}{
			Msg:  "hello world",
			Msg1: "apiserver-from-scratch says '👋 hello world 👋'",
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
	"apiVersion": "v1",
	"kind": "Status",
	"metadata": {},
	"status": "Failure",
	"message": "%s",
	"reason": "%s",
	"details": {"group": "hello.zeng.dev", "kind": "foos", "name": "%s"},
	"code": %d
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
	w.WriteHeader(status)
	w.Write([]byte(errStatus))
}

var fooCol = []metav1.TableColumnDefinition{
	{Name: "Name", Type: "string", Format: "name", Description: metav1.ObjectMeta{}.SwaggerDoc()["name"]},
	{Name: "Age", Type: "string", Description: metav1.ObjectMeta{}.SwaggerDoc()["creationTimestamp"]},
	{Name: "Message", Type: "string", Format: "message", Description: "foo message"},
	{Name: "Message1", Type: "string", Format: "message1", Description: "foo message plus", Priority: 1}, // kubectl -o wide
}

func foo2TableRow(f *Foo) []metav1.TableRow {
	ts := "<unknown>"
	if timestamp := f.CreationTimestamp; !timestamp.IsZero() {
		ts = durationutil.HumanDuration(time.Since(timestamp.Time))
	}
	return []metav1.TableRow{
		{
			Object: kruntime.RawExtension{Object: f}, // get -A show NAMESPACE column
			Cells:  []interface{}{f.Name, ts, f.Spec.Msg, f.Spec.Msg1},
		},
	}
}

func fooList2TableRows(f *FooList) (ret []metav1.TableRow) {
	for idx := range f.Items {
		ret = append(ret, foo2TableRow(&f.Items[idx])...)
	}
	return
}

// application/json;as=Table;v=v1;g=meta.k8s.io,application/json;as=Table;v=v1beta1;g=meta.k8s.io,application/json
func tryConvert2Table(obj interface{}, acceptedContentType string) interface{} {
	if strings.Contains(acceptedContentType, "application/json") && strings.Contains(acceptedContentType, "as=Table") {
		switch typedObj := obj.(type) {
		case Foo:
			return metav1.Table{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Table",
					APIVersion: "meta.k8s.io/v1",
				},
				ColumnDefinitions: fooCol,
				Rows:              foo2TableRow(&typedObj),
			}
		case FooList:
			return metav1.Table{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Table",
					APIVersion: "meta.k8s.io/v1",
				},
				ColumnDefinitions: fooCol,
				Rows:              fooList2TableRows(&typedObj),
			}
		default:
			return obj
		}
	}
	return obj
}

// GetFoo swag doc
// @Summary      Get one Foo Object
// @Description  Get one Foo by Resource Name
// @Tags         foos
// @Produce      json
// @Param        namespace	path	string  true  "Namepsace"
// @Param        name	path	string  true  "Resource Name"
// @Success      200  {object}  Foo
// @Router       /apis/hello.zeng.dev/v1/namespaces/{namespace}/foos/{name} [get]
func GetFoo(w http.ResponseWriter, r *http.Request, name string) {
	ns := r.Context().Value(ctxkey("namespace"))
	nsname := fmt.Sprintf("%s/%s", ns, name)

	x.RLock()
	defer x.RUnlock()

	f, ok := foos[nsname]
	if !ok {
		writeErrStatus(w, nsname, http.StatusNotFound, "")
		return
	}
	renderJSON(w, tryConvert2Table(f, r.Header.Get("Accept")))
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

	x.RLock()
	defer x.RUnlock()

	for _, f := range foos {
		if f.Namespace == r.Context().Value(ctxkey("namespace")) {
			flist.Items = append(flist.Items, f)
		}
	}
	renderJSON(w, tryConvert2Table(flist, r.Header.Get("Accept")))
}

// GetAllFoos swag doc
// @Summary      List all Foos
// @Description  List all Foos
// @Tags         foos
// @Produce      json
// @Success      200  {object}  FooList
// @Router       /apis/hello.zeng.dev/v1/foos [get]
func GetAllFoos(w http.ResponseWriter, r *http.Request) {
	flist := FooList{
		TypeMeta: metav1.TypeMeta{Kind: "FooList", APIVersion: "hello.zeng.dev/v1"},
	}

	x.RLock()
	defer x.RUnlock()

	for _, f := range foos {
		flist.Items = append(flist.Items, f)
	}
	renderJSON(w, tryConvert2Table(flist, r.Header.Get("Accept")))
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
	f.CreationTimestamp = metav1.Now()

	ns := r.Context().Value(ctxkey("namespace"))
	nsname := fmt.Sprintf("%s/%s", ns, f.Name)

	x.Lock()
	defer x.Unlock()

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

	x.Lock()
	defer x.Unlock()

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

	x.Lock()
	defer x.Unlock()

	var old, ok = foos[nsname]
	if !ok { // not exists
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var originalBytes, _ = json.Marshal(old)
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
		var schema Foo = old
		var originalMap, _ = kruntime.DefaultUnstructuredConverter.ToUnstructured(&old)
		var patchMap map[string]interface{}

		if err = json.Unmarshal(patchBytes, &patchMap); err != nil {
			writeErrStatus(w, nsname, http.StatusBadRequest, err.Error())
			return
		}

		if patchedMap, err := kstrategicpatch.StrategicMergeMapPatch(originalMap, patchMap, schema); err != nil {
			writeErrStatus(w, nsname, http.StatusBadRequest, err.Error())
			return
		} else {
			var theFoo Foo
			if err = kruntime.DefaultUnstructuredConverter.FromUnstructuredWithValidation(patchedMap, &theFoo, false); err != nil {
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

	x.Lock()
	defer x.Unlock()

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
