##### kubeasz信息

> 地址：https://github.com/easzlab/kubeasz
>
> kubeasz版本和k8s版本对照表
>
> ![1722861171865](./images/1722861171865.png)    



##### 安装步骤

###### 1.ssh免密登录

```shell
# ssh-keygen
Generating public/private rsa key pair.
Enter file in which to save the key (/root/.ssh/id_rsa):
Created directory '/root/.ssh'.
Enter passphrase (empty for no passphrase):
Enter same passphrase again:
Your identification has been saved in /root/.ssh/id_rsa.
Your public key has been saved in /root/.ssh/id_rsa.pub.
The key fingerprint is:
SHA256:XUZtM6IQrPJZGHanvQFT0JIefliqLPzJfQg8Oofh7Zw root@ops
The key's randomart image is:
+---[RSA 2048]----+
|       .o*. ..   |
|      o X +.. =  |
|     . B & .oo o |
|    . o B.=o     |
|   . = +S..o     |
|    + O   .      |
|   . O = .       |
|    =.*.o .      |
|     +E  .       |
+----[SHA256]-----+

# ssh-copy-id 10.0.1.201
# ssh-copy-id 10.0.1.202
# ssh-copy-id 10.0.1.203
# ssh-copy-id 10.0.1.204

```



###### 2.安装ansible

```shell
root@ops:/etc/kubeasz# apt-get install ansible -y

root@ops:/etc/kubeasz# ansible --version
ansible 2.10.8
  config file = /etc/kubeasz/ansible.cfg
  configured module search path = ['/root/.ansible/plugins/modules', '/usr/share/ansible/plugins/modules']
  ansible python module location = /usr/lib/python3/dist-packages/ansible
  executable location = /usr/bin/ansible
  python version = 3.10.12 (main, Nov 20 2023, 15:14:05) [GCC 11.4.0]
```



###### 3.下载工具脚本ezdown

```shell
root@ops:~# export release=3.2.0
root@ops:~# wget https://github.com/easzlab/kubeasz/releases/download/${release}/ezdown
root@ops:~# chmod +x ./ezdown
root@ops:~# ls -l
total 2376820
-rwxr-xr-x  1 root root      13660 Dec  7  2021 ezdown

# 下载kubeasz代码、二进制、默认容器镜像（更多关于ezdown的参数，运行./ezdown 查看）
# 如果下载不下来的话，可以买个香港的云服务器进行下载
root@ops:~# # 国内环境
root@ops:~# ./ezdown -D
root@ops:~# # 海外环境
root@ops:~# #./ezdown -D -m standard

## 上述脚本运行成功后，所有文件（kubeasz代码、二进制、离线镜像）均已整理好放入目录 /etc/kubeasz
```



###### 4.创建集群配置实例

```shell
# 容器化运行kubeasz
root@ops:~# ./ezdown -S
2024-08-05 20:47:51 INFO Action begin: start_kubeasz_docker
f1417ff83b31: Loading layer [==================================================>]  7.338MB/7.338MB
afe664e55619: Loading layer [==================================================>]  2.729MB/2.729MB
db35aecc3002: Loading layer [==================================================>]  33.12MB/33.12MB
d5584830d725: Loading layer [==================================================>]   5.12kB/5.12kB
c45401fc392c: Loading layer [==================================================>]  11.64MB/11.64MB
16d2ad51cdea: Loading layer [==================================================>]  108.8MB/108.8MB
49589bcfb288: Loading layer [==================================================>]  2.909MB/2.909MB
Loaded image: easzlab/kubeasz:3.6.2
2024-08-05 20:47:54 INFO try to run kubeasz in a container
2024-08-05 20:47:54 DEBUG get host IP: 10.0.1.21
786e933221041d82bbed112ddc8d345ed3131df22b09f77feba22f8e12cbf203
2024-08-05 20:47:55 INFO Action successed: start_kubeasz_docker

root@ops:~# docker ps
CONTAINER ID        IMAGE                   COMMAND               CREATED             STATUS              PORTS               NAMES
786e93322104        easzlab/kubeasz:3.6.2   "tail -f /dev/null"   36 seconds ago      Up 34 seconds                           kubeasz

# 创建新集群 k8s-01
root@ops:~# docker exec -it kubeasz ezctl new k8s-01
2024-08-05 20:48:52 DEBUG generate custom cluster files in /etc/kubeasz/clusters/k8s-01
2024-08-05 20:48:52 DEBUG set versions
2024-08-05 20:48:52 DEBUG cluster k8s-01: files successfully created.
2024-08-05 20:48:52 INFO next steps 1: to config '/etc/kubeasz/clusters/k8s-01/hosts'
2024-08-05 20:48:52 INFO next steps 2: to config '/etc/kubeasz/clusters/k8s-01/config.yml'

## 然后根据提示配置'/etc/kubeasz/clusters/k8s-01/hosts' 和 '/etc/kubeasz/clusters/k8s-01/config.yml'：根据前面节点规划修改hosts 文件和其他集群层面的主要配置选项；其他集群组件等配置项可以在config.yml 文件中修改
root@ops:~# pwd
/etc/kubeasz/clusters/k8s-01
root@ops:~# ls -l
total 12
-rw-r--r-- 1 root root 7414 Aug  5 20:48 config.yml
-rw-r--r-- 1 root root 2374 Aug  5 20:48 hosts
```



