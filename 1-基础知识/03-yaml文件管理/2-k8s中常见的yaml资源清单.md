### yaml基础

> ```
> YAML是专门用来写配置文件的语言，非常简洁和强大，使用比json更方便，它实质上是一种通用的数据串行化格式
> 
> 
> YAML语法规则
> 大小写敏感
> 使用缩进表示层级关系
> 缩进时不允许使用Tal键，只允许使用空格
> 缩进的空格数目不重要，只要相同层级的元素左侧对齐即可
> ”#” 表示注释，从这个字符一直到行尾，都会被解析器忽略　　
> 
> 
> 在k8s中，只需要知道两种结构类型即可：
> Lists
> Maps
> 
> 
> ## YAML Maps
> 1. Map顾名思义指的是字典，即一个Key:Value 的键值对信息
> apiVersion: v1
> kind: Pod
> 注：---为可选的分隔符，当需要在一个文件中定义多个结构的时候需要使用
> 上述内容表示有两个键apiVersion和kind，分别对应的值为v1和Pod
> {
>   "apiVersion": "v1",
>   "kind": "Pod"
> }
> 
> 
> 2. Maps的value既能够对应字符串也能够对应一个Maps
> apiVersion: v1
> kind: Pod
> metadata:
>   name: kube100-site
>   labels:
>     app: web
> 上述的YAML文件中，metadata这个KEY对应的值为一个Maps，而嵌套的labels这个KEY的值又是一个Map
> 使用两个空格作为缩进，但空格的数据量并不重要，只是至少要求一个空格并且所有缩进保持一致的空格数
> name和labels是相同缩进级别，因此YAML处理器知道他们属于同一map
> YAML处理器知道app是lables的值因为app的缩进更大
> 在YAML文件中绝对不要使用tab键
> {
>   "apiVersion": "v1",
>   "kind": "Pod",
>   "metadata": {
>     "name": "kube100-site",
>     "labels": {
>       "app": "web"
>     }
>   }
> }
> 
> 
> 
> ## YAML Lists
> List即列表，说白了就是数组，例如：
> args:
> - beijing: v1
>   shanghai: v2
> {
>   "args": [{"beijing": "v1", "shanghai": "v2"}]
> }  
> 
> args:
>   - beijing: v1
>     shanghai: v2
> {
>   "args": [{"beijing": "v1", "shanghai": "v2"}]
> }
> 
> 
> args:
> - beijing: v1
> - shanghai: v2
> {
>   "args": [{"beijing": "v1"}, {"shanghai": "v2"}]
> }
> 
> args:
>   - beijing: v1
>   - shanghai: v2
> {
>   "args": [{"beijing": "v1"}, {"shanghai": "v2"}]
> }
> ```



### pod

> ```yaml
> apiVersion: v1
> kind: Pod
> metadata:
>   #可以添加 metadata.annotations 字段用于存放其他注释信息
>   annotations:
>     description: This is a sample Pod configuration file.  # 描述该配置文件的注释
>   name: my-pod  # Pod 的名称
>   labels:
>     app: my-app  # Pod 的标签，可用于选择器和服务发现
> spec:
>   containers:
>   - name: my-container  # 容器的名称
>     image: nginx  # 使用的容器镜像
>     ports:
>     - containerPort: 80  # 容器内部暴露的端口号
>       protocol: TCP  # 端口的协议类型
>     resources:  # 容器所需的资源限制和需求
>       limits:
>         cpu: "0.5"  # CPU 使用上限
>         memory: "256Mi"  # 内存使用上限
>       requests:
>         cpu: "0.2"  # CPU 最小需求
>         memory: "128Mi"  # 内存最小需求
>     env:  # 容器的环境变量
>     - name: ENV_VAR1
>       value: value1
>     - name: ENV_VAR2
>       valueFrom:
>         secretKeyRef:
>           name: my-secret  # 引用的 Secret 对象的名称
>           key: secret-key  # 从 Secret 对象中获取的键名
>     volumeMounts:  # 容器的挂载路径
>     - name: data-volume  # 挂载的卷的名称
>       mountPath: /data  # 挂载的路径
>   volumes:  # Pod 的卷
>   - name: data-volume  # 卷的名称
>     emptyDir: {}  # 空的卷，生命周期与 Pod 相关联
> ```



