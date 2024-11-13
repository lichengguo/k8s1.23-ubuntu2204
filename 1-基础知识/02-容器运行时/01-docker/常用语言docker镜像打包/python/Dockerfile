#FROM registry.cn-hangzhou.aliyuncs.com/alnktest/python:3.9.17-slim-bullseye
FROM harbor.alnk.com/public/python:3.9.17-slim-bullseye
MAINTAINER 1029612787@qq.com

WORKDIR /app

COPY requirements.txt .

RUN  sed -i 's/deb.debian.org/ftp.cn.debian.org/g' /etc/apt/sources.list \
  && sed -i 's/security.debian.org/ftp.cn.debian.org/g' /etc/apt/sources.list \
  && apt-get update -y \
  && apt-get install -y wget gcc libsm6 libxext6 libglib2.0-0 libxrender1 make \
  && apt-get clean && apt-get autoremove -y && rm -rf /var/lib/apt/lists/*

# 安装python代码依赖
RUN pip install --no-cache-dir -i https://mirrors.aliyun.com/pypi/simple -r requirements.txt \
    && rm requirements.txt

# 指定多个pip安装源
# RUN pip install -r requirements.txt -i http://mirrors.aliyun.com/pypi/simple/ \
# --extra-index-url https://pypi.tuna.tsinghua.edu.cn/simple \
# --trusted-host mirrors.aliyun.com \
# --trusted-host pypi.tuna.tsinghua.edu.cn \

COPY . .

EXPOSE 5000
HEALTHCHECK CMD curl --fail http://localhost:5000 || exit 1

ENTRYPOINT ["gunicorn", "app:app", "-c", "gunicorn_config.py"]


##打包python项目的镜像并运行测试
# cd python
# docker build -t harbor.alnk.com/public/hello-python-flask:v0.0.1 .
# docker push harbor.alnk.com/public/hello-python-flask:v0.0.1
# docker run -d --name hello-python-flask -p 9999:5000 harbor.alnk.com/public/hello-python-flask:v0.0.1

# docker exec -it hello-python-flask bash
# root@ace8cc5407a7:/app# cat /etc/os-release
# PRETTY_NAME="Debian GNU/Linux 11 (bullseye)"
# NAME="Debian GNU/Linux"
# VERSION_ID="11"
# VERSION="11 (bullseye)"
# VERSION_CODENAME=bullseye
# ID=debian
# HOME_URL="https://www.debian.org/"
# SUPPORT_URL="https://www.debian.org/support"
# BUG_REPORT_URL="https://bugs.debian.org/"
