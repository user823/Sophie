# Project Structure

项目结构：
Sophie
├── api
│    ├── domain         ------领域模型定义、存储层交互等
│    ├── thrift         ------IDL 定义
│    ├── consts.go      ------通用常量
│    └── meta.go      
├── bin                 ------二进制文件
├── cmd                 ------各个模块的入口
│    ├── file
│    ├── gateway
│    ├── gen
│    ├── logstash
│    ├── schedule
│    ├── system
│    └── README.md
├── configs             ------项目配置
│    ├── cert           ------提供https服务的证书
│    ├── cfssl          ------自签证书工具
│    ├── templates      ------各个服务组件的配置文件模板
│    ├── README.md
│    ├── data.sql       
│    ├── es-init.sh
│    ├── goss.yml
│    └── sophie.sql
├── deployments         ------部署相关
│    ├── docker         ------服务依赖组件的docker部署、各个服务组件的镜像制作等
│    └── README.md
├── docs                ------文档相关
│    ├── api            ------各个组件RESTful api
│    ├── devel          ------开发文档
│    ├── guide          ------部署手册及项目补充信息
│    ├── images       
│    ├── swagger        
│    └── README.md
├── githooks            
├── internal            ------组件源码
│    ├── file
│    ├── gateway
│    ├── gen
│    ├── logstash
│    ├── pkg            ------Sophie项目内部通用库
│    ├── schedule
│    └── system
├── pkg                 ------通用库
│    ├── app            ------app 构建框架库
│    ├── core           
│    ├── db             ------存储层交互相关库
│    ├── ds             ------数据结构库
│    ├── errors         ------错误处理库
│    ├── log            ------Sophie-Log 日志系统
│    ├── mq             ------消息队列库
│    ├── shutdown       ------优雅关停库
│    ├── test           
│    ├── utils          ------通用工具库
│    └── validators     ------扩展validator自定义验证器
├── scripts             ------项目运行、部署相关脚本
├── templates           ------代码生成服务 的模板引擎扫描目录
└── tools               ------第三方工具
├── canal
└── README.md

