FROM golang:1.22.5 as builder
WORKDIR /go/src/github.com/AliyunContainerService/image-syncer
COPY ./ ./
ENV GOPROXY=https://proxy.golang.com.cn,direct
RUN CGO_ENABLED=0 GOOS=linux make

FROM alpine:latest
WORKDIR /bin/
COPY --from=builder /root/image-syncer ./
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN chmod +x ./image-syncer
RUN apk add -U --no-cache ca-certificates && rm -rf /var/cache/apk/* && mkdir -p /etc/ssl/certs \
  && update-ca-certificates --fresh
ENTRYPOINT ["image-syncer"]
CMD ["-c", "/root/sync.yaml"]
