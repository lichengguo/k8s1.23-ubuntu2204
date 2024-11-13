# 使用定制的php镜像
FROM harbor.alnk.com/public/php8-swoole-yaf-redis-amqp:latest

# 安装 Nginx
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \
    apk --no-cache add nginx tzdata ca-certificates && \
    cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    mkdir -p /run/nginx/ && \
    rm -rf /var/cache/apk/*

# 复制 Nginx 配置文件
COPY hello.conf /etc/nginx/conf.d/hello.conf

# 复制 PHP 代码到工作目录
COPY src/ /var/www/html

# 设置工作目录
WORKDIR /var/www/html

# Expose 80 port
EXPOSE 80

# debian镜像启动 Nginx 和 PHP-FPM7.4
# CMD service nginx start && php-fpm7.4

# alpine镜像启动nginx和php-fpm8.0
CMD nginx  && php-fpm

## 构建镜像
# docker build -t my-php-nginx-app .
## 启动容器
# docker run -d -p 8888:80 --name my-php-nginx-container my-php-nginx-app