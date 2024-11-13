### 基于go语言的cicd流水线

> Go程序需要连接Nacos获取配置，然后发布到k8s中
>
> ```text
> 前置条件
> 1.k8s集群
> 2.harbor仓库
> 3.jenkins
> 4.gitlab
> 5.docker
> 6.mysql
> ```



#### docker部署nacos:v2.2.0

###### 部署nacos

```shell
【10.0.1.21】
## 镜像准备
# docker pull swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/nacos/nacos-server:v2.2.0
# docker tag  swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/nacos/nacos-server:v2.2.0   harbor.alnk.com/public/nacos-server:v2.2.0

# docker push harbor.alnk.com/public/nacos-server:v2.2.0

## 目录准备
# mkdir -p /data/nacos/logs
# mkdir -p /data/nacos/init.d
# mkdir -p /data/nacos/data

## nacos依赖的数据
# cd /data/nacos/
# vi nacos.sql

## 安装数据库
#  apt install mariadb-server -y
# 数据库初始化设置
# mysql_secure_installation 


## 设置远程登录
# vi /etc/mysql/mariadb.conf.d/50-server.cnf
bind-address            = 0.0.0.0
# systemctl restart mariadb.service

# mysql -uroot -proot123
## 赋予root用户远程登录的权限
> GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY 'root123' WITH GRANT OPTION;
## 刷新权限使更改生效
> FLUSH PRIVILEGES;

##  导入数据
# mysql -uroot -proot123
> CREATE DATABASE nacos_config;
> use nacos_config;
> source nacos.sql


## docker启动nacos
# docker启动nacos
docker run -d \
--name nacos \
-p 8848:8848 \
-p 9848:9848 \
-p 9849:9849 \
--privileged=true \
--restart=always \
-e JVM_XMS=256m \
-e JVM_XMX=256m \
-e MODE=standalone \
-e PREFER_HOST_MODE=hostname \
-e SPRING_DATASOURCE_PLATFORM=mysql \
-e MYSQL_SERVICE_HOST=10.0.1.21 \
-e MYSQL_SERVICE_PORT=3306 \
-e MYSQL_SERVICE_DB_NAME=nacos_config \
-e MYSQL_SERVICE_USER=root \
-e MYSQL_SERVICE_PASSWORD=root123 \
-v /data/nacos/logs:/home/nacos/logs \
-v /data/nacos/init.d/custom.properties:/etc/nacos/init.d/custom.properties \
-v /data/nacos/data:/home/nacos/data \
harbor.alnk.com/public/nacos-server:v2.2.0

# 参数解释
docker run：运行 Docker 容器的命令。
-d：将容器放入后台并运行。
--name nacos：将容器命名为“nacos”。
-p 8848:8848 -p 9848:9848 -p 9849:9849：将主机的端口映射到容器中。主机上的端口 8848、9848 和 9849 分别映射到容器中相同的端口。
--privileged=true：赋予容器扩展权限。
--restart=always：如果容器停止运行，自动重新启动容器。
-e：在容器内设置环境变量。
JVM_XMS=256m：设置 JVM 的初始堆大小为 256 MB。
JVM_XMX=256m：设置 JVM 的最大堆大小为 256 MB。
MODE=standalone：将 Nacos 的模式设置为独立模式。
PREFER_HOST_MODE=hostname：设置主机模式为主机名。
SPRING_DATASOURCE_PLATFORM=mysql：设置数据源平台为MySQL。
MYSQL_SERVICE_HOST=10.4.7.11：设置 MySQL 服务的主机地址。
MYSQL_SERVICE_PORT=3306：设置 MySQL 服务的端口。
MYSQL_SERVICE_DB_NAME=nacos_config：设置 Nacos 使用的 MySQL 数据库名称。
MYSQL_SERVICE_USER=root：设置 MySQL 的用户名。
MYSQL_SERVICE_PASSWORD=123456：设置 MySQL 的密码。
-v：挂载主机上的卷到容器中。
/data/nacos/logs:/home/nacos/logs：将主机上 Nacos 日志目录挂载到容器中。
/data/nacos/init.d/custom.properties:/etc/nacos/init.d/custom.properties：将主机上的 custom.properties 文件挂载到容器中。
/data/nacos/data:/home/nacos/data：将主机上 Nacos 数据目录挂载到容器中。
harbor.od.com/public/nacos-server:v2.2.0：指定容器使用的 Docker 镜像，这里是 Nacos 服务器镜像。

注意：
-e JVM_XMS=256m
-e JVM_XMX=256m 可以不设置，但nacos默认值会占用1G左右内存，内存不够用的同学最好设置一下
MYSQL_SERVICE_DB_NAME=nacos_config：必须与之前创建的数据库同名。
```