### Deployment

> ```yaml
> apiVersion: apps/v1
> kind: Deployment
> metadata:
>   name: my-deployment  # Deployment 的名称
> spec:
>   replicas: 3  # 副本数
>   selector:
>     matchLabels:
>       app: my-app  # 用于选择 Pod 的标签
>   template:
>     metadata:
>       labels:
>         app: my-app  # Pod 的标签，用于与选择器匹配
>     spec:
>       containers:
>       - name: my-container  # 容器的名称
>         image: nginx  # 使用的容器镜像
>         ports:
>         - containerPort: 80  # 容器内部暴露的端口号
>           protocol: TCP  # 端口的协议类型
>         resources:  # 容器所需的资源限制和需求
>           limits:
>             cpu: "0.5"  # CPU 使用上限
>             memory: "256Mi"  # 内存使用上限
>           requests:
>             cpu: "0.2"  # CPU 最小需求
>             memory: "128Mi"  # 内存最小需求
>         env:  # 容器的环境变量
>         - name: ENV_VAR1
>           value: value1
>         - name: ENV_VAR2
>           valueFrom:
>             secretKeyRef:
>               name: my-secret  # 引用的 Secret 对象的名称
>               key: secret-key  # 从 Secret 对象中获取的键名
>         volumeMounts:  # 容器的挂载路径
>         - name: data-volume  # 挂载的卷的名称
>           mountPath: /data  # 挂载的路径
>       volumes:  # Pod 的卷
>       - name: data-volume  # 卷的名称
>         emptyDir: {}  # 空的卷，生命周期与 Pod 相关联
> ```



### Service

> ```text
> 在Kubernetes中，服务是一种抽象，用于公开Pod或一组Pod提供的应用程序或功能
> 通过服务，其他组件可以无需了解具体Pod的IP地址和端口而直接与应用程序进行通信
> 服务可以通过多种方式公开，包括ClusterIP、NodePort、LoadBalancer等类型
> ```



#### ClusterIP

> ```yaml
> kind: Service
> apiVersion: v1
> metadata:
>   name: my-service
> spec:
>   type: ClusterIP
>   ports:
>     - port: 8080  # 集群内部的服务监听端口
>       targetPort: 80  # Pod 内部的端口
>   selector:
>     app: my-app
> ```
>
> ```text
> 当服务的类型为ClusterIP时，它会为集群内部的其他组件分配一个虚拟IP地址
> 并使用服务监听端口来接收来自其他组件的请求
> 这样，其他Pod或服务就可以通过服务的虚拟IP地址和监听端口与服务进行通信
> 
> 服务名称为my-service，类型为ClusterIP，它监听端口8080，并将请求转发到Pod内部的80端口
> 其他组件可以通过访问my-service:8080来与服务进行通信
> 需要注意的是，集群内部的服务监听端口通常是在集群内部使用的端口，不直接暴露给外部请求
> 如需要将服务公开到集群外部，可以考虑使用其他类型的服务，如NodePort或LoadBalance
> ```



#### NodePort

> ```yaml
> kind: Service
> apiVersion: v1
> metadata:
>   name: my-service
> spec:
>   type: NodePort
>   ports:
>     - port: 80
>       targetPort: 8080
>       nodePort: 30000
>   selector:
>     app: my-app
> ```
>
> ```text
> 在NodePort的YAML文件中，ports下的port字段表示服务在集群内部使用的端口号
> 这是服务对内提供服务的端口，其他Pod可以通过该端口与服务进行通信
> 
> nodePort字段表示服务在节点上公开的端口号
> 当服务类型为NodePort时，k8s随机分配一个未使用的端口号，并将该端口号映射到每个节点上
> 这样，可以通过节点的IP地址和nodePort端口号来访问服务
> 
> targetPort字段是服务所指向的Pod的端口
> 当请求到达nodePort监听的端口时，它将被转发到目标端口（targetPort）上的Pod进行处理
> 
> 该示例中，服务名称为my-service，通过NodePort类型将集群中的应用程序暴露到节点的30000端口上
> 集群内部的服务监听端口为80，将请求转发到Pod上的8080端口
> 
> 请注意，NodePort的端口范围是30000-32767，确保选择一个未被占用的端口号
> ```



