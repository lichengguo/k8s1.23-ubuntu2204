### Node模式日志收集方案

> ```
> 收集日志的框架流程
> 
> 模拟业务服务   日志收集       日志存储     消费日志         消费目的服务
> tomcat ---> log-pilot ---> kafka ---> logstash ---> elasticsearch ---> kibana
> 
> 日志量很大的时候，怕es扛不住可以使用上面架构
> 如果日志量不大，可以直接使用log-pilot --- > elasticsearch架构
> ```



#### log-pilot部署

> ```shell
> 参考文档地址：https://github.com/AliyunContainerService/log-pilot
> 
> 【10.0.1.21】
> ## 准备镜像
> # docker pull registry.cn-hangzhou.aliyuncs.com/acs/log-pilot:0.9.7-filebeat
> # docker tag registry.cn-hangzhou.aliyuncs.com/acs/log-pilot:0.9.7-filebeat harbor.alnk.com/public/log-pilot:0.9.7-filebeat
> # docker push harbor.alnk.com/public/log-pilot:0.9.7-filebeat
> 
> ## 创建目录
> # mkdir -p /data/k8s-yaml/elk && cd /data/k8s-yaml/elk
> # vi log-pilot.yaml
> 
> ## 应用
> # kubectl apply -f log-pilot.yaml
> # kubectl -n ns-elastic get pod
> ```
>
> `log-pilot.yaml`
>
> ```yaml
> ---
> apiVersion: v1
> kind: ConfigMap
> metadata:
>   name: log-pilot2-configuration
>   namespace: elk
> data:
>   logging_output: "kafka" # 采集日志输出到kafka
>   kafka_brokers: "10.0.1.21:9092" # kafka地址,多个地址可以写成"kafka1:9092,kafka2:9092"
>   kafka_version: "0.10.0" # 指定版本，生产中实测kafka 2.12-2.5 及以下，这里都配置为"0.10.0" 就可以了
>   # 当禁止自动创建topic是，这里配置kakfa中的有效topic
>   # 对应的业务容器收集日志需要传入这主题
>   kafka_topics: "tomcat-syslog,tomcat-access"
> 
> ---
> apiVersion: apps/v1
> kind: DaemonSet
> metadata:
>   name: log-pilot2
>   namespace: elk
>   labels:
>     k8s-app: log-pilot2
> spec:
>   selector:
>     matchLabels:
>       k8s-app: log-pilot2
>   updateStrategy:
>     type: RollingUpdate
>   template:
>     metadata:
>       labels:
>         k8s-app: log-pilot2
>     spec:
>       tolerations:
>       - key: node-role.kubernetes.io/master
>         effect: NoSchedule
>       containers:
>       - name: log-pilot2
>         image: harbor.alnk.com/public/log-pilot:0.9.7-filebeat
>         env:
>           - name: "LOGGING_OUTPUT"
>             valueFrom:
>               configMapKeyRef:
>                 name: log-pilot2-configuration
>                 key: logging_output
>           - name: "KAFKA_BROKERS"
>             valueFrom:
>               configMapKeyRef:
>                 name: log-pilot2-configuration
>                 key: kafka_brokers
>           - name: "KAFKA_VERSION"
>             valueFrom:
>               configMapKeyRef:
>                 name: log-pilot2-configuration
>                 key: kafka_version
>           - name: "NODE_NAME"
>             valueFrom:
>               fieldRef:
>                 fieldPath: spec.nodeName
>         volumeMounts:
>         - name: sock
>           mountPath: /var/run/docker.sock
>         - name: logs
>           mountPath: /var/log/filebeat
>         - name: state
>           mountPath: /var/lib/filebeat
>         - name: root
>           mountPath: /host
>           readOnly: true
>         - name: localtime
>           mountPath: /etc/localtime
>         # configure all valid topics in kafka
>         # when disable auto-create topic
>         - name: config-volume
>           mountPath: /etc/filebeat/config
>         securityContext:
>           capabilities:
>             add:
>             - SYS_ADMIN
>       terminationGracePeriodSeconds: 30
>     
>       volumes:
>       - name: sock
>         hostPath:
>           path: /var/run/docker.sock
>           type: Socket
>       - name: logs
>         hostPath:
>           path: /var/log/filebeat
>           type: DirectoryOrCreate
>       - name: state
>         hostPath:
>           path: /var/lib/filebeat
>           type: DirectoryOrCreate
>       - name: root
>         hostPath:
>           path: /
>           type: Directory
>       - name: localtime
>         hostPath:
>           path: /etc/localtime
>           type: File
>       # kubelet sync period
>       - name: config-volume
>         configMap:
>           name: log-pilot2-configuration
>           items:
>           - key: kafka_topics
>             path: kafka_topics
> ```



