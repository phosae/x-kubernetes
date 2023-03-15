LOCAL_IP_LIST=$(ip a | grep inet |  awk '{print $2}' | cut -d/ -f1 | paste -sd "," -)

PKIDIR ?= $(shell pwd)/pki

pkidir:
	[ -d $(PKIDIR) ] || mkdir -p $(PKIDIR)

apicert: pkidir
	go run ./hack/gen_cert.go --host "kubernetes,kubernetes.default,kubernetes.default.svc,kubernetes.default.svc.cluster.local,localhost,apiserver,kube-apiserver,172.20.0.1,$LOCAL_IP_LIST"  --ecdsa-curve P256 --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h \
	&& mv tls.crt ./pki/apiserver.crt && mv ./tls.key ./pki/apiserver.key

sa-pki: pkidir
	openssl ecparam -name prime256v1 -genkey -noout -out $(PKIDIR)/sa-ecdsa.key \
	&& openssl ec -in $(PKIDIR)/sa-ecdsa.key -pubout -out $(PKIDIR)/sa-ecdsa.pub

token-users:
	printf 'admin-token,admin,,"system:masters"' > $(PKIDIR)/token-users.csv

pki: apicert sa-pki token-users

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