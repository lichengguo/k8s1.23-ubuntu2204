# 使用php8.0-fpm-alpine作为基础镜像
FROM harbor.alnk.com/public/php:8.0.1-fpm-alpine

# 签名
LABEL maintainer="Alnk 1029612787@qq.com"

# 设置apk源为国内镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update 
COPY ./conf.d/ $PHP_INI_DIR/conf.d/

# 时区
ENV TZ "Asia/Shanghai"
ENV TERM xterm

# 默认关闭opcode
ENV OPCODE 0

# 添加用户
RUN addgroup -g 1000 -S www && adduser -s /sbin/nologin -S -D -u 1000 -G www www

# 指定扩展版本
ENV PHPREDIS_VERSION 5.3.4
ENV PHPXDEBUG_VERSION 2.6.0
ENV PHPZIP_VERSION 1.19.3
ENV PHPSWOOLE_VERSION 4.6.7
ENV PHPYAF_VERSION 3.3.2
ENV PHPAMQP_VERSION 1.11.0beta
ENV PHPMONGODB_VERSION 1.9.1

# PHPIZE_DEPS 包含 gcc g++ 等编译辅助类库，完成编译后删除
RUN apk add --no-cache $PHPIZE_DEPS \
    && apk add --no-cache libstdc++ libzip-dev vim\
    && apk add rabbitmq-c-dev \
    && apk update \
    && pecl install redis-$PHPREDIS_VERSION \
    && pecl install zip-$PHPZIP_VERSION \
    && pecl install yaf-$PHPYAF_VERSION \
    && pecl install swoole-$PHPSWOOLE_VERSION \
    && pecl install mongodb-$PHPMONGODB_VERSION \
    && pecl install amqp-$PHPAMQP_VERSION \
    && docker-php-ext-enable redis zip yaf swoole mongodb amqp \
    && apk del $PHPIZE_DEPS

# 安装扩展
RUN apk add --no-cache \
       freetype \
       libpng \
       libjpeg-turbo \
       freetype-dev \
       libpng-dev \
       jpeg-dev \
       libjpeg \
       libjpeg-turbo-dev \
       libwebp \
       libwebp-dev \
    && NUMPROC=$(grep -c ^processor /proc/cpuinfo 2>/dev/null || 1) \
    && docker-php-ext-install -j${NUMPROC} gd \
    && docker-php-ext-install -j${NUMPROC} pdo_mysql \
    && docker-php-ext-install -j${NUMPROC} opcache \
    && docker-php-ext-install -j${NUMPROC} bcmath \ 
    && docker-php-ext-install -j${NUMPROC} mysqli

# 拷贝配置文件
COPY  www.conf /usr/local/etc/php-fpm.d/www.conf

RUN mv "$PHP_INI_DIR/php.ini-production" "$PHP_INI_DIR/php.ini"