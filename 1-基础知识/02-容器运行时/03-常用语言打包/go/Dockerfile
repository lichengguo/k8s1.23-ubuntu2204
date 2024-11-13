# stage 1: build src code to binary
#FROM registry.cn-hangzhou.aliyuncs.com/alnktest/golang:1.19.11-alpine3.18 as builder
FROM harbor.alnk.com/public/golang:1.19.11-alpine3.18 as builder

ENV GOPROXY https://goproxy.cn

COPY . /app/

RUN cd /app \
    && go mod init github.com/alnk/golangtest \
    && go mod tidy \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o hellogo .


## stage 2: use alpine as base image
#FROM registry.cn-hangzhou.aliyuncs.com/alnktest/alpine:3.18
FROM harbor.alnk.com/public/alpine:3.18

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \
    apk --no-cache add tzdata ca-certificates && \
    cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    # apk del tzdata && \
    rm -rf /var/cache/apk/*

COPY --from=builder /app/hellogo /hellogo

EXPOSE 3000

CMD ["/hellogo"] 


# 创建golang目录
# cd golang
# docker build -t harbor.alnk.com/public/hello-go:v0.0.1 .
# docker push harbor.alnk.com/public/hello-go:v0.0.1
# docker run -d --name hello-go -p 9999:3000 harbor.alnk.com/public/hello-go:v0.0.1

# root@ops:/data/docker-build/go# docker exec -it hello-go sh
# / # cat /etc/os-release
# NAME="Alpine Linux"
# ID=alpine
# VERSION_ID=3.18.6
# PRETTY_NAME="Alpine Linux v3.18"
# HOME_URL="https://alpinelinux.org/"
# BUG_REPORT_URL="https://gitlab.alpinelinux.org/alpine/aports/-/issues"
