##### docker和传统虚拟化比较

> 传统虚拟机技术是虚拟出一套硬件后，在其上运行一个完整操作系统，在该系统上再运行所需应用进程
>
> 容器内的应用进程直接运行于宿主的内核，容器内没有自己的内核，而且也没有进行硬件虚拟。因此容器要比传统虚拟机更为轻便。

![1730475613956](images\1730475613956.png)

![1730475660250](images\1730475660250.png)

##### 为什么要使用docker

> - 更高效的利用系统资源
> - 更快速的启动时间
> - 一致的运行环境
> - 持续交付和部署
> - 更轻松的迁移
> - 更轻松的维护和扩展

| 特性       | 容器               | 虚拟机        |
| ---------- | ------------------ | ------------- |
| 启动       | 秒级               | 分钟级        |
| 硬盘使用   | 一般为 `MB`      | 一般为 `GB` |
| 性能       | 接近原生           | 弱于          |
| 系统支持量 | 单机支持上千个容器 | 一般几十个    |

##### docker核心概念或者组件

> - 镜像（`Image`）
> - 容器（`Container`）
> - 仓库（`Repository`）
> - Dockerfile

```shell
### 登陆docker镜像仓库
#docker login "仓库地址" -u "仓库用户名" -p "仓库密码"

### 从仓库下载镜像
docker pull "仓库地址"/"仓库命名空间"/"镜像名称":"版本号"

### 基于Dockerfile构建本地镜像
## 简单命令
# docker build -t "仓库地址"/"仓库命名空间"/"镜像名称":"镜像版本号" .

## 复杂命令
# docker build --network host --build-arg PYPI_IP="xx.xx.xx.xx" --cache-from "仓库地址"/"仓库命名空间"/"镜像名称":latest --tag "仓库地址"/"仓库命名空间"/"镜像名称":"镜像版本号" --tag "仓库地址"/"仓库命名空间"/"镜像名称":"版本号" .

### 将构建好的本地镜像推到远端镜像仓库里面
#docker push "仓库地址"/"仓库命名空间"/"镜像名称":"镜像版本号"

## 减少镜像大小方法，加速镜像构建
# 选择轻量级的基础镜像
# 编译型语言（go、java），打包镜像和运行镜像分开，分层构建
# 利用缓存，注意dockerfile中run指令尽量在一层，以便于后面构建可以直接使用缓存
```



##### docker镜像国内不能访问解决办法

> 1.可以在华为云购买香港服务器进行下载，按量按需按时间付费不贵
>
> 2.github开源工具https://github.com/tech-shrimp/docker_image_pusher

```shell
###阿里云docker镜像仓库资料
##空间名称：alnktest
##账号：1029612787@qq.com
##链接：registry.cn-hangzhou.aliyuncs.com
##密码：*******

##登录阿里云docker镜像仓库命令
# sudo docker login --username=1029612787@qq.com registry.cn-hangzhou.aliyuncs.com

##下载阿里云docker镜像
# docker pull registry.cn-hangzhou.aliyuncs.com/alnktest/alpine:3.19.1

```