#### 部署tomcat模拟业务

> ```shell
> 【10.0.1.21】
> ## 创建目录
> # mkdir -p /data/k8s-yaml/elk && cd /data/k8s-yaml/elk
> # vi tomcat.yaml
> 
> ## 准备镜像
> # docker pull registry.cn-hangzhou.aliyuncs.com/alnktest/tomcat:7.0
> # docker tag registry.cn-hangzhou.aliyuncs.com/alnktest/tomcat:7.0 harbor.alnk.com/public/tomcat:7.0
> # docker push harbor.alnk.com/public/tomcat:7.0
> 
> ## 应用
> # kubectl apply -f tomcat.yaml
> 
> # 查看
> # kubectl -n elk get pod
> 
> ```
>
> `tomcat.yaml`
>
> ```yaml
> apiVersion: apps/v1
> kind: Deployment
> metadata:
>   labels:
>     app: tomcat
>   name: tomcat
>   namespace: elk
> spec:
>   replicas: 1
>   selector:
>     matchLabels:
>       app: tomcat
>   template:
>     metadata:
>       labels:
>         app: tomcat
>     spec:
>       containers:
>       - name: tomcat
>         image: harbor.alnk.com/public/tomcat:7.0
>         # 添加相应的环境变量
>         # 下面收集了两块日志1、stdout 2、/usr/local/tomcat/logs/catalina.*.log
>         env:    
>         # 如日志发送到es，那index名称为tomcat-syslog,如日志发送到kafka，那topic则为tomcat-syslog
>         - name: aliyun_logs_tomcat-syslog   
>           value: "stdout"
>         # 如日志发送到es，那index名称为tomcat-access，如日志发送到kafka，那topic则为tomcat-access
>         - name: aliyun_logs_tomcat-access   
>           value: "/usr/local/tomcat/logs/catalina.*.log"
>         # 对pod内要收集的业务日志目录需要进行共享，可以收集多个目录下的日志文件
>         volumeMounts:   
>           - name: tomcat-log
>             mountPath: /usr/local/tomcat/logs
>       volumes:
>         - name: tomcat-log
>           emptyDir: {}
> ```
>
> `tomcat日志收集成功`
>
> ![1731461197383](images/1731461197383.png)  



#### 新业务上线收集日志

