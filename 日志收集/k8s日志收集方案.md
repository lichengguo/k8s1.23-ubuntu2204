### k8s日志收集方案

#### k8s环境

> k8s部署工具：kubeasz
>
> k8s 版本：1.23.1
>
> 容器运行时： docker
>
> 操作系统 ：ubuntu2204
>
> ##### 服务器规划
>
> | **IP**         | **主机名称** | **角色**        | **系统**       | **软件**                                                     | **配置**  |
> | -------------- | ------------ | --------------- | -------------- | ------------------------------------------------------------ | --------- |
> | **10.0.1.21**  | **ops**      | **运维机**      | ubuntu2204     | **harbor仓库、kubeasz**                                      | **2c/4g** |
> | **10.0.1.100** | **虚拟IP**   | **/**           | **/**          | **流量入口、负载均衡、高可用、七层反向代理**                 | **/**     |
> | **10.0.1.101** | **ha-1**     | **反向代理**    | **ubuntu2204** | **nginx、keepalived**                                        | **1c/2g** |
> | **10.0.1.102** | **ha-2**     | **反向代理**    | **ubuntu2204** | **nginx、keepalived**                                        | **1c/2g** |
> | **10.0.1.200** | **虚拟IP**   | **/**           | **/**          | **apiserver高可用、4层反向代理**                             | **/**     |
> | **10.0.1.201** | **master-1** | **k8s主节点**   | **ubuntu2204** | **apiserver、controller、scheduler、etcd、keepalived、nginx(l4lb)** | **2c/4g** |
> | **10.0.1.202** | **master-2** | **k8s主节点**   | **ubuntu2204** | **apiserver、controller、scheduler、etcd、keepalived、nginx(l4lb)** | **2c/4g** |
> | **10.0.1.203** | **node-1**   | **k8s工作节点** | **ubuntu2204** | **kubelet、kube-proxy**、**etcd**                            | **2c/8g** |
> | **10.0.1.204** | **node-2**   | **k8s工作节点** | **ubuntu2204** | **kubelet、kube-proxy**                                      | **2c/8g** |
>
> `如果使用公有云、10.0.1.100、10.0.1.101、10.0.1.102可直接使用SLB即可`



#### 方案介绍

> k8s官网日志收集方案地址 ：https://kubernetes.io/zh/docs/concepts/cluster-administration/logging/
>
> ```text
> 1-边车模式
> 使用sidcar容器(一个pod多容器)收集当前pod内一个或者多个业务容器的日志(通常基于emptyDir实现业务容器与sidcar之间的日志共享)
> 该方式的优点：对当前收集日志的pod能做更详细的处理，按照预定义规则处理好的日志，再输出到kafka或者filebeat
> 该方式的缺点：性能可能没那么好，因为一个pod就需要一个sidecar，如果有100个pod，那就需要100个sidecar，那么这个日志收集的agent就必须选用一个轻量级的工具，不然会占用很多额外资源
> 
> 2-node节点收集
> 基于daemonset部署日志收集进程，实现json-file类型(标准输出/dev/stdout、错误输出/dev/stderr)日志收集；
> 该方式的优点：日志收集架构简单，易部署、易维护
> 该方式的缺点：node节点产生的日志、路径、类型、日志内容都不同，很难做到把每个容器的日志都完整的处理好
> 
> 3-在业务容器中内置日志收集服务进程
> 就是除了业务容器，再另外起一个日志收集的进程，然后直接把日志推送到远端存储
> ```
>
> ![1731227966562](images\1731227966562.png)  
>
> ![1731227996239](images\1731227996239.png)  





#### operator部署elasticsearch和kibana

> 使用operator部署生产级别的Elasticserach集群和Kibana
>
> 可参考 https://www.cuiliangblog.cn/detail/article/71



##### 查看版本支持

> https://www.elastic.co/guide/en/cloud-on-k8s/2.7/k8s-quickstart.html
>
> ![1731228390645](images\1731228390645.png)  



##### 部署eck-operator

```shell
【10.0.1.21】
## 默认情况下eck-operator运行在elastic-system名称空间

## eck镜像准备（github开源项目把镜像同步到阿里云仓库）
# docker pull registry.cn-hangzhou.aliyuncs.com/alnktest/eck-operator:2.7.0
# docker tag registry.cn-hangzhou.aliyuncs.com/alnktest/eck-operator:2.7.0 harbor.alnk.com/public/eck-operator:2.7.0 
# docker push harbor.alnk.com/public/eck-operator:2.7.0 

## 下载yaml文件
# mkdir /data/k8s-yaml/eck -p
# cd /data/k8s-yaml/eck
# wget https://download.elastic.co/downloads/eck/2.7.0/crds.yaml
# wget https://download.elastic.co/downloads/eck/2.7.0/operator.yaml

## 替换镜像
# sed -ri 's+docker.elastic.co/.*/+harbor.alnk.com/public/+g' operator.yaml

## 部署
# kubectl create -f crds.yaml
# kubectl apply -f operator.yaml

# 查看资源
# kubectl -n elastic-system get all
pod/elastic-operator-0   1/1     Running   0          68s
service/elastic-webhook-server   ClusterIP   10.68.159.142   <none>        443/TCP   68s
statefulset.apps/elastic-operator   1/1     68s

## 查看日志
# kubectl -n elastic-system logs -f statefulset.apps/elastic-operator
# kubectl -n elastic-system logs -f elastic-operator-0
```



