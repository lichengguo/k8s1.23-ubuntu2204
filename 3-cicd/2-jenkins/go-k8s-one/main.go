package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"go-k8s-one/nacoscfg"

	"github.com/gin-gonic/gin"
)

// NaConf nacos上的配置
type NaConf struct {
	Port int64  `json:"port"`
	Mes  string `json:"mes"`
}

func main() {
	if err := os.MkdirAll("./logs", 0777); err != nil {
		log.Fatal(err)
	}

	logFile, err := os.OpenFile("./logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

	r := gin.Default()

	var n NaConf

	err = json.Unmarshal(nacoscfg.LoadNacos(), &n)
	if err != nil {
		panic("解析失败")
	}

	r.GET("/", func(c *gin.Context) {
		// 增加一行日志
		logger.Println("path: ", c.FullPath(), ",resData: ", n)

		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"mes":  "hello world! 开发了很多新功能 version:v0.0.2",
			"data": n,
		})
	})

	_ = r.Run(fmt.Sprintf(":%d", n.Port))
}