###### nacos数据库文件

`nacos.sql`

```sql
/*
 * Copyright 1999-2018 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/******************************************/
/*   表名称 = config_info                  */
/******************************************/
CREATE TABLE `config_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `data_id` varchar(255) NOT NULL COMMENT 'data_id',
  `group_id` varchar(128) DEFAULT NULL COMMENT 'group_id',
  `content` longtext NOT NULL COMMENT 'content',
  `md5` varchar(32) DEFAULT NULL COMMENT 'md5',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  `src_user` text COMMENT 'source user',
  `src_ip` varchar(50) DEFAULT NULL COMMENT 'source ip',
  `app_name` varchar(128) DEFAULT NULL COMMENT 'app_name',
  `tenant_id` varchar(128) DEFAULT '' COMMENT '租户字段',
  `c_desc` varchar(256) DEFAULT NULL COMMENT 'configuration description',
  `c_use` varchar(64) DEFAULT NULL COMMENT 'configuration usage',
  `effect` varchar(64) DEFAULT NULL COMMENT '配置生效的描述',
  `type` varchar(64) DEFAULT NULL COMMENT '配置的类型',
  `c_schema` text COMMENT '配置的模式',
  `encrypted_data_key` text NOT NULL COMMENT '密钥',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_configinfo_datagrouptenant` (`data_id`,`group_id`,`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='config_info';

/******************************************/
/*   表名称 = config_info_aggr             */
/******************************************/
CREATE TABLE `config_info_aggr` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `data_id` varchar(255) NOT NULL COMMENT 'data_id',
  `group_id` varchar(128) NOT NULL COMMENT 'group_id',
  `datum_id` varchar(255) NOT NULL COMMENT 'datum_id',
  `content` longtext NOT NULL COMMENT '内容',
  `gmt_modified` datetime NOT NULL COMMENT '修改时间',
  `app_name` varchar(128) DEFAULT NULL COMMENT 'app_name',
  `tenant_id` varchar(128) DEFAULT '' COMMENT '租户字段',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_configinfoaggr_datagrouptenantdatum` (`data_id`,`group_id`,`tenant_id`,`datum_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='增加租户字段';


/******************************************/
/*   表名称 = config_info_beta             */
/******************************************/
CREATE TABLE `config_info_beta` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `data_id` varchar(255) NOT NULL COMMENT 'data_id',
  `group_id` varchar(128) NOT NULL COMMENT 'group_id',
  `app_name` varchar(128) DEFAULT NULL COMMENT 'app_name',
  `content` longtext NOT NULL COMMENT 'content',
  `beta_ips` varchar(1024) DEFAULT NULL COMMENT 'betaIps',
  `md5` varchar(32) DEFAULT NULL COMMENT 'md5',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  `src_user` text COMMENT 'source user',
  `src_ip` varchar(50) DEFAULT NULL COMMENT 'source ip',
  `tenant_id` varchar(128) DEFAULT '' COMMENT '租户字段',
  `encrypted_data_key` text NOT NULL COMMENT '密钥',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_configinfobeta_datagrouptenant` (`data_id`,`group_id`,`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='config_info_beta';

/******************************************/
/*   表名称 = config_info_tag              */
/******************************************/
CREATE TABLE `config_info_tag` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `data_id` varchar(255) NOT NULL COMMENT 'data_id',
  `group_id` varchar(128) NOT NULL COMMENT 'group_id',
  `tenant_id` varchar(128) DEFAULT '' COMMENT 'tenant_id',
  `tag_id` varchar(128) NOT NULL COMMENT 'tag_id',
  `app_name` varchar(128) DEFAULT NULL COMMENT 'app_name',
  `content` longtext NOT NULL COMMENT 'content',
  `md5` varchar(32) DEFAULT NULL COMMENT 'md5',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  `src_user` text COMMENT 'source user',
  `src_ip` varchar(50) DEFAULT NULL COMMENT 'source ip',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_configinfotag_datagrouptenanttag` (`data_id`,`group_id`,`tenant_id`,`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='config_info_tag';

