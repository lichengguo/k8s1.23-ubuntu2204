#### 编写golang程序

`main.go`代码

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// 初始化gin框架路由
	r := gin.Default()

	// 创建路由
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello Alnk!")
	})

	// 监听端口
	r.Run(":8080")
}
```



#### 编译代码

> 此处是编译好以后在上传镜像，也可以使用容器进行编译

```shell
# 安装代码依赖
lichengguo@MacBook-Pro hello % go mod init
lichengguo@MacBook-Pro hello % go mod tidy

# 编译成linux包
lichengguo@MacBook-Pro hello % CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o hello-linux
lichengguo@MacBook-Pro hello % ls -lh
total 39672
-rw-r--r--  1 lichengguo  staff    63B Aug 27 14:56 go.mod
-rw-r--r--  1 lichengguo  staff   4.9K Aug 27 14:56 go.sum
-rwxr-xr-x  1 lichengguo  staff   9.7M Aug 27 14:57 hello-linux
-rw-r--r--  1 lichengguo  staff   268B Aug 27 14:56 main.go
```



#### 编写Dockerfile文件

```dockerfile
# 测试直接使用centos镜像了，生产中建议使用alpine这种小的镜像
FROM centos

MAINTAINER alnk<1029612787@qq.com>

COPY hello-linux /hello-linux

EXPOSE 8080

# centos官方镜像默认的工作目录是根目录，所以这里不需要设置工作目录，直接运行即可
CMD ["./hello-linux"]
```



#### 上传程序和dockerfile文件到服务器

```shell
[root@alnk hello]# pwd
/root/hello
[root@alnk hello]# ll
-rw-r--r-- 1 root root       79 8月  27 15:11 Dockerfile
-rwxrwxrwx 1 root root 10150772 8月  27 14:57 hello-linux
```



#### 打包镜像

```shell
[root@alnk hello]# docker build -t alnk_app .
Sending build context to Docker daemon  10.15MB
Step 1/4 : FROM centos
....
Successfully built 9181f20cdbfd
Successfully tagged alnk_app:latest

[root@alnk hello]# docker images
REPOSITORY   TAG                IMAGE ID       CREATED         SIZE
alnk_app     latest             9181f20cdbfd   2 minutes ago   219MB
```



#### 启动容器

```shell
[root@alnk hello]# docker run -d -p 8080:8080 --name alnk_hello_01 alnk_app
ecd4d7e6686b9e464ffa043b9e8d96752eed53a5e52b74266d18a7f4257bb38f

[root@alnk hello]# docker ps 
CONTAINER ID   IMAGE      COMMAND           CREATED         STATUS         PORTS                                       NAMES
ecd4d7e6686b   alnk_app   "./hello-linux"   3 seconds ago   Up 3 seconds   0.0.0.0:8080->8080/tcp, :::8080->8080/tcp   alnk_hello_01
```



#### 测试

```shell
# 进入容器测试查看
[root@alnk hello]# docker exec -it alnk_hello_01 /bin/bash

[root@ecd4d7e6686b /]# ls -l
total 9964
lrwxrwxrwx   1 root root        7 Nov  3  2020 bin -> usr/bin
drwxr-xr-x   5 root root      340 Aug 27 07:19 dev
drwxr-xr-x   1 root root     4096 Aug 27 07:19 etc
-rwxrwxrwx   1 root root 10150772 Aug 27 06:57 hello-linux
......

[root@ecd4d7e6686b /]# ps -ef|grep hello-linux
root         1     0  0 07:19 ?        00:00:00 ./hello-linux
root        49    32  0 08:41 pts/0    00:00:00 grep --color=auto hello-linux


# 宿主机curl一下
[root@alnk hello]# curl localhost:8080/hello
hello Alnk!

# web界面
#IP:8080/hello
```