> `main.go`
>
> ```go
> package main
> 
> import (
> 	"fmt"
> 	"log"
> 	"net/http"
> 	"os"
> 	"path/filepath"
> 
> 	"github.com/gin-gonic/gin"
> )
> 
> var LogToFileAndStdout func(format string, v ...interface{})
> 
> func main() {
> 	r := gin.Default()
> 
> 	r.GET("/hello", func(c *gin.Context) {
> 		log.Printf("这是直接打印到std的日志, 请求URL: %s, 请求方法Method: %s, 请求主机地址: %s\n", c.Request.URL, c.Request.Method, c.Request.Host)
> 		LogToFileAndStdout("请求URL: %s, 请求方法Method: %s, 请求主机地址: %s\n", c.Request.URL, c.Request.Method, c.Request.Host)
> 		c.String(http.StatusOK, "hello, Gin! 新的业务上线了 version: 0.1")
> 	})
> 
> 	r.Run(":3000")
> }
> 
> func init() {
> 	// 当前程序执行的目录
> 	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
> 	if err != nil {
> 		panic(err)
> 	}
> 
> 	// 判断日志目录是否存在
> 	logDirPath := fmt.Sprintf("%s/%s", dir, "log")
> 	_, err = os.Stat(logDirPath)
> 	if err != nil {
> 		// 不存在则创建目录
> 		os.MkdirAll(dir, 0755)
> 	}
> 
> 	// 创建一个文件，用于写入日志
> 	file, err := os.OpenFile(fmt.Sprintf("%s/%s", logDirPath, "go-gin-log.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
> 	if err != nil {
> 		panic(err)
> 	}
> 	// defer file.Close()
> 
> 	// 创建一个日志写入文件的Logger
> 	fileLogger := log.New(file, "", log.LstdFlags)
> 
> 	// 同时创建一个标准Logger，用于输出屏幕
> 	// consoleLogger := log.New(os.Stdout, "", log.LstdFlags)
> 
> 	// 定义一个写入日志的函数，同时向文件和屏幕输出
> 	LogToFileAndStdout = func(format string, v ...interface{}) {
> 		fileLogger.Printf(format, v...)
> 		// consoleLogger.Printf(format, v...)
> 	}
> }
> 
> ```
>
> `Dockerfile`
>
> ```dockerfile
> FROM harbor.alnk.com/public/golang:1.22.8 as builder
> ENV GOPROXY https://goproxy.cn
> COPY . /app/
> RUN cd /app \
>     && go mod init go-gin-log \
>     && go mod tidy \
>     && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o go-gin-log .
> 
> FROM harbor.alnk.com/public/alpine:3.18
> RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
>     apk update && \
>     apk --no-cache add tzdata ca-certificates && \
>     cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
>     # apk del tzdata && \
>     rm -rf /var/cache/apk/*
> COPY --from=builder /app/go-gin-log /app/go-gin-log
> WORKDIR /app/
> EXPOSE 3000
> CMD ["./go-gin-log"] 
> 
> # docker build -t harbor.alnk.com/public/go-gin-log:0.4 .
> # docker push harbor.alnk.com/public/go-gin-log:0.4
> # docker run --rm -it harbor.alnk.com/public/go-gin-log:0.4 sh
> ```
>
> ```shell
> 【10.0.1.21】
> ### 创建目录
> # mkdir /data/k8s-yaml/go-gin-log -p && cd /data/k8s-yaml/go-gin-log
> # vi go-gin-log.yaml
> 
> ## 应用
> # kubectl apply -f go-gin-log.yaml
> ```
>
> `go-gin-log.yaml`
>
> ```yaml
> ---
> # deployment.yaml
> kind: Deployment
> apiVersion: apps/v1
> metadata:
>   name: go-gin-log-k8s
>   namespace: prod
>   labels:
>     name: go-gin-log-k8s
> spec:
>   replicas: 1
>   selector:
>     matchLabels:
>       name: go-gin-log-k8s
>   template:
>     metadata:
>       labels:
>         app: go-gin-log-k8s
>         name: go-gin-log-k8s
>     spec:
>       containers:
>       - name: go-gin-log-k8s
>         image: harbor.alnk.com/public/go-gin-log:0.4
>         # 添加相应的环境变量
>         # 下面收集了两块日志
>         # go-gin-std:屏幕打印日志
>         # go-gin-log:/app目录下所有*.log文件日志
>         # 如日志发送到es，那index名称为go-gin-std,如日志发送到kafka，那topic则为go-gin-std
>         # 如日志发送到es，那index名称为go-gin-log,如日志发送到kafka，那topic则为go-gin-log
>         env: 
>         - name: aliyun_logs_go-gin-std
>           value: "stdout"
>         - name: aliyun_logs_go-gin-log
>           value: "/app/log/*.log"
>         ports:
>         - containerPort: 3000
>           protocol: TCP
>         terminationMessagePath: /dev/termination-log
>         terminationMessagePolicy: File
>         imagePullPolicy: IfNotPresent
>         # 对pod内要收集的业务日志目录需要进行共享，可以收集多个目录下的日志文件
>         volumeMounts:
>           - name: gin-log
>             mountPath: /app/log/
>       volumes: 
>         - name: gin-log
>           emptyDir: {}
>     
>       imagePullSecrets:
>       - name: harbor
>       restartPolicy: Always
>       terminationGracePeriodSeconds: 30
>       securityContext:
>         runAsUser: 0
>       schedulerName: default-scheduler
>   strategy:
>     type: RollingUpdate
>     rollingUpdate:
>       maxUnavailable: 1
>       maxSurge: 1
>   revisionHistoryLimit: 7
>   progressDeadlineSeconds: 600
> 
> ---
> # service.yaml
> kind: Service
> apiVersion: v1
> metadata:
>   name: go-gin-log-k8s
>   namespace: prod
>   labels:
>     go-app: go-gin-log-k8s
> spec:
>   ports:
>   - protocol: TCP
>     port: 80
>     targetPort: 3000
>     name: http
>   selector:
>     app: go-gin-log-k8s
> 
> ---
> # ingress.yaml
> apiVersion: networking.k8s.io/v1
> kind: Ingress
> metadata:
>   namespace: prod
>   name: go-gin-log-k8s
> spec:
>   rules:
>   - host: go-gin-log-k8s.alnk.com
>     http:
>       paths:
>       - backend:
>           service:
>             name: go-gin-log-k8s
>             port:
>               number: 80
>         path: /
>         pathType: Prefix
> ```
>
> `修改log-pilot的配置文件`
>
> ```yaml
> apiVersion: v1
> kind: ConfigMap
> metadata:
>   name: log-pilot2-configuration
>   namespace: elk
> data:
>   logging_output: "kafka" # 采集日志输出到kafka
>   kafka_brokers: "10.0.1.21:9092" # kafka地址,多个地址可以写成"kafka1:9092,kafka2:9092"
>   kafka_version: "0.10.0" # 指定版本，生产中实测kafka 2.12-2.5 及以下，这里都配置为"0.10.0" 就可以了
>   # 当禁止自动创建topic是，这里配置kakfa中的有效topic
>   # 对应的业务容器收集日志需要传入这主题
>   kafka_topics: "tomcat-syslog,tomcat-access,go-gin-std,go-gin-log"
> ```
>
> `重启log-pilot服务,然后去kafka-ui中页面看是否有相应的topic`
>
> ![1731471703297](images/1731471703297.png)  
>
> `修改logstash配置，把日志从kafka存储到es中去`
>
> `修改logstash的configmap文件，添加go-gin-log的解析，重启logstash服务`
>
> ```yaml
> apiVersion: v1
> kind: ConfigMap
> metadata:
>   namespace: elk
>   name: logstash-configmap
> data:
>   logstash.conf: |
>     input {
>       kafka {
>           bootstrap_servers => "10.0.1.21:9092" # kafka地址
>           auto_offset_reset => "latest"  # 从最新的偏移量开始消费
>           consumer_threads => 1 
>           # 此属性会将当前topic、offset、group、partition等信息也带到message中
>           decorate_events => true  
>           topics_pattern  => "tomcat-.*" # 匹配以tomcat开头的topic
>           codec => "json"
>           group_id => "logstash" # 消费组id，如果需要重新从头消费的话，可更换id
>       }
>       ### 以下为添加内容 ###
>       # 匹配go-gin-log的日志
>       kafka {
>           bootstrap_servers => "10.0.1.21:9092" # kafka地址
>           auto_offset_reset => "latest"  # 从最新的偏移量开始消费
>           consumer_threads => 1 
>           # 此属性会将当前topic、offset、group、partition等信息也带到message中
>           decorate_events => true  
>           topics_pattern  => "go-gin-.*" # 匹配以go-gin开头的topic
>           codec => "json"
>           group_id => "go-gin-logstash" # 消费组id，如果需要重新从头消费的话，可更换id
>       }
>       ### 以上为添加内容 ###
>     }
>   
>     filter {
>     }
> 
>     output {
>         elasticsearch {
>           index => "%{[@metadata][kafka][topic]}-%{+YYYY-MM-dd}" 
>           hosts => "http://quickstart-es-http:9200" # es的地址，这里直接用svc名称
>           user => "elastic" # es账号
>           password => "${ELASTICSEARCH_PASSWORD}" # es密码
> 
>         }
>         stdout {
>           codec => rubydebug # 往控制台也打印收集到的日志
>         }
>     }
> 
> 
> ```
>
> `查看logstash的日志输出，已经解析到go-gin-log`
>
> ![1731505222435](images/1731505222435.png)  
>
> `kafka-ui上查看，已经多出一个go-gin-logstash的消费者`
>
> ![1731505255229](images/1731505255229.png)
>
> `kibana重建索引`
>
> ![1731505357599](images/1731505357599.png) 
>
> ![1731505411012](images/1731505411012.png)  
>
> ![1731505445612](./images/1731505445612.png)  

