package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/alertmanager/notify/webhook"
	"github.com/prometheus/alertmanager/template"
)

func main() {
	r := gin.Default()

	// 钉钉
	r.POST("/webhook/dingding", alertReceiveDingding)

	r.Run(":8085")
}

// alertReceive 解析alert消息发送至钉钉
func alertReceiveDingding(c *gin.Context) {
	// 打印Prometheus发送过来的原始信息
	// ioutil.ReadAll读取到的是[]byte,读完body就没有了
	// body, err := ioutil.ReadAll(c.Request.Body)
	// 使用ioutil.NopCloser重新赋值给body
	// c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}
	log.Printf("prometheus发过来的原始信息: %s\n", string(b))
	c.Request.Body = io.NopCloser(bytes.NewReader(b))

	// ------------------------------------------------------------------------------------------
	var msg webhook.Message
	if err := c.BindJSON(&msg); err != nil {
		c.JSON(400, errors.New("invalid args"))
		return
	}

	// fmt.Println("==================================")
	// fmt.Printf("收到Prometheus告警信息: %#v\n", &msg)
	// fmt.Println("==================================")

	baseMsg := fmt.Sprintf("[状态：%s][报警条数:%d]", msg.Status, len(msg.Alerts))
	log.Printf("[alertReceive][baseMsg:%+v]", baseMsg)

	for i := 0; i < len(msg.Alerts); i++ {
		alert := msg.Alerts[i]
		bs, _ := buildDDContent(alert)

		log.Printf("[detail][%d/%d][alert:%+v]", i+1, len(msg.Alerts), alert)
		sendToDing(bs)
	}

	c.JSON(200, "ok")
}

// dingMsg 钉钉消息格式
type dingMsg struct {
	Msgtype string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	At struct {
		AtMobiles []string `json:"atMobiles"`
	} `json:"at"`
}

// buildDDContent 拼接钉钉信息的函数
func buildDDContent(msg template.Alert) ([]byte, error) {
	recM := map[string]string{"firing": "已触发", "resolved": "已恢复"}

	// msgTpl := fmt.Sprintf(
	// 		"钉钉这个老六需要一个关键字，不然不给告警，关键字:alnk"+
	// 		"[规则名称：%s]\n"+
	// 		"[是否已恢复：%s]\n"+
	// 		"[告警级别：%s]\n"+
	// 		"[触发时间：%s]\n"+
	// 		"[看图连接：%s]\n"+
	// 		"[当前值：%s]\n"+
	// 		"[标签组：%s]",
	// 	msg.Labels["alertname"],
	// 	recM[msg.Status],
	// 	msg.Labels["severity"],
	// 	// prometheus使用utc时间，转换为当前时间
	// 	msg.StartsAt.In(time.Local).Format(time.DateTime),
	// 	msg.GeneratorURL,
	// 	msg.Annotations["value"],
	// 	msg.Labels.SortedPairs(),
	// )

	msgTpl := fmt.Sprintf(
		"[钉钉这个老六需要一个关键字，不然不给告警，关键字:alnk]\n"+
			"[当前状态: %s]\n"+
			"[告警规则名称alertname: %s]\n"+
			"[实例instance: %s]\n"+
			"[告警级别severity: %s]\n"+
			"[告警简要信息summary: %s]\n"+
			"[告警详细信息description: %s]\n"+
			"[告警时间startsAt: %s]\n"+
			"[看图连接externalURL: %s]\n",
		recM[msg.Status],
		msg.Labels["alertname"],
		msg.Labels["instance"],
		msg.Labels["severity"],
		msg.Annotations["summary"],
		msg.Annotations["description"],
		// prometheus使用utc时间，转换为当前时间
		msg.StartsAt.In(time.Local).Format(time.DateTime),
		msg.GeneratorURL,
	)

	dm := dingMsg{Msgtype: "text"}
	dm.Text.Content = msgTpl
	bs, err := json.Marshal(dm)
	return bs, err
}

// sendToDing 发送消息到钉钉
func sendToDing(jsonByte []byte) {
	apiUrl := "https://oapi.dingtalk.com/robot/send?access_token=ab21f0ebc52254c61c0e8356b96e536cb95341cab716efe6acd47285527b6d08"

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonByte))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[http.post.request.err][url:%v][err:%v]", apiUrl, err)
		return
	}
	defer resp.Body.Close()

	log.Printf("response Status:%v", resp.Status)
	log.Printf("response Headers:%v", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	log.Printf("response Body:%v", string(body))
}
