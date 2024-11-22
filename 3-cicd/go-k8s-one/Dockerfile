# stage 1: build src code to binary
FROM harbor.alnk.com/public/golang:1.22.8 as builder

ENV GOPROXY https://goproxy.cn

COPY . /app/

RUN cd  /app \
&& go mod init go-k8s-one \
&& go mod tidy \
&& CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o go-k8s-one .

# stage 2: use alpine as base image
FROM harbor.od.com/public/alpine:3.18

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
apk update && \
apk --no-cache add tzdata ca-certificates && \
cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
# apk del tzdata && \
rm -rf /var/cache/apk/*

COPY --from=builder /app/go-k8s-one /app/go-k8s-one
COPY --from=builder /app/conf/config.ini /app/conf/config.ini

CMD ["/app/go-k8s-one"]


#创建golang目录
#cd golang
#docker build -t harbor.od.com/app/go-k8s-one:v0.0.1 .
#docker run -d -p9999:9999 harbor.od.com/app/go-k8s-one:v0.0.1