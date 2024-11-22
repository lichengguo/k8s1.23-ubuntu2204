### 基于go语言的cicd流水线

> Go程序需要连接Nacos获取配置，然后在把服务发布到k8s集群中
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

##### 部署nacos

> ```
> 【10.0.1.21】
> ## 镜像准备
> # docker pull swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/nacos/nacos-server:v2.2.0
> # docker tag  swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/nacos/nacos-server:v2.2.0   harbor.alnk.com/public/nacos-server:v2.2.0
> # docker push harbor.alnk.com/public/nacos-server:v2.2.0
> 
> ## 目录准备
> # mkdir -p /data/nacos/logs
> # mkdir -p /data/nacos/init.d
> # mkdir -p /data/nacos/data
> 
> ## nacos依赖的数据
> # cd /data/nacos/
> # vi nacos.sql
> 
> ## 安装数据库
> # apt install mariadb-server -y
> # 数据库初始化设置
> # mysql_secure_installation 
> 
> 
> ## 设置远程登录
> # vi /etc/mysql/mariadb.conf.d/50-server.cnf
> bind-address            = 0.0.0.0
> # systemctl restart mariadb.service
> 
> # mysql -uroot -proot123
> ## 赋予root用户远程登录的权限
> > GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY 'root123' WITH GRANT OPTION;
> ## 刷新权限使更改生效
> > FLUSH PRIVILEGES;
> 
> ##  导入数据
> # mysql -uroot -proot123
> > CREATE DATABASE nacos_config;
> > use nacos_config;
> > source nacos.sql
> 
> 
> ## docker启动nacos
> # docker启动nacos
> docker run -d \
> --name nacos \
> -p 8848:8848 \
> -p 9848:9848 \
> -p 9849:9849 \
> --privileged=true \
> --restart=always \
> -e JVM_XMS=256m \
> -e JVM_XMX=256m \
> -e MODE=standalone \
> -e PREFER_HOST_MODE=hostname \
> -e SPRING_DATASOURCE_PLATFORM=mysql \
> -e MYSQL_SERVICE_HOST=10.0.1.21 \
> -e MYSQL_SERVICE_PORT=3306 \
> -e MYSQL_SERVICE_DB_NAME=nacos_config \
> -e MYSQL_SERVICE_USER=root \
> -e MYSQL_SERVICE_PASSWORD=root123 \
> -v /data/nacos/logs:/home/nacos/logs \
> -v /data/nacos/init.d/custom.properties:/etc/nacos/init.d/custom.properties \
> -v /data/nacos/data:/home/nacos/data \
> harbor.alnk.com/public/nacos-server:v2.2.0
> 
> # 参数解释
> docker run：运行 Docker 容器的命令。
> -d：将容器放入后台并运行。
> --name nacos：将容器命名为“nacos”。
> -p 8848:8848 -p 9848:9848 -p 9849:9849：将主机的端口映射到容器中。主机上的端口 8848、9848 和 9849 分别映射到容器中相同的端口。
> --privileged=true：赋予容器扩展权限。
> --restart=always：如果容器停止运行，自动重新启动容器。
> -e：在容器内设置环境变量。
> JVM_XMS=256m：设置 JVM 的初始堆大小为 256 MB。
> JVM_XMX=256m：设置 JVM 的最大堆大小为 256 MB。
> MODE=standalone：将 Nacos 的模式设置为独立模式。
> PREFER_HOST_MODE=hostname：设置主机模式为主机名。
> SPRING_DATASOURCE_PLATFORM=mysql：设置数据源平台为MySQL。
> MYSQL_SERVICE_HOST=10.4.7.11：设置 MySQL 服务的主机地址。
> MYSQL_SERVICE_PORT=3306：设置 MySQL 服务的端口。
> MYSQL_SERVICE_DB_NAME=nacos_config：设置 Nacos 使用的 MySQL 数据库名称。
> MYSQL_SERVICE_USER=root：设置 MySQL 的用户名。
> MYSQL_SERVICE_PASSWORD=123456：设置 MySQL 的密码。
> -v：挂载主机上的卷到容器中。
> /data/nacos/logs:/home/nacos/logs：将主机上 Nacos 日志目录挂载到容器中。
> /data/nacos/init.d/custom.properties:/etc/nacos/init.d/custom.properties：将主机上的 custom.properties 文件挂载到容器中。
> /data/nacos/data:/home/nacos/data：将主机上 Nacos 数据目录挂载到容器中。
> harbor.od.com/public/nacos-server:v2.2.0：指定容器使用的 Docker 镜像，这里是 Nacos 服务器镜像。
> 
> 注意：
> -e JVM_XMS=256m
> -e JVM_XMX=256m 可以不设置，但nacos默认值会占用1G左右内存，内存不够用的同学最好设置一下
> MYSQL_SERVICE_DB_NAME=nacos_config：必须与之前创建的数据库同名。
> ```

