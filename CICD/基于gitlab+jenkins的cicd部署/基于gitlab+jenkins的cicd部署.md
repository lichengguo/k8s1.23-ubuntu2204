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



#### 安装部署Jenkins

```

```

