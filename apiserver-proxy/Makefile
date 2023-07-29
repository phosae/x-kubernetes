deploy: cert kubeconfig
	kustomize build | KOCACHE=/tmp/ko KO_DOCKER_REPO=localhost:5000 KO_DEFAULTBASEIMAGE=debian:bullseye-20230502 ko apply -B -f -

undeploy:
	kustomize build | kubectl delete -f -

cert:
	mkdir -p apiserver.local.config/certificates
	docker cp artifacts/gen-apiserver-cert.sh kind-control-plane:/etc/kubernetes/pki
	docker exec -it -w /etc/kubernetes/pki kind-control-plane /bin/sh -c ./gen-apiserver-cert.sh
	docker cp kind-control-plane:/etc/kubernetes/pki/apiserver-proxy.crt apiserver.local.config/certificates/apiserver.crt
	docker cp kind-control-plane:/etc/kubernetes/pki/apiserver-proxy.key apiserver.local.config/certificates/apiserver.key

kubeconfig:
	cp $(HOME)/.kube/config ./config.pxy
	os_name=$$(uname); \
	sed_opts=""; \
	if [ "$$os_name" = "Darwin" ]; then \
	    sed_opts="-i ''"; \
	elif [ "$$os_name" = "Linux" ]; then \
	    sed_opts="-i"; \
	fi; \
	sed $$sed_opts 's~server: https://[^:]*:[0-9]*~server: https://172.18.0.2:30443~g' config.pxy
	