#### LoadBalancer

> ```yaml
> apiVersion: v1
> kind: Service
> metadata:
>   name: my-service  # Service 的名称
> spec:
>   selector:
>     app: my-app  # 选择要路由到的 Pod 的标签
>   ports:
>   - name: http  # 端口的名称
>     protocol: TCP  # 端口的协议类型
>     port: 80  # Service 暴露的端口号
>     targetPort: 8080  # 路由到的 Pod 的端口号
>   type: LoadBalancer  # Service 的类型，可以是 ClusterIP、NodePort 或者 LoadBalancer
> ```



### Ingress

> ```yaml
> apiVersion: networking.k8s.io/v1
> kind: Ingress
> metadata:
>   name: my-ingress  # Ingress 的名称
>   annotations:
>     nginx.ingress.kubernetes.io/rewrite-target: /$1  # 添加 NGINX Ingress 控制器的注解，用于重写 URL
> spec:
>   rules:
>   - host: example.com  # 定义要匹配的域名
>     http:
>       paths:
>       - path: /appA  # URL 路径
>         pathType: Prefix  # 路径匹配类型，可以是 Prefix 或 Exact
>         backend:
>           service:
>             name: appA-service  # 要路由到的 Service 的名称
>             port:
>               number: 80  # 路由到的 Service 的端口号
>       - path: /appB(/|$)(.*)  # 使用正则表达式匹配 URL 路径
>         pathType: Prefix  # 路径匹配类型，可以是 Prefix 或 Exact
>         backend:
>           service:
>             name: appB-service  # 要路由到的 Service 的名称
>             port:
>               number: 80  # 路由到的 Service 的端口号
> ```



### ConfigMap

> ```yaml
> apiVersion: v1
> kind: ConfigMap
> metadata:
>   name: my-configmap  # ConfigMap 的名称
> data:
>   server.conf: |
>     # Server 配置文件
>     port=8080
>     host=localhost
>   client.conf: |
>     # Client 配置文件
>     timeout=5000
>     retries=3
> ```
>
> ```text
> ConfigMap的配置文件可以包含多个键值对，键名和对应的值可以是任意字符串类型
> 例如文件内容、环境变量、命令行参数等
> 
> 在使用ConfigMap时，可以将其挂载到Pod中的容器内，从而使容器可以轻松地访问配置信息
> 为了更好地管理和维护ConfigMap，建议使用有意义的名称和注释对其进行命名和描述
> ```



### Secret

> ```yaml
> apiVersion: v1
> kind: Secret
> metadata:
>   name: my-secret  # Secret 的名称
> type: Opaque  # Secret 类型（Opaque 表示任意类型）
> data:
>   username: dXNlcm5hbWU=  # 加密后的用户名
>   password: cGFzc3dvcmQ=  # 加密后的密码
> ```
>
> ```text
> Secret的配置文件包含敏感信息，如用户名、密码等，需要进行加密处理
> 在YAML文件中，可以将敏感信息以base64编码的方式保存在data字段中，以保证安全性
> ```



### Volume

