doc: swag
	swag init --outputTypes json --parseDependency true

swag: # find or download swag if necessary
ifeq (, $(shell which ko))
	go install github.com/swaggo/swag/cmd/swag@v1.8.12
endif

deploy:
	KOCACHE=/tmp/ko KO_DOCKER_REPO=localhost:5000 ko apply -B -f ./deploy.yml

undeploy:
	kubectl delete -f ./deploy.yml