/******************************************/
/*   表名称 = config_tags_relation         */
/******************************************/
CREATE TABLE `config_tags_relation` (
  `id` bigint(20) NOT NULL COMMENT 'id',
  `tag_name` varchar(128) NOT NULL COMMENT 'tag_name',
  `tag_type` varchar(64) DEFAULT NULL COMMENT 'tag_type',
  `data_id` varchar(255) NOT NULL COMMENT 'data_id',
  `group_id` varchar(128) NOT NULL COMMENT 'group_id',
  `tenant_id` varchar(128) DEFAULT '' COMMENT 'tenant_id',
  `nid` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'nid, 自增长标识',
  PRIMARY KEY (`nid`),
  UNIQUE KEY `uk_configtagrelation_configidtag` (`id`,`tag_name`,`tag_type`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='config_tag_relation';

/******************************************/
/*   表名称 = group_capacity               */
/******************************************/
CREATE TABLE `group_capacity` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `group_id` varchar(128) NOT NULL DEFAULT '' COMMENT 'Group ID，空字符表示整个集群',
  `quota` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '配额，0表示使用默认值',
  `usage` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '使用量',
  `max_size` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '单个配置大小上限，单位为字节，0表示使用默认值',
  `max_aggr_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '聚合子配置最大个数，，0表示使用默认值',
  `max_aggr_size` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '单个聚合数据的子配置大小上限，单位为字节，0表示使用默认值',
  `max_history_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '最大变更历史数量',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_group_id` (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='集群、各Group容量信息表';

/******************************************/
/*   表名称 = his_config_info              */
/******************************************/
CREATE TABLE `his_config_info` (
  `id` bigint(20) unsigned NOT NULL COMMENT 'id',
  `nid` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'nid, 自增标识',
  `data_id` varchar(255) NOT NULL COMMENT 'data_id',
  `group_id` varchar(128) NOT NULL COMMENT 'group_id',
  `app_name` varchar(128) DEFAULT NULL COMMENT 'app_name',
  `content` longtext NOT NULL COMMENT 'content',
  `md5` varchar(32) DEFAULT NULL COMMENT 'md5',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  `src_user` text COMMENT 'source user',
  `src_ip` varchar(50) DEFAULT NULL COMMENT 'source ip',
  `op_type` char(10) DEFAULT NULL COMMENT 'operation type',
  `tenant_id` varchar(128) DEFAULT '' COMMENT '租户字段',
  `encrypted_data_key` text NOT NULL COMMENT '密钥',
  PRIMARY KEY (`nid`),
  KEY `idx_gmt_create` (`gmt_create`),
  KEY `idx_gmt_modified` (`gmt_modified`),
  KEY `idx_did` (`data_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='多租户改造';


/******************************************/
/*   表名称 = tenant_capacity              */
/******************************************/
CREATE TABLE `tenant_capacity` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `tenant_id` varchar(128) NOT NULL DEFAULT '' COMMENT 'Tenant ID',
  `quota` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '配额，0表示使用默认值',
  `usage` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '使用量',
  `max_size` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '单个配置大小上限，单位为字节，0表示使用默认值',
  `max_aggr_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '聚合子配置最大个数',
  `max_aggr_size` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '单个聚合数据的子配置大小上限，单位为字节，0表示使用默认值',
  `max_history_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '最大变更历史数量',
  `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='租户容量信息表';


CREATE TABLE `tenant_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `kp` varchar(128) NOT NULL COMMENT 'kp',
  `tenant_id` varchar(128) default '' COMMENT 'tenant_id',
  `tenant_name` varchar(128) default '' COMMENT 'tenant_name',
  `tenant_desc` varchar(256) DEFAULT NULL COMMENT 'tenant_desc',
  `create_source` varchar(32) DEFAULT NULL COMMENT 'create_source',
  `gmt_create` bigint(20) NOT NULL COMMENT '创建时间',
  `gmt_modified` bigint(20) NOT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tenant_info_kptenantid` (`kp`,`tenant_id`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='tenant_info';

CREATE TABLE `users` (
	`username` varchar(50) NOT NULL PRIMARY KEY COMMENT 'username',
	`password` varchar(500) NOT NULL COMMENT 'password',
	`enabled` boolean NOT NULL COMMENT 'enabled'
);

CREATE TABLE `roles` (
	`username` varchar(50) NOT NULL COMMENT 'username',
	`role` varchar(50) NOT NULL COMMENT 'role',
	UNIQUE INDEX `idx_user_role` (`username` ASC, `role` ASC) USING BTREE
);

CREATE TABLE `permissions` (
    `role` varchar(50) NOT NULL COMMENT 'role',
    `resource` varchar(128) NOT NULL COMMENT 'resource',
    `action` varchar(8) NOT NULL COMMENT 'action',
    UNIQUE INDEX `uk_role_permission` (`role`,`resource`,`action`) USING BTREE
);

INSERT INTO users (username, password, enabled) VALUES ('nacos', '$2a$10$EuWPZHzz32dJN7jexM34MOeYirDdFAZm2kuWj7VEOJhhZkDrxfvUu', TRUE);

INSERT INTO roles (username, role) VALUES ('nacos', 'ROLE_ADMIN');
```



###### 访问设置nacos

`访问nacos`

```
#浏览器访问：http://10.0.1.21:8848/nacos/#/login
账号：nacos
密码：nacos
```

![1731080941941](images/1731080941941.png)  

![1731081419117](images/1731081419117.png)  



![1731081613134](images/1731081613134.png)  



![1731081678221](images/1731081678221.png)  

![1731081700254](images/1731081700254.png)  





#### gitlab上传代码

###### 上传go项目代码

`main.go`

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-k8s-one/nacoscfg"
	"log"
	"net/http"
	"os"
)

// NaConf nacos上的配置
type NaConf struct {
	Port int64  `json:"port"`
	Mes  string `json:"mes"`
}

func main() {
	if err := os.MkdirAll("./logs", 0777); err != nil {
		log.Fatal(err)
	}

	logFile, err := os.OpenFile("./logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

	r := gin.Default()

	var n NaConf

	err = json.Unmarshal(nacoscfg.LoadNacos(), &n)
	if err != nil {
		panic("解析失败")
	}

	r.GET("/", func(c *gin.Context) {
		// 增加一行日志
		logger.Println("path: ", c.FullPath(), ",resData: ", n)

		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"mes":  "hello world! 开发了很多新功能 version:v0.0.2",
			"data": n,
		})
	})

	_ = r.Run(fmt.Sprintf(":%d", n.Port))
}