##### 部署elasticsearch

```shell
【10.0.1.21】
## elasticsearch镜像准备（github开源项目把镜像同步到阿里云仓库） ##
# docker pull registry.cn-hangzhou.aliyuncs.com/alnktest/elasticsearch:6.8.6
# docker tag registry.cn-hangzhou.aliyuncs.com/alnktest/elasticsearch:6.8.6 harbor.alnk.com/public/elasticsearch:6.8.6
# docker push harbor.alnk.com/public/elasticsearch:6.8.6

## yaml文件 ##
# mkdir /data/k8s-yaml/elk -p
# cd /data/k8s-yaml/elk
# vi elasticsearch.yaml

【10.0.1.21】
## 部署应用 ##
# kubectl create ns elk
# kubectl apply -f elasticsearch.yaml

# kubectl -n elk get pod

# 查看es集群是否健康 ##
# kubectl -n elk  get elasticsearches.elasticsearch.k8s.elastic.co quickstart
# kubectl -n elk  get es quickstart
quickstart   green    1       8.11.3    Ready   62s

## 查看es密码 ##
# kubectl -n elk get secret quickstart-es-elastic-user -o go-template='{{.data.elastic | base64decode}}'
25y69WWe1FC8F6Gh6iVHVt53

## 查看es的svc地址 ##
# kubectl get svc -n elk
quickstart-es-http            ClusterIP   10.68.120.223   <none>        9200/TCP   56s

【10.0.1.201】
## 访问es集群 ##
# curl -k http://elastic:25y69WWe1FC8F6Gh6iVHVt53@10.68.120.223:9200
## 查看索引列表
# curl -k http://elastic:25y69WWe1FC8F6Gh6iVHVt53@10.68.120.223:9200/_cat/indices?v

```

`elasticsearch.yaml`

```yaml
apiVersion: elasticsearch.k8s.elastic.co/v1
kind: Elasticsearch
metadata:
  name: quickstart # 名称
  namespace: elk # 名称空间
spec:
  http: # 关闭https，使用http对外提供服务，否则需要配置证书
    tls:
      selfSignedCertificate:
        disabled: true
  version: 6.8.6 # es版本
  image: harbor.alnk.com/public/elasticsearch:6.8.6 #自定义镜像地址，不指定则从官方镜像仓库拉取
  nodeSets:  
  - name: default # 节点名称
    count: 1 # 集群数量，建议奇数个1,3,5
    config:
      node.store.allow_mmap: false
    podTemplate:
      spec:
        containers:
        - name: elasticsearch
          env:
          - name: ES_JAVA_OPTS # jvm参数
            value: -Xms1g -Xmx1g # 指定节点JVM大小 
          resources: # 根据实际情况调整资源限制
            requests:     # 资源请求值，通常与JVM保持一致
              memory: 1Gi
              cpu: 0.5
            limits:        # 资源限制值，通常为JVM的2倍
              memory: 2Gi
              cpu: 1
        initContainers:
        - name: sysctl
          securityContext:
            privileged: true
          command: ['sh', '-c', 'sysctl -w vm.max_map_count=262144']
    volumeClaimTemplates:
    - metadata:
        name: elasticsearch-data
      spec:
        accessModes:
        - ReadWriteOnce
        resources:
          requests:
            storage: 10Gi # 存储空间，根据实际情况调整
        storageClassName: nfs-alnk # 建议固态硬盘
  #http:
    #tls:
      #selfSignedCertificate:
        #disabled: true 
        
---
##测试为了方便加上ingress，生产环境中可以不加
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: elastic
  namespace: elk
  annotations:
    #nginx.ingress.kubernetes.io/force-ssl-redirect: "true" #不使用https访问
    kubernetes.io/ingress.class: nginx
spec:
  #tls:
  #- hosts:
    #- kibana.alnk.com
    #secretName: alnk-com-tls
  rules:
  - host: elastic.alnk.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: quickstart-es-http
            port:
              number: 9200
```

![1731394188425](images\1731394188425.png)  

  `http://elastic:25y69WWe1FC8F6Gh6iVHVt53@elastic.alnk.com/`

 ![1731420011582](images\1731420011582.png)  



##### 部署kibana

```shell
【10.0.1.21】
##kibana镜像准备 github开源项目把镜像同步到阿里云仓库##
# docker pull registry.cn-hangzhou.aliyuncs.com/alnktest/kibana:6.8.6
# docker tag registry.cn-hangzhou.aliyuncs.com/alnktest/kibana:6.8.6 harbor.alnk.com/public/kibana:6.8.6
# docker push harbor.alnk.com/public/kibana:6.8.6

## yaml文件 ##
# mkdir /data/k8s-yaml/elk -p
# cd /data/k8s-yaml/elk
# vi kibana.yaml

## 应用 ##
# kubectl apply -f kibana.yaml

## 查看 ##
# kubectl -n elk get pod
# kubectl -n elk get pod -w
# kubectl -n elk get svc
# kubectl -n elk get ingress
```

