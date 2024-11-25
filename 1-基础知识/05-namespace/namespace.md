### 命名空间（namespace）

#### 什么是命名空间

> ```text
> namespace命名空间，简称ns
> 在K8s上面，大部分资源都受ns的限制，来做资源的隔离
> 一个Namespace可以看作是一个虚拟的集群，它将物理集群划分为多个逻辑部分
> 每个部分都有自己的一组资源如Pod、Service、ConfigMap等
> 
> 
> K8s启动时会创建四个初始名字空间
> 1. default：你的service和app默认被创建于此
> 2. kube-system：kubernetes系统组件使用
> 3. kube-public：公共资源使用。但实际上现在并不常用
> 4. kube-node-lease：该空间包含用于与各个节点关联的Lease租约对象，节点租约允许kubelet发送心跳， 由此控制面能够检测到节点故障
> 
> 
> k8s命名空间主要用于隔离集群资源、隔离容器等，为集群提供了一种虚拟隔离的策略
> 1. 资源隔离：可为不同的团队、用户或项目提供虚拟的集群空间，共享同一个k8s集群的资源，使用ResourceQuota与Resource LimitRange来指定与限制各个namesapce的资源分配与使用
> 2. 权限控制：可以指定某个namespace哪些用户可以访问，哪些用户不能访问
> ```



#### 命名空间的查看和创建

> ```shell
> ### 1. 查看所有命名空间
> # kubectl get namespaces
> # kubectl get ns               
> # kubectl get ns --show-labels # 显示namespace的label
> 
> ### 2. 查看命名空间的详细信息
> # kubectl describe namespace kube-system
> 
> ### 3. 创建命名空间
> # kubectl create namespace alnk-test
> 
> ### 4.查看某个命名空间下的pod
> ## 如果不指定-n命名空间，会默认查看default命名空间里的pod，
> ## 创建pod的时候不指定命名空间，只会将pod创建在default命名空间里
> # kubectl get pod -n kube-system
> 
> ### 5. 删除命名空间
> # kubectl delete namespace alnk-test
> ```



#### 跨命名空间通信

> ```
> 命名空间彼此是透明的，但它们并不是绝对相互隔绝但
> 
> 一个namespace中service可以和另一个namespace中的service通信
> 比如团队的一个service要和另外一个团队的service通信
> 而你们的service都在各自namespace中
> 通常会把mysql,redis,rabbitmq,mongodb这些公用组件放在一个namespace里
> 或者每个公用组件都有自己的namespace，而业务组件会统一放在自己的namespace里
> 这时就涉及到了跨namespace的数据通讯问题
> 
> 当应用app要访问k8s的service，可以使用内置的DNS服务发现并把你的app指到Service的名称
> 然而可以在多个namespace中创建同名的service
> 解决这个问题，就用到DNS地址的扩展形式
> 在k8s中，Service通过一个DNS模式来暴露endpoint
> 这个模式类似 <Service Name>.<Namespace Name>.svc.cluster.local
> 一般情况下，只需要service的名称，DNS会自动解析到它的全地址
> 如果要访问其他namespace中的service，那么就需要同时使用service名称和namespace名称
> 例如想访问namespace为test中的“database”服务，可以使用下面的地址
> database.test
> ```



#### 资源限制

> ```
> 在默认情况下，k8s不会对pod进行CPU和内存限制，如果某个Pod发生内存泄露那么将是一个非常糟糕的事情
> 
> 一. pod级别的资源限制
> 在部署Pod的时候把Requests和limits加上，配置文件示例如下
> apiVersion: apps/v1
> kind: Deployment
> metadata:
>   name: ng-deploy
> spec:
>   selector:
>     matchLables:
>       app: ng-demo
>     replicas: 2
>     template:
>       metadata:
>         labels:
>           app: ng-demo
>       spec:
>         containers:
>         - name: ng-demo
>           image: nginx
>           imagePullPolicy:IfNotPresent
>           resources:
>             requests:
>               cpu: 100m
>               memory: 200Mi
>             limits:
>               cpu: 200m
>               memory: 400Mi          
> 
> 二. 名称空间级别的限制
> 如果Pod多并只需要相同的限制，这样一个一个设置就比较麻烦了
> 这时可以通过LimitRange做一个全局限制
> 如果在部署pod的时候指定了requests和limits，则指定的生效，反之则全局给pod设置默认的限制
> apiVersion: v1
> kind: LimitRange
> metadata:
>   name: alnk-test
> spec:
>   limits:
>   - type: Container       #资源类型
>     max:
>       cpu: "1"            #限定最大CPU
>       memory: "1Gi"       #限定最大内存
>     min:
>       cpu: "100m"         #限定最小CPU
>       memory: "100Mi"     #限定最小内存
>     default:
>       cpu: "900m"         #默认CPU限定
>       memory: "800Mi"     #默认内存限定
>     defaultRequest:
>       cpu: "200m"         #默认CPU请求
>       memory: "200Mi"     #默认内存请求
>     maxLimitRequestRatio:
>       cpu: 2              #限定CPU limit/request比值最大为2  
>       memory: 1.5         #限定内存limit/request比值最大为1.5
>   - type: Pod
>     max:
>       cpu: "2"            #限定Pod最大CPU
>       memory: "2Gi"       #限定Pod最大内存
>   - type: PersistentVolumeClaim
>     max:
>       storage: 2Gi        #限定PVC最大的requests.storage
>     min:
>       storage: 1Gi        #限定PVC最小的requests.storage
>       
> ## 注释
> 该文件定义了在namespace alnk-test中，容器、Pod、PVC的资源限制，在该namesapce中，只有满足如下条件，对象才能创建成功
> 1. 容器的resources.limits部分CPU必须在100m-1之间，内存必须在100Mi-1Gi之间，否则创建失败
> 2. 容器的resources.limits部分CPU与resources.requests部分CPU的比值最大为2，memory比值最大为1.5，否则创建失败
> 3. Pod内所有容器的resources.limits部分CPU总和最大为2，内存总和最大为2Gi，否则创建失败
> 4. PVC的resources.requests.storage最大为2Gi，最小为1Gi，否则创建失败
> 
> 如果容器定义了resources.requests没有定义resources.limits，则LimitRange中的default部分将作为limit注入到容器中；如果容器定义了resources.limits却没有定义resources.requests，则将requests值也设置为limits的值；如果容器两者都没有定义，则使用LimitRange中default作为limits，defaultRequest作为requests值
> ```



