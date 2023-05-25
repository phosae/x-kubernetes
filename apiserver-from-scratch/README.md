# apiserver from scratch

## play 
```shell
go run main.go
```

```shell
$ kubectl -s localhost:8000 get foo

NAME   AGE   MESSAGE       MESSAGE1
bar    3s    hello world 
```

```bash
{
kubectl -s localhost:8000 apply -f ../api/artifacts/samples/hello-foo.yml 
kubectl -s localhost:8000  patch foo/test --patch-file ../api/artifacts/samples/patch-hello-foo.yml
kubectl -s localhost:8000 get foo
}
```

outputs
```shell
foo.hello.zeng.dev/myfoo created
foo.hello.zeng.dev/myfoo created
foo.hello.zeng.dev/test created
foo.hello.zeng.dev/test patched
NAME    AGE    MESSAGE        MESSAGE1
bar     2m8s   hello world    
myfoo   0s     my first foo   
test    0s     testmsg        hey, there, 👋
```

If serve on HTTPS, play with `kubectl -s https://localhost:6443 --certificate-authority /path/to/ca`

## aggregate with kube-apiserver
use `make deploy-apiserver-from-scratch` in parent directory
```bash
cd ..
make deploy-apiserver-from-scratch
```
## OpenAPI v2 json

gen swag(i.e openapi/v2)
```bash
make doc
```