`kibana.yaml`

```yaml
apiVersion: kibana.k8s.elastic.co/v1
kind: Kibana
metadata:
  name: quickstart # 名称
  namespace: elk # 名称空间
spec:
  #增加该配置，不然出bug连接不上es
  #https://github.com/elastic/cloud-on-k8s/issues/5952
  #config:
    #monitoring.ui.ccs.enabled: false
  version: 6.8.6
  image: harbor.alnk.com/public/kibana:6.8.6
  count: 1 # kibana集群个数
  elasticsearchRef: # 关联es集群名称,与Elasticsearch资源名称匹配
    name: quickstart
  http: # 关闭https，使用http对外提供服务，否则需要配置证书
    tls:
      selfSignedCertificate:
        disabled: true
  podTemplate:
    spec:
      containers:
      - name: kibana
        env:
          - name: NODE_OPTIONS
            value: "--max-old-space-size=2048"
          - name: SERVER_PUBLICBASEURL
            value: "https://kibana.alnk.com"
          - name: I18N_LOCALE # 中文配置
            value: "zh-CN"
        resources:
          requests:
            memory: 1Gi
            cpu: 0.5
          limits:
            memory: 2Gi
            cpu: 2  
        
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: kibana
  namespace: elk
  annotations:
    #nginx.ingress.kubernetes.io/force-ssl-redirect: "true" #不使用https访问
    kubernetes.io/ingress.class: nginx
spec:
  #tls:
  #- hosts:
    #- kibana.alnk.com
    #secretName: alnk-com-tls
  rules:
  - host: kibana.alnk.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: quickstart-kb-http
            port:
              number: 5601
```

`http://kibana.alnk.com/`

`账号密码 elastic:25y69WWe1FC8F6Gh6iVHVt53`

  ![1731442773697](images\1731442773697.png)      



#### Node模式日志收集方案

> 收集日志的框架流程
>
> ```text
> 模拟业务服务   日志收集       日志存储     消费日志         消费目的服务
> tomcat ---> log-pilot ---> kafka ---> logstash ---> elasticsearch ---> kibana
> ```



##### docker-compose部署kafka

> Kafka是一种流行的分布式消息传递系统，设计用于处理大规模的数据流。它在日志传输和流数据处理领域尤为重要。以下是Kafka的一些关键特性，以及它们如何适用于日志传输体系：
>
> 高吞吐量和可伸缩性
>
> - Kafka能够处理高速数据流，支持大量数据的读写操作。它通过分区（partitions）和分布式架构来实现高吞吐量和水平扩展
> - 在日志传输中，这意味着Kafka可以高效地收集和分发大量的日志数据，满足大型系统的需求
>
> 持久性和可靠性
>
> - Kafka将数据存储在磁盘上，并支持数据的复制，确保信息不会因为单点故障而丢失
> - 对于日志管理来说，这确保了日志数据的安全和完整性，即使在系统故障的情况下也能保持日志的可用性
>
> 实时处理
>
> - Kafka支持实时数据流处理，使得数据在产生后立即可用
> - 在日志系统中，这允许实时监控和分析日志数据，有助于快速故障排除和系统性能监控
>
> 分布式架构
>
> - Kafka的分布式特性意味着它可以跨多个服务器和数据中心运行，提高了容错能力和系统的稳定性
> - 在分布式日志系统中，这有助于管理跨地域或多个服务的日志数据，提供更好的可扩展性和冗余
>
> 灵活的消费模型
>
> - Kafka支持发布-订阅和队列两种消息模式，使消费者能够灵活地读取数据
> - 日志传输中，这使得不同的系统和应用可以根据需要订阅特定类型的日志数据
>
> 低延迟
>
> - Kafka设计用于低延迟消息传递，这对于需要快速响应的日志分析和监控系统至关重要
>
> 结合日志传输体系，Kafka的这些特性使其成为一个理想的中央日志聚合点。它不仅能够高效地收集来自不同源的日志数据，还能将这些数据实时分发给各种日志分析和监控工具，支持大数据分析、实时监控和快速故障检测
>
> 在Apache Kafka中，几个核心概念包括主题（Topic）、分区（Partition）和消费组（Consumer Group）。这些概念是Kafka高效消息处理和分布式数据流管理的基石

