#### gitlab pipeline语法表

> script	  运行的 Shell 命令或脚本
> image	 使用 docker 映像
> services	使用 docker 服务映像
> before_script	在作业运行前运行脚本
> after_script	在作业运行后运行脚本
> stages		定义管道中的阶段，运行顺序
> stage		为job定义一个阶段，可选，未指定默认为test阶段
> only		限制创建作业的条件
> except	限制未创建作业的条件
> rules		条件列表，用于评估和确定作业的选定属性，以及是否创建该作业。不能 only与/ except一起使用
> tags					用于选择 Runner 的标签列表
> allow_failure					允许作业失败，失败的 job 不会影响提交状态
> when							什么时候开始运行工作
> environment				作业部署到的环境的名称
> cache							在后续运行之间应缓存的文件列表
> artifacts					成功时附加到作业的文件和目录列表
> dependencies		通过提供要从中获取工件的作业列表，限制将哪些工件传递给特定作业
> retry						发生故障时可以自动重试作业的时间和次数
> timeout				  定义自定义作业级别的超时，该超时优先于项目范围的设置
> parallel				   多 个作业并行运行
> trigger	  			  定义下游管道触发
> include	              允许此作业包括外部 YAML 文件
> extends	            该作业将要继承的配置条目
> pages	               上载作业结果以用于 GitLab 页面
> variables	           在作业级别上定义作业变量

##### 语法检测

![1730993482732](images/1730993482732.png)

![1730993531567](images/1730993531567.png)

#### pipline语法基础

##### job

> 在每个项目中，使用名为 `.gitlab-ci.yml `的 YAML 文件配置 GitLab CI / CD 管道
>
> - 可以定义一个或多个作业(job)
> - 每个作业必须具有唯一的名称（不能使用关键字）
> - 每个作业是独立执行的
> - 每个作业至少要包含一个 script
>
> ```yaml
> job1:
>   script: "execute-script-for-job1"
>  
> job2:
>   script: "execute-script-for-job2"
> ```
>
> 注释： 这里在 pipeline 中定义了两个作业，每个作业运行不同的命令。命令可以是 shell 或脚本

##### script

> 每个作业（job）至少要包含一个 script
>
> ```yaml
> job:
>   script:
>     - uname -a
>     - bundle exec rspec
> ```
>
> 注意：有时，script 命令将需要用单引号或双引号引起来. 例如，包含冒号命令（ : ）需要加引号，以便被包裹的 YAML 解析器知道来解释整个事情作为一个字符串，而不是一个"键：值"对。使用特殊字符时要小心：: ， { ， } ， [ ， ] ， , ， & ， * ， # ， ? ， | ， - ， < ， > ， = ! ， % ， @ .

##### before_script

> 用于定义一个命令，该命令在每个作业之前运行。必须是一个数组。指定的 script 与主脚本中指定的任何脚本串联在一起，并在单个 shell 中一起执行。
>
> before_script 失败会导致整个作业失败，其他作业将不再执行。作业失败不会影响 after_script 运行

##### after_script

> 用于定义将在每个作业（包括失败的作业）之后运行的命令。这必须是一个数组。指定的脚本在新的 shell 中执行，与任何 `before_script `或 `script `脚本分开。可以在全局定义，也可以在 job 中定义，在 job 中定义会覆盖全局
>
> ```yaml
> before_script:
>   - echo "before-script!!"
>  
> variables:
>   DOMAIN: example.com
>  
> stages:
>   - build
>   - deploy
>  
>  
> build:
>   before_script:
>     - echo "before-script in job"
>   stage: build
>   script:
>     - echo "mvn clean "
>     - echo "mvn install"
>   after_script:
>     - echo "after script in job"
>  
>  
> deploy:
>   stage: deploy
>   script:
>     - echo "hello deploy"
>   
> after_script:
>   - echo "after-script"
> ```
>
> after_script 失败不会影响作业失败

##### stages

