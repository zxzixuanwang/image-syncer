FROM golang:1.23 as builder
WORKDIR /root
COPY ./ ./
ENV GOPROXY=https://goproxy.cn,direct
RUN CGO_ENABLED=0 GOOS=linux make

FROM alpine:latest
WORKDIR /bin/
COPY --from=builder /root/image-syncer ./
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN chmod +x ./image-syncer
RUN sed -i 's#https\?://dl-cdn.alpinelinux.org/alpine#https://mirrors.tuna.tsinghua.edu.cn/alpine#g' /etc/apk/repositories
RUN apk add -U --no-cache ca-certificates && rm -rf /var/cache/apk/* && mkdir -p /etc/ssl/certs \
  && update-ca-certificates --fresh
ENTRYPOINT ["image-syncer"]
CMD ["-c", "/root/sync.yaml"]
