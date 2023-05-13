build:
	go build -o ./bin/hello-apiserver cmd/hello-apiserver/main.go

run: # make sure you've set up some local cluster like KinD cluster
	go run cmd/hello-apiserver/main.go --etcd-servers http://localhost:2379 --authentication-kubeconfig /root/.kube/config --authorization-kubeconfig /root/.kube/config --kubeconfig /root/.kube/config