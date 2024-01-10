localenv: kind kubectl ko kustomize
	./hack/setup-kind-with-registry.sh

kind: # find or download kind if necessary
ifeq (, $(shell which kind))
	GOBIN=/usr/local/bin/ go install sigs.k8s.io/kind@v0.20.0
endif

kubectl: # find or download kubectl if necessary
ifeq (, $(shell which kubectl))
	curl -LO https://dl.k8s.io/release/v1.29.0/bin/linux/amd64/kubectl
	sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
	rm kubectl
endif

ko: # find or download ko if necessary
ifeq (, $(shell which ko))
	GOBIN=/usr/local/bin/ go install github.com/google/ko@v0.15.1
endif

kustomize: # find or download kustomize if necessary
ifeq (, $(shell which kustomize))
	GOBIN=/usr/local/bin/ go install -ldflags="-X 'sigs.k8s.io/kustomize/api/provenance.version=v5.3.0'" sigs.k8s.io/kustomize/kustomize/v5@v5.3.0
endif