###### 5.修改配置文件

```shell
## 修改/etc/kubeasz/clusters/k8s-01/hosts
root@ops:~# cd /etc/kubeasz/clusters/k8s-01/
root@ops:~# cat hosts
# 'etcd' cluster should have odd member(s) (1,3,5,...)
[etcd]
10.0.1.201
10.0.1.202
10.0.1.203

# master node(s)
[kube_master]
10.0.1.201
10.0.1.202

# work node(s)
[kube_node]
10.0.1.203
10.0.1.204

# [optional] harbor server, a private docker registry
# 'NEW_INSTALL': 'true' to install a harbor server; 'false' to integrate with existed one
[harbor]
#192.168.1.8 NEW_INSTALL=false

# [optional] loadbalance for accessing k8s from outside
[ex_lb]
10.0.1.201 LB_ROLE=backup EX_APISERVER_VIP=10.0.1.200 EX_APISERVER_PORT=8443
10.0.1.202 LB_ROLE=master EX_APISERVER_VIP=10.0.1.200 EX_APISERVER_PORT=8443

# [optional] ntp server for the cluster
[chrony]
#192.168.1.1

[all:vars]
# --------- Main Variables ---------------
# Secure port for apiservers
SECURE_PORT="6443"

# Cluster container-runtime supported: docker, containerd
CONTAINER_RUNTIME="docker"

# Network plugins supported: calico, flannel, kube-router, cilium, kube-ovn
CLUSTER_NETWORK="flannel"

# Service proxy mode of kube-proxy: 'iptables' or 'ipvs'
PROXY_MODE="ipvs"

# K8S Service CIDR, not overlap with node(host) networking
SERVICE_CIDR="10.68.0.0/16"

# Cluster CIDR (Pod CIDR), not overlap with node(host) networking
CLUSTER_CIDR="172.20.0.0/16"

# NodePort Range
NODE_PORT_RANGE="30000-32767"

# Cluster DNS Domain
CLUSTER_DNS_DOMAIN="cluster.local"

# -------- Additional Variables (don't change the default value right now) ---
# Binaries Directory
bin_dir="/opt/kube/bin"

# Deploy Directory (kubeasz workspace)
base_dir="/etc/kubeasz"

# Directory for a specific cluster
cluster_dir="{{ base_dir }}/clusters/k8s-01"

# CA and other components cert/key Directory
ca_dir="/etc/kubernetes/ssl"


```

