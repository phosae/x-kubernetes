# apiserver using official libraries

## run
just store objects in memory

```
go run main.go
```

persist objects in etcd as json

```bash
go run main.go --enable-etcd-storage --etcd-servers http://localhost:2379
```

persist objects in etcd as Protocol Buffers
```bash
go run main.go --enable-etcd-storage --etcd-servers http://localhost:2379 --storage-media-type application/vnd.kubernetes.protobuf
```

## play

```bash
$ kubectl apply -f ../api/artifacts/samples/hello-foo.yml 
$ kubectl delete fo test
$ kubectl edit fo myfoo
---
$ kubectl get fo -w --output-watch-events 
```