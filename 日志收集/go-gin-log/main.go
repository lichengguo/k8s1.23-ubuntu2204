package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var LogToFileAndStdout func(format string, v ...interface{})

func main() {
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		log.Printf("这是直接打印到std的日志, 请求URL: %s, 请求方法Method: %s, 请求主机地址: %s\n", c.Request.URL, c.Request.Method, c.Request.Host)
		LogToFileAndStdout("请求URL: %s, 请求方法Method: %s, 请求主机地址: %s\n", c.Request.URL, c.Request.Method, c.Request.Host)
		c.String(http.StatusOK, "hello, Gin! 新的业务上线了 version: 0.1")
	})

	r.Run(":3000")
}

func init() {
	// 当前程序执行的目录
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	// 判断日志目录是否存在
	logDirPath := fmt.Sprintf("%s/%s", dir, "log")
	_, err = os.Stat(logDirPath)
	if err != nil {
		// 不存在则创建目录
		os.MkdirAll(dir, 0755)
	}

	// 创建一个文件，用于写入日志
	file, err := os.OpenFile(fmt.Sprintf("%s/%s", logDirPath, "go-gin-log.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	// defer file.Close()

	// 创建一个日志写入文件的Logger
	fileLogger := log.New(file, "", log.LstdFlags)

	// 同时创建一个标准Logger，用于输出屏幕
	// consoleLogger := log.New(os.Stdout, "", log.LstdFlags)

	// 定义一个写入日志的函数，同时向文件和屏幕输出
	LogToFileAndStdout = func(format string, v ...interface{}) {
		fileLogger.Printf(format, v...)
		// consoleLogger.Printf(format, v...)
	}
}
