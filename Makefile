.PHONY: api-aggregation-simple
api-aggregation-simple:
	KOCACHE=/tmp/ko KO_DOCKER_REPO=zengxu ko build --platform linux/amd64,linux/arm64 -B github.com/phosae/x-kubernetes/api-aggregation-simple

deploy-api-aggregation-simple:
	KOCACHE=/tmp/ko KO_DOCKER_REPO=zengxu ko apply -B -f ./api-aggregation-simple/deploy.yml