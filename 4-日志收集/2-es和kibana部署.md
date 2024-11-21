### operator部署elasticsearch和kibana

> ```
> 使用operator部署生产级别的Elasticserach集群和Kibana
> 可参考 https://www.cuiliangblog.cn/detail/article/71
> 
> 查看版本支持
> https://www.elastic.co/guide/en/cloud-on-k8s/2.7/k8s-quickstart.html
> ```
>
> ![1731228390645](images/1731228390645.png)  



#### 部署eck-operator

> ```
> 【10.0.1.21】
> ## 默认情况下eck-operator运行在elastic-system名称空间
> 
> ## eck镜像准备（github开源项目把镜像同步到阿里云仓库）
> # docker pull registry.cn-hangzhou.aliyuncs.com/alnktest/eck-operator:2.7.0
> # docker tag registry.cn-hangzhou.aliyuncs.com/alnktest/eck-operator:2.7.0 harbor.alnk.com/public/eck-operator:2.7.0 
> # docker push harbor.alnk.com/public/eck-operator:2.7.0 
> 
> ## 下载yaml文件
> # mkdir /data/k8s-yaml/eck -p
> # cd /data/k8s-yaml/eck
> # wget https://download.elastic.co/downloads/eck/2.7.0/crds.yaml
> # wget https://download.elastic.co/downloads/eck/2.7.0/operator.yaml
> 
> ## 替换镜像
> # sed -ri 's+docker.elastic.co/.*/+harbor.alnk.com/public/+g' operator.yaml
> 
> ## 部署
> # kubectl create -f crds.yaml
> # kubectl apply -f operator.yaml
> 
> # 查看资源
> # kubectl -n elastic-system get all
> pod/elastic-operator-0   1/1     Running   0          68s
> service/elastic-webhook-server   ClusterIP   10.68.159.142   <none>        443/TCP   68s
> statefulset.apps/elastic-operator   1/1     68s
> 
> ## 查看日志
> # kubectl -n elastic-system logs -f statefulset.apps/elastic-operator
> # kubectl -n elastic-system logs -f elastic-operator-0
> ```



#### 部署elasticsearch

> ```
> 【10.0.1.21】
> ## elasticsearch镜像准备（github开源项目把镜像同步到阿里云仓库） ##
> # docker pull registry.cn-hangzhou.aliyuncs.com/alnktest/elasticsearch:6.8.6
> # docker tag registry.cn-hangzhou.aliyuncs.com/alnktest/elasticsearch:6.8.6 harbor.alnk.com/public/elasticsearch:6.8.6
> # docker push harbor.alnk.com/public/elasticsearch:6.8.6
> 
> ## yaml文件 ##
> # mkdir /data/k8s-yaml/elk -p
> # cd /data/k8s-yaml/elk
> # vi elasticsearch.yaml
> 
> 【10.0.1.21】
> ## 部署应用 ##
> # kubectl create ns elk
> # kubectl apply -f elasticsearch.yaml
> 
> # kubectl -n elk get pod
> 
> # 查看es集群是否健康 ##
> # kubectl -n elk  get elasticsearches.elasticsearch.k8s.elastic.co quickstart
> # kubectl -n elk  get es quickstart
> quickstart   green    1       8.11.3    Ready   62s
> 
> ## 查看es密码 ##
> # kubectl -n elk get secret quickstart-es-elastic-user -o go-template='{{.data.elastic | base64decode}}'
> 25y69WWe1FC8F6Gh6iVHVt53
> 
> ## 查看es的svc地址 ##
> # kubectl get svc -n elk
> quickstart-es-http            ClusterIP   10.68.120.223   <none>        9200/TCP   56s
> 
> 【10.0.1.201】
> ## 访问es集群 ##
> # curl -k http://elastic:25y69WWe1FC8F6Gh6iVHVt53@10.68.120.223:9200
> ## 查看索引列表
> # curl -k http://elastic:25y69WWe1FC8F6Gh6iVHVt53@10.68.120.223:9200/_cat/indices?v
> ```
>
> `elasticsearch.yaml`
>
> ```yaml
> apiVersion: elasticsearch.k8s.elastic.co/v1
> kind: Elasticsearch
> metadata:
>   name: quickstart # 名称
>   namespace: elk # 名称空间
> spec:
>   http: # 关闭https，使用http对外提供服务，否则需要配置证书
>     tls:
>       selfSignedCertificate:
>         disabled: true
>   version: 6.8.6 # es版本
>   image: harbor.alnk.com/public/elasticsearch:6.8.6 #自定义镜像地址，不指定则从官方镜像仓库拉取
>   nodeSets:  
>   - name: default # 节点名称
>     count: 1 # 集群数量，建议奇数个1,3,5
>     config:
>       node.store.allow_mmap: false
>     podTemplate:
>       spec:
>         containers:
>         - name: elasticsearch
>           env:
>           - name: ES_JAVA_OPTS # jvm参数
>             value: -Xms1g -Xmx1g # 指定节点JVM大小 
>           resources: # 根据实际情况调整资源限制
>             requests:     # 资源请求值，通常与JVM保持一致
>               memory: 1Gi
>               cpu: 0.5
>             limits:        # 资源限制值，通常为JVM的2倍
>               memory: 2Gi
>               cpu: 1
>         initContainers:
>         - name: sysctl
>           securityContext:
>             privileged: true
>           command: ['sh', '-c', 'sysctl -w vm.max_map_count=262144']
>     volumeClaimTemplates:
>     - metadata:
>         name: elasticsearch-data
>       spec:
>         accessModes:
>         - ReadWriteOnce
>         resources:
>           requests:
>             storage: 10Gi # 存储空间，根据实际情况调整
>         storageClassName: nfs-alnk # 建议固态硬盘
>   #http:
>     #tls:
>       #selfSignedCertificate:
>         #disabled: true 
>       
> ---
> ##测试为了方便加上ingress，生产环境中可以不加
> apiVersion: networking.k8s.io/v1
> kind: Ingress
> metadata:
>   name: elastic
>   namespace: elk
>   annotations:
>     #nginx.ingress.kubernetes.io/force-ssl-redirect: "true" #不使用https访问
>     kubernetes.io/ingress.class: nginx
> spec:
>   #tls:
>   #- hosts:
>     #- kibana.alnk.com
>     #secretName: alnk-com-tls
>   rules:
>   - host: elastic.alnk.com
>     http:
>       paths:
>       - path: /
>         pathType: Prefix
>         backend:
>           service:
>             name: quickstart-es-http
>             port:
>               number: 9200
> ```
>
> ![1731394188425](images/1731394188425.png)  
>
>  `http://elastic:25y69WWe1FC8F6Gh6iVHVt53@elastic.alnk.com/`
>
> ![1731420011582](images/1731420011582.png)  