> 用于定义作业可以使用的阶段，并且是全局定义的。同一阶段的作业并行运行，不同阶段按顺序执行
>
> ```yaml
> stages：
>   - build
>   - test
>   - deploy
> ```
>
> 这里定义了三个阶段，首先 build 阶段并行运行，然后 test 阶段并行运行，最后 deploy 阶段并行运行。deploy 阶段运行成功后将提交状态标记为 passed 状态。如果任何一个阶段运行失败，最后提交状态为 failed

##### 未定义 stages

> 全局定义的 stages 是来自于每个 job。如果 job 没有定义 stage 则默认是 test 阶段。如果全局未定义 stages，则按顺序运行 build,test,deploy

##### 定义 stages 控制 stage 运行顺序

> 一个标准的 yaml 文件中是需要定义 stages，可以帮助我们对每个 stage 进行排序

##### .pre & .post

> .pre 始终是整个管道的第一个运行阶段，.post 始终是整个管道的最后一个运行阶段。 用户定义的阶段都在两者之间运行。`.pre `和 `.post `的顺序无法更改。如果管道仅包含 `.pre `或 `.post`阶段的作业，则不会创建管道

##### stage

> 是按 JOB 定义的，并且依赖于全局定义的 stages 。它允许将作业分为不同的阶段，并且同一 `stage `作业可以并行执行（取决于特定条件 ）
>
> ```yaml
> unittest:
>   stage: test
>   script:
>     - echo "run test"
>   
> interfacetest:
>   stage: test
>   script:
>     - echo "run test"
> ```
>
> 可能遇到下面问题： 阶段并没有并行运行
>
> 在这里我把这两个阶段在同一个 runner 运行了，所以需要修改 runner 每次运行的作业数量。默认是 1，改为 10
>
> ![1730994605581](images/1730994605581.png)

##### variables

> 定义变量，pipeline 变量、job 变量、Runner 变量。job 变量优先级最大

##### 综合实例（一）

```yaml
before_script:
  - echo "before-script!!"
 
variables:
  DOMAIN: example.com
  
stages:
  - build
  - test
  - codescan
  - deploy
 
build:
  before_script:
    - echo "before-script in job"
  stage: build
  script:
    - echo "mvn clean "
    - echo "mvn install"
    - echo "$DOMAIN"
  after_script:
    - echo "after script in buildjob"
 
unittest:
  stage: test
  script:
    - echo "run test"
 
deploy:
  stage: deploy
  script:
    - echo "hello deploy"
    - sleep 2;
  
codescan:
  stage: codescan
  script:
    - echo "codescan"
    - sleep 5;
 
after_script:
  - echo "after-script"
```

![1730994738276](images/1730994738276.png)

> 可能遇到的问题：pipeline 卡主，为降低复杂性目前没有学习 tags，所以流水线是在共享的runner 中运行的。需要设置共享的 runner 运行没有 tag 的作业
>
> ![1730994797507](images/1730994797507.png)
>
> ![1730994816426](images/1730994816426.png)  

##### tags

> 用于从允许运行该项目的所有 Runner 列表中选择特定的 Runner，在 Runner 注册期间，您可以指定 Runner 的标签
>
> `tags `可让您使用指定了标签的runner来运行作业，此 runner 具有 ruby 和 postgres 标签
>
> ```yaml
> job:
>   tags:
>     - ruby
>     - postgres
> ```
>
> 给定带有 `deploy `标签的 deploy Runner 和带有 `build `标签的 build Runner，以下作业将在各自的平台上运行
>
> ```yaml
> build job:
>   stage:
>     - build
>   tags:
>     - build
>   script:
>     - echo Hello, %USERNAME%!
>  
> deploy job:
>   stage:
>     - build
>   tags:
>     - deploy
>   script:
>     - echo "Hello, $USER!"
> ```

##### allow_failure

> allow_failure 允许作业失败，默认值为 false 。启用后，如果作业失败，该作业将在用户界面中显示橙色警告。但是，管道的逻辑流程将认为作业成功/通过，并且不会被阻塞。假设所有其他作业均成功，则该作业的阶段及其管道将显示相同的橙色警告。但是，关联的提交将被标记为"通过"，而不会发出警告
>
> ```yaml
> job1:
>   stage: test
>   script:
>     - execute_script_that_will_fail
>   allow_failure: true
> ```

