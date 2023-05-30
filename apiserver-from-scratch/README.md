# apiserver from scratch

## play 
```shell
go run main.go
```

```shell
kubectl -s localhost:8000 get foo -o wide
NAME   AGE   MESSAGE       MESSAGE1
bar    16s   hello world   apiserver-from-scratch says '👋 hello world 👋'

kubectl -s localhost:8000 patch fo/bar --type merge -p '{"spec": {"msg": "hi👋", "msg1": ""}}'

kubectl -s localhost:8000 get foo -o wide
NAME   AGE     MESSAGE   MESSAGE1
bar    2m49s   hi👋

kubectl -s localhost:8000 delete fo/bar
foo.hello.zeng.dev "bar" deleted
```

```bash
{
kubectl -s localhost:8000 apply -f ../api/artifacts/samples/hello-foo.yml 
kubectl -s localhost:8000 patch foo/test --patch-file ../api/artifacts/samples/patch-hello-foo.yml
kubectl -s localhost:8000 get foo
}
```

outputs
```shell
foo.hello.zeng.dev/myfoo created
foo.hello.zeng.dev/myfoo created
foo.hello.zeng.dev/test created
foo.hello.zeng.dev/test patched
NAME    AGE   MESSAGE
test    1s    hey there
bar     43s   hello world
myfoo   0s    my first foo
```

If serve on HTTPS, play with `kubectl -s https://localhost:6443 --certificate-authority /path/to/ca`

## aggregate with kube-apiserver

```bash
../hack/setup-kind-with-registry.sh
make deploy
```
## OpenAPI v2 json

gen swag(i.e openapi/v2)
```bash
make doc
```