```shell
## 修改/etc/kubeasz/clusters/k8s-01/config.yml
root@ops:~# cd /etc/kubeasz/clusters/k8s-01/
root@ops:~# cat config.yml
############################
# prepare
############################
# 可选离线安装系统软件包 (offline|online)
INSTALL_SOURCE: "online"

# 可选进行系统安全加固 github.com/dev-sec/ansible-collection-hardening
OS_HARDEN: false

# 设置时间源服务器【重要：集群内机器时间必须同步】
ntp_servers:
  - "ntp1.aliyun.com"
  - "time1.cloud.tencent.com"
  - "0.cn.pool.ntp.org"

# 设置允许内部时间同步的网络段，比如"10.0.0.0/8"，默认全部允许
local_network: "0.0.0.0/0"


############################
# role:deploy
############################
# default: ca will expire in 100 years
# default: certs issued by the ca will expire in 50 years
CA_EXPIRY: "876000h"
CERT_EXPIRY: "438000h"

# kubeconfig 配置参数
CLUSTER_NAME: "cluster1"
CONTEXT_NAME: "context-{{ CLUSTER_NAME }}"

# k8s version
K8S_VER: "1.23.1"

############################
# role:etcd
############################
# 设置不同的wal目录，可以避免磁盘io竞争，提高性能
ETCD_DATA_DIR: "/var/lib/etcd"
ETCD_WAL_DIR: ""


############################
# role:runtime [containerd,docker]
############################
# ------------------------------------------- containerd
# [.]启用容器仓库镜像
ENABLE_MIRROR_REGISTRY: false

# [containerd]基础容器镜像
SANDBOX_IMAGE: "easzlab/pause:3.6"

# [containerd]容器持久化存储目录
CONTAINERD_STORAGE_DIR: "/var/lib/containerd"

# ------------------------------------------- docker
# [docker]容器存储目录
DOCKER_STORAGE_DIR: "/var/lib/docker"

# [docker]开启Restful API
ENABLE_REMOTE_API: false

# [docker]信任的HTTP仓库
INSECURE_REG: '["127.0.0.1/8"]'


############################
# role:kube-master
############################
# k8s 集群 master 节点证书配置，可以添加多个ip和域名（比如增加公网ip和域名）
MASTER_CERT_HOSTS:
  - "10.1.1.1"
  - "k8s.test.io"
  #- "www.test.com"

# node 节点上 pod 网段掩码长度（决定每个节点最多能分配的pod ip地址）
# 如果flannel 使用 --kube-subnet-mgr 参数，那么它将读取该设置为每个节点分配pod网段
# https://github.com/coreos/flannel/issues/847
NODE_CIDR_LEN: 24


############################
# role:kube-node
############################
# Kubelet 根目录
KUBELET_ROOT_DIR: "/var/lib/kubelet"

# node节点最大pod 数
MAX_PODS: 110

# 配置为kube组件（kubelet,kube-proxy,dockerd等）预留的资源量
# 数值设置详见templates/kubelet-config.yaml.j2
KUBE_RESERVED_ENABLED: "no"

# k8s 官方不建议草率开启 system-reserved, 除非你基于长期监控，了解系统的资源占用状况；
# 并且随着系统运行时间，需要适当增加资源预留，数值设置详见templates/kubelet-config.yaml.j2
# 系统预留设置基于 4c/8g 虚机，最小化安装系统服务，如果使用高性能物理机可以适当增加预留
# 另外，集群安装时候apiserver等资源占用会短时较大，建议至少预留1g内存
SYS_RESERVED_ENABLED: "no"

# haproxy balance mode
BALANCE_ALG: "roundrobin"


############################
# role:network [flannel,calico,cilium,kube-ovn,kube-router]
############################
# ------------------------------------------- flannel
# [flannel]设置flannel 后端"host-gw","vxlan"等
FLANNEL_BACKEND: "vxlan"
DIRECT_ROUTING: false

# [flannel] flanneld_image: "quay.io/coreos/flannel:v0.10.0-amd64"
flannelVer: "v0.15.1"
flanneld_image: "easzlab/flannel:{{ flannelVer }}"

# [flannel]离线镜像tar包
flannel_offline: "flannel_{{ flannelVer }}.tar"

# ------------------------------------------- calico
# [calico]设置 CALICO_IPV4POOL_IPIP=“off”,可以提高网络性能，条件限制详见 docs/setup/calico.md
CALICO_IPV4POOL_IPIP: "Always"

# [calico]设置 calico-node使用的host IP，bgp邻居通过该地址建立，可手工指定也可以自动发现
IP_AUTODETECTION_METHOD: "can-reach={{ groups['kube_master'][0] }}"

# [calico]设置calico 网络 backend: brid, vxlan, none
CALICO_NETWORKING_BACKEND: "brid"

# [calico]更新支持calico 版本: [v3.3.x] [v3.4.x] [v3.8.x] [v3.15.x]
calico_ver: "v3.19.3"

# [calico]calico 主版本
calico_ver_main: "{{ calico_ver.split('.')[0] }}.{{ calico_ver.split('.')[1] }}"

# [calico]离线镜像tar包
calico_offline: "calico_{{ calico_ver }}.tar"

# ------------------------------------------- cilium
# [cilium]CILIUM_ETCD_OPERATOR 创建的 etcd 集群节点数 1,3,5,7...
ETCD_CLUSTER_SIZE: 1

# [cilium]镜像版本
cilium_ver: "v1.4.1"

# [cilium]离线镜像tar包
cilium_offline: "cilium_{{ cilium_ver }}.tar"

# ------------------------------------------- kube-ovn
# [kube-ovn]选择 OVN DB and OVN Control Plane 节点，默认为第一个master节点
OVN_DB_NODE: "{{ groups['kube_master'][0] }}"

# [kube-ovn]离线镜像tar包
kube_ovn_ver: "v1.5.3"
kube_ovn_offline: "kube_ovn_{{ kube_ovn_ver }}.tar"

# ------------------------------------------- kube-router
# [kube-router]公有云上存在限制，一般需要始终开启 ipinip；自有环境可以设置为 "subnet"
OVERLAY_TYPE: "full"

# [kube-router]NetworkPolicy 支持开关
FIREWALL_ENABLE: "true"

# [kube-router]kube-router 镜像版本
kube_router_ver: "v0.3.1"
busybox_ver: "1.28.4"

# [kube-router]kube-router 离线镜像tar包
kuberouter_offline: "kube-router_{{ kube_router_ver }}.tar"
busybox_offline: "busybox_{{ busybox_ver }}.tar"


############################
# role:cluster-addon
############################
# coredns 自动安装
dns_install: "yes"
corednsVer: "1.8.6"
ENABLE_LOCAL_DNS_CACHE: true
dnsNodeCacheVer: "1.21.1"
# 设置 local dns cache 地址
LOCAL_DNS_CACHE: "169.254.20.10"

# metric server 自动安装
metricsserver_install: "yes"
metricsVer: "v0.5.2"

# dashboard 自动安装
dashboard_install: "yes"
dashboardVer: "v2.4.0"
dashboardMetricsScraperVer: "v1.0.7"

# ingress 自动安装
ingress_install: "no"
ingress_backend: "traefik"
traefik_chart_ver: "10.3.0"

# prometheus 自动安装
prom_install: "no"
prom_namespace: "monitor"
prom_chart_ver: "12.10.6"

# nfs-provisioner 自动安装
nfs_provisioner_install: "no"
nfs_provisioner_namespace: "kube-system"
nfs_provisioner_ver: "v4.0.2"
nfs_storage_class: "managed-nfs-storage"
nfs_server: "192.168.1.10"
nfs_path: "/data/nfs"

############################
# role:harbor
############################
# harbor version，完整版本号
HARBOR_VER: "v2.1.3"
HARBOR_DOMAIN: "harbor.yourdomain.com"
HARBOR_TLS_PORT: 8443

# if set 'false', you need to put certs named harbor.pem and harbor-key.pem in directory 'down'
HARBOR_SELF_SIGNED_CERT: true

# install extra component
HARBOR_WITH_NOTARY: false
HARBOR_WITH_TRIVY: false
HARBOR_WITH_CLAIR: false
HARBOR_WITH_CHARTMUSEUM: true

```



