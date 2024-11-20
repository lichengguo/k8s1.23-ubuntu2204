### jenkins

#### jenkins构建方式

> ```
> ### 在k8s中部署的jenkins
>
> 方式1（jenkins本地打包）
> 1. jenkins本地需要安装相关语言的编译环境进行支持，例如java需要安装maven环境，面对多版本的maven环境，安装需要留意；go需要安装go编译环境等
> 2. jenkins本地需要安装docker客户端，然后调用宿主机的docker服务端（docker in docker），把编译后的的代码包拷贝到业务运行容器中去，再进行业务镜像构建
> 3. 这种方式的Dockerfile在其他地方构建时需要有编译后的代码包支持
>
> 方式2（jenkins启用docker打包）
> 1. 不需要在jenkins本地安装相关编译环境
> 2. jenkins安装docker客户端，然后调用宿主机的docker服务端（docker in docker），把相关的代码拷贝到构建容器进行构建，最后用另外的业务容器运行
> 3. 这种方式的Dockerfile在其他地方构建时需要有相关的代码支持
>
> ## 方式2 案例
> # from一个go的构建镜像
> FROM harbor.alnk.com/public/golang:1.22.8 as builder
> ENV GOPROXY https://goproxy.cn
> COPY . /app/
> RUN cd  /app \
>     && go mod init go-k8s-one \
>     && go mod tidy \
>     && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o go-k8s-one .
> # 使用其他镜像运行业务
> FROM harbor.od.com/public/alpine:3.18
> RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
>     apk update && \
>     apk --no-cache add tzdata ca-certificates && \
>     cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
>     # apk del tzdata && \
>     rm -rf /var/cache/apk/*
> COPY --from=builder /app/go-k8s-one /app/go-k8s-one
> COPY --from=builder /app/conf/config.ini /app/conf/config.ini
> CMD ["/app/go-k8s-one"]
> ```

#### 准备jenkins镜像

> ```dockerfile
> FROM harbor.alnk.com/public/jenkins:2.483
> USER root
> RUN rm -rf /etc/apt/sources.list.d/* &&\
>     echo "deb http://mirrors.aliyun.com/debian bookworm main non-free contrib" > /etc/apt/sources.list.d/sources.list &&\
>     echo "deb http://mirrors.aliyun.com/debian bookworm-updates main non-free contrib" >> /etc/apt/sources.list.d/sources.list &&\
>     echo "deb http://mirrors.aliyun.com/debian bookworm-backports main non-free contrib" >> /etc/apt/sources.list.d/sources.list &&\
>    	echo "deb-src http://mirrors.aliyun.com/debian bookworm main non-free contrib" >> /etc/apt/sources.list.d/sources.list &&\
>     echo "deb-src http://mirrors.aliyun.com/debian bookworm-updates main non-free contrib" >> /etc/apt/sources.list.d/sources.list &&\
>     echo "deb-src http://mirrors.aliyun.com/debian bookworm-backports main non-free contrib" >> /etc/apt/sources.list.d/sources.list &&\
>     apt update && \
>     apt install -y apt-transport-https ca-certificates curl software-properties-common gnupg lsb-release && \
>     curl -fsSL https://download.docker.com/linux/debian/gpg|gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg && \
>     echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null && \
>     apt update && \
>     apt install -y docker-ce docker-ce-cli containerd.io
>
> ADD daemon.json /etc/docker/daemon.json
> ADD id_rsa /root/.ssh/id_rsa
> ADD config.json /root/.docker/config.json
>
> # 增加maven环境
> ADD ./apache-maven-3.9.0-bin.tar.gz /usr/local/ 
> ENV MAVEN_HOME=/usr/local/apache-maven-3.9.0
> ENV PATH=$JAVA_HOME/bin:$MAVEN_HOME/bin:$PATH
>
> #ADD get-docker.sh /get-docker.sh
> #RUN apt-get update -y
> RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime &&\ 
>     echo 'Asia/Shanghai' >/etc/timezone
> RUN echo "    StrictHostKeyChecking no" >> /etc/ssh/ssh_config 
> #RUN /get-docker.sh
> ```
> ```
> ### 重新打包上传harbor仓库
> # docker build -t harbor.alnk.com/public/jenkins:2.483-docker-maven .
> # docker push harbor.alnk.com/public/jenkins:2.483-docker-maven
>
> ### 更新jenkins镜像
> ```
