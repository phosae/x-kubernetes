# Generate TLS Certificate for Kubernetes Webhook Server

## Usage

```bash
$ docker run -v $PWD:/output -w /output -e HOSTS=zeng.dev,example.com,1.1.1.1 -e NS=default,kube-system zengxu/genselfcert

# or copy from std
docker run -e HOSTS=example.com -e NS=test zengxu/genselfcert
```

output

```bash
+ NS=default,kube-system
+ SUFFIX=cluster.local
++ printf default,kube-system
++ sed 's/\([^,]\+\)/*.\1.svc/g'
+ NSHOSTS='*.default.svc,*.kube-system.svc'
++ printf '*.default.svc,*.kube-system.svc'
++ sed 's/\([^,]\+\)/\1.cluster.local/g'
+ SUFFIXED_NSHOSTS='*.default.svc.cluster.local,*.kube-system.svc.cluster.local'
+ /gencert --host 'example.com,*.default.svc,*.kube-system.svc,*.default.svc.cluster.local,*.kube-system.svc.cluster.local,127.0.0.1,::1,fe80::1' --ecdsa-curve P256 --ca --start-date 'Jan 1 00:00:00 1970' --duration=1000000h
2023/04/25 13:44:09 wrote tls.crt
2023/04/25 13:44:09 wrote tls.key
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

## Build

```bash
$ docker buildx build --platform linux/amd64,linux/arm64 -t zengxu/genselfcert --push .
```