#### gitlab-runner简介

> GitLab Runner 是一个开源项目，用于运行作业并将结果发送回GitLab。它与 Gitlab CI 结合使用， Gitlab CI 是 Gitlab 随附的用于协调作业的开源持续集成服务



##### 要求

> GitLab Runner 是用 Go 编写的，可以作为一个二进制文件运行。它旨在在 GNU/Linux，macOS 和 Windows 操作系统上运行。只要可以在其他操作系统上编译 Go 二进制文件，其他操作系统就可能会运行。
>
> 如果要使用 Docker，请安装最新版本。GitLab Runner 需要最少的 Docker v1.13.0。
>
> GitLab Runner 版本应与 GitLab 版本同步。尽管较旧的 Runner 仍可以使用较新的 GitLab 版本，反之亦然，但在某些情况下，如果版本存在差异，则功能可能不可用或无法正常工作。在次要版本更新之间可以保证向后兼容性，但是请注意，GitLab 的次要版本更新会引入新功能，这些新功能将要求 Runner 在同一次要版本上使用
>



##### 特点

> 允许运行：
>
> 同时执行多个作业
>
> 对多个服务器（甚至每个项目）使用多个令牌
>
> 限制每个令牌的并行作业数
>
> 
>
> 可以运行作业：
>
> 在本地
>
> 使用 Docker 容器
>
> 使用 Docker 容器并通过 SSH 执行作业
>
> 使用 Docker 容器在不同的云和虚拟化管理程序上自动缩放
>
> 连接到远程 SSH 服务器
>
> 
>
> 用 Go 编写并以单个二进制文件的形式分发，而没有其他要求
>
> 支持 Bash，Windows Batch 和 Windows PowerShell
>
> 在 GNU / Linux，macOS 和 Windows（几乎可以在任何可以运行 Docker 的地方）上运行
>
> 允许自定义作业运行环境
>
> 自动重新加载配置，无需重启
>
> 易于使用的设置，并支持 Docker，Docker-SSH，Parallels 或 SSH 运行环境
>
> 启用 Docker 容器的缓存
>
> 易于安装，可作为 GNU / Linux，macOS 和 Windows 的服务
>
> 嵌入式 Prometheus 指标 HTTP 服务器
>
> 裁判工作者监视 Prometheus 度量标准和其他特定于工作的数据并将其传递给 GitLab



#### gitlab-runner安装

> 建议生产环境不要和gitlab安装在同一台机器

##### 使用dep软件包安装

```shell
【10.0.1.21】
## 清华源地址： https://mirrors.tuna.tsinghua.edu.cn/

# cd /opt/src
# wget https://mirrors.tuna.tsinghua.edu.cn/gitlab-runner/ubuntu/pool/jammy/main/g/gitlab-runner/gitlab-runner_15.5.0_amd64.deb
# apt install ./gitlab-runner_15.5.0_amd64.deb


# 启动/停止 服务
# systemctl start gitlab-runner
# systemctl stop gitlab-runner

# 开机启动
# systemctl enable gitlab-runner

```





#### gitlab-runner注册

##### gitlab-runner类型

> - shared：运行整个平台项目的作业（gitlab）
> - group：运行特定 group 下的所有项目的作业（group）
> - specific: 运行指定的项目作业（project）
> - locked：无法运行项目作业
> - paused：不会运行作业



##### 获取runner token

###### **获取 shared 类型 runner token**  

> 需要管理员权限

![1730983771301](images\1730983771301.png)  

![1730983834006](images\1730983834006.png)  





###### 获取 group 类型的 runner token

![1730985365658](images\1730985365658.png)    

![1730985402131](images\1730985402131.png)  

![1730985437530](images\1730985437530.png)  

![1730985475244]( images\1730985475244.png)  



###### 获取 specific 类型的 runner token

![1730985516803](images\1730985516803.png)  

![1730985543715](images\1730985543715.png)  

![1730985664768](images\1730985664768.png)  





##### 进行注册

> 注册一个group类型的runner

```shell
【10.0.1.21】
# gitlab-runner register
Enter the GitLab instance URL (for example, https://gitlab.com/):
http://10.0.1.21:6666

Enter the registration token:
GR1348941kiYwWxsTttC5yN2KMxPp

Enter a description for the runner:
[ops]: test

Enter tags for the runner (comma-separated):
build

Enter optional maintenance note for the runner:
this is d test

Enter an executor: ssh, virtualbox, docker-ssh+machine, instance, custom, docker, shell, docker+machine, kubernetes, docker-ssh, parallels:
shell

```

![1730987025394](images\1730987025394.png)  



> 注册一个specific（项目）类型的runner

```shell
# sudo gitlab-runner register --url http://10.0.1.21:6666/ --registration-token GR13489411Ye2q-BfyqaEVqNTNgv-

[http://10.0.1.21:6666/]:
Enter the registration token:
[GR13489411Ye2q-BfyqaEVqNTNgv-]:
Enter a description for the runner:
[ops]: test
Enter tags for the runner (comma-separated):
this is a test
Enter optional maintenance note for the runner:
this is a test test

Enter an executor: ssh, docker+machine, docker-ssh+machine, instance, docker-ssh, parallels, shell, kubernetes, custom, docker, virtualbox:
shell
```

![1730986996553](images\1730986996553.png)  



#### gitlab-runner命令

##### 启动命令

