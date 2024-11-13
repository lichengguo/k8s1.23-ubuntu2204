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