###### 6.查看命令帮助

```shell
root@ops:/etc/kubeasz# ./ezctl --help
Usage: ezctl COMMAND [args]
-------------------------------------------------------------------------------------
Cluster setups:
    list                             to list all of the managed clusters
    checkout    <cluster>            to switch default kubeconfig of the cluster
    new         <cluster>            to start a new k8s deploy with name 'cluster'
    setup       <cluster>  <step>    to setup a cluster, also supporting a step-by-step way
    start       <cluster>            to start all of the k8s services stopped by 'ezctl stop'
    stop        <cluster>            to stop all of the k8s services temporarily
    upgrade     <cluster>            to upgrade the k8s cluster
    destroy     <cluster>            to destroy the k8s cluster
    backup      <cluster>            to backup the cluster state (etcd snapshot)
    restore     <cluster>            to restore the cluster state from backups
    start-aio                        to quickly setup an all-in-one cluster with 'default' settings

Cluster ops:
    add-etcd    <cluster>  <ip>      to add a etcd-node to the etcd cluster
    add-master  <cluster>  <ip>      to add a master node to the k8s cluster
    add-node    <cluster>  <ip>      to add a work node to the k8s cluster
    del-etcd    <cluster>  <ip>      to delete a etcd-node from the etcd cluster
    del-master  <cluster>  <ip>      to delete a master node from the k8s cluster
    del-node    <cluster>  <ip>      to delete a work node from the k8s cluster

Extra operation:
    kcfg-adm    <cluster>  <args>    to manage client kubeconfig of the k8s cluster

Use "ezctl help <command>" for more information about a given command.


root@ops:/etc/kubeasz# ./ezctl setup --help
Usage: ezctl setup <cluster> <step>
available steps:
    01  prepare            to prepare CA/certs & kubeconfig & other system settings
    02  etcd               to setup the etcd cluster
    03  container-runtime  to setup the container runtime(docker or containerd)
    04  kube-master        to setup the master nodes
    05  kube-node          to setup the worker nodes
    06  network            to setup the network plugin
    07  cluster-addon      to setup other useful plugins
    90  all                to run 01~07 all at once
    10  ex-lb              to install external loadbalance for accessing k8s from outside
    11  harbor             to install a new harbor server or to integrate with an existed one

examples: ./ezctl setup test-k8s 01  (or ./ezctl setup test-k8s prepare)
          ./ezctl setup test-k8s 02  (or ./ezctl setup test-k8s etcd)
          ./ezctl setup test-k8s all
          ./ezctl setup test-k8s 04 -t restart_master


```



###### 7.部署

```shell
# 可以分步骤安装
root@ops:/etc/kubeasz# ./ezctl setup k8s-01 01
root@ops:/etc/kubeasz# ./ezctl setup k8s-01 02
root@ops:/etc/kubeasz# ./ezctl setup k8s-01 03
root@ops:/etc/kubeasz# ./ezctl setup k8s-01 04
root@ops:/etc/kubeasz# ./ezctl setup k8s-01 05
root@ops:/etc/kubeasz# ./ezctl setup k8s-01 06
root@ops:/etc/kubeasz# ./ezctl setup k8s-01 07
root@ops:/etc/kubeasz# ./ezctl setup k8s-01 10

# 或者一条命令安装
root@ops:/etc/kubeasz# ./ezctl setup k8s-01 all
root@ops:/etc/kubeasz# ./ezctl setup k8s-01 10
```



###### 8.修改./kube/config文件