```shell
gitlab-runner --debug <command>   # 调试模式排查错误特别有用。
gitlab-runner <command> --help    # 获取帮助信息
gitlab-runner run                 # 普通用户模式  配置文件位置 ~/.gitlab-runner/config.toml
sudo gitlab-runner run            # 超级用户模式  配置文件位置 /etc/gitlab-runner/config.toml
```



##### 注册命令

```shell
gitlab-runner register      # 默认交互模式下使用，非交互模式添加 --non-interactive
gitlab-runner list          # 此命令列出了保存在配置文件中的所有运行程序
gitlab-runner verify        # 此命令检查注册的 runner 是否可以连接，但不验证 GitLab 服务是否正在使用 runner。 --delete 删除
gitlab-runner unregister    # 该命令使用 GitLab 取消已注册的 runner。
 
 
# 使用令牌注销
gitlab-runner unregister --url http://xxx/ --token t0kxx
 
# 使用名称注销（同名删除第一个）
gitlab-runner unregister --name test-runner
 
# 注销所有
gitlab-runner unregister --all-runners
```



##### 服务管理

```shell
gitlab-runner install --user=gitlab-runner --working-directory=/home/gitlab-runner
 
# --user 指定将用于执行构建的用户
# --working-directory  指定将使用 Shell executor 运行构建时所有数据将存储在其中的根目录
 
gitlab-runner uninstall # 该命令停止运行并从服务中卸载 GitLab Runner。
 
gitlab-runner start     # 该命令启动 GitLab Runner 服务。
 
gitlab-runner stop      # 该命令停止 GitLab Runner 服务。
 
gitlab-runner restart   # 该命令将停止，然后启动 GitLab Runner 服务。
 
gitlab-runner status    # 此命令显示 GitLab Runner 服务的状态。当服务正在运行时，退出代码为零；而当服务未运行时，退出代码为非零。
 
# 也可以是使用 systemctl 管理 runner
```



#### 运行流水线任务

##### 编写yaml文件

> 在 gitlab 仓库中项目根目录添加一个 .gitlab-ci.yml 文件，文件内容如下
>
> ```yaml
> #这个流水线共包含两个 job，分别是 build 和 deploy
> stages:
>   - build
>   - deploy
> 
> #build job包含一个 stage build
> #build stage配置了在具有build标签的runne 中运行，限制为master分支提交，运行构建命令
> build:
>   stage: build
>   tags:
>     - build
>   only:
>     - master
>   script:
>     - echo "mvn clean "
>     - echo "mvn install"
> 
> #deploy job包含一个stage deploy 
> #deploy stage配置了在具有deploy标签的runner中运行，限制为master分支提交，运行发布命令
> deploy:
>   stage: deploy
>   tags:
>     - deploy
>   only:
>     - master
>   script:
>     - echo "hello deploy"
> 
> ```





##### 测试流水线

> 修改一下项目runner的标签

![1730987475379](images\1730987475379.png)  

![1730987446045](images\1730987446045.png)



> 新建一个demo群组
>
> ![1730988269112](images\1730988269112.png) 
>
> 在demo群组下再创建一个runner
>
> ```shell
> 【10.0.1.21】
> # gitlab-runner register
> Enter the GitLab instance URL (for example, https://gitlab.com/):
> http://10.0.1.21:6666
> 
> Enter the registration token:
> GR13489412k2PMwGhEdBdG8DHDWss
> 
> Enter a description for the runner:
> [ops]: test
> 
> Enter tags for the runner (comma-separated):
> build
> 
> Enter optional maintenance note for the runner:
> this is demo group test
> 
> Enter an executor: ssh, virtualbox, docker-ssh+machine, instance, custom, docker, shell, docker+machine, kubernetes, docker-ssh, parallels:
> shell
> 
> ```
>
> ![1730988463976](images\1730988463976.png)  





> 创建一个代码项目
>
> ![1730987599612](images\1730987599612.png)  
>
> ![1730987638980](images\1730987638980.png)  
>
> 修改下项目所属组
>
> ![1730988083787](images\1730988083787.png)  
>
> ![1730988127715](images\1730988127715.png) 
>
> 在该项目下创建一个runner
>
> ![1730988598698](images\1730988598698.png)  
>
> ```shell
> # sudo gitlab-runner register --url http://10.0.1.21:6666/ --registration-token GR134894181jDkfGxBZfaJ4A4t5Cd
> 
> # sudo gitlab-runner register --url http://10.0.1.21:6666/ --registration-token GR134894181jDkfGxBZfaJ4A4t5Cd
> 
> Enter the GitLab instance URL (for example, https://gitlab.com/):
> [http://10.0.1.21:6666/]:
> 
> Enter the registration token:
> [GR134894181jDkfGxBZfaJ4A4t5Cd]:
> 
> Enter a description for the runner:
> [ops]: demo-maven-service
> 
> Enter tags for the runner (comma-separated):
> deploy
> 
> Enter optional maintenance note for the runner:
> this is a demo test
> 
> docker+machine, instance, docker, shell, ssh, docker-ssh+machine:
> shell
> 
> ```
>
> ![1730988736692](images\1730988736692.png)  



模拟提交代码

![1730988802840]( images\1730988802840.png)  

  ![1730988909419](images\1730988909419.png)  

![1730992878207](images\1730992878207.png)  

  



  

  