```shell
【10.0.1.21】
### 准备镜像
## kafka 2.2.0
# docker pull registry.cn-hangzhou.aliyuncs.com/alnktest/kafka:2.2.0
# docker tag  registry.cn-hangzhou.aliyuncs.com/alnktest/kafka:2.2.0  harbor.alnk.com/public/kafka:2.2.0
# docker push harbor.alnk.com/public/kafka:2.2.0

## zookeeper 3.4.13
# docker pull registry.cn-hangzhou.aliyuncs.com/alnktest/zookeeper:3.4.13
# docker tag registry.cn-hangzhou.aliyuncs.com/alnktest/zookeeper:3.4.13 harbor.alnk.com/public/zookeeper:3.4.13
# docker push harbor.alnk.com/public/zookeeper:3.4.13

## kafka-ui
# docker pull registry.cn-hangzhou.aliyuncs.com/alnktest/kafka-ui:v0.7.2
# docker tag registry.cn-hangzhou.aliyuncs.com/alnktest/kafka-ui:v0.7.2 harbor.alnk.com/public/kafka-ui:v0.7.2 
# docker push harbor.alnk.com/public/kafka-ui:v0.7.2 

### 创建目录
## 创建数据目录
# mkdir /data/zookeeper/data -p
# mkdir /data/zookeeper/datalog -p
# mkdir /data/zookeeper/logs -p
# mkdir /data/kafka1/data -p

## docker-compose文件目录
# mkdir -p /opt/kafka-2.2.0 
# cd /opt/kafka-2.2.0
# vi docker-compose.yml
version: '3.3'
services:
  zookeeper:
    image: harbor.alnk.com/public/zookeeper:3.4.13
    container_name: zookeeper
    ports:
      - 2181:2181
    volumes:
      - /data/zookeeper/data:/data
      - /data/zookeeper/datalog:/datalog
      - /data/zookeeper/logs:/logs
    restart: unless-stopped

  kafka1:
    image: harbor.alnk.com/public/kafka:2.2.0
    depends_on:
      - zookeeper
    container_name: kafka1
    ports:
      - 9092:9092
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: 10.0.1.21:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://10.0.1.21:9092
      KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9092
      KAFKA_LOG_DIRS: /data/kafka-data
      KAFKA_LOG_RETENTION_HOURS: 24
      ALLOW_PLAINTEXT_LISTENER: yes
    volumes:
      - /data/kafka1/data:/data/kafka-data
    restart: unless-stopped
  kafka-ui:
    container_name: kafka-ui
    image: harbor.alnk.com/public/kafka-ui:v0.7.2
    ports:
      - "8000:8080"
    environment:
      DYNAMIC_CONFIG_ENABLED: true
    restart: unless-stopped
# KAFKA_ZOOKEEPER_CONNECT: zookeeper服务地址
# KAFKA_ADVERTISED_LISTENERS: kafka服务地址
# KAFKA_BROKER_ID: kafka节点ID，不可重复
# KAFKA_LOG_DIRS: kafka数据文件地址(非必选。下面的volumes会将容器内文件挂载到宿主机上)
# KAFKA_LOG_RETENTION_HOURS: 数据文件保留时间(非必选。缺省168小时)
# kafka2:
#     image: harbor.alnk.com/public/kafka:2.2.0
#     depends_on:
#       - zookeeper
#     container_name: kafka2
#     ports:
#       - 9093:9093
#     environment:
#       KAFKA_BROKER_ID: 2
#       KAFKA_ZOOKEEPER_CONNECT: 10.0.1.21:2181
#       KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://10.0.1.21:9093
#       KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9093
#       KAFKA_LOG_DIRS: /data/kafka-data
#       KAFKA_LOG_RETENTION_HOURS: 24
#     volumes:
#       - /data/kafka2/data:/data/kafka-data
#     restart: unless-stopped
#   kafka3:
#     image: harbor.alnk.com/public/kafka:2.2.0
#     depends_on:
#       - zookeeper
#     container_name: kafka3
#     ports:
#       - 9094:9094
#     environment:
#       KAFKA_BROKER_ID: 3
#       KAFKA_ZOOKEEPER_CONNECT: 10.0.1.21:2181
#       KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://10.0.1.21:9094
#       KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9094
#       KAFKA_LOG_DIRS: /data/kafka-data
#       KAFKA_LOG_RETENTION_HOURS: 24
#     volumes:
#       - /data/kafka3/data:/data/kafka-data
#     restart: unless-stopped

## 启动kafka
# docker-compose up -d

## 查看结果
# docker-compose ps
                  
## 创建生产者
# docker exec -it kafka1 bash
$ kafka-console-producer.sh --broker-list 10.0.1.21:9092 --topic tomcat-syslog
>aaa
>badfa
>
$ kafka-console-producer.sh --broker-list 10.0.1.21:9092 --topic tomcat-access

## 创建消费者
# docker exec -it kafka bash
$ kafka-console-consumer.sh --bootstrap-server 10.0.1.21:9092 --from-beginning --topic tomcat-syslog
aaa
badfa

```

`http://10.0.1.21:8000/`

![1731414846763](images\1731414846763.png)    

​    ![1731443122183](images\1731443122183.png)  



​    

##### 部署logstash

