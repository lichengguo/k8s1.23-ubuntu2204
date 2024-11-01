#### k8s环境

> k8s部署工具：kubeasz
>
> k8s 版本：1.27.5
>
> 容器运行时： containerd
>
> 操作系统 ：ubuntu2204
>
> ##### 服务器规划
>
> | IP         | 主机名称 | 角色        | 系统       | 软件                                                         | 配置  |
> | :--------- | -------- | ----------- | ---------- | ------------------------------------------------------------ | ----- |
> | 10.0.1.21  | ops      | 运维机      | ubuntu2204 | harbor仓库、kubeasz                                          | 1c/2g |
> | 10.0.1.100 | 虚拟IP   | /           | /          | 流量入口、负载均衡、高可用、七层反向代理                     | /     |
> | 10.0.1.101 | ha-1     | 反向代理    | ubuntu2204 | nginx、keepalived                                            | 1c/2g |
> | 10.0.1.102 | ha-2     | 反向代理    | ubuntu2204 | nginx、keepalived、etcd                                      | 1c/2g |
> | 10.0.1.200 | 虚拟IP   | /           | /          | apiserver高可用、4层反向代理                                 | /     |
> | 10.0.1.201 | master-1 | k8s主节点   | ubuntu2204 | apiserver、controller、scheduler、etcd、keepalived、nginx(l4lb) | 1c/2g |
> | 10.0.1.202 | master-2 | k8s主节点   | ubuntu2204 | apiserver、controller、scheduler、etcd、keepalived、nginx(l4lb) | 1c/2g |
> | 10.0.1.203 | node-1   | k8s工作节点 | ubuntu2204 | kubelet、kube-proxy                                          | 2c/8g |
> | 10.0.1.204 | node-2   | k8s工作节点 | ubuntu2204 | kubelet、kube-proxy                                          |       |



#### Prometheus Operator的架构示意图

![1730298915954](.\images\1730298915954.png)



#### Prometheus Operator能做什么

> 要了解Prometheus Operator能做什么，其实就是要了解Prometheus Operator为我们提供了哪些自定义的Kubernetes资源，列出了Prometheus Operator目前提供的️4类资源：
>
> - Prometheus：声明式创建和管理Prometheus Server实例；
> - ServiceMonitor：负责声明式的管理监控配置；
> - PrometheusRule：负责声明式的管理告警配置；
> - Alertmanager：声明式的创建和管理Alertmanager实例。
>
> 简言之，Prometheus Operator能够帮助用户自动化的创建以及管理Prometheus Server以及其相应的配置



#### 部署

> https://github.com/prometheus-operator/kube-prometheus/releases