##### nacos数据库文件

> `nacos.sql`
>
> ```
> /*
>  * Copyright 1999-2018 Alibaba Group Holding Ltd.
>  *
>  * Licensed under the Apache License, Version 2.0 (the "License");
>  * you may not use this file except in compliance with the License.
>  * You may obtain a copy of the License at
>  *
>  *      http://www.apache.org/licenses/LICENSE-2.0
>  *
>  * Unless required by applicable law or agreed to in writing, software
>  * distributed under the License is distributed on an "AS IS" BASIS,
>  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
>  * See the License for the specific language governing permissions and
>  * limitations under the License.
>  */
> 
> /******************************************/
> /*   表名称 = config_info                  */
> /******************************************/
> CREATE TABLE `config_info` (
>   `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
>   `data_id` varchar(255) NOT NULL COMMENT 'data_id',
>   `group_id` varchar(128) DEFAULT NULL COMMENT 'group_id',
>   `content` longtext NOT NULL COMMENT 'content',
>   `md5` varchar(32) DEFAULT NULL COMMENT 'md5',
>   `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
>   `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
>   `src_user` text COMMENT 'source user',
>   `src_ip` varchar(50) DEFAULT NULL COMMENT 'source ip',
>   `app_name` varchar(128) DEFAULT NULL COMMENT 'app_name',
>   `tenant_id` varchar(128) DEFAULT '' COMMENT '租户字段',
>   `c_desc` varchar(256) DEFAULT NULL COMMENT 'configuration description',
>   `c_use` varchar(64) DEFAULT NULL COMMENT 'configuration usage',
>   `effect` varchar(64) DEFAULT NULL COMMENT '配置生效的描述',
>   `type` varchar(64) DEFAULT NULL COMMENT '配置的类型',
>   `c_schema` text COMMENT '配置的模式',
>   `encrypted_data_key` text NOT NULL COMMENT '密钥',
>   PRIMARY KEY (`id`),
>   UNIQUE KEY `uk_configinfo_datagrouptenant` (`data_id`,`group_id`,`tenant_id`)
> ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='config_info';
> 
> /******************************************/
> /*   表名称 = config_info_aggr             */
> /******************************************/
> CREATE TABLE `config_info_aggr` (
>   `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
>   `data_id` varchar(255) NOT NULL COMMENT 'data_id',
>   `group_id` varchar(128) NOT NULL COMMENT 'group_id',
>   `datum_id` varchar(255) NOT NULL COMMENT 'datum_id',
>   `content` longtext NOT NULL COMMENT '内容',
>   `gmt_modified` datetime NOT NULL COMMENT '修改时间',
>   `app_name` varchar(128) DEFAULT NULL COMMENT 'app_name',
>   `tenant_id` varchar(128) DEFAULT '' COMMENT '租户字段',
>   PRIMARY KEY (`id`),
>   UNIQUE KEY `uk_configinfoaggr_datagrouptenantdatum` (`data_id`,`group_id`,`tenant_id`,`datum_id`)
> ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='增加租户字段';
> 
> 
> /******************************************/
> /*   表名称 = config_info_beta             */
> /******************************************/
> CREATE TABLE `config_info_beta` (
>   `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
>   `data_id` varchar(255) NOT NULL COMMENT 'data_id',
>   `group_id` varchar(128) NOT NULL COMMENT 'group_id',
>   `app_name` varchar(128) DEFAULT NULL COMMENT 'app_name',
>   `content` longtext NOT NULL COMMENT 'content',
>   `beta_ips` varchar(1024) DEFAULT NULL COMMENT 'betaIps',
>   `md5` varchar(32) DEFAULT NULL COMMENT 'md5',
>   `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
>   `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
>   `src_user` text COMMENT 'source user',
>   `src_ip` varchar(50) DEFAULT NULL COMMENT 'source ip',
>   `tenant_id` varchar(128) DEFAULT '' COMMENT '租户字段',
>   `encrypted_data_key` text NOT NULL COMMENT '密钥',
>   PRIMARY KEY (`id`),
>   UNIQUE KEY `uk_configinfobeta_datagrouptenant` (`data_id`,`group_id`,`tenant_id`)
> ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='config_info_beta';
> 
> /******************************************/
> /*   表名称 = config_info_tag              */
> /******************************************/
> CREATE TABLE `config_info_tag` (
>   `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
>   `data_id` varchar(255) NOT NULL COMMENT 'data_id',
>   `group_id` varchar(128) NOT NULL COMMENT 'group_id',
>   `tenant_id` varchar(128) DEFAULT '' COMMENT 'tenant_id',
>   `tag_id` varchar(128) NOT NULL COMMENT 'tag_id',
>   `app_name` varchar(128) DEFAULT NULL COMMENT 'app_name',
>   `content` longtext NOT NULL COMMENT 'content',
>   `md5` varchar(32) DEFAULT NULL COMMENT 'md5',
>   `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
>   `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
>   `src_user` text COMMENT 'source user',
>   `src_ip` varchar(50) DEFAULT NULL COMMENT 'source ip',
>   PRIMARY KEY (`id`),
>   UNIQUE KEY `uk_configinfotag_datagrouptenanttag` (`data_id`,`group_id`,`tenant_id`,`tag_id`)
> ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='config_info_tag';
> 
> /******************************************/
> /*   表名称 = config_tags_relation         */
> /******************************************/
> CREATE TABLE `config_tags_relation` (
>   `id` bigint(20) NOT NULL COMMENT 'id',
>   `tag_name` varchar(128) NOT NULL COMMENT 'tag_name',
>   `tag_type` varchar(64) DEFAULT NULL COMMENT 'tag_type',
>   `data_id` varchar(255) NOT NULL COMMENT 'data_id',
>   `group_id` varchar(128) NOT NULL COMMENT 'group_id',
>   `tenant_id` varchar(128) DEFAULT '' COMMENT 'tenant_id',
>   `nid` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'nid, 自增长标识',
>   PRIMARY KEY (`nid`),
>   UNIQUE KEY `uk_configtagrelation_configidtag` (`id`,`tag_name`,`tag_type`),
>   KEY `idx_tenant_id` (`tenant_id`)
> ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='config_tag_relation';
> 
> /******************************************/
> /*   表名称 = group_capacity               */
> /******************************************/
> CREATE TABLE `group_capacity` (
>   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
>   `group_id` varchar(128) NOT NULL DEFAULT '' COMMENT 'Group ID，空字符表示整个集群',
>   `quota` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '配额，0表示使用默认值',
>   `usage` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '使用量',
>   `max_size` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '单个配置大小上限，单位为字节，0表示使用默认值',
>   `max_aggr_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '聚合子配置最大个数，，0表示使用默认值',
>   `max_aggr_size` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '单个聚合数据的子配置大小上限，单位为字节，0表示使用默认值',
>   `max_history_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '最大变更历史数量',
>   `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
>   `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
>   PRIMARY KEY (`id`),
>   UNIQUE KEY `uk_group_id` (`group_id`)
> ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='集群、各Group容量信息表';
> 
> /******************************************/
> /*   表名称 = his_config_info              */
> /******************************************/
> CREATE TABLE `his_config_info` (
>   `id` bigint(20) unsigned NOT NULL COMMENT 'id',
>   `nid` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'nid, 自增标识',
>   `data_id` varchar(255) NOT NULL COMMENT 'data_id',
>   `group_id` varchar(128) NOT NULL COMMENT 'group_id',
>   `app_name` varchar(128) DEFAULT NULL COMMENT 'app_name',
>   `content` longtext NOT NULL COMMENT 'content',
>   `md5` varchar(32) DEFAULT NULL COMMENT 'md5',
>   `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
>   `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
>   `src_user` text COMMENT 'source user',
>   `src_ip` varchar(50) DEFAULT NULL COMMENT 'source ip',
>   `op_type` char(10) DEFAULT NULL COMMENT 'operation type',
>   `tenant_id` varchar(128) DEFAULT '' COMMENT '租户字段',
>   `encrypted_data_key` text NOT NULL COMMENT '密钥',
>   PRIMARY KEY (`nid`),
>   KEY `idx_gmt_create` (`gmt_create`),
>   KEY `idx_gmt_modified` (`gmt_modified`),
>   KEY `idx_did` (`data_id`)
> ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='多租户改造';
> 
> 
> /******************************************/
> /*   表名称 = tenant_capacity              */
> /******************************************/
> CREATE TABLE `tenant_capacity` (
>   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
>   `tenant_id` varchar(128) NOT NULL DEFAULT '' COMMENT 'Tenant ID',
>   `quota` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '配额，0表示使用默认值',
>   `usage` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '使用量',
>   `max_size` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '单个配置大小上限，单位为字节，0表示使用默认值',
>   `max_aggr_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '聚合子配置最大个数',
>   `max_aggr_size` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '单个聚合数据的子配置大小上限，单位为字节，0表示使用默认值',
>   `max_history_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '最大变更历史数量',
>   `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
>   `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
>   PRIMARY KEY (`id`),
>   UNIQUE KEY `uk_tenant_id` (`tenant_id`)
> ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='租户容量信息表';
> 
> 
> CREATE TABLE `tenant_info` (
>   `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
>   `kp` varchar(128) NOT NULL COMMENT 'kp',
>   `tenant_id` varchar(128) default '' COMMENT 'tenant_id',
>   `tenant_name` varchar(128) default '' COMMENT 'tenant_name',
>   `tenant_desc` varchar(256) DEFAULT NULL COMMENT 'tenant_desc',
>   `create_source` varchar(32) DEFAULT NULL COMMENT 'create_source',
>   `gmt_create` bigint(20) NOT NULL COMMENT '创建时间',
>   `gmt_modified` bigint(20) NOT NULL COMMENT '修改时间',
>   PRIMARY KEY (`id`),
>   UNIQUE KEY `uk_tenant_info_kptenantid` (`kp`,`tenant_id`),
>   KEY `idx_tenant_id` (`tenant_id`)
> ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='tenant_info';
> 
> CREATE TABLE `users` (
> 	`username` varchar(50) NOT NULL PRIMARY KEY COMMENT 'username',
> 	`password` varchar(500) NOT NULL COMMENT 'password',
> 	`enabled` boolean NOT NULL COMMENT 'enabled'
> );
> 
> CREATE TABLE `roles` (
> 	`username` varchar(50) NOT NULL COMMENT 'username',
> 	`role` varchar(50) NOT NULL COMMENT 'role',
> 	UNIQUE INDEX `idx_user_role` (`username` ASC, `role` ASC) USING BTREE
> );
> 
> CREATE TABLE `permissions` (
>     `role` varchar(50) NOT NULL COMMENT 'role',
>     `resource` varchar(128) NOT NULL COMMENT 'resource',
>     `action` varchar(8) NOT NULL COMMENT 'action',
>     UNIQUE INDEX `uk_role_permission` (`role`,`resource`,`action`) USING BTREE
> );
> 
> INSERT INTO users (username, password, enabled) VALUES ('nacos', '$2a$10$EuWPZHzz32dJN7jexM34MOeYirDdFAZm2kuWj7VEOJhhZkDrxfvUu', TRUE);
> 
> INSERT INTO roles (username, role) VALUES ('nacos', 'ROLE_ADMIN');
> ```

##### 访问设置nacos

> `访问nacos`
>
> ```
> #浏览器访问：http://10.0.1.21:8848/nacos/#/login
> 账号：nacos
> 密码：nacos
> ```
>
> ![1731080941941](images/1731080941941.png)  
>
> ![1731081419117](images/1731081419117.png)  
>
> ![1731081613134](images/1731081613134.png)  
>
> ![1731081678221](images/1731081678221.png)  
>
> ![1731081700254](images/1731081700254.png)  



#### gitlab上传代码

> `代码目录`
>
> ![1732239860597](images/1732239860597.png)  
>
> `新建项目`
> 
> ![1731082143074](images/1731082143074.png)  
> 
> `创建test分支`
> 
> ![1731082226406](images/1731082226406.png)  
> 
> `上传代码`
> 
> ![1731082466485](images/1731082466485.png)  



#### 配置jenkins流水线，test分支发布服务到k8s test环境中  

> `新建流水线`
>
> ![1732236297010](images/1732236297010.png)    
>
> ![1731083603900](images/1731083603900.png)  
>
> ![1731083651543](images/1731083651543.png)  
>
> ![1731083764308](images/1731083764308.png) 
>
> ![1731083839547](images/1731083839547.png)  
>
> `这里直接选择SCM，保证在git clone下的目录进行流水线的操作，他会自动拉取整个项目，项目包含了Jenkinsfile，然后会按照Jenkinsfile的pipeline去执行，直接构建打包推送发布一键完成` 
>
>  ![1732239938028](images/1732239938028.png)  
>
> ```
> 然后进行构建，构建之前现在k8s环境中创建相关的名称空间
> kubectl create ns hello-yewu
> ```
>
> 
>
> `到蓝海查看输出`
>
> `1. 拉取代码，切换分支`
>
> ![1732240344875](images/1732240344875.png)  
>
> `2.构建镜像 `
>
> ![1732240397389](images/1732240397389.png)  
>
> `3.推送镜像到harbor仓库，并且删除掉本地的镜像，节约磁盘空间 `
>
> ![1732240445470](images/1732240445470.png)  
>
> `4. 替换yaml文件中的镜像，然后进行发布到k8s集群中`
>
> ![1732240494941](images/1732240494941.png)  
>
>  
>
>  
>
> ![1732240607794](images/1732240607794.png)  
>
> ![1732240621402](images/1732240621402.png)  
>
> ![1732240636536](images/1732240636536.png)  
>
> 









  

#### 配置gitlab提交代码触发jenkins

##### jenkins 安装gitlab插件

> `jenkins安装GitLab插件`
>
> ![1731152828884](images/1731152828884.png)

##### jenkins修改构建器

> `jenkins修改构建触发器`
>
> http://jenkins.alnk.com/project/go-k8s-one
>
> ![1731154093132](images/1731154093132.png)  
>
> `选择test分支，只允许test分支进行构建`
>
> ![1731154203804](images/1731154203804.png)  
>
> `允许访问/project，去掉√`
>
> ![1731155267310](images/1731155267310.png)  

##### gitlab修改webhooks

> `gitlab上修改webhooks`
>
> ![1731154875298](images/1731154875298.png)  
>
> ![1731154314493](images/1731154314493.png)  
>
> ```
> 【10.0.1.21】
> ## 添加解析
> # cat /etc/hosts
> 10.0.1.100 jenkins.alnk.com
> 
> ```
>
> `gitlab测试`
>
> ![1731155339092](images/1731155339092.png)  

##### gitlab test分支提交代码进行测试

> `gitlab上test分支提交代码`
>
> ![1731155427126](images/1731155427126.png)  
>
> `可以看到已经触发`
>
> ![1731155465474](images/1731155465474.png)  
>
> ![1731155599502](images/1731155599502.png)  





#### 测试发布服务到k8s生产环境

##### jenkins流水线复制与修改

> 复制一条流水线，然后修改
>
> ![1732241676355](images/1732241676355.png)  
>
> ![1732241721359](images/1732241721359.png)  
>
> ​    
>
> `修改分支版本默认值`
>
> ![1731156481449](images/1731156481449.png)  
>
> ![1732241780296](images/1732241780296.png)  
>
> `记录webhook地址，等下gitlab需要用到`
>
> http://jenkins.alnk.com/project/prod-go-k8s-one
>
> ![1731156571077](images/1731156571077.png)  
>
> `修改触发的分支为main`
>
> ![1731156603617](images/1731156603617.png)  



##### gitlab合并代码，修改配置文件和main文件，jenkinsfile文件,yaml文件等

> `修改jenkinsfile`
>
> `修改为prod`
>
> ![1732242102855](images/1732242102855.png)  
>
> ![1732242134082](images/1732242134082.png)  
>
> ![1732242210744](images/1732242210744.png)  
>
> 
>
> 
>
> `修改nacos的名称空间为prod的ID`
>
> ![1732242277913](images/1732242277913.png)  
>
>   
>
> ![1732242381041](images/1732242381041.png)    

##### jenkins构建

> ```
> 构建之前，在prod 的k8s环境中创建名称空间
> kubectl create ns hello-yewu
> ```
>
> `可以看到通过不同kube config文件，把不同分支的代码发布到了不同的k8s环境中`
>
> ![1732242697178](images/1732242697178.png)  



##### gitlab配置webhooks

> 参考上面的操作