```shell
【10.0.1.21】
### 准备镜像
## logstash的版本需要和es的版本一样
# docker pull registry.cn-hangzhou.aliyuncs.com/alnktest/logstash:6.8.6
# docker tag  registry.cn-hangzhou.aliyuncs.com/alnktest/logstash:6.8.6  harbor.alnk.com/public/logstash:6.8.6
# docker push harbor.alnk.com/public/logstash:6.8.6

### 做个简单的测试，看链路 kafka  --> logstash --> elastic --> kibana 是否打通
# docker run --name logstash-test harbor.alnk.com/public/logstash:6.8.6 tail -f /dev/null
## 进入logstash-test容器
# docker exec -it -u root logstash-test bash
## 以下操作是在容器进行
# cd /root
# vi logstatsh-test.conf
input {
  kafka {
      bootstrap_servers => "10.0.1.21:9092"
      auto_offset_reset => "latest" 
      consumer_threads => 1 
      decorate_events => true 
      topics_pattern  => "tomcat-.*"
      codec => "json"
      group_id => "logstash"
      decorate_events => true
  }
}

filter {
}

output {
    elasticsearch {
      index => "%{[@metadata][kafka][topic]}-%{+YYYY-MM-dd}" 
      hosts => "http://elastic.alnk.com:80" 
      user => "elastic"
      password => "25y69WWe1FC8F6Gh6iVHVt53"

    }
    stdout {
      codec => rubydebug
    }
} 
## 添加下解析，不然docker容器内访问不到域名
# vi /etc/hosts
10.0.1.100 elastic.alnk.com

## 启动logstash测试
# logstash -f ./logstatsh-test.conf
```

> `直接在kafka ui 创建消息，看es是否能收到`
>
> ![1731443709537](images\1731443709537.png)  
>
> ![1731443740702](images\1731443740702.png)  
>
> `访问kibana`
>
> ![1731443786504](images\1731443786504.png)  
>
> `可以看到已经有数据传递过来了`
>
> ![1731443835960](images\1731443835960.png)  
>
> `创建索引`
>
> ![1731443904901](images\1731443904901.png)    
>
> `匹配所有以tomcat开头的记录`
>
> ![1731443928524](images\1731443928524.png)
>
> ![1731443961728](images\1731443961728.png)  
>
> `Discover:可以看到已经收集到日志`
>
> ![1731444012176](images\1731444012176.png) 



```shell
【10.0.1.21】
### 准备yaml文件，把logstash以deployment的形式部署到k8s集群中去
# mkdir -p /data/k8s-yaml/elk -p
# cd /data/k8s-yaml/elk
# vi logstash.yaml

## 应用
# kubectl apply -f logstash.yaml
```

`logstash.yaml`

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  namespace: elk
  name: logstash-configmap
data:
  logstash.conf: |
    input {
      kafka {
          bootstrap_servers => "10.0.1.21:9092" # kafka地址
          auto_offset_reset => "latest"  # 从最新的偏移量开始消费
          consumer_threads => 1 
          # 此属性会将当前topic、offset、group、partition等信息也带到message中
          decorate_events => true  
          topics_pattern  => "tomcat-.*" # 匹配以tomcat开头的topic
          codec => "json"
          group_id => "logstash" # 消费组id，如果需要重新从头消费的话，可更换id
      }
    }

    filter {
    }

    output {
        elasticsearch {
          index => "%{[@metadata][kafka][topic]}-%{+YYYY-MM-dd}" 
          hosts => "http://quickstart-es-http:9200" # es的地址，这里直接用svc名称
          user => "elastic" # es账号
          password => "${ELASTICSEARCH_PASSWORD}" # es密码

        }
        stdout {
          codec => rubydebug # 往控制台也打印收集到的日志
        }
    }

---
# es密码
apiVersion: v1
kind: Secret
metadata:
  name: elasticsearch-pw-elastic
  namespace: elk
type: Opaque
data:
  # https://www.jb51.net/tools/base64.htm
  # 这里是es密码25y69WWe1FC8F6Gh6iVHVt53通过base64后加密的字符串
  password: MjV5NjlXV2UxRkM4RjZHaDZpVkhWdDUz

---
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: elk
  name: logstash
spec:
  replicas: 1
  selector:
    matchLabels:
      app: logstash
  template:
    metadata:
      labels:
        app: logstash
    spec:
      containers:
      - name: logstash
        image: harbor.alnk.com/public/logstash:6.8.6
        resources:
          limits:
            cpu: "1"
            memory: 1Gi
        args:
        - -f
        - /usr/share/logstash/config/logstash.conf
        ports:
        - containerPort: 9600
          name: logstash
        env:
        #- name: XPACK_MONITORING_ENABLED
          #value: "false"
        - name: ELASTICSEARCH_PASSWORD
          valueFrom:
            secretKeyRef:
              name: elasticsearch-pw-elastic
              key: password
        #- name: ELASTICSEARCH_URL
          #value: "http://quickstart-es-http:9200"
        #- name: ELASTICSEARCH_USER
          #value: "elastic"
        #- name: ELASTICSEARCH_PASSWORD
          #value: "37qSPrb68X9Oj641R6qw3nYr"
        volumeMounts:
        - name: logstash-conf
          mountPath: /usr/share/logstash/config/logstash.conf
          subPath: logstash.conf       
        - name: localtime
          mountPath: /etc/localtime
      volumes:
      - name: logstash-conf
        configMap:
          name: logstash-configmap
      - name: localtime
        hostPath:
          path: /etc/localtime
          type: File
