FROM node:12.6.0-alpine
MAINTAINER 1029612787@qq.com

WORKDIR /app
COPY package.json .

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \
    yarn  config set registry https://registry.npm.taobao.org && \
    yarn install

RUN yarn build

COPY . .

EXPOSE 6868

ENTRYPOINT ["yarn", "start"]