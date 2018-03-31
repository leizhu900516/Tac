# Tac
##### golang版本的后台任务部署平台：支持web形式的后台任务新建。
#### 实现功能
- 新建后台任务 (已完成)
- 新建计划任务(类crontab) (未开发)

#### 效果图
![](/example/images/index.png '主页预览图')
![](/example/images/newAddTask.png '新建任务')

####安装及使用
- agent部署
  - 配置文件修改/agent/parserconfig/config
  - agent目录运行 go run rpcServer.go 文件
- web端部署
  - 基于beego v1.9.1框架实现
  - clone代码到服务器。配置mysql服务。配置文件conf/app.conf
  - 运行项目 go run main.go