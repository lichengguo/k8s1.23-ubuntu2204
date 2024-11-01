#### k8s环境

> k8s部署工具：kubeasz
>
> k8s 版本：1.27.5
>
> 容器运行时： containerd
>
> 操作系统 ：ubuntu2204
>
> 
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



#### k8s部署gitlab

> 部署gitlab私有代码仓库需要的数据库postgresql和redis



##### 部署持久化存储nfs

```
# nfs共享
# 在10.0.1.201上安装（在生产中，大家要提供作好NFS-SERVER环境的规划）

# 安装nfs软件
# sudo apt install nfs-kernel-server -y

# 创建NFS挂载目录
# mkdir -p /nfs_dir
# chown -R nobody:nogroup /nfs_dir

# 修改NFS-SERVER配置
# echo '/nfs_dir *(rw,sync,no_root_squash)' > /etc/exports

# 重启服务
# systemctl restart rpcbind.service
# systemctl restart nfs-utils.service 
# systemctl restart nfs-server.service 

# 查看
# showmount -e 10.0.1.201
Export list for 10.0.1.201:
/nfs_dir *

# 开机启动
# systemctl enable nfs-kernel-server.service
```



##### 部署postgresql

```
【OPS机器10.0.1.21】
# 准备镜像
# docker pull swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/postgres:12.19-alpine3.20

# docker tag  swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/postgres:12.19-alpine3.20  harbor.alnk.com/public/postgres:12.19-alpine3.20

# docker push harbor.alnk.com/public/postgres:12.19-alpine3.20

# 创建yaml文件目录
# mkdir -p /data/k8s-yaml/postgresql
# cd /data/k8s-yaml/postgresql

# vi postgresql.yaml
---
# pv
apiVersion: v1
kind: PersistentVolume
metadata:
  #namespace: gitlab-ver130806
  name: gitlab-postgresql-data-ver130806
  labels:
    type: gitlab-postgresql-data-ver130806
spec:
  capacity:
    storage: 10Gi
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  storageClassName: nfs
  nfs:
    path: /nfs_dir/gitlab_postgresql_data_ver130806
    server: 10.0.1.201

---
# pvc
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  namespace: gitlab-ver130806
  name: gitlab-postgresql-data-ver130806-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: nfs
  selector:
    matchLabels:
      type: gitlab-postgresql-data-ver130806

---
apiVersion: v1
kind: Service
metadata:
  namespace: gitlab-ver130806
  name: postgresql
  labels:
    app: gitlab
    tier: postgreSQL
spec:
  ports:
    - port: 5432
  selector:
    app: gitlab
    tier: postgreSQL

---
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: gitlab-ver130806
  name: postgresql
  labels:
    app: gitlab
    tier: postgreSQL
spec:
  replicas: 1
  selector:
    matchLabels:
      app: gitlab
      tier: postgreSQL
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: gitlab
        tier: postgreSQL
    spec:
      #nodeSelector:
      #  gee/disk: "500g"
      containers:
        - image: harbor.alnk.com/public/postgres:12.19-alpine3.20
          name: postgresql
          env:
            - name: POSTGRES_USER
              value: gitlab
            - name: POSTGRES_DB
              value: gitlabhq_production
            - name: POSTGRES_PASSWORD
              value: bogeusepg
            - name: TZ
              value: Asia/Shanghai
          ports:
            - containerPort: 5432
              name: postgresql
          livenessProbe:
            exec:
              command:
              - sh
              - -c
              - exec pg_isready -U gitlab -h 127.0.0.1 -p 5432 -d gitlabhq_production
            initialDelaySeconds: 110
            timeoutSeconds: 5
            failureThreshold: 6
          readinessProbe:
            exec:
              command:
              - sh
              - -c
              - exec pg_isready -U gitlab -h 127.0.0.1 -p 5432 -d gitlabhq_production
            initialDelaySeconds: 20
            timeoutSeconds: 3
            periodSeconds: 5
#          resources:
#            requests:
#              cpu: 100m
#              memory: 512Mi
#            limits:
#              cpu: "1"
#              memory: 1Gi
          volumeMounts:
            - name: postgresql
              mountPath: /var/lib/postgresql/data
      volumes:
        - name: postgresql
          persistentVolumeClaim:
            claimName: gitlab-postgresql-data-ver130806-pvc
            

【master节点10.0.1.201上】
#  mkdir -p /nfs_dir/{gitlab_etc_ver130806,gitlab_log_ver130806,gitlab_opt_ver130806,gitlab_postgresql_data_ver130806}

# kubectl create ns gitlab-ver130806
# kubectl apply -f http://k8s-yaml.alnk.com/postgresql/postgresql.yaml

# kubectl -n gitlab-ver130806 get pod
NAME                          READY   STATUS    RESTARTS   AGE
postgresql-6f97f7dbbb-g5s9b   1/1     Running   0          51s

```





##### 部署redis

