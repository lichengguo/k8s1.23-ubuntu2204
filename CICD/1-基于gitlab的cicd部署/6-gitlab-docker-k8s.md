

环境

10.0.1.21已经安装了docker、dockerhub、gitlab、gitlab-runner



创建一个go项目，属于demo组，因为该组已经绑定了一个runner

![1730997140534](images\1730997140534.png)  



上传一个`main.go`文件

```go
package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	WebRequestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "web_reqeust_total",
		Help: "Number of hello requests in total",
	}, []string{"method", "path"})

	WebRequestDurationHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Histogram of the duration of HTTP requests",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})
)

func init() {
	// 注册计数器到 Prometheus
	prometheus.MustRegister(WebRequestTotal)
	// 注册直方图到 Prometheus
	prometheus.MustRegister(WebRequestDurationHistogram)
}

func main() {
	r := gin.Default()

	r.Use(func(ctx *gin.Context) {
		startTime := time.Now().UnixNano()
		// 处理请求
		ctx.Next()

		//记录请求次数
		WebRequestTotal.WithLabelValues(ctx.Request.Method, ctx.Request.URL.Path).Inc()

		//记录http方法和路径对应的耗时
		endTime := time.Now().UnixNano()
		seconds := float64((endTime - startTime) / 1e9) // s
		// Milliseconds := float64((endTime - startTime) / 1e6) // ms
		// nanoSeconds := float64(endTime - startTime)          // ns
		WebRequestDurationHistogram.WithLabelValues(ctx.Request.Method, ctx.Request.URL.Path).Observe(seconds)
	})

	// 将Prometheus的metrics接口挂载到Gin的路由上
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 健康检查
	r.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})

	// 其他业务逻辑
	r.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello, Go! hello, 业务!")
	})

	r.Run(":9999")
}

// 交叉编译
// 设置Go交叉编译环境变量
// $env:GOOS = "linux" 
// $env:GOARCH = "amd64"
// go build -o hello-go-linux


```



上传一个`Dockerfile`文件

```dockerfile
FROM harbor.alnk.com/public/golang:1.22.8 as builder
ENV GOPROXY https://goproxy.cn
COPY . /app/
RUN cd /app \
    && go mod init github.com/alnk/go-hello-prometheus-k8s \
    && go mod tidy \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o go-hello-prometheus-k8s .

FROM harbor.alnk.com/public/alpine:3.18
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \
    apk --no-cache add tzdata ca-certificates && \
    cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    # apk del tzdata && \
    rm -rf /var/cache/apk/*
COPY --from=builder /app/go-hello-prometheus-k8s /go-hello-prometheus-k8s
EXPOSE 9999
CMD ["/go-hello-prometheus-k8s"] 

```





编写`.gitlab-ci.yml`文件

```yaml
variables:
  IMAGE_NAME: go-hello-prometheus-k8s
  HARBOR_URL: harbor.alnk.com


stages:
  - build_iamge
  - tag_image
  - push_image


build_iamge: 
  stage: build_iamge
  script: 
    - IMAGE_VERSION="$(date +%Y%m%d%H%M%S)" 
    - docker build -t  "$HARBOR_URL"/public/"$IMAGE_NAME":"$IMAGE_VERSION" .
    - echo $IMAGE_VERSION > docker_images_version.env # 暴露局部变量IMAGE_VERSION
  artifacts: # 暴露局部变量IMAGE_VERSION
    paths:
      - docker_images_version.env

tag_image:
  stage: tag_image
  script:
    - IMAGE_VERSION=$(cat docker_images_version.env)  # 获取局部变量IMAGE_VERSION
    - docker push "$HARBOR_URL"/public/"$IMAGE_NAME":"$IMAGE_VERSION"

push_image:
  stage: push_image
  script:
    - IMAGE_VERSION=$(cat docker_images_version.env)  # 获取局部变量IMAGE_VERSION
    - docker rmi "$HARBOR_URL"/public/"$IMAGE_NAME":"$IMAGE_VERSION"

```



注意给下权限，其实就是gitlab-runner这个用户在操作docker build这个命令

```shell
# docker socket
root@ops:/var/run# ll docker.sock
srw-rw-rw- 1 root docker 0 Nov  7 20:31 docker.sock=

root@ops:/var/run# id gitlab-runner
uid=993(gitlab-runner) gid=993(gitlab-runner) groups=993(gitlab-runner)

root@ops:/var/run# chmod 666 docker.sock


root@ops:/var/run# ll /root/.docker/config.json
-rw------- 1 root root 183 Oct 31 23:39 /root/.docker/config.json

#docker harbor凭证
root@ops:/var/run# mkdir -p /home/gitlab-runner/.docker
root@ops:/var/run# cp ~/.docker/config.json /home/gitlab-runner/.docker/
root@ops:/var/run# cd /home/gitlab-runner/
root@ops:/home/gitlab-runner# chown gitlab-runner.gitlab-runner -R .docker/



```

![1731000667304](images\1731000667304.png)

![1730999907607](images\1730999907607.png)  