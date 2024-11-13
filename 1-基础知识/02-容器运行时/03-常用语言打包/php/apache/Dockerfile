# 使用PHP 8.0的官方Apache镜像作为基础镜像
FROM registry.cn-hangzhou.aliyuncs.com/alnktest/php:8.0-apache
 
# 将应用程序文件复制到镜像内的Apache服务目录
COPY . /var/www/html/
 
# # 安装软件
# RUN sed -i 's/deb.debian.org/ftp.cn.debian.org/g' /etc/apt/sources.list &&\
#     sed -i 's/security.debian.org/ftp.cn.debian.org/g' /etc/apt/sources.list &&\
#     apt-get update && apt-get install -y \
#     libzip-dev \
#     libpng-dev \
#     libjpeg62-turbo-dev \
#     libfreetype6-dev \
#     libonig-dev \
#     libwebp-dev \
#     libxml2-dev \
#     libxslt-dev \
#     vim \
#     net-tools \
#     iproute2 \
#     && docker-php-ext-configure gd --with-jpeg \
#     && docker-php-ext-install -j$(nproc) gd pdo_mysql zip bcmath exif mysqli pcntl \
#     && pecl install redis \
#     && docker-php-ext-enable redis \
#     && apt-get install -y nginx && rm -rf /var/lib/apt/lists/*

# # （可选）配置Apache，例如修改DocumentRoot
# RUN a2enmod rewrite && \
#     sed -ri -e 's!/var/www/html!/var/www/html/public!g' /etc/apache2/sites-available/* && \
#     sed -ri -e 's!/var/www/html/!/var/www/html/public!g' /etc/apache2/sites-enabled/*

## 构建镜像
# docker build -t my-php-nginx-app .
## 启动容器
# docker run -d -p 8888:80 --name my-php-nginx-container my-php-nginx-app