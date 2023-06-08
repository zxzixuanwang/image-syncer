FROM golang:1.20 as builder
WORKDIR /root
RUN  go env -w GOPROXY=https://goproxy.cn,direct
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make

FROM alpine:latest
WORKDIR /bin/
COPY --from=builder /root/image-syncer ./
RUN chmod +x ./image-syncer
RUN apk add -U --no-cache ca-certificates && rm -rf /var/cache/apk/* && mkdir -p /etc/ssl/certs \
  && update-ca-certificates --fresh
ENTRYPOINT ["image-syncer"]
CMD ["-c", "/root/sync.yaml"]