```





##### 部署log-pilot

> 参考文档地址：https://github.com/AliyunContainerService/log-pilot

```shell
【10.0.1.21】
## 准备镜像
# docker pull registry.cn-hangzhou.aliyuncs.com/acs/log-pilot:0.9.7-filebeat
# docker tag registry.cn-hangzhou.aliyuncs.com/acs/log-pilot:0.9.7-filebeat harbor.alnk.com/public/log-pilot:0.9.7-filebeat
# docker push harbor.alnk.com/public/log-pilot:0.9.7-filebeat

## 创建目录
# mkdir -p /data/k8s-yaml/elk && cd /data/k8s-yaml/elk
# vi log-pilot.yaml

## 应用
# kubectl apply -f log-pilot.yaml
# kubectl -n ns-elastic get pod
```

`log-pilot.yaml`

```yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: log-pilot2-configuration
  namespace: elk
data:
  logging_output: "kafka" # 采集日志输出到kafka
  kafka_brokers: "10.0.1.21:9092" # kafka地址,多个地址可以写成"kafka1:9092,kafka2:9092"
  kafka_version: "0.10.0" # 指定版本，生产中实测kafka 2.12-2.5 及以下，这里都配置为"0.10.0" 就可以了
  # 当禁止自动创建topic是，这里配置kakfa中的有效topic
  # 对应的业务容器收集日志需要传入这主题
  kafka_topics: "tomcat-syslog,tomcat-access"

---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: log-pilot2
  namespace: elk
  labels:
    k8s-app: log-pilot2
spec:
  selector:
    matchLabels:
      k8s-app: log-pilot2
  updateStrategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        k8s-app: log-pilot2
    spec:
      tolerations:
      - key: node-role.kubernetes.io/master
        effect: NoSchedule
      containers:
      - name: log-pilot2
        image: harbor.alnk.com/public/log-pilot:0.9.7-filebeat
        env:
          - name: "LOGGING_OUTPUT"
            valueFrom:
              configMapKeyRef:
                name: log-pilot2-configuration
                key: logging_output
          - name: "KAFKA_BROKERS"
            valueFrom:
              configMapKeyRef:
                name: log-pilot2-configuration
                key: kafka_brokers
          - name: "KAFKA_VERSION"
            valueFrom:
              configMapKeyRef:
                name: log-pilot2-configuration
                key: kafka_version
          - name: "NODE_NAME"
            valueFrom:
              fieldRef:
                fieldPath: spec.nodeName
        volumeMounts:
        - name: sock
          mountPath: /var/run/docker.sock
        - name: logs
          mountPath: /var/log/filebeat
        - name: state
          mountPath: /var/lib/filebeat
        - name: root
          mountPath: /host
          readOnly: true
        - name: localtime
          mountPath: /etc/localtime
        # configure all valid topics in kafka
        # when disable auto-create topic
        - name: config-volume
          mountPath: /etc/filebeat/config
        securityContext:
          capabilities:
            add:
            - SYS_ADMIN
      terminationGracePeriodSeconds: 30
      
      volumes:
      - name: sock
        hostPath:
          path: /var/run/docker.sock
          type: Socket
      - name: logs
        hostPath:
          path: /var/log/filebeat
          type: DirectoryOrCreate
      - name: state
        hostPath:
          path: /var/lib/filebeat
          type: DirectoryOrCreate
      - name: root
        hostPath:
          path: /
          type: Directory
      - name: localtime
        hostPath:
          path: /etc/localtime
          type: File
      # kubelet sync period
      - name: config-volume
        configMap:
          name: log-pilot2-configuration
          items:
          - key: kafka_topics
            path: kafka_topics
```



##### 部署tomcat模拟业务

```shell
【10.0.1.21】
## 创建目录
# mkdir -p /data/k8s-yaml/elk && cd /data/k8s-yaml/elk
# vi tomcat.yaml

## 准备镜像
# docker pull registry.cn-hangzhou.aliyuncs.com/alnktest/tomcat:7.0
# docker tag registry.cn-hangzhou.aliyuncs.com/alnktest/tomcat:7.0 harbor.alnk.com/public/tomcat:7.0
# docker push harbor.alnk.com/public/tomcat:7.0

## 应用
# kubectl apply -f tomcat.yaml

# 查看
# kubectl -n ns-elastic get pod

```

`tomcat.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: tomcat
  name: tomcat
  namespace: elk