> ```yaml
> apiVersion: v1
> kind: Pod
> metadata:
>   name: my-pod  # Pod 的名称
> spec:
>   containers:
>     - name: my-container  # 容器的名称
>       image: nginx  # 容器的镜像
>       volumeMounts:
>         - name: data-volume  # 挂载卷的名称
>           mountPath: /data  # 挂载到容器中的路径
>   volumes:
>     - name: data-volume  # 卷的名称
>       emptyDir: {}  # 空目录卷
> ```
>
> ```text
> Volume的配置文件可以包含多个卷定义，
> 每个卷可以是不同类型的卷（如 emptyDir、hostPath、persistentVolumeClaim 等）
> 在上述示例中，使用的是emptyDir类型的卷，它会在Pod运行时创建一个空目录，并将其挂载到容器内的指定路径
> 
> 访问模式
> 在k8s中，访问模式（Access Modes）是用来定义持久化存储卷（Persistent Volume）的访问方式的
> 下面是Kubernetes支持的三种访问模式：
> 1. ReadWriteOnce（RWO）：该访问模式表示该存储卷可以被单个节点以读写方式挂载。这意味着同一时间内只能有一个Pod能够挂载并对存储卷进行读写操作。当存储卷被某个节点上的Pod挂载时，它将成为该节点的专属卷，在其他节点上不可见
> 2. ReadOnlyMany（ROX）：该访问模式表示该存储卷可以以只读方式被多个节点挂载。多个Pod可以共享对存储卷的只读访问权限，但不能进行写入操作
> 3. ReadWriteMany（RWX）：该访问模式表示该存储卷可以以读写方式被多个节点挂载。多个Pod可以同时挂载并对存储卷进行读写操作，即具有读写共享的功能
> 
> ```





### StatefulSet

> ```yaml
> apiVersion: apps/v1
> kind: StatefulSet
> metadata:
>   name: my-statefulset  # StatefulSet 的名称
> spec:
>   selector: 
>     matchLabels:
>       app: my-app  # 匹配标签，用于选择要管理的 Pod
>   serviceName: my-service  # Headless Service 的名称
>   replicas: 3  # 副本数
>   template:
>     metadata:
>       labels:
>         app: my-app  # Pod 的标签
>     spec:
>       containers:
>         - name: my-container  # 容器的名称
>           image: nginx  # 容器的镜像
>           ports:
>             - containerPort: 80  # 容器监听的端口号
>           volumeMounts:
>             - name: data-volume  # 挂载卷的名称
>               mountPath: /data  # 挂载到容器中的路径
>   volumeClaimTemplates:
>     - metadata:
>         name: data-volume  # 持久化存储卷模板的名称
>       spec:
>         accessModes:
>           - ReadWriteOnce  # 访问模式
>         resources:
>           requests:
>             storage: 1Gi  # 存储容量
> ```
>
> ```text
> StatefulSet是用于管理有状态应用程序的控制器，它保证Pod的唯一性和稳定性，并按照序号进行命名
> 
> 在上述示例中，创建了一个包含3个副本的StatefulSet，每个副本都会被命名为my-statefulset-{0…2}
> ```



### DaemonSet

> ```yaml
> apiVersion: apps/v1
> kind: DaemonSet
> metadata:
>   name: my-daemonset  # DaemonSet 的名称
> spec:
>   selector:
>     matchLabels:
>       app: my-app  # 匹配标签，用于选择要管理的节点上的 Pod
>   template:
>     metadata:
>       labels:
>         app: my-app  # Pod 的标签
>     spec:
>       containers:
>         - name: my-container  # 容器的名称
>           image: nginx  # 容器的镜像
>           ports:
>             - containerPort: 80  # 容器监听的端口号
>           volumeMounts:
>             - name: data-volume  # 挂载卷的名称
>               mountPath: /data  # 挂载到容器中的路径
>       nodeSelector:
>         disktype: ssd  # 节点的标签选择器，用于选择带有指定标签的节点
>   updateStrategy:
>     type: RollingUpdate  # 更新策略为滚动更新
>     rollingUpdate:
>       maxUnavailable: 1  # 在更新期间最多允许一个 Pod 不可用
>   volumeClaimTemplates:
>     - metadata:
>         name: data-volume  # 持久化存储卷模板的名称
>       spec:
>         accessModes:
>           - ReadWriteOnce  # 访问模式
>         resources:
>           requests:
>             storage: 1Gi  # 存储容量
> ```
>
> ```
> DaemonSet是用于在每个节点上运行一个Pod的控制器，它保证每个节点上都有一个唯一的Pod进行运行，并自动适应节点变化。
> 在上述示例中，创建了一个DaemonSet，每个节点上都会运行一个Pod，该Pod将被命名为 my-daemonset-{node-name}
> 
> DaemonSet还包含了volumeClaimTemplates字段，用于定义持久化存储卷模板
> 在上述示例中，使用了一个名为data-volume的持久化存储卷模板，访问模式为ReadWriteOnce，存储容量为 1GB。Pod将根据该模板创建一个与之对应的持久化存储卷
> 
> 通过使用DaemonSet和持久化存储卷，可以在集群的每个节点上运行一个Pod，并确保应用程序的数据持久化和可靠性存储。
> 
> 由于DaemonSet会自动适应节点变化，因此在增加或删除节点时，应用程序的数据不会丢失或受影响
> ```



