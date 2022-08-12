FROM node:15.5 as front
WORKDIR /work
ADD . .
RUN make ui

FROM golang:1.18 as builder
ENV GOPROXY=https://goproxy.io,direct
ENV GO111MODULE=on

WORKDIR /go/cache
ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /work
ADD . .
COPY --from=front /work/static/ui /work/static/ui
RUN GOOS=linux CGO_ENABLED=0 make build

FROM gruebel/upx:latest as upx
WORKDIR /work
COPY --from=builder /work/arceus /work/arceus.pre
RUN upx --best --lzma -o /work/arceus arceus.pre

FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates

FROM alpine:3.6
MAINTAINER zc
LABEL \
    maintainer="zc" \
    email="zc2638@qq.com" \
    version="0.9"

ENV TZ="Asia/Shanghai"
ENV ARCEUS_BASE_PATH="/etc/arceus"

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=upx /work/arceus /bin/arceus
COPY --from=builder /work/static/kubernetes /etc/arceus/static/kubernetes
COPY --from=builder /work/static/patch /etc/arceus/static/patch
COPY --from=builder /work/static/custom /etc/arceus/static/custom
COPY --from=builder /work/static/rule /etc/arceus/static/rule
COPY --from=builder /work/static/template /etc/arceus/static/template

WORKDIR /work
CMD ["arceus", "server"]