```
【10.0.1.201】
# mkdir -p /data/operator/prometheus
# cd /data/operator/prometheus

# wget https://github.com/prometheus-operator/kube-prometheus/archive/refs/tags/v0.14.0.zip

# wget https://github.com/prometheus-operator/kube-prometheus/archive/refs/tags/v0.13.0.zip

# unzip v0.13.0.zip
# cd kube-prometheus-0.13.0/

# find ./ -type f |xargs egrep 'image: quay.io|image: registry.k8s.io|image: grafana|image: docker.io'|awk '{print $3}'|sort|uniq
#---------------------------------------------------------------------------------------
# 注意：这两个镜像配置比较特殊，上面命令过滤不出来
quay.io/prometheus-operator/prometheus-config-reloader:v0.67.1  
docker.io/jimmidyson/configmap-reload:v0.5.0

grafana/grafana:9.5.3
docker.io/cloudnativelabs/kube-router
quay.io/brancz/kube-rbac-proxy:v0.14.2
quay.io/prometheus/alertmanager:v0.26.0
quay.io/prometheus/blackbox-exporter:v0.24.0
quay.io/prometheus/node-exporter:v1.6.1
quay.io/prometheus-operator/prometheus-operator:v0.67.1
quay.io/prometheus/prometheus:v2.46.0
registry.k8s.io/kube-state-metrics/kube-state-metrics:v2.9.2
registry.k8s.io/prometheus-adapter/prometheus-adapter:v0.11.1
quay.io/fabxc/prometheus_demo_service
#---------------------------------------------------------------------------------------

#-------------镜像上传到阿里云--------------------------------------------------------------
registry.cn-hangzhou.aliyuncs.com/alnktest/prometheus-config-reloader:v0.67.1
registry.cn-hangzhou.aliyuncs.com/alnktest/configmap-reload:v0.5.0

registry.cn-hangzhou.aliyuncs.com/alnktest/kube-router
registry.cn-hangzhou.aliyuncs.com/alnktest/grafana:9.5.3
registry.cn-hangzhou.aliyuncs.com/alnktest/kube-rbac-proxy:v0.14.2
registry.cn-hangzhou.aliyuncs.com/alnktest/alertmanager:v0.26.0
registry.cn-hangzhou.aliyuncs.com/alnktest/blackbox-exporter:v0.24.0
registry.cn-hangzhou.aliyuncs.com/alnktest/node-exporter:v1.6.1
registry.cn-hangzhou.aliyuncs.com/alnktest/prometheus-operator:v0.67.1
registry.cn-hangzhou.aliyuncs.com/alnktest/prometheus:v2.46.0
registry.cn-hangzhou.aliyuncs.com/alnktest/kube-state-metrics:v2.9.2
registry.cn-hangzhou.aliyuncs.com/alnktest/prometheus-adapter:v0.11.1
quay.io/fabxc/prometheus_demo_service
#---------------------------------------------------------------------------------------

#-------------镜像上传到本地----------------------------------------------------------------
harbor.alnk.com/public/prometheus-config-reloader:v0.67.1
harbor.alnk.com/public/configmap-reload:v0.5.0

harbor.alnk.com/public/kube-router
harbor.alnk.com/public/grafana:9.5.3
harbor.alnk.com/public/kube-rbac-proxy:v0.14.2
harbor.alnk.com/public/alertmanager:v0.26.0
harbor.alnk.com/public/blackbox-exporter:v0.24.0
harbor.alnk.com/public/node-exporter:v1.6.1
harbor.alnk.com/public/prometheus-operator:v0.67.1
harbor.alnk.com/public/prometheus:v2.46.0
harbor.alnk.com/public/kube-state-metrics:v2.9.2
harbor.alnk.com/public/prometheus-adapter:v0.11.1
harbor.alnk.com/public/prometheus_demo_service
#---------------------------------------------------------------------------------------


## 替换
find ./ -type f |xargs  sed -ri 's+quay.io/.*/+harbor.alnk.com/public/+g'
find ./ -type f |xargs  sed -ri 's+docker.io/cloudnativelabs/+harbor.alnk.com/public/+g'
find ./ -type f |xargs  sed -ri 's+grafana/+harbor.alnk.com/public/+g'
find ./ -type f |xargs  sed -ri 's+registry.k8s.io/.*/+harbor.alnk.com/public/+g'
##
# vi blackboxExporter-deployment.yaml
这里面还有个镜像没过滤到 docker.io/jimmidyson/configmap-reload:v0.5.0
替换成 harbor.alnk.com/public/configmap-reload:v0.5.0

# 开始创建所有服务
kubectl create -f manifests/setup
kubectl create -f manifests/

# 查看创建结果：
kubectl -n monitoring get all
kubectl -n monitoring get pod -w

# 附：清空上面部署的prometheus所有服务：
kubectl delete --ignore-not-found=true -f manifests/ -f manifests/setup

```

#### 访问prometheus的UI

```
# mkdir /data/operator/prometheus/kube-prometheus-0.13.0/zidingyi
# cd /data/operator/prometheus/kube-prometheus-0.13.0/zidingyi
# vim prometheus-ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  namespace: monitoring
  name: prometheus
spec:
  rules:
  - host: prometheus-operator.alnk.com
    http:
      paths:
      - backend:
          service:
            name: prometheus-k8s
            port:
              number: 9090
        path: /
        pathType: Prefix

# kubectl apply -f prometheus-ingress.yaml
```

#### grafana ingress创建

