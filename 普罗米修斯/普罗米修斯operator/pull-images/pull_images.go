package main

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func main() {

	imagesURL := [...]string{
		"registry.cn-hangzhou.aliyuncs.com/alnktest/cpvpa-amd64:v0.8.1",
		"registry.cn-hangzhou.aliyuncs.com/alnktest/metrics-server-amd64:v0.2.0",
		"registry.cn-hangzhou.aliyuncs.com/alnktest/grafana:8.5.5",
		"registry.cn-hangzhou.aliyuncs.com/alnktest/configmap-reload:v0.5.0",
		"registry.cn-hangzhou.aliyuncs.com/alnktest/kube-state-metrics:v2.5.0",
		"registry.cn-hangzhou.aliyuncs.com/alnktest/prometheus-adapter:v0.9.1",
		"registry.cn-hangzhou.aliyuncs.com/alnktest/kube-rbac-proxy:v0.12.0",
		"registry.cn-hangzhou.aliyuncs.com/alnktest/alertmanager:v0.24.0",
		"registry.cn-hangzhou.aliyuncs.com/alnktest/blackbox-exporter:v0.21.0",
		"registry.cn-hangzhou.aliyuncs.com/alnktest/node-exporter:v1.3.1",
		"registry.cn-hangzhou.aliyuncs.com/alnktest/prometheus-operator:v0.57.0",
		"registry.cn-hangzhou.aliyuncs.com/alnktest/prometheus:v2.36.1",
		"registry.cn-hangzhou.aliyuncs.com/alnktest/thanos:v0.19.0",
		"quay.io/fabxc/prometheus_demo_service",
	}

	for i := 0; i < len(imagesURL); i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			// pull
			pull := fmt.Sprintf("/usr/bin/docker pull %s", imagesURL[i])
			pullCmd := exec.Command("/bin/sh", "-c", pull)
			out, err := pullCmd.CombinedOutput()
			if err != nil {
				fmt.Println("pull err: ", imagesURL[i])
				return
			}
			fmt.Println("pull image: ", string(out))

			// tag
			temp := strings.Split(imagesURL[i], "/")
			imageVerson := temp[len(temp)-1]
			tag := fmt.Sprintf("/usr/bin/docker tag %s %s", imagesURL[i], "harbor.alnk.com/public/"+imageVerson)
			tagCmd := exec.Command("/bin/sh", "-c", tag)
			_, err = tagCmd.CombinedOutput()
			if err != nil {
				fmt.Println("tag err: ", imageVerson)
			}

			// push
			push := fmt.Sprintf("/usr/bin/docker push %s", "harbor.alnk.com/public/"+imageVerson)
			pushCmd := exec.Command("/bin/sh", "-c", push)
			_, err = pushCmd.CombinedOutput()
			if err != nil {
				fmt.Println("push err: ", imageVerson)
			}
		}()

	}

	wg.Wait()

	fmt.Println("执行完成")
}