```



`nacoscfg/nacos.go`

```go
package nacoscfg

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"go-k8s-one/conf"
)

func LoadNacos() []byte {
	//create ServerConfig
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(conf.IP, conf.Port, constant.WithContextPath("/nacos")),
	}

	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(conf.NameSpaceID),
		//constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		//constant.WithLogDir("/tmp/nacos/log"),
		//constant.WithCacheDir("/tmp/nacos/cache"),
		//constant.WithLogLevel("error"),
	)

	// create config client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	//get config
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: conf.DataID,
		Group:  conf.Group,
	})
	if err != nil {
		panic(err)
	}

	return []byte(content)
}

```



`conf/config.go`

```go
package conf

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"path/filepath"
)

var (
	IP          string
	Port        uint64
	NameSpaceID string
	DataID      string
	Group       string
)

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Print("获取当前程序执行目录失败")
		return
	}

	file, err := ini.Load(fmt.Sprintf("%s/%s", dir, "conf/config.ini"))
	if err != nil {
		fmt.Println("配置文件读取有误,请检查配置文件.", err)
		return
	}

	LoadNacos(file)
}

func LoadNacos(file *ini.File) {
	s := file.Section("nacos")
	IP = s.Key("IP").String()
	Port = s.Key("Port").MustUint64(8848)
	NameSpaceID = s.Key("NameSpaceID").String()
	DataID = s.Key("DataID").String()
	Group = s.Key("Group").String()
}

