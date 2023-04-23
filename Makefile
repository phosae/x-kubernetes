wild-svc-cert:
	go run ./hack/gen_cert.go --host "*.kube-system,*.kube-system.svc,*.kube-system.svc.cluster.local,*.default,*.default.svc,*.default.svc.cluster.local,localhost,127.0.0.1,::1"  --ecdsa-curve P256 --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h \
	&& mv tls.crt ./api-aggregation-simple/tls.crt && mv ./tls.key ./api-aggregation-simple/tls.key

.PHONY: api-aggregation-simple
api-aggregation-simple:
	KOCACHE=/tmp/ko KO_DOCKER_REPO=zengxu ko build --platform linux/amd64,linux/arm64 -B github.com/phosae/x-kubernetes/api-aggregation-simple

deploy-api-aggregation-simple:
	KOCACHE=/tmp/ko KO_DOCKER_REPO=zengxu ko apply -B -f ./api-aggregation-simple/deploy.yml

rsync-remote:
	rsync -a ~/go/src/github.com/phosae/x-kubernetes root@122.228.207.19:/root