spec:
  replicas: 1
  selector:
    matchLabels:
      app: tomcat
  template:
    metadata:
      labels:
        app: tomcat
    spec:
      containers:
      - name: tomcat
        image: harbor.alnk.com/public/tomcat:7.0
        # 添加相应的环境变量
        # 下面收集了两块日志1、stdout 2、/usr/local/tomcat/logs/catalina.*.log
        env:      
        # 如日志发送到es，那index名称为tomcat-syslog,如日志发送到kafka，那topic则为tomcat-syslog
        - name: aliyun_logs_tomcat-syslog   
          value: "stdout"
        # 如日志发送到es，那index名称为tomcat-access，如日志发送到kafka，那topic则为tomcat-access
        - name: aliyun_logs_tomcat-access   
          value: "/usr/local/tomcat/logs/catalina.*.log"
        # 对pod内要收集的业务日志目录需要进行共享，可以收集多个目录下的日志文件
        volumeMounts:   
          - name: tomcat-log
            mountPath: /usr/local/tomcat/logs
      volumes:
        - name: tomcat-log
          emptyDir: {}
```

`tomcat日志收集成功`

![1731461197383](images\1731461197383.png)  



##### 模拟新业务上线收集日志

`main.go`

```go
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

```

`Dockerfile`

```yaml
FROM harbor.alnk.com/public/golang:1.22.8 as builder
ENV GOPROXY https://goproxy.cn
COPY . /app/
RUN cd /app \
    && go mod init go-gin-log \
    && go mod tidy \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o go-gin-log .

