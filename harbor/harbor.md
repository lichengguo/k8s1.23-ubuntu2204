```
# 安装docker
[root@harbor harbor]# sudo yum install -y yum-utils \
    device-mapper-persistent-data \
    lvm2
    
# 添加阿里云 Docker CE 镜像源
[root@harbor harbor]# sudo yum-config-manager --add-repo \
    https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo

# 更新 yum 缓存
[root@harbor harbor]# sudo yum makecache fast

# 查看可用的 Docker 版本
[root@harbor harbor]# yum list docker-ce --showduplicates | sort -r

# 安装最新版本的 Docker
[root@harbor harbor]# sudo yum install -y docker-ce docker-ce-cli containerd.io

# 启动 Docker 服务
[root@harbor harbor]# sudo systemctl start docker

# 设置开机自启
[root@harbor harbor]# sudo systemctl enable docker

# 检查服务状态
[root@harbor harbor]# sudo systemctl status docker
```



```
# 创建工作目录
[root@harbor harbor]# mkdir -p /harbor && cd /harbor

# 创建证书目录
[root@harbor harbor]# mkdir ssl & cd ssl

# 生成根证书私钥
[root@harbor harbor]# openssl genrsa -out ca.key 4096

# 生成根证书（有效期10年）
[root@harbor harbor]# openssl req -x509 -new -nodes -sha512 -days 3650 \
  -subj "/C=CN/ST=Beijing/L=Beijing/O=alnk/OU=IT/CN=alnk.com" \
  -key ca.key \
  -out ca.crt

# 生成服务器私钥
[root@harbor harbor]# openssl genrsa -out alnk.com.key 4096

# 创建证书签名请求（CSR）
[root@harbor harbor]# openssl req -sha512 -new \
  -subj "/C=CN/ST=Beijing/L=Beijing/O=alnk/OU=IT/CN=*.alnk.com" \
  -key alnk.com.key \
  -out alnk.com.csr

# 创建扩展配置文件
[root@harbor harbor]# cat > v3.ext <<EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1=*.alnk.com
DNS.2=harbor.alnk.com
DNS.3=gitlab.alnk.com
DNS.4=jenkins.alnk.com
DNS.5=argocd.alnk.com
DNS.6=jumpserver.alnk.com
DNS.7=rancher.alnk.com
EOF

# 生成证书（使用根证书签名）
[root@harbor harbor]# openssl x509 -req -sha512 -days 3650 \
  -extfile v3.ext \
  -CA ca.crt -CAkey ca.key -CAcreateserial \
  -in alnk.com.csr \
  -out alnk.com.crt

# 转换证书格式供 Docker 使用
[root@harbor harbor]# openssl x509 -inform PEM -in alnk.com.crt -out alnk.com.cert

# 文件
[root@harbor harbor]# ls -l
-rw-r--r-- 1 root root 2179 Nov  9 13:37 alnk.com.cert
-rw-r--r-- 1 root root 2179 Nov  9 13:36 alnk.com.crt
-rw-r--r-- 1 root root 1691 Nov  9 13:36 alnk.com.csr
-rw-r--r-- 1 root root 3243 Nov  9 13:36 alnk.com.key
-rw-r--r-- 1 root root 1992 Nov  9 13:36 ca.crt
-rw-r--r-- 1 root root 3243 Nov  9 13:36 ca.key
-rw-r--r-- 1 root root   17 Nov  9 13:36 ca.srl
-rw-r--r-- 1 root root  367 Nov  9 13:36 v3.ext

```



```
# 下载docker-compose
[root@harbor harbor]# wget https://github.com/docker/compose/releases/download/v2.40.3/docker-compose-linux-x86_64

[root@harbor harbor]# docker-compose version
Docker Compose version v2.40.3
```





```
# 下载最新版 Harbor
[root@harbor harbor]# cd /harbor
[root@harbor harbor]# wget https://github.com/goharbor/harbor/releases/download/v2.14.0/harbor-offline-installer-v2.14.0.tgz
[root@harbor harbor]# tar xvf harbor-offline-installer-v2.14.0.tgz
[root@harbor harbor]# cd harbor

[root@harbor harbor]# pwd
/harbor/harbor
[root@harbor harbor]# ls -l
total 656308
-rw-r--r-- 1 root root      3646 Sep  9 19:44 common.sh
-rw-r--r-- 1 root root 672014938 Sep  9 19:44 harbor.v2.14.0.tar.gz
-rw-r--r-- 1 root root     14688 Sep  9 19:44 harbor.yml.tmpl
-rwxr-xr-x 1 root root      1975 Sep  9 19:44 install.sh
-rw-r--r-- 1 root root     11347 Sep  9 19:44 LICENSE
-rwxr-xr-x 1 root root      2211 Sep  9 19:44 prepare


[root@harbor harbor]# cp harbor.yml.tmpl harbor.yml
[root@harbor harbor]# cat harbor.yml|grep -v "#" |grep -v "^$"
hostname: harbor.alnk.com
http:
  port: 80
https:
  port: 443
  certificate: /harbor/ssl/alnk.com.crt
  private_key: /harbor/ssl/alnk.com.key
harbor_admin_password: Harbor12345
database:
  password: root123
  max_idle_conns: 100
  max_open_conns: 900
  conn_max_lifetime: 5m
  conn_max_idle_time: 0
data_volume: /harbor/data
trivy:
  ignore_unfixed: false
  skip_update: false
  skip_java_db_update: false
  offline_scan: false
  security_check: vuln
  insecure: false
  timeout: 5m0s
jobservice:
  max_job_workers: 10
  max_job_duration_hours: 24
  job_loggers:
    - STD_OUTPUT
    - FILE
notification:
  webhook_job_max_retry: 3
log:
  level: info
  local:
    rotate_count: 50
    rotate_size: 200M
    location: /harbor/log/harbor
_version: 2.14.0
proxy:
  http_proxy:
  https_proxy:
  no_proxy:
  components:
    - core
    - jobservice
    - trivy
upload_purging:
  enabled: true
  age: 168h
  interval: 24h
  dryrun: false
cache:
  enabled: false
  expire_hours: 24
  
  
# 或者临时关闭防火墙（测试环境）
[root@harbor harbor]# sudo systemctl stop firewalld
[root@harbor harbor]# sudo systemctl disable firewalld  
  
# 执行安装脚本
[root@harbor harbor]# sudo ./install.sh


# 停止服务
[root@harbor harbor]# docker-compose down -v

# 启动服务
[root@harbor harbor]# docker-compose up -d

# 卸载 Harbor
[root@harbor harbor]# sudo ./uninstall.sh
```



```
# 自签SSL证书问题处理
[root@harbor harbor]# echo "10.0.1.20 harbor.alnk.com" >> /etc/hosts
[root@harbor harbor]# docker login harbor.alnk.com -u admin -p Harbor12345
WARNING! Using --password via the CLI is insecure. Use --password-stdin.
Error response from daemon: Get "https://harbor.alnk.com/v2/": tls: failed to verify certificate: x509: certificate signed by unknown authority


[root@harbor harbor]# sudo mkdir -p /etc/docker/certs.d/harbor.alnk.com
[root@harbor harbor]# sudo cp /harbor/ssl/alnk.com.cert /etc/docker/certs.d/harbor.alnk.com/ca.crt
[root@harbor harbor]# systemctl restart docker


[root@harbor harbor]# docker login harbor.alnk.com -u admin -p Harbor12345
WARNING! Using --password via the CLI is insecure. Use --password-stdin.
WARNING! Your password will be stored unencrypted in /root/.docker/config.json.
Configure a credential helper to remove this warning. See
https://docs.docker.com/engine/reference/commandline/login/#credentials-store

Login Succeeded
```