```
【OPS机器10.0.1.21】
# 准备镜像
# docker pull swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/redis:6.2.4-alpine

# docker tag  swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/redis:6.2.4-alpine  harbor.alnk.com/public/redis:6.2.4-alpine

# docker push harbor.alnk.com/public/redis:6.2.4-alpine

# 创建yaml文件目录
# mkdir -p /data/k8s-yaml/redis6.2
# cd /data/k8s-yaml/redis6.2

# vi redis.yaml
---
apiVersion: v1
kind: Service
metadata:
  namespace: gitlab-ver130806
  name: redis
  labels:
    app: gitlab
    tier: backend
spec:
  ports:
    - port: 6379
      targetPort: 6379
  selector:
    app: gitlab
    tier: backend

---
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: gitlab-ver130806
  name: redis
  labels:
    app: gitlab
    tier: backend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: gitlab
      tier: backend
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: gitlab
        tier: backend
    spec:
      #nodeSelector:
      #  gee/disk: "500g"
      containers:
        - image: harbor.alnk.com/public/redis:6.2.4-alpine
          name: redis
          command:
            - "redis-server"
          args:
            - "--requirepass"
            - "bogeuseredis"
#          resources:
#            requests:
#              cpu: "1"
#              memory: 2Gi
#            limits:
#              cpu: "1"
#              memory: 2Gi
          ports:
            - containerPort: 6379
              name: redis
          livenessProbe:
            exec:
              command:
              - sh
              - -c
              - "redis-cli ping"
            initialDelaySeconds: 30
            periodSeconds: 10
            timeoutSeconds: 5
            successThreshold: 1
            failureThreshold: 3
          readinessProbe:
            exec:
              command:
              - sh
              - -c
              - "redis-cli ping"
            initialDelaySeconds: 5
            periodSeconds: 10
            timeoutSeconds: 1
            successThreshold: 1
            failureThreshold: 3
      initContainers:
      - command:
        - /bin/sh
        - -c
        - |
          ulimit -n 65536
          mount -o remount rw /sys
          echo never > /sys/kernel/mm/transparent_hugepage/enabled
          mount -o remount rw /proc/sys
          echo 2000 > /proc/sys/net/core/somaxconn
          echo 1 > /proc/sys/vm/overcommit_memory
        image: registry.cn-beijing.aliyuncs.com/acs/busybox:v1.29.2
        imagePullPolicy: IfNotPresent
        name: init-redis
        resources: {}
        securityContext:
          privileged: true
          procMount: Default
          
【master节点10.0.1.201上】
# kubectl apply -f http://k8s-yaml.alnk.com/redis6.2/redis.yaml

# kubectl -n gitlab-ver130806 get pod redis-77498b78f-6bt82
NAME                    READY   STATUS    RESTARTS   AGE
redis-77498b78f-6bt82   1/1     Running   0          54s

```





##### 部署gitlab

> 先定制一下镜像

```
【OPS机器10.0.1.21】
# 准备镜像
# docker pull swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/gitlab/gitlab-ce:17.5.1-ce.0

# docker tag  swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/gitlab/gitlab-ce:17.5.1-ce.0  docker.io/gitlab/gitlab-ce:17.5.1-ce.0


# mkdir -p /data/k8s-yaml/gitlab13
# cd /data/k8s-yaml/gitlab13/

# vi Dockerfile
FROM  docker.io/gitlab/gitlab-ce:17.5.1-ce.0

RUN rm /etc/apt/sources.list \
    && echo "deb http://apt.postgresql.org/pub/repos/apt/ jammy-pgdg main" > /etc/apt/sources.list.d/pgdg.list \
    && wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
COPY sources.list /etc/apt/sources.list

RUN apt-get update -yq && \
    apt-get install -y libpq5 vim iproute2 net-tools iputils-ping curl wget software-properties-common unzip postgresql-client-12 && \
    rm -rf /var/cache/apt/archives/*

RUN ln -svf /usr/bin/pg_dump /opt/gitlab/embedded/bin/pg_dump

# vi sources.list
deb http://mirrors.aliyun.com/ubuntu/dists xenial main
deb-src http://mirrors.aliyun.com/ubuntu/dists xenial main
deb http://mirrors.aliyun.com/ubuntu/dists xenial-updates main
deb-src http://mirrors.aliyun.com/ubuntu/dists xenial-updates main
deb http://mirrors.aliyun.com/ubuntu/dists xenial universe
deb-src http://mirrors.aliyun.com/ubuntu/dists xenial universe
deb http://mirrors.aliyun.com/ubuntu/dists xenial-updates universe
deb-src http://mirrors.aliyun.com/ubuntu/dists xenial-updates universe
deb http://mirrors.aliyun.com/ubuntu/dists xenial-security main
deb-src http://mirrors.aliyun.com/ubuntu/dists xenial-security main
deb http://mirrors.aliyun.com/ubuntu/dists xenial-security universe
deb-src http://mirrors.aliyun.com/ubuntu/dists xenial-security universe

# 制作镜像
# docker build -t harbor.alnk.com/public/gitlab-ce:17.5.1-ce.0 .
# docker push harbor.alnk.com/public/gitlab-ce:17.5.1-ce.0
```





```

```