### Job

> ```yaml
> apiVersion: batch/v1
> kind: Job
> metadata:
>   name: my-job  # Job 的名称
> spec:
>   completions: 1  # 完成的任务数
>   parallelism: 1  # 并行运行的 Pod 数量
>   template:
>     metadata:
>       name: my-pod  # Pod 的名称
>     spec:
>       restartPolicy: Never  # 不重启容器
>       containers:
>         - name: my-container  # 容器的名称
>           image: nginx  # 容器的镜像
>           command: ["echo", "Hello, world!"]  # 容器启动命令
>       volumes:
>         - name: data-volume  # 卷的名称
>           emptyDir: {}  # 空的临时卷
>   backoffLimit: 3  # 重试次数上限
> ```
>
> ```text
> Job是用于在Kubernetes中运行一次性任务的控制器
> 
> Job的配置还包括一些重要的字段，如completions（完成的任务数）和parallelism（并行运行的Pod 数量）completions 字段指定了job 完成的任务数，一旦达到该数量，Job 就会被标记为成功。parallelism 字段指定了同时运行的 Pod 数量，可以控制并行执行任务的速度。
> 
> 另外，还有一些其他常用的字段，例如 restartPolicy（容器的重启策略），volumes（卷的定义）和 backoffLimit（重试次数上限）。这些参数可以根据业务需求进行调整和配置
> 
> 通过使用 Job 控制器，可以在 Kubernetes 中运行一次性任务，并确保任务的完成和可靠性运行。由于 Job 可以指定任务数量和并行度，可以很好地适应不同规模和要求的任务场景。
> 
> 
> 
> 容器的重启策略
> Always（默认）：当容器失败、终止或退出时，总是自动重启容器。这是默认的重启策略
> OnFailure：只有当容器以非零状态码退出时，才会自动重启容器。如果容器正常终止并退出（即零状态码），则不会自动重启容器
> Never：当容器失败、终止或退出时，不会自动重启容器。如果容器异常终止，它将保持在该状态，直到手动重启
> 
> 
> Pod的restartPolicy字段适用于整个Pod中的所有容器
> 如果多个容器在同一个Pod中运行，并且希望它们具有不同的重启策略，可以考虑将它们放置在不同的Pod中
> ```



### ConJob

> ```yaml
> apiVersion: batch/v1beta1  # 使用的 API 版本
> kind: CronJob  # CronJob 类型
> metadata:
>   name: my-cronjob  # CronJob 的名称
> spec:
>   schedule: "*/1 * * * *"  # Cron 表达式，用于定义作业执行的时间表
>   jobTemplate:  # 作业模板，指定 CronJob 创建的作业配置
>     spec:
>       template:
>         metadata:
>           name: my-job  # 作业的名称
>         spec:
>           restartPolicy: OnFailure  # 容器的重启策略
>           containers:
>             - name: my-container  # 容器的名称
>               image: nginx  # 容器的镜像
>               command: ["echo", "Hello, world!"]  # 容器启动命令
>   successfulJobsHistoryLimit: 5  # 历史成功作业保存的数量上限
>   failedJobsHistoryLimit: 5  # 历史失败作业保存的数量上限
> 
> ```
>
> ```text
> CronJob是Kubernetes中的一种控制器，用于定期运行作业
> 在上述示例中，创建了一个CronJob，其中定义了一个Cron表达式 */1 * * * *，表示每分钟执行一次作业
> 
> successfulJobsHistoryLimit 和 failedJobsHistoryLimit字段分别指定历史成功和失败作业保存的数量上限，可以根据需求进行调整。
> ```