##### when

> - `on_success：`前面阶段中的所有作业都成功（或由于标记为 `allow_failure `而被视为成功）时才执行作业。 这是默认值。
> - `on_failure：`当前面阶段出现失败则执行。
> - `always`：执行作业，而不管先前阶段的作业状态如何，放到最后执行。总是执行

##### manual 手动

> `manual` 手动执行作业，不会自动执行，需要由用户显式启动。手动操作的示例用法是部署到生产环境。可以从管道，作业，环境和部署视图开始手动操作
>
> 此时在 deploy 阶段添加 manual，则流水线运行到 deploy 阶段为锁定状态，需要手动点击按钮才能运行 deploy 阶段

##### delayed 延迟

> `delayed` 延迟一定时间后执行作业（在 GitLab 11.14中已添加）
>
> 有效值 `'5', 10 seconds, 30 minutes, 1 day, 1 week`

##### 综合实例（二）

> ```yaml
> before_script:
>   - echo "before-script!!"
>  
> variables:
>   DOMAIN: example.com
>   
> stages:
>   - build
>   - test
>   - codescan
>   - deploy
>  
> build:
>   before_script:
>     - echo "before-script in job"
>   stage: build
>   script:
>     - echo "mvn clean "
>     - echo "mvn install"
>     - echo "$DOMAIN"
>   after_script:
>     - echo "after script in buildjob"
>  
> unittest:
>   stage: test
>   script:
>     - ech "run test"
>   when: delayed
>   start_in: '30'
>   allow_failure: true
>   
>  
> deploy:
>   stage: deploy
>   script:
>     - echo "hello deploy"
>     - sleep 2;
>   when: manual
>   
> codescan:
>   stage: codescan
>   script:
>     - echo "codescan"
>     - sleep 5;
>   when: on_success
>  
> after_script:
>   - echo "after-script"
> ```

##### retry-重试

> 配置在失败的情况下重试作业的次数。
>
> 当作业失败并配置了 retry ，将再次处理该作业，直到达到 retry 关键字指定的次数。如果retry 设置为 2，并且作业在第二次运行成功（第一次重试），则不会再次重试；retry 值必须是一个正整数，等于或大于 0，但小于或等于 2（最多两次重试，总共运行 3 次）
>
> ```yaml
> unittest:
>   stage: test
>   retry: 2
>   script:
>     - ech "run test"
> ```
>
> 默认情况下，将在所有失败情况下重试作业。为了更好地控制 `retry `哪些失败，可以是具有以下键的哈希值：
>
> - `max` ：最大重试次数.
> - `when` ：重试失败的案例.
>
> 根据错误原因设置重试的次数
>
> ```text
> always ：在发生任何故障时重试（默认）.
> unknown_failure ：当失败原因未知时。
> script_failure ：脚本失败时重试。
> api_failure ：API失败重试。
> stuck_or_timeout_failure ：作业卡住或超时时。
> runner_system_failure ：运行系统发生故障。
> missing_dependency_failure: 如果依赖丢失。
> runner_unsupported ：Runner不受支持。
> stale_schedule ：无法执行延迟的作业。
> job_execution_timeout ：脚本超出了为作业设置的最大执行时间。
> archived_failure ：作业已存档且无法运行。
> unmet_prerequisites ：作业未能完成先决条件任务。
> scheduler_failure ：调度程序未能将作业分配给运行scheduler_failure。
> data_integrity_failure ：检测到结构完整性问题。
> ```
>
> ##### 实验
>
> 定义当出现脚本错误重试两次，也就是会运行三次。
>
> ```yaml
> unittest:
>   stage: test
>   tags:
>     - build
>   only:
>     - master
>   script:
>     - ech "run test"
>   retry:
>     max: 2
>     when:
>       - script_failure
> ```

