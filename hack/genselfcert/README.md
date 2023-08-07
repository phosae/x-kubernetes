# Generate TLS Certificate for Kubernetes Webhook Server

## Usage

```bash
$ docker run -v $PWD:/output -w /output -e HOSTS=zeng.dev,example.com,1.1.1.1 -e NAMESPACE=default,kube-system zengxu/genselfcert

# or copy from std
docker run -e HOSTS=example.com,1.1.1.1 -e NAMESPACE=test,default zengxu/genselfcert
```

output

```bash
+ NS=test,default
+ SUFFIX=cluster.local
++ printf test,default
++ sed 's/\([^,]\+\)/*.\1.svc/g'
+ NSHOSTS='*.test.svc,*.default.svc'
++ printf '*.test.svc,*.default.svc'
++ sed 's/\([^,]\+\)/\1.cluster.local/g'
+ SUFFIXED_NSHOSTS='*.test.svc.cluster.local,*.default.svc.cluster.local'
+ /gencert --host 'example.com,1.1.1.1,*.test.svc,*.default.svc,*.test.svc.cluster.local,*.default.svc.cluster.local,127.0.0.1,::1' --ecdsa-curve P256 --ca --start-date 'Jan 1 00:00:00 1970' --duration=1000000h
2023/04/26 03:20:31 wrote tls.crt
2023/04/26 03:20:31 wrote tls.key
+ cat tls.crt
-----BEGIN CERTIFICATE-----
MIICKDCCAc6gAwIBAgIQQnM+BmfxM1X69ysKGU
...
9LGp5KW951Q57iVJG9ws/xZP8mihProFx7MB4Q==
-----END CERTIFICATE-----
+ cat tls.key
-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgCFRhF
...
gwdjXO8+
-----END PRIVATE KEY-----
```
Other ENV Args
- `ORG` Organization in Certificate
- `SVC_SUFFIX` change default service suffic `cluster.local`

## Build

```bash
$ docker buildx build --platform linux/amd64,linux/arm64 -t zengxu/genselfcert --push .
```