```shell
# 使用kubeasz安装完成以后，会在/root/.kube/目录下生成config文件
# 修改config文件中server地址为vip地址10.0.1.200:8443

root@ops:~# ll /root/.kube/config
-r-------- 1 root root 6198 Aug  5 21:24 /root/.kube/config
root@ops:~# chmod u+w /root/.kube/config
root@ops:~# ll /root/.kube/config
-rw------- 1 root root 6198 Aug  5 21:24 /root/.kube/config
root@ops:~# vi /root/.kube/config

root@ops:~# cat /root/.kube/config
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURtakNDQW9LZ0F3SUJBZ0lVS08xU0x0cWlwcUZDWUVBQWJaMXFGMXJWdXJJd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1pERUxNQWtHQTFVRUJoTUNRMDR4RVRBUEJnTlZCQWdUQ0VoaGJtZGFhRzkxTVFzd0NRWURWUVFIRXdKWQpVekVNTUFvR0ExVUVDaE1EYXpoek1ROHdEUVlEVlFRTEV3WlRlWE4wWlcweEZqQVVCZ05WQkFNVERXdDFZbVZ5CmJtVjBaWE10WTJFd0lCY05NalF3T0RBMU1UTXlNREF3V2hnUE1qRXlOREEzTVRJeE16SXdNREJhTUdReEN6QUoKQmdOVkJBWVRBa05PTVJFd0R3WURWUVFJRXdoSVlXNW5XbWh2ZFRFTE1Ba0dBMVVFQnhNQ1dGTXhEREFLQmdOVgpCQW9UQTJzNGN6RVBNQTBHQTFVRUN4TUdVM2x6ZEdWdE1SWXdGQVlEVlFRREV3MXJkV0psY201bGRHVnpMV05oCk1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBeFFmVElQV3JjT2d4Q1VWQzFxL0gKeEEwUXVKVjJGZWZFdXZrUUNtQ3gyY1BiajhuallZR0hJVnFSU214MDUyc0JqYk5pdUgvZjQ1azhINWNUSVZYcgo0aWdEeVZkdHdQNHMxV2JUVXN0TEV5cnU2YVdBOTV3TXEzWlNMb3R4TWthTDRQMXQ4VWFYcTA3Qk1wN1F6c2NOCkxyaHhaOXhKYmsyb2YxeDlqcFhpZWg3T2hySktaSCtUbHlXcHRvVXhsOU5tclZWSWpMcG1XYzU5UTZSRkNadGsKSldoM0Mwdm1saWwwUzdjV2YyNTdZeHF6Y2JkTHpuYWpCNmNNNjdyMTRpc1RDTENOZGN3VFVRMlh6T2FFYk14dApQVUZXSWNiNFBPREZPYUdjWlpmVHVtKzJNQjExZVUydkdlQ1VCeEhsRkZzZmEwcnlqY2dMM3FiQjlFU2RqME56CkJRSURBUUFCbzBJd1FEQU9CZ05WSFE4QkFmOEVCQU1DQVFZd0R3WURWUjBUQVFIL0JBVXdBd0VCL3pBZEJnTlYKSFE0RUZnUVVnL05LRnU1cVBOWTR2WTBvMlJpaXkyc1h5TkV3RFFZSktvWklodmNOQVFFTEJRQURnZ0VCQUFPLwpQNHM2bWtFeFRpYVZLM2dkaWZJUmptL2c5dy9ucmgrTENRSEcvMDA0dkxFNXMwdW1RQTNtTzVxRlNIT3ZnQ3E0CktjWkxhaHozTllCYzFFZGpBdXpObjdpOHV3TW9HR2lIQ3NWcjloVTY2eWRPVkgwcTlkK1Fqcm1RQkFhMW54UVYKVmFNT21ZWjBXR29OUmtub25hNDFTQnZvWVMxOFJDaDJyTnlyS1owUERKaFgvVncwdDRINEdWN2dhSkNJeCs2ZQpBdzFRblpleDByTUhXS1g2ajgxUzFnb0lXTXZPeUJjRHJDc0d2SU5OVzdCaVFsQVpHOHEvNGdlWlJrNFpqZEVFCm83ay9QQTlENHVGUS83cDdRNXhtK0kzY1hGNWFhWm9pUTQxcDRsVVdjV21zbzBwTHJyRGlCT0JQUU9UbmUxYjUKbUdlZ3hIUEwzb2JEbFFiRlRiRT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
    server: https://10.0.1.200:8443
  name: cluster1
contexts:
- context:
    cluster: cluster1
    user: admin
  name: context-cluster1
current-context: context-cluster1
kind: Config
preferences: {}
users:
- name: admin
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUQyakNDQXNLZ0F3SUJBZ0lVT1Jod1BtalhmbHhVWWxXRzNwUkJkbm9mQ3FRd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1pERUxNQWtHQTFVRUJoTUNRMDR4RVRBUEJnTlZCQWdUQ0VoaGJtZGFhRzkxTVFzd0NRWURWUVFIRXdKWQpVekVNTUFvR0ExVUVDaE1EYXpoek1ROHdEUVlEVlFRTEV3WlRlWE4wWlcweEZqQVVCZ05WQkFNVERXdDFZbVZ5CmJtVjBaWE10WTJFd0lCY05NalF3T0RBMU1UTXlNREF3V2hnUE1qQTNOREEzTWpReE16SXdNREJhTUdjeEN6QUoKQmdOVkJBWVRBa05PTVJFd0R3WURWUVFJRXdoSVlXNW5XbWh2ZFRFTE1Ba0dBMVVFQnhNQ1dGTXhGekFWQmdOVgpCQW9URG5ONWMzUmxiVHB0WVhOMFpYSnpNUTh3RFFZRFZRUUxFd1pUZVhOMFpXMHhEakFNQmdOVkJBTVRCV0ZrCmJXbHVNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQXFJREtQdThmUytyMnVaRnMKT29rb3NmWVFGWXRoZ1lyQWpqK1d0bU5jNS9pTHdrQWI1SVpidElGQzJoT1Z6dzR1VDJtUTlIeUQya1dpdUhSZgpBTWs1Y2s4dWluQi9SMWlHdktIRkxDRmUxTTgwbXhUWjgxZTArMnZScFNOZWRkdnVWK0gwbnFBVVhrVVNGbDRaClFVZDlLeEZjeWkrK3VvZTVaZStCK3ViTm5sOCtZQ21ZYW9JbWhwMVZuMG4zYzlHeDJ5dHIrbWhvdFIzMmNRd0MKMUs2N2xUcW9RQm1qbXFJZ2QrR1ZNbXFtbUFTVmcwakhOY242VlkxNlJJajlaZzZaTmhsckNNdVRzRmlGQnlaZQpWSkNHUGl5Znp6TDREeGViaWVSUkhVaHNmOEd4b2RsK0V1Uk9VdFgxM2VaTTZZemcxdXcvNXM0OXNiUFMyYmh1Ci9YNkVqUUlEQVFBQm8zOHdmVEFPQmdOVkhROEJBZjhFQkFNQ0JhQXdIUVlEVlIwbEJCWXdGQVlJS3dZQkJRVUgKQXdFR0NDc0dBUVVGQndNQ01Bd0dBMVVkRXdFQi93UUNNQUF3SFFZRFZSME9CQllFRk1OY0ZHemp3aSt2ZzJzYgo1eTlTREl4OGt2b0ZNQjhHQTFVZEl3UVlNQmFBRklQelNoYnVhanpXT0wyTktOa1lvc3RyRjhqUk1BMEdDU3FHClNJYjNEUUVCQ3dVQUE0SUJBUUFUV3pONWdsb3J5a1NSaDBkb08za2tocjhVcmJSQW5ySk5EQWo5eXZBa0U1aTMKa2EwaFR2TTNROXVWc3htODBTeW1sM3lLYjB1Q2Y4d3FvdzZHM24zQ2RlYUwvOG51eHUrSHBOR21Kd01Tejk3QgpGeDZ6eFcwNTNCVTlPaFA4dHZPUzdFenlMNE54K2pHV21Nbzd5NTNZYzBYNFdCUnJxS3hWN1ZDeDNSbzB6dnZOCjlaemgwTnk1Yy9rZkVIdHNZd3lRM0k1U2V3cnc0U29aSkNpYVREb3ZhZ2lVU2hxMW5ad1BSV0FKT3lxOHltOG8KRlNsSEZqQW8rb0UzU1Y2bmtRV0hici94WEswVCtmM1J3aUhaYndHZHR2ckFrTjNMRXN5WmxodFVTY3AzbmI2eAp5T3o5V3JkWUJjZ2R5N0xKbDk1YjBMSDY4RFh1TmNQeFNDeXVQd2t4Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBcUlES1B1OGZTK3IydVpGc09va29zZllRRll0aGdZckFqaitXdG1OYzUvaUx3a0FiCjVJWmJ0SUZDMmhPVnp3NHVUMm1ROUh5RDJrV2l1SFJmQU1rNWNrOHVpbkIvUjFpR3ZLSEZMQ0ZlMU04MG14VFoKODFlMCsydlJwU05lZGR2dVYrSDBucUFVWGtVU0ZsNFpRVWQ5S3hGY3lpKyt1b2U1WmUrQit1Yk5ubDgrWUNtWQphb0ltaHAxVm4wbjNjOUd4Mnl0cittaG90UjMyY1F3QzFLNjdsVHFvUUJtam1xSWdkK0dWTW1xbW1BU1ZnMGpICk5jbjZWWTE2UklqOVpnNlpOaGxyQ011VHNGaUZCeVplVkpDR1BpeWZ6ekw0RHhlYmllUlJIVWhzZjhHeG9kbCsKRXVST1V0WDEzZVpNNll6ZzF1dy81czQ5c2JQUzJiaHUvWDZFalFJREFRQUJBb0lCQVFDTzJva2I2OHdEcDlLKwpKZ3kwWDUzeHFlT2U4RWljODQ0bHlzRXlKWEVaZEx5LzFXd1FSTnd3WHJuWGVtMHlXaHBtMXQrK2RtV3VRZ1VmCkRnR1JBQWZFNGw1Wm5lQTZFMUkvVDRLQzFGTzZaV1I2NmFXUlNWVTFKZ1paVTIvOHlaVDZTcVVCYWtONlhHSm0KTmlhQkVtT0toTEMrOU5Way9EWHd0cW5pbFpQUzVCbVFHcEJVYW04a25Nc1QzT0xKYUhMRlEwTUhFbTBGMzVGRgpPeWo5djRJRjVVUWs3L3V2RGFqUFJnWmZzZ3VSQkRVaWpUSnB0Sm5ONlhiRThseXZjRkk0TWlDenZWdElxUGxlCjI5dG5OUFo2RlArSnJKZXYxelFuL09CV011V0Q1eFhjVEJ1MzNCUmN1MmxkWjhDNjd4aXJRQUIvUGExRThwRWsKSk1pbUkzZmRBb0dCQU53NUtEL0xkK09VM1FZSXA3MGxQWXFqT3o4UXpLSWtlNFA5WE9rdytUUE90bGlGem5lZgpzWHh1LzQ2TCt2czVpdFowS1phQTg2WTRZL0NNUjcxLzR2YUt4ZXZSbVZrM0RjY2tTQ2FjQkdJc3FwVVhxS0NHCnBHc1ZOd0o3YnhxcmhoNFA5elBLMlE3NmFTbUo4enFlTWRtTUpoR0ZrbnBNMmVhSWtYTU9wTkIzQW9HQkFNUGcKcGJ5bGlJQkxLNUpKODFBaHJWVW5DWnpFdVJWMWFrQzBncHo5N1YyZWY3czV1VmxuNkMxbExwMUE5bFZBOFRYTAo5TUZCUGpRU0ZzNGQ3ZFlVRmFzUDB3eFU1RDdjTCt2S3Z4alptdHF3T05XVE16ZFVMK0hab05MbGF1TWlOQmhKCjNRNVFkNnZYR01VVXd0TkFMdDN4UGVWS0JrRGgzQUtLWU4xMnNiZ2JBb0dBRFNITGFLSjFiN2k4eFZOV3pVeWYKTXRrdyt6M0JOaG4rMDR3VU1rT2RXSjJHK2hoZ2kzbVdWOWsybkFWMDNlNDhmVFZJRlpWeThnS0MweUZLVmQ1KwpaajA0T0N1emZVSnZLK1RaK0pOdEgzMlNYbm1lc0pQVzBodmR2K1FrWis2NmZLaHZFVU9UVmZWUXVBMWwxNlQvClMvMnpkM0FEb0E5ZEh3WWR4a0tsU1ZrQ2dZQlpCazhOY0VhYjJJNVRESjB6UERzbFNucko3M2NYVTZnWkJIR2cKbktBM1BvUmJPWjhPRFhXdXZCLzFoTUx3ZUhXb3Q2dmo4WjB0MlZMWUZ5NHpjQ2x3OTk0NTZwTmFKb1Q1SzhxeQpwcVFFNUxiUUN2anFHcTh3Zk5MbFJ6UFBTNHBWeDZ4YWh5UDh5K1FNSHFWMWtlUTdKeHUwakhKUEp0ZnhwNmJpCndNR0JKUUtCZ1FDb0w5S0NOUlBGZEhKZnBUZWdlTGRHY2VJSU9kTE9NeU5jK3hvVEI1UTFQUVRXaGU1WldjQysKaVZuVEMyZndjS01vQXpSa01LMnBzdFl1VDBGTTJ5eWhESVpTb0I1VzFmMUdPNXhQK2pMY2xzai92RUV1NjFpZgprRHZjOHpsYXZoWlNjbTdzYXV0cEZIcWtTK0NBVW9qdkxzdDR1R0JDRDV2V1JhekpTYTNQUUE9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
```