#### 部署kibana

> ```
> 【10.0.1.21】
> ##kibana镜像准备 github开源项目把镜像同步到阿里云仓库##
> # docker pull registry.cn-hangzhou.aliyuncs.com/alnktest/kibana:6.8.6
> # docker tag registry.cn-hangzhou.aliyuncs.com/alnktest/kibana:6.8.6 harbor.alnk.com/public/kibana:6.8.6
> # docker push harbor.alnk.com/public/kibana:6.8.6
> 
> ## yaml文件 ##
> # mkdir /data/k8s-yaml/elk -p
> # cd /data/k8s-yaml/elk
> # vi kibana.yaml
> 
> ## 应用 ##
> # kubectl apply -f kibana.yaml
> 
> ## 查看 ##
> # kubectl -n elk get pod
> # kubectl -n elk get pod -w
> # kubectl -n elk get svc
> # kubectl -n elk get ingress
> ```
>
> `kibana.yaml`
>
> ```yaml
> apiVersion: kibana.k8s.elastic.co/v1
> kind: Kibana
> metadata:
>   name: quickstart # 名称
>   namespace: elk # 名称空间
> spec:
>   #增加该配置，不然出bug连接不上es
>   #https://github.com/elastic/cloud-on-k8s/issues/5952
>   #config:
>     #monitoring.ui.ccs.enabled: false
>   version: 6.8.6
>   image: harbor.alnk.com/public/kibana:6.8.6
>   count: 1 # kibana集群个数
>   elasticsearchRef: # 关联es集群名称,与Elasticsearch资源名称匹配
>     name: quickstart
>   http: # 关闭https，使用http对外提供服务，否则需要配置证书
>     tls:
>       selfSignedCertificate:
>         disabled: true
>   podTemplate:
>     spec:
>       containers:
>       - name: kibana
>         env:
>           - name: NODE_OPTIONS
>             value: "--max-old-space-size=2048"
>           - name: SERVER_PUBLICBASEURL
>             value: "https://kibana.alnk.com"
>           - name: I18N_LOCALE # 中文配置
>             value: "zh-CN"
>         resources:
>           requests:
>             memory: 1Gi
>             cpu: 0.5
>           limits:
>             memory: 2Gi
>             cpu: 2  
>       
> ---
> apiVersion: networking.k8s.io/v1
> kind: Ingress
> metadata:
>   name: kibana
>   namespace: elk
>   annotations:
>     #nginx.ingress.kubernetes.io/force-ssl-redirect: "true" #不使用https访问
>     kubernetes.io/ingress.class: nginx
> spec:
>   #tls:
>   #- hosts:
>     #- kibana.alnk.com
>     #secretName: alnk-com-tls
>   rules:
>   - host: kibana.alnk.com
>     http:
>       paths:
>       - path: /
>         pathType: Prefix
>         backend:
>           service:
>             name: quickstart-kb-http
>             port:
>               number: 5601
> ```
>
> ```
> http://kibana.alnk.com/
> 账号密码 elastic:25y69WWe1FC8F6Gh6iVHVt53
> ```
>
> ![1731442773697](images/1731442773697.png)  