##### timeout 超时

> 特定作业配置超时，作业级别的超时可以超过项目级别的超时，但不能超过 Runner 特定的超时。
>
> ```yaml
> build:
>   script: build.sh
>   timeout: 3 hours 30 minutes
>  
> test:
>   script: rspec
>   timeout: 3h 30m
> ```
>
> ##### 项目设置流水线超时时间
>
>     超时定义了作业可以运行的最长时间（以分钟为单位）。 这可以在项目的**"设置">" CI / CD">"常规管道"设置下进行配置** 。 默认值为 60 分钟

##### runner 超时时间

> 此类超时（如果小于项目定义的超时 ）将具有优先权。此功能可用于通过设置大超时（例如一个星期）来防止 Shared Runner 被项目占用。未配置时，Runner 将不会覆盖项目超时
>
> 此功能如何工作：
>
> 示例1-运行程序超时大于项目超时
>
> runner 超时设置为 24 小时，项目的 CI / CD 超时设置为 2 小时。该工作将在 2 小时后超时。
>
> 示例2-未配置运行程序超时
>
> runner 不设置超时时间，项目的 CI / CD 超时设置为2 小时。该工作将在 2 小时后超时。
>
> 示例3-运行程序超时小于项目超时
>
> runner 超时设置为 30 分钟，项目的 CI / CD 超时设置为 2 小时。工作在 30 分钟后将超时

##### parallel-并行作业

> 配置要并行运行的作业实例数，此值必须大于或等于 2 并且小于或等于 50。
>
> 这将创建 N 个并行运行的同一作业实例。它们从 `job_name 1/N `到 `job_name N/N `依次命名
>
> ```yaml
> codescan:
>   stage: codescan
>   tags:
>     - build
>   only:
>     - master
>   script:
>     - echo "codescan"
>     - sleep 5;
>   parallel: 5
> ```

##### 综合实例（三）

> ```yaml
> before_script:
>   - echo "before-script!!"
>  
> variables:
>   DOMAIN: example.com
>   
> stages:
>   - build
>   - test
>   - codescan
>   - deploy
>  
> build:
>   before_script:
>     - echo "before-script in job"
>   stage: build
>   script:
>     - echo "mvn clean "
>     - echo "mvn install"
>     - echo "$DOMAIN"
>   after_script:
>     - echo "after script in buildjob"
>  
> unittest:
>   stage: test
>   script:
>     - ech "run test"
>   when: delayed
>   start_in: '5'
>   allow_failure: true
>   retry:
>     max: 1
>     when:
>       - script_failure
>   timeout: 1 hours 10 minutes
>   
> deploy:
>   stage: deploy
>   script:
>     - echo "hello deploy"
>     - sleep 2;
>   when: manual
>   
> codescan:
>   stage: codescan
>   script:
>     - echo "codescan"
>     - sleep 5;
>   when: on_success
>   parallel: 5
>  
> after_script:
>   - echo "after-script"
>   - ech
> ```

##### only & except-限制分支标签

> only 和 except 是两个参数用分支策略来限制 jobs 构建：
>
> 1. `only `定义哪些分支和标签的 git 项目将会被 job 执行。
> 2. `except `定义哪些分支和标签的 git 项目将不会被 job 执行
>
> ```yaml
> job:
>   # use regexp
>   only:
>     - /^issue-.*$/
>   # use special keyword
>   except:
>     - branches
> ```

##### rules-构建规则