```
# mkdir /data/operator/prometheus/kube-prometheus-0.13.0/zidingyi
# cd /data/operator/prometheus/kube-prometheus-0.13.0/zidingyi
# vim grafana-ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  namespace: monitoring
  name: grafana
spec:
  rules:
  - host: grafana-operator.alnk.com
    http:
      paths:
      - backend:
          service:
            name: grafana
            port:
              number: 3000
        path: /
        pathType: Prefix


# kubectl apply -f grafana-ingress.yaml

# 注意：删除自带的网络策略，否则访问服务都会被阻塞
https://github.com/prometheus-operator/kube-prometheus/issues/1763#issuecomment-1139553506
# kubectl -n monitoring delete networkpolicies.networking.k8s.io --all

```



#### 监控kube-controller-manager和kube-scheduler

> 访问prometheus后台，点击上方`菜单栏`-`Status` — `Targets` ，发现kube-controller-manager和kube-scheduler未发现

```
# ss -tlnp|egrep 'controller|schedule'
LISTEN 0      32768              *:10259            *:*    users:(("kube-scheduler",pid=849,fd=3))
LISTEN 0      32768              *:10257            *:*    users:(("kube-controller",pid=1154,fd=3))

# 因为K8s的这两上核心组件我们是以二进制形式部署的，为了能让K8s上的prometheus能发现，需要来创建相应的service和endpoints来将其关联起来

# mkdir /data/operator/prometheus/kube-prometheus-0.13.0/zidingyi
# cd /data/operator/prometheus/kube-prometheus-0.13.0/zidingyi
# vi repair-prometheus.yaml
apiVersion: v1
kind: Service
metadata:
  namespace: kube-system
  name: kube-controller-manager
  labels:
    app.kubernetes.io/name: kube-controller-manager
spec:
  type: ClusterIP
  clusterIP: None
  ports:
  - name: https-metrics
    port: 10257
    targetPort: 10257
    protocol: TCP

---
apiVersion: v1
kind: Endpoints
metadata:
  labels:
    app.kubernetes.io/name: kube-controller-manager
  name: kube-controller-manager
  namespace: kube-system
subsets:
- addresses:
  - ip: 10.0.1.201
  - ip: 10.0.1.202
  ports:
  - name: https-metrics
    port: 10257
    protocol: TCP

---
apiVersion: v1
kind: Service
metadata:
  namespace: kube-system
  name: kube-scheduler
  labels:
    app.kubernetes.io/name: kube-scheduler
spec:
  type: ClusterIP
  clusterIP: None
  ports:
  - name: https-metrics
    port: 10259
    targetPort: 10259
    protocol: TCP

---
apiVersion: v1
kind: Endpoints
metadata:
  labels:
    app.kubernetes.io/name: kube-scheduler
  name: kube-scheduler
  namespace: kube-system
subsets:
- addresses:
  - ip: 10.0.1.201
  - ip: 10.0.1.202
  ports:
  - name: https-metrics
    port: 10259
    protocol: TCP
    
# kubectl apply -f repair-prometheus.yaml
```



#### 监控etcd

> 作为K8s所有资源存储的关键服务ETCD，也有必要把它给监控起来
>
> 完整的演示一次利用Prometheus来监控非K8s集群服务的步骤
>
> 在前面部署K8s集群的时候，是用二进制的方式部署的ETCD集群，并且利用自签证书来配置访问ETCD现在关键的服务基本都会留有指标metrics接口支持prometheus的监控
>
> 利用下面命令，可以看到ETCD都暴露出了哪些监控指标出来

```
【10.0.1.21】
# curl --cacert /etc/kubeasz/clusters/k8s-01/ssl/ca.pem --cert /etc/kubeasz/clusters/k8s-01/ssl/etcd.pem  --key /etc/kubeasz/clusters/k8s-01/ssl/etcd-key.pem https://10.0.1.201:2379/metrics
```