FROM harbor.alnk.com/public/alpine:3.18
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \
    apk --no-cache add tzdata ca-certificates && \
    cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    # apk del tzdata && \
    rm -rf /var/cache/apk/*
COPY --from=builder /app/go-gin-log /app/go-gin-log
WORKDIR /app/
EXPOSE 3000
CMD ["./go-gin-log"] 

# docker build -t harbor.alnk.com/public/go-gin-log:0.4 .
# docker push harbor.alnk.com/public/go-gin-log:0.4
# docker run --rm -it harbor.alnk.com/public/go-gin-log:0.4 sh
```



```shell
【10.0.1.21】
### 创建目录
# mkdir /data/k8s-yaml/go-gin-log -p && cd /data/k8s-yaml/go-gin-log
# vi go-gin-log.yaml

## 应用
# kubectl apply -f go-gin-log.yaml
```

`go-gin-log.yaml`

```yaml
---
# deployment.yaml
kind: Deployment
apiVersion: apps/v1
metadata:
  name: go-gin-log-k8s
  namespace: prod
  labels:
    name: go-gin-log-k8s
spec:
  replicas: 1
  selector:
    matchLabels:
      name: go-gin-log-k8s
  template:
    metadata:
      labels:
        app: go-gin-log-k8s
        name: go-gin-log-k8s
    spec:
      containers:
      - name: go-gin-log-k8s
        image: harbor.alnk.com/public/go-gin-log:0.4
        # 添加相应的环境变量
        # 下面收集了两块日志
        # go-gin-std:屏幕打印日志
        # go-gin-log:/app目录下所有*.log文件日志
        # 如日志发送到es，那index名称为go-gin-std,如日志发送到kafka，那topic则为go-gin-std
        # 如日志发送到es，那index名称为go-gin-log,如日志发送到kafka，那topic则为go-gin-log
        env: 
        - name: aliyun_logs_go-gin-std
          value: "stdout"
        - name: aliyun_logs_go-gin-log
          value: "/app/log/*.log"
        ports:
        - containerPort: 3000
          protocol: TCP
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        imagePullPolicy: IfNotPresent
        # 对pod内要收集的业务日志目录需要进行共享，可以收集多个目录下的日志文件
        volumeMounts:
          - name: gin-log
            mountPath: /app/log/
      volumes: 
        - name: gin-log
          emptyDir: {}
      
      imagePullSecrets:
      - name: harbor
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      securityContext:
        runAsUser: 0
      schedulerName: default-scheduler
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 1
      maxSurge: 1
  revisionHistoryLimit: 7
  progressDeadlineSeconds: 600

---
# service.yaml
kind: Service
apiVersion: v1
metadata:
  name: go-gin-log-k8s
  namespace: prod
  labels:
    go-app: go-gin-log-k8s
spec:
  ports:
  - protocol: TCP
    port: 80
    targetPort: 3000
    name: http
  selector:
    app: go-gin-log-k8s

---
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  namespace: prod
  name: go-gin-log-k8s
spec:
  rules:
  - host: go-gin-log-k8s.alnk.com
    http:
      paths:
      - backend:
          service:
            name: go-gin-log-k8s
            port:
              number: 80
        path: /
        pathType: Prefix
```



`修改log-pilot的配置文件`

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: log-pilot2-configuration
  namespace: elk
data:
  logging_output: "kafka" # 采集日志输出到kafka
  kafka_brokers: "10.0.1.21:9092" # kafka地址,多个地址可以写成"kafka1:9092,kafka2:9092"
  kafka_version: "0.10.0" # 指定版本，生产中实测kafka 2.12-2.5 及以下，这里都配置为"0.10.0" 就可以了
  # 当禁止自动创建topic是，这里配置kakfa中的有效topic
  # 对应的业务容器收集日志需要传入这主题
  kafka_topics: "tomcat-syslog,tomcat-access,go-gin-std,go-gin-log"
```

`重启log-pilot服务,然后去kafka-ui中页面看是否有相应的topic`

![1731471703297](images\1731471703297.png)  



`修改logstash配置，把日志从kafka存储到es中去`

`修改logstash的configmap文件，添加go-gin-log的解析，重启logstash服务`

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  namespace: elk
  name: logstash-configmap
data:
  logstash.conf: |
    input {
      kafka {
          bootstrap_servers => "10.0.1.21:9092" # kafka地址
          auto_offset_reset => "latest"  # 从最新的偏移量开始消费
          consumer_threads => 1 
          # 此属性会将当前topic、offset、group、partition等信息也带到message中
          decorate_events => true  
          topics_pattern  => "tomcat-.*" # 匹配以tomcat开头的topic
          codec => "json"
          group_id => "logstash" # 消费组id，如果需要重新从头消费的话，可更换id
      }
      ### 以下为添加内容 ###
      # 匹配go-gin-log的日志
      kafka {
          bootstrap_servers => "10.0.1.21:9092" # kafka地址
          auto_offset_reset => "latest"  # 从最新的偏移量开始消费
          consumer_threads => 1 
          # 此属性会将当前topic、offset、group、partition等信息也带到message中
          decorate_events => true  
          topics_pattern  => "go-gin-.*" # 匹配以go-gin开头的topic
          codec => "json"
          group_id => "go-gin-logstash" # 消费组id，如果需要重新从头消费的话，可更换id
      }
      ### 以上为添加内容 ###
    }
    
    filter {
    }

    output {
        elasticsearch {
          index => "%{[@metadata][kafka][topic]}-%{+YYYY-MM-dd}" 
          hosts => "http://quickstart-es-http:9200" # es的地址，这里直接用svc名称
          user => "elastic" # es账号
          password => "${ELASTICSEARCH_PASSWORD}" # es密码

        }
        stdout {
          codec => rubydebug # 往控制台也打印收集到的日志
        }
    }


```



`查看logstash的日志输出，已经解析到go-gin-log`

![1731505222435](images\1731505222435.png)  



`kafka-ui上查看，已经多出一个go-gin-logstash的消费者`

![1731505255229](images\1731505255229.png)  

`kibana重建索引`

![1731505357599](images\1731505357599.png)  

![1731505411012](images\1731505411012.png)  

![1731505445612](./images/1731505445612.png)  





#### 边车模式日志收集方案

> 收集日志的框架流程
>
> ```text
> 模拟业务服务/日志收集       日志存储     消费日志         消费目的服务
> tomcat/filebeat     ---> kafka ---> logstash ---> elasticsearch ---> kibana
> ```



##### 模拟业务上线收集日志



```yaml
---
# deployment.yaml
kind: Deployment
apiVersion: apps/v1
metadata:
  name: go-gin-log-k8s
  namespace: prod
  labels:
    name: go-gin-log-k8s
spec:
  replicas: 1
  selector:
    matchLabels:
      name: go-gin-log-k8s
  template:
    metadata:
      labels:
        app: go-gin-log-k8s
        name: go-gin-log-k8s
    spec:
      containers:
      - name: go-gin-log-k8s
        image: harbor.alnk.com/public/go-gin-log:0.4
        # 添加相应的环境变量
        # 下面收集了两块日志
        # go-gin-std:屏幕打印日志
        # go-gin-log:/app目录下所有*.log文件日志
        # 如日志发送到es，那index名称为go-gin-std,如日志发送到kafka，那topic则为go-gin-std
        # 如日志发送到es，那index名称为go-gin-log,如日志发送到kafka，那topic则为go-gin-log
        env: 
        - name: aliyun_logs_go-gin-std
          value: "stdout"
        - name: aliyun_logs_go-gin-log
          value: "/app/log/*.log"
        ports:
        - containerPort: 3000
          protocol: TCP
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        imagePullPolicy: IfNotPresent
        # 对pod内要收集的业务日志目录需要进行共享，可以收集多个目录下的日志文件
        volumeMounts:
          - name: gin-log
            mountPath: /app/log/
      volumes: 
        - name: gin-log
          emptyDir: {}
      
      imagePullSecrets:
      - name: harbor
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      securityContext:
        runAsUser: 0
      schedulerName: default-scheduler
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 1
      maxSurge: 1
  revisionHistoryLimit: 7
  progressDeadlineSeconds: 600

---
# service.yaml
kind: Service
apiVersion: v1
metadata:
  name: go-gin-log-k8s
  namespace: prod
  labels:
    go-app: go-gin-log-k8s
spec:
  ports:
  - protocol: TCP
    port: 80
    targetPort: 3000
    name: http
  selector:
    app: go-gin-log-k8s

---
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  namespace: prod
  name: go-gin-log-k8s
spec:
  rules:
  - host: go-gin-log-k8s.alnk.com
    http:
      paths:
      - backend:
          service:
            name: go-gin-log-k8s
            port:
              number: 80
        path: /
        pathType: Prefix
```





























































