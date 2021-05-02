FROM golang:1.16-alpine as builder
ENV GOPROXY=https://goproxy.io,direct
ENV GO111MODULE=on

WORKDIR /go/cache
ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /work
ADD . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o arceus github.com/zc2638/arceus/cmd

FROM gruebel/upx:latest as upx
WORKDIR /work
COPY --from=builder /work/arceus /work/arceus.pre
RUN upx --best --lzma -o /work/arceus arceus.pre

FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates

FROM alpine
LABEL maintainer="zc" \
    email="zc2638@qq.com" \
    version="v0.3.0"

ENV TS="Asia/Shanghai"
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=upx /work/arceus /work/arceus

WORKDIR /work
CMD ["./arceus"]
