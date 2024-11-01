```shell
### 在命令行使用ctr拉取带有账号密码的仓库镜像
# ctr -n k8s.io images pull  -u 账号:密码 -k registry.cn-hangzhou.aliyuncs.com/alnktest/hello:v0.0.1
# ctr -n k8s.io images pull -u 账号:密码 -k registry.cn-hangzhou.aliyuncs.com/alnktest/nginx:1.21.6

## containerd配置harbor仓库
# vi /etc/containerd/config.toml
    [plugins."io.containerd.grpc.v1.cri".registry]
      [plugins."io.containerd.grpc.v1.cri".registry.auths]
      [plugins."io.containerd.grpc.v1.cri".registry.configs]
        [plugins."io.containerd.grpc.v1.cri".registry.configs."harbor.alnk.com".tls]
          insecure_skip_verify = true
      [plugins."io.containerd.grpc.v1.cri".registry.headers]
      [plugins."io.containerd.grpc.v1.cri".registry.mirrors]
        [plugins."io.containerd.grpc.v1.cri".registry.mirrors."harbor.alnk.com"]
          endpoint = ["http://harbor.alnk.com"]
   
# systemctl restart containerd.service
# crictl pull harbor.alnk.com/library/go-hello:v0.0.1

```