```



`conf/config.ini`

```ini
[nacos]
IP = 10.0.1.21
Port = 8848
# test名称空间
NameSpaceID = c19d7abe-0683-4168-af1d-c30f1183854f
DataID = go-k8s-one
Group = DEFAULT_GROUP
```



`Dockerfile`

```Dockerfile
# stage 1: build src code to binary
FROM harbor.alnk.com/public/golang:1.22.8 as builder

ENV GOPROXY https://goproxy.cn

COPY . /app/

RUN cd  /app \
    && go mod init go-k8s-one \
    && go mod tidy \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o go-k8s-one .

# stage 2: use alpine as base image
FROM harbor.od.com/public/alpine:3.18

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \
    apk --no-cache add tzdata ca-certificates && \
    cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    # apk del tzdata && \
    rm -rf /var/cache/apk/*

COPY --from=builder /app/go-k8s-one /app/go-k8s-one
COPY --from=builder /app/conf/config.ini /app/conf/config.ini

CMD ["/app/go-k8s-one"]


#创建golang目录
#cd golang
#docker build -t harbor.od.com/app/go-k8s-one:v0.0.1 .
#docker run -d -p9999:9999 harbor.od.com/app/go-k8s-one:v0.0.1
```



`新建项目`

![1731082143074](images/1731082143074.png)  

`创建test分支`

![1731082226406](images/1731082226406.png)  



`创建目录`

![1731082273663](images/1731082273663.png)  



 ![1731082307507](images/1731082307507.png) 



`把文件上传到对应的目录中`

![1731082466485](images/1731082466485.png) 





#### jenkins配置

###### jenkins插件设置国内源

`将jenkins插件安装配置国内源`

> https://mirrors.tuna.tsinghua.edu.cn/jenkins/updates/update-center.json

![1731153786118](images/1731153786118.png)



###### 新建jenkins流水线任务

![1731082525727](images/1731082525727.png)  

![1731083603900](images/1731083603900.png)  

![1731083651543](images/1731083651543.png)  

![1731083764308](images/1731083764308.png)  

![1731083839547](images/1731083839547.png)  

![1731083951259](images/1731083951259.png)  



###### pipeline流水线    

```pipeline
// pipeline流水线
pipeline {
    agent any

    environment {
        // 获取系统当前时间作为镜像的版本号
        BUILD_TIME = new Date().format('yyyyMMdd_HHmmss')

        // harbor仓库地址
        HARBOR_URL = "harbor.alnk.com"
    }

    // 为了让远程SSH也获取到变量
    parameters {
        // 项目名称，deployment名称
        string(name: 'app_name', defaultValue: 'go-k8s-one')
        // 镜像名称
        string(name: 'image_name', defaultValue: 'public/go-k8s-one')
        // git仓库地址
        string(name: 'git_repo', defaultValue: 'git@10.0.1.21:demo/go-k8s-one.git')
        // git仓库分支
        string(name: 'git_ver', defaultValue: 'test')
    }

    stages {     
        stage('pull') { 
            steps {
                sh "git clone ${git_repo} ${app_name}/${env.BUILD_NUMBER} && cd ${app_name}/${env.BUILD_NUMBER} && git checkout ${git_ver}"
            }
        }
        
        stage('build') { 
            steps {
                sh "cd ${app_name}/${env.BUILD_NUMBER}  && docker build -t ${HARBOR_URL}/${image_name}:${BUILD_TIME} ."
            }
        }
        
        stage('push image') { 
            steps {
                sh "docker push ${HARBOR_URL}/${image_name}:${BUILD_TIME}"
                sh "docker rmi ${HARBOR_URL}/${image_name}:${BUILD_TIME}"
            }
        }           
    }
}
```



###### jenkins构建镜像推送到harbor仓库

![1731084029322](images/1731084029322.png)  



![1731084083821](images/1731084083821.png)  

![1731084106231](images/1731084106231.png)  

![1731084151476](images/1731084151476.png)  

![1731086448713](images/1731086448713.png)  

![1731086491776](images/1731086491776.png)  



#### 发布服务到k8s测试环境test

##### 编辑yaml文件，发布服务

```shell
【10.0.1.21】
# mkdir -p /data/k8s-yaml/go-k8s-one
# cd /data/k8s-yaml/go-k8s-one
# vi go-k8s-one.yaml

【10.0.1.201】
# kubectl create ns test
# kubectl apply -f http://k8s-yaml.alnk.com/go-k8s-one/go-k8s-one.yaml
```

`go-k8s-one.yaml`

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: go-k8s-one-cm
  namespace: test
data:
  config.ini: |
    [nacos]
    IP = 10.0.1.21
    Port = 8848
    # test名称空间
    NameSpaceID = c19d7abe-0683-4168-af1d-c30f1183854f
    DataID = go-k8s-one
    Group = DEFAULT_GROUP
---
kind: Deployment
apiVersion: apps/v1
metadata:
  name: go-k8s-one
  namespace: test
  labels:
    name: go-k8s-one
spec:
  replicas: 1
  selector:
    matchLabels:
      name: go-k8s-one
  template:
    metadata:
      labels:
        app: go-k8s-one
        name: go-k8s-one
    spec:
      volumes:
      - name: configmap-volume
        configMap:
          name: go-k8s-one-cm
      containers:
      - name: go-k8s-one
        image: harbor.alnk.com/public/go-k8s-one:20241109_011408
        ports:
        - containerPort: 9999
          protocol: TCP
        volumeMounts:
        - name: configmap-volume
          mountPath: app/conf
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        imagePullPolicy: IfNotPresent
      imagePullSecrets:
      - name: harbor
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      securityContext:
        runAsUser: 0
      schedulerName: default-scheduler
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 1
      maxSurge: 1
  revisionHistoryLimit: 7
  progressDeadlineSeconds: 600
  
---
kind: Service
apiVersion: v1
metadata:
  name: go-k8s-one
  namespace: test
spec:
  ports:
  - protocol: TCP
    port: 9999
    targetPort: 9999
  selector:
    app: go-k8s-one

---         
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  namespace: test
  name: go-k8s-one
spec:
  rules:
  - host: go-k8s-one-test.alnk.com
    http:
      paths:
      - backend:
          service:
            name: go-k8s-one
            port:
              number: 9999
        path: /
        pathType: Prefix
```

![1731139328854](images/1731139328854.png)  





##### 配置jenkins远程ssh登录linux服务器

`安装插件3个插件` 

> Publish Over SSH
>
> SSH Agent Plugin
>
> SSH Pipeline Steps

![1731143629417](images/1731143629417.png)    



`在需要被远程执行的linux服务器上执行`

```shell
【10.0.1.201】
## 生成密钥
# ssh-keygen
## 添加公钥
# cat id_rsa.pub >> authorized_keys

## 查看密钥
# cat id_rsa
-----BEGIN OPENSSH PRIVATE KEY-----
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
-----END OPENSSH PRIVATE KEY-----

```



 `创建凭证`

![1731143909832](images/1731143909832.png)  



##### 修改上面的jenkins流水线，实现构建与发布一步完成，发布到测试环境

```pipeline
pipeline {
    agent any

    environment {
        // 获取系统当前时间作为镜像的版本号
        BUILD_TIME = new Date().format('yyyyMMdd_HHmmss')

        // harbor仓库地址
        HARBOR_URL = "harbor.alnk.com"
    }

    // 为了让远程SSH也获取到变量
    parameters {
        // 项目名称，deployment名称
        string(name: 'app_name', defaultValue: 'go-k8s-one')
        // 镜像名称
        string(name: 'image_name', defaultValue: 'public/go-k8s-one')
        // git仓库地址
        string(name: 'git_repo', defaultValue: 'git@10.0.1.21:demo/go-k8s-one.git')
        // git仓库分支
        string(name: 'git_ver', defaultValue: 'test')
    }

    stages {     
        stage('pull') { 
            steps {
                sh "git clone ${git_repo} ${app_name}/${env.BUILD_NUMBER} && cd ${app_name}/${env.BUILD_NUMBER} && git checkout ${git_ver}"
            }
        }
        
        stage('build') { 
            steps {
                sh "cd ${app_name}/${env.BUILD_NUMBER}  && docker build -t ${HARBOR_URL}/${image_name}:${BUILD_TIME} ."
            }
        }
        
        stage('push image') { 
            steps {
                sh "docker push ${HARBOR_URL}/${image_name}:${BUILD_TIME}"
                sh "docker rmi ${HARBOR_URL}/${image_name}:${BUILD_TIME}"
            }
        }
            
        stage('deploy') { 
            steps {
                sshagent(credentials: ['10.0.1.201']) {					
                sh 'ssh root@10.0.1.201 kubectl -n test set image deployment/${app_name} ${app_name}=${HARBOR_URL}/${image_name}:${BUILD_TIME} --record '					
                }
            }
        }
    }
}
```

`再次构建`

![1731150405589](images/1731150405589.png)  



![1731150431734](images/1731150431734.png)  



##### gitlab提交代码触发jenkins

##### gitlab插件安装

`jenkins安装GitLab插件`

![1731152828884](images/1731152828884.png)  



##### jenkins修改构建器

`jenkins修改构建触发器`

> http://jenkins.alnk.com/project/go-k8s-one

![1731154093132](images/1731154093132.png)  

`选择test分支，只允许test分支进行构建`

![1731154203804](images/1731154203804.png)  

`允许访问/project`

![1731155267310](images/1731155267310.png)  



##### gitlab修改webhooks

`gitlab上修改webhooks`

![1731154875298](images/1731154875298.png)  

![1731154314493](images/1731154314493.png)  

```shell
【10.0.1.21】
## 添加解析
# cat /etc/hosts
10.0.1.100 jenkins.alnk.com

```

`gitlab测试`

![1731155339092](images/1731155339092.png)  





##### gitlab测试提交test分支代码

`gitlab上test分支提交代码`

![1731155427126](images/1731155427126.png)  



`可以看到已经触发`

![1731155465474](images/1731155465474.png)  





![1731155599502](images/1731155599502.png)  



#### 测试发布服务到k8s生产环

##### jenkins流水线复制与修改

> 复制一条流水线，然后修改

![1731156332531](images/1731156332531.png)  

`修改分支版本默认值`

![1731156481449](images/1731156481449.png)  

![1731156501752](images/1731156501752.png)  



`修改pipeline发布k8s应用时的名称空间，这里因为只有一个k8s集群，所以测试和生产是分名称空间区别的，正常情况下测试和生产是有两个集群的`

![1731156731240](images/1731156731240.png) 



`pipeline`

```
pipeline {
    agent any

    environment {
        // 获取系统当前时间作为镜像的版本号
        BUILD_TIME = new Date().format('yyyyMMdd_HHmmss')

        // harbor仓库地址
        HARBOR_URL = "harbor.alnk.com"
    }

    // 为了让远程SSH也获取到变量
    parameters {
        // 项目名称，deployment名称
        string(name: 'app_name', defaultValue: 'go-k8s-one')
        // 镜像名称
        string(name: 'image_name', defaultValue: 'public/go-k8s-one')
        // git仓库地址
        string(name: 'git_repo', defaultValue: 'git@10.0.1.21:demo/go-k8s-one.git')
        // git仓库分支
        string(name: 'git_ver', defaultValue: 'main')
    }

    stages {     
        stage('pull') { 
            steps {
                sh "git clone ${git_repo} ${app_name}/${env.BUILD_NUMBER} && cd ${app_name}/${env.BUILD_NUMBER} && git checkout ${git_ver}"
            }
        }
        
        stage('build') { 
            steps {
                sh "cd ${app_name}/${env.BUILD_NUMBER}  && docker build -t ${HARBOR_URL}/${image_name}:${BUILD_TIME} ."
            }
        }
        
        stage('push image') { 
            steps {
                sh "docker push ${HARBOR_URL}/${image_name}:${BUILD_TIME}"
                sh "docker rmi ${HARBOR_URL}/${image_name}:${BUILD_TIME}"
            }
        }
            
        stage('deploy') { 
            steps {
                sshagent(credentials: ['10.0.1.201']) {					
                sh 'ssh root@10.0.1.201 kubectl -n prod set image deployment/${app_name} ${app_name}=${HARBOR_URL}/${image_name}:${BUILD_TIME} --record '					
                }
            }
        }
    }
}
```



`记录webhook地址，等下gitlab需要用到`

> http://jenkins.alnk.com/project/go-k8s-one-prod

![1731156571077](images/1731156571077.png)  



`修改触发的分支为main`

![1731156603617](images/1731156603617.png)    





##### gitlab合并代码，修改配置文件和main文件

`修改nacos的名称空间为prod的ID`

![1731156921710](images/1731156921710.png) 

![1731157123471](images/1731157123471.png)  



##### jenkins构建

`构建一次，这次构建会成功，但是发布到k8s会失败，只是为了打个镜像`

```shell
# ssh root@10.0.1.201 kubectl -n prod set image deployment/go-k8s-one go-k8s-one=harbor.alnk.com/public/go-k8s-one:20241109_205326 --record
Flag --record has been deprecated, --record will be removed in the future
Error from server (NotFound): namespaces "prod" not found
```



##### 第一次服务发布到k8s（第一次发布需要手动发布应用）

```shell
【10.0.1.21】
# mkdir -p /data/k8s-yaml/go-k8s-one
# cd /data/k8s-yaml/go-k8s-one
# vi go-k8s-one-prod.yaml

【10.0.1.201】
# kubectl create ns prod
# kubectl apply -f http://k8s-yaml.alnk.com/go-k8s-one/go-k8s-one-prod.yaml
```



`go-k8s-one-prod`

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: go-k8s-one-cm
  namespace: prod # 修改为prod名称空间
data:
  config.ini: |
    [nacos]
    IP = 10.0.1.21
    Port = 8848
    # 修改为nacos的prod名称空间ID
    NameSpaceID = eab40093-96bd-48ce-b88e-e13b268efd45
    DataID = go-k8s-one
    Group = DEFAULT_GROUP
---
kind: Deployment
apiVersion: apps/v1
metadata:
  name: go-k8s-one
  namespace: prod # 修改为prod名称空间
  labels:
    name: go-k8s-one
spec:
  replicas: 1
  selector:
    matchLabels:
      name: go-k8s-one
  template:
    metadata:
      labels:
        app: go-k8s-one
        name: go-k8s-one
    spec:
      volumes:
      - name: configmap-volume
        configMap:
          name: go-k8s-one-cm
      containers:
      - name: go-k8s-one
        image: harbor.alnk.com/public/go-k8s-one:20241109_205326 # 修改为上面的镜像
        ports:
        - containerPort: 8888 # 修改服务端口
          protocol: TCP
        volumeMounts:
        - name: configmap-volume
          mountPath: app/conf
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        imagePullPolicy: IfNotPresent
      imagePullSecrets:
      - name: harbor
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      securityContext:
        runAsUser: 0
      schedulerName: default-scheduler
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 1
      maxSurge: 1
  revisionHistoryLimit: 7
  progressDeadlineSeconds: 600
  
---
kind: Service
apiVersion: v1
metadata:
  name: go-k8s-one
  namespace: prod # 修改为prod名称空间
spec:
  ports:
  - protocol: TCP
    port: 8888 # 修改service的端口为8888
    targetPort: 8888 # 修改容器的端口为8888
  selector:
    app: go-k8s-one

---         
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  namespace: prod # 修改为prod名称空间
  name: go-k8s-one
spec:
  rules:
  - host: go-k8s-one-prod.alnk.com # 修改域名
    http:
      paths:
      - backend:
          service:
            name: go-k8s-one
            port:
              number: 8888 # 修改service的端口为888
        path: /
        pathType: Prefix
```

![1731157607741](images/1731157607741.png)  





##### gitlab配置webhooks

![1731157774449](images/1731157774449.png)  



`测试`

![1731157815682](images/1731157815682.png)   





##### 模拟main主分支提交代码

![1731157869361](images/1731157869361.png)  



`jenkins已经开始自动构建，main主分支这条流水开始自动构建了，test分支流水线没有构建`

![1731157965159](images/1731157965159.png)



`可以看到成功发布`

![1731158128553](images/1731158128553.png)  