###### 9.拷贝文件和kubectl命令

```shell
# 拷贝config文件到主控制节点10.0.1.201/10.0.1.202
root@ops:~# scp /root/.kube/config root@10.0.1.201:/root/.kube/
root@ops:~# scp /root/.kube/config root@10.0.1.202:/root/.kube/

# 拷贝kubectl命令
# scp /etc/kubeasz/bin/kubectl root@10.0.1.201:/usr/bin/
# scp /etc/kubeasz/bin/kubectl root@10.0.1.202:/usr/bin/

## kubectl命令tab键补全
# sudo apt-get install bash-completion
# source /usr/share/bash-completion/bash_completion
# kubectl completion bash | sudo tee /etc/bash_completion.d/kubectl > /dev/null
```



###### 10.查看k8s集群状态

```shell
【10.0.1.201】
root@master-1:~# ll /root/.kube/config
-rw------- 1 root root 6198 Aug  6 00:22 /root/.kube/config

root@master-1:~# kubectl version
WARNING: This version information is deprecated and will be replaced with the output from kubectl version --short.  Use --output=yaml|json to get the full version.
Client Version: version.Info{Major:"1", Minor:"27", GitVersion:"v1.27.5", GitCommit:"93e0d7146fb9c3e9f68aa41b2b4265b2fcdb0a4c", GitTreeState:"clean", BuildDate:"2023-08-24T00:48:26Z", GoVersion:"go1.20.7", Compiler:"gc", Platform:"linux/amd64"}
Kustomize Version: v5.0.1
Server Version: version.Info{Major:"1", Minor:"27", GitVersion:"v1.27.5", GitCommit:"93e0d7146fb9c3e9f68aa41b2b4265b2fcdb0a4c", GitTreeState:"clean", BuildDate:"2023-08-24T00:42:11Z", GoVersion:"go1.20.7", Compiler:"gc", Platform:"linux/amd64"}

root@master-1:~# kubectl get componentstatuses
Warning: v1 ComponentStatus is deprecated in v1.19+
NAME                 STATUS    MESSAGE                         ERROR
controller-manager   Healthy   ok
scheduler            Healthy   ok
etcd-2               Healthy   {"health":"true","reason":""}
etcd-0               Healthy   {"health":"true","reason":""}
etcd-1               Healthy   {"health":"true","reason":""}

root@master-1:~# kubectl get cs
Warning: v1 ComponentStatus is deprecated in v1.19+
NAME                 STATUS    MESSAGE   ERROR
controller-manager   Healthy   ok
scheduler            Healthy   ok
etcd-2               Healthy
etcd-1               Healthy
etcd-0               Healthy

root@master-1:~# kubectl cluster-info
Kubernetes control plane is running at https://10.0.1.200:8443
CoreDNS is running at https://10.0.1.200:8443/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy
KubeDNSUpstream is running at https://10.0.1.200:8443/api/v1/namespaces/kube-system/services/kube-dns-upstream:dns/proxy
kubernetes-dashboard is running at https://10.0.1.200:8443/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy
To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.

root@master-1:~# kubectl get node -o wide
NAME       STATUS                     ROLES    AGE   VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION       CONTAINER-RUNTIME
master-1   Ready,SchedulingDisabled   master   11h   v1.27.5   10.0.1.201    <none>        Ubuntu 22.04.3 LTS   5.15.0-112-generic   containerd://1.6.23
master-2   Ready,SchedulingDisabled   master   11h   v1.27.5   10.0.1.202    <none>        Ubuntu 22.04.3 LTS   5.15.0-112-generic   containerd://1.6.23
node-1     Ready                      node     11h   v1.27.5   10.0.1.203    <none>        Ubuntu 22.04.3 LTS   5.15.0-112-generic   containerd://1.6.23
node-2     Ready                      node     11h   v1.27.5   10.0.1.204    <none>        Ubuntu 22.04.3 LTS   5.15.0-112-generic   containerd://1.6.23

root@master-1:~# kubectl get pod -A
NAMESPACE     NAME                                         READY   STATUS    RESTARTS        AGE
kube-system   calico-kube-controllers-67c67b9b5f-qscn9     1/1     Running   2 (9m46s ago)   11h
kube-system   calico-node-9rfzs                            1/1     Running   2 (9m44s ago)   11h
kube-system   calico-node-q8jbh                            1/1     Running   2 (9m46s ago)   11h
kube-system   calico-node-wdcwn                            1/1     Running   2 (9m46s ago)   11h
kube-system   calico-node-xv9dk                            1/1     Running   2 (9m41s ago)   11h
kube-system   coredns-65bc7b648d-zwx4x                     1/1     Running   2 (9m41s ago)   11h
kube-system   dashboard-metrics-scraper-5c876f54bd-2l5ms   1/1     Running   2 (9m41s ago)   11h
kube-system   kubernetes-dashboard-89b5448d6-dhksg         1/1     Running   3 (8m57s ago)   11h
kube-system   metrics-server-56774d6954-pdkf9              1/1     Running   2 (9m46s ago)   11h
kube-system   node-local-dns-q86jv                         1/1     Running   2 (9m44s ago)   11h
kube-system   node-local-dns-txkdg                         1/1     Running   2 (9m41s ago)   11h
kube-system   node-local-dns-vc9qr                         1/1     Running   2 (9m46s ago)   11h
kube-system   node-local-dns-w7xvt                         1/1     Running   2 (9m46s ago)   11h

root@master-1:~# kubectl get svc -A
NAMESPACE     NAME                        TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                  AGE
default       kubernetes                  ClusterIP   10.68.0.1       <none>        443/TCP                  11h
kube-system   dashboard-metrics-scraper   ClusterIP   10.68.67.77     <none>        8000/TCP                 11h
kube-system   kube-dns                    ClusterIP   10.68.0.2       <none>        53/UDP,53/TCP,9153/TCP   11h
kube-system   kube-dns-upstream           ClusterIP   10.68.165.209   <none>        53/UDP,53/TCP            11h
kube-system   kubernetes-dashboard        NodePort    10.68.14.76     <none>        443:30552/TCP            11h
kube-system   metrics-server              ClusterIP   10.68.50.197    <none>        443/TCP                  11h
kube-system   node-local-dns              ClusterIP   None            <none>        9253/TCP                 11h
```
