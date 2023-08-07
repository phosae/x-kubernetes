FROM --platform=$BUILDPLATFORM golang:1.18-bullseye as builder
ARG TARGETOS TARGETARCH
WORKDIR /workspace
COPY gencert.go gencert.go
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -a -o gencert-$TARGETARCH gencert.go

FROM bash:5.1-alpine3.17
ARG TARGETARCH
WORKDIR /
COPY run.sh /run.sh
COPY --from=builder /workspace/gencert-$TARGETARCH /gencert

ENTRYPOINT [ "/run.sh" ]