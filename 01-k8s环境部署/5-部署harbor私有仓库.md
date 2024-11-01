##### 部署harbor仓库

```shell
#### 注意docker版本，如果版本过低，安装harbor2.9.5 会出问题
#### 
root@ops:/opt/harbor# docker version
Client:
 Version:           24.0.7
 API version:       1.43
 Go version:        go1.21.1
 Git commit:        24.0.7-0ubuntu2~22.04.1
 Built:             Wed Mar 13 20:23:54 2024
 OS/Arch:           linux/amd64
 Context:           default

Server:
 Engine:
  Version:          24.0.7
  API version:      1.43 (minimum version 1.12)
  Go version:       go1.21.1
  Git commit:       24.0.7-0ubuntu2~22.04.1
  Built:            Wed Mar 13 20:23:54 2024
  OS/Arch:          linux/amd64
  Experimental:     false
 containerd:
  Version:          1.7.12
  GitCommit:
 runc:
  Version:          1.1.12-0ubuntu2~22.04.1
  GitCommit:
 docker-init:
  Version:          0.19.0
  GitCommit:

## ops运维机
root@ops:~# mkdir /opt/src && cd /opt/src

# harbor仓库在github上的地址
# https://github.com/goharbor/harbor

## 下载安装包
root@ops:/opt/src# wget https://github.com/goharbor/harbor/releases/download/v2.9.5/harbor-offline-installer-v2.9.5.tgz

## 解压 && 软连接
root@ops:/opt/src# tar -xf harbor-offline-installer-v2.9.5.tgz -C /opt/
root@ops:/opt/src# cd /opt/
root@ops:/opt# mv harbor/ harbor-v2.9.5
root@ops:/opt# ln -s harbor-v2.9.5/ harbor


## 修改配置
root@ops:/opt# cd harbor
root@ops:/opt/harbor# cp harbor.yml.tmpl harbor.yml
root@ops:/opt/harbor# vi harbor.yml
# 1. 修改端口
http:
  port: 180   #原80
# 2. 注释https
#https:
  # https port for harbor, default is 443
  #  port: 443
  # The path of cert and key files for nginx
  #certificate: /your/certificate/path
  #private_key: /your/private/key/path
# 3. 修改存储目录
data_volume: /data/harbor
location: /data/harbor/logs
# 4. 修改主机名
hostname: harbor.alnk.com


## 创建目录
root@ops:/opt/harbor# mkdir -p /data/harbor/logs


## kubeasz已经下载了docker-compose执行文件
root@ops:/opt/harbor# cp /etc/kubeasz/bin/docker-compose /usr/bin/
root@ops:/opt/harbor# docker-compose --version
Docker Compose version v2.20.3


## 安装
root@ops:/opt/harbor# ./install.sh


## 查看
root@ops:/opt/harbor# docker-compose ps
root@ops:/opt/harbor# docker ps


## 设置开机启动
# 停止harbor
root@ops:/opt/harbor# docker-compose down -v
# 创建服务文件
root@ops:/opt/harbor# vi /etc/systemd/system/harbor.service
[Unit]
Description=Harbor
After=docker.service systemd-networkd.service systemd-resolved.service
Requires=docker.service
Documentation=http://github.com/vmware/harbor

[Service]
Type=simple
Restart=on-failure
RestartSec=5
ExecStart=/usr/bin/docker-compose -f /opt/harbor/docker-compose.yml up
ExecStop=/usr/bin/docker-compose -f /opt/harbor/docker-compose.yml down

[Install]
WantedBy=multi-user.target

# 启动，并设置开机启动harbor
root@ops:/opt/harbor# systemctl daemon-reload
root@ops:/opt/harbor# systemctl start harbor.service
root@ops:/opt/harbor# systemctl status harbor
root@ops:/opt/harbor# systemctl enable harbor
```

##### 设置nginx反向代理

```shell
## 安装nginx
root@ops:/opt/harbor# sudo apt install -y nginx

root@ops:/opt/harbor# systemctl status nginx
root@ops:/opt/harbor# systemctl enable nginx

root@ops:/opt/harbor# vi /etc/nginx/conf.d/harbor.alnk.com.conf
server {
    listen       80;
    server_name  harbor.alnk.com;

    client_max_body_size 1000m;

    location / {
        proxy_pass http://127.0.0.1:180;
    }
}

root@ops:/opt/harbor# nginx -t
nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful

root@ops:/opt/harbor# systemctl reload nginx



```

##### docker登录harbor仓库

```shell
## 修改hosts
## 10.0.1.21、10.0.1.201、10.0.1.202、10.0.1.203、10.0.1.204都做一下
# vi /etc/hosts
10.0.1.21 harbor.alnk.com

# vi /etc/docker/daemon.json
{
  "insecure-registries" : ["harbor.alnk.com"]
}

# systemctl daemon-reload
# systemctl restart docker
# docker login harbor.alnk.com
Username: admin
Password: Harbor12345

# cat /root/.docker/config.json
{
        "auths": {
                "harbor.alnk.com": {
                        "auth": "YWRtaW46SGFyYm9yMTIzNDU="
                },
                "registry.cn-hangzhou.aliyuncs.com": {
                        "auth": "xxxxxxxxxxxxxxxx"
                }
        }

```

##### 浏览器访问

修改windows系统的hosts文件

![1722910056126](.\images\1722910056126.png)

```
10.0.1.21 harbor.alnk.com
```

默认账号密码 admin/Harbor12345

![1722910127716](.\images\1722910127716.png)

![1722910237821](.\images\1722910237821.png)
