FROM harbor.alnk.com/public/jenkins:2.483
USER root
RUN rm -rf /etc/apt/sources.list.d/* &&\
 echo "deb http://mirrors.aliyun.com/debian bookworm main non-free contrib" > /etc/apt/sources.list.d/sources.list &&\
 echo "deb http://mirrors.aliyun.com/debian bookworm-updates main non-free contrib" >> /etc/apt/sources.list.d/sources.list &&\
 echo "deb http://mirrors.aliyun.com/debian bookworm-backports main non-free contrib" >> /etc/apt/sources.list.d/sources.list &&\
 echo "deb-src http://mirrors.aliyun.com/debian bookworm main non-free contrib" >> /etc/apt/sources.list.d/sources.list &&\
 echo "deb-src http://mirrors.aliyun.com/debian bookworm-updates main non-free contrib" >> /etc/apt/sources.list.d/sources.list &&\
 echo "deb-src http://mirrors.aliyun.com/debian bookworm-backports main non-free contrib" >> /etc/apt/sources.list.d/sources.list &&\
 apt update

## 安装docker客户端
#  apt update && \
#  apt install -y apt-transport-https ca-certificates curl software-properties-common gnupg lsb-release && \
#  curl -fsSL https://download.docker.com/linux/debian/gpg|gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg && \
#  echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null && \
#  apt update && \
#  apt install -y docker-ce docker-ce-cli containerd.io

## docker相关设置 ##
# harbor仓库不使用https以及harbor仓库账号密码
ADD daemon.json /etc/docker/daemon.json
ADD config.json /root/.docker/config.json
# 添加docker命令
ADD docker /usr/bin/

# 添加kubectl命令
Add kubectl /usr/bin/

# gitlab密钥
ADD id_rsa /root/.ssh/id_rsa

# 增加maven环境
ADD ./apache-maven-3.9.0-bin.tar.gz /usr/local/ 
ENV MAVEN_HOME=/usr/local/apache-maven-3.9.0
ENV PATH=$JAVA_HOME/bin:$MAVEN_HOME/bin:$PATH

RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime &&\ 
    echo 'Asia/Shanghai' >/etc/timezone

RUN echo "    StrictHostKeyChecking no" >> /etc/ssh/ssh_config

## 镜像构建
## docker build -t harbor.alnk.com/public/jenkins:2.483-docker .
## docker push harbor.alnk.com/public/jenkins:2.483-docker