> rules 允许按顺序评估单个规则对象的列表，直到一个匹配并为作业动态提供属性；请注意， rules 不能 only/except 与 only/except 组合使用。
>
> 可用的规则条款包括：
>
> if （类似于 only:variables ）
>
> changes （ only:changes 相同）
>
> exists
>
> rules:if
> 如果 DOMAIN 的值匹配，则需要手动运行。不匹配 on_success。 条件判断从上到下，匹配即停止。多条件匹配可以使用 && ||
>
> ```yaml
> variables:
>   DOMAIN: example.com
>  
> codescan:
>   stage: codescan
>   tags:
>     - build
>   script:
>     - echo "codescan"
>     - sleep 5;
>   #parallel: 5
>   rules:
>     - if: '$DOMAIN == "example.com"'
>       when: manual
>     - when: on_success
> ```
>
> ##### rules:changes
>
> 接受文件路径数组。 如果提交中 `Jenkinsfile `文件发生的变化则为 true。
>
> ```yaml
> codescan:
>   stage: codescan
>   tags:
>     - build
>   script:
>     - echo "codescan"
>     - sleep 5;
>   #parallel: 5
>   rules:
>     - changes:
>       - Jenkinsfile
>       when: manual
>     - if: '$DOMAIN == "example.com"'
>       when: on_success
>     - when: on_success
> ```
>
> ##### rules:exists
>
> 接受文件路径数组。当仓库中存在指定的文件时操作。
>
> ```yaml
> codescan:
>   stage: codescan
>   tags:
>     - build
>   script:
>     - echo "codescan"
>     - sleep 5;
>   #parallel: 5
>   rules:
>     - exists:
>       - Jenkinsfile
>       when: manual 
>     - changes:
>       - Jenkinsfile
>       when: on_success
>     - if: '$DOMAIN == "example.com"'
>       when: on_success
>     - when: on_success
> ```
>
> ##### rules:allow_failure
>
> 使用[ allow_failure: true](http://s0docs0gitlab0com.icopy.site/12.9/ee/ci/yaml/README.html#allow_failure)  `rules:`在不停止管道本身的情况下允许作业失败或手动作业等待操作.
>
> ```yaml
> job:
>   script: "echo Hello, Rules!"
>   rules:
>     - if: '$CI_MERGE_REQUEST_TARGET_BRANCH_NAME == "master"'
>       when: manual
>       allow_failure: true
> ```
>
> 在此示例中，如果第一个规则匹配，则作业将具有以下 `when: manual `和 `allow_failure: true`。
>
> workflow:rules
>
>     顶级`workflow:`关键字适用于整个管道，并将确定是否创建管道。[when](http://s0docs0gitlab0com.icopy.site/12.9/ee/ci/yaml/README.html#when) ：可以设置为 `always`或 `never`；如果未提供，则默认值 `always`。
>
> ```yaml
> variables:
>   DOMAIN: example.com
>  
> workflow:
>   rules:
>     - if: '$DOMAIN == "example.com"'
>     - when: always
> ```

##### 综合实例（四）

```yaml
before_script:
  - echo "before-script!!"
 
variables:
  DOMAIN: example.com
  
workflow:
  rules:
    - if: '$DOMAIN == "example.com"'
      when: always
    - when: never
  
stages:
  - build
  - test
  - codescan
  - deploy
 
build:
  before_script:
    - echo "before-script in job"
  stage: build
  script:
    - echo "mvn clean "
    - echo "mvn install"
    - ech "$DOMAIN"
  after_script:
    - echo "after script in buildjob"
  rules:
    - exists:
      - Dockerfile
      when: on_success 
      allow_failure: true
 
    - changes:
      - Dockerfile
      when: manual
    - when: on_failure
 
unittest:
  stage: test
  script:
    - ech "run test"
  when: delayed
  start_in: '5'
  allow_failure: true
  retry:
    max: 1
    when:
      - script_failure
  timeout: 1 hours 10 minutes
  
  
 
deploy:
  stage: deploy
  script:
    - echo "hello deploy"
    - sleep 2;
  rules:
    - if: '$DOMAIN == "example.com"'
      when: manual
    - if: '$DOMAIN == "aexample.com"'
      when: delayed
      start_in: '5'
    - when: on_failure
  
codescan:
  stage: codescan
  script:
    - echo "codescan"
    - sleep 5;
  when: on_success
  parallel: 5
 
after_script:
  - echo "after-script"
  - ech
```

更多语法可以参考官网

推荐博客：https://blog.csdn.net/weixin_46560589/category_12381994.html?spm=1001.2014.3001.5482
