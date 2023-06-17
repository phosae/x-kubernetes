deploy:
	kustomize build --load-restrictor LoadRestrictionsNone | KOCACHE=/tmp/ko KO_DOCKER_REPO=localhost:5000 KO_DEFAULTBASEIMAGE=localhost:5000/debian:bullseye-20230502 ko apply -B -f -

undeploy:
	kustomize build --load-restrictor LoadRestrictionsNone | kubectl delete -f -