```
【10.0.1.201】
# 拷贝10.0.1.21上的证书
# pwd
/data/operator/prometheus/kube-prometheus-0.13.0/zidingyi
# scp root@10.0.1.21:/etc/kubeasz/clusters/k8s-01/ssl/ca.pem .
# scp root@10.0.1.21:/etc/kubeasz/clusters/k8s-01/ssl/etcd.pem .
# scp root@10.0.1.21:/etc/kubeasz/clusters/k8s-01/ssl/etcd-key.pem .


# 首先把ETCD的证书创建为secret
# kubectl -n monitoring create secret generic etcd-certs --from-file=./etcd.pem   --from-file=./etcd-key.pem   --from-file=./ca.pem

# 接着在prometheus里面引用这个secrets
kubectl -n monitoring edit prometheus k8s 

spec:
...
  secrets:
  - etcd-certs

# 保存退出后，prometheus会自动重启服务pod以加载这个secret配置，过一会，我们进pod来查看下是不是已经加载到ETCD的证书了
# kubectl -n monitoring exec -it prometheus-k8s-0 -c prometheus  -- sh 
/prometheus $ ls /etc/prometheus/secrets/etcd-certs/
ca.pem        etcd-key.pem  etcd.pem


# 创建service、endpoints以及ServiceMonitor的yaml配置
# vim prometheus-etcd.yaml 
apiVersion: v1
kind: Service
metadata:
  name: etcd-k8s
  namespace: monitoring
  labels:
    k8s-app: etcd
spec:
  type: ClusterIP
  clusterIP: None
  ports:
  - name: api
    port: 2379
    protocol: TCP

---
apiVersion: v1
kind: Endpoints
metadata:
  name: etcd-k8s
  namespace: monitoring
  labels:
    k8s-app: etcd
subsets:
- addresses:
  - ip: 10.0.1.201
  - ip: 10.0.1.202
  - ip: 10.0.1.102
  ports:
  - name: api
    port: 2379
    protocol: TCP

---
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: etcd-k8s
  namespace: monitoring
  labels:
    k8s-app: etcd-k8s
spec:
  jobLabel: k8s-app
  endpoints:
  - port: api
    interval: 30s
    scheme: https
    tlsConfig:
      caFile: /etc/prometheus/secrets/etcd-certs/ca.pem
      certFile: /etc/prometheus/secrets/etcd-certs/etcd.pem
      keyFile: /etc/prometheus/secrets/etcd-certs/etcd-key.pem
      #use insecureSkipVerify only if you cannot use a Subject Alternative Name
      insecureSkipVerify: true 
  selector:
    matchLabels:
      k8s-app: etcd
  namespaceSelector:
    matchNames:
    - monitoring

# kubectl apply -f prometheus-etcd.yaml 
```



#### grafana来展示被监控的ETCD指标

```
1. 在grafana官网模板中心搜索etcd，下载这个json格式的模板文件
https://grafana.com/grafana/dashboards/3070-etcd/
# download json

2.然后打开自己先部署的grafana首页，
点击左上边菜单栏HOME --- Data source --- Add data source --- 选择 Prometheus

查看prometheus的详细地址 并编辑进去保存：
# kubectl -n monitoring get secrets grafana-datasources -o yaml
# 然后把Secret解码一下
# # echo 'ewogICAgImFwaVZlcnNpb24iOiAxLAogICAgImRhdGFzb3VyY2VzIjogWwogICAgICAgIHsKICAgICAgICAgICAgImFjY2VzcyI6ICJwcm94eSIsCiAgICAgICAgICAgICJlZGl0YWJsZSI6IGZhbHNlLAogICAgICAgICAgICAibmFtZSI6ICJwcm9tZXRoZXVzIiwKICAgICAgICAgICAgIm9yZ0lkIjogMSwKICAgICAgICAgICAgInR5cGUiOiAicHJvbWV0aGV1cyIsCiAgICAgICAgICAgICJ1cmwiOiAiaHR0cDovL3Byb21ldGhldXMtazhzLm1vbml0b3Jpbmcuc3ZjOjkwOTAiLAogICAgICAgICAgICAidmVyc2lvbiI6IDEKICAgICAgICB9CiAgICBdCn0='|base64 -d
{
    "apiVersion": 1,
    "datasources": [
        {
            "access": "proxy",
            "editable": false,
            "name": "prometheus",
            "orgId": 1,
            "type": "prometheus",
            "url": "http://prometheus-k8s.monitoring.svc:9090",
            "version": 1
        }
    ]
}



再点击右上角 +^ Import dashboard --- 
点击Upload .json File 按钮，上传上面下载好的json文件 3070_rev3.json，
点击Import，即可显示etcd集群的图形监控信息

```



























