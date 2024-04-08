# Sophie 部署指南

二进制部署方式, 本项目主要采用[docker部署方式](../../deployments/README.md)

## 环境准备

部署前需要的准备工作：
1. 提供https 服务需要用到的公钥和私钥（pem）
2. mysql
3. etcd
4. redis
5. jaeger（可选）, 用于链路追踪 和 集成promethus的服务监控
6. minio(可选), 用于对象存储，如果未部署可替换华为云、阿里云等其他对象存储
7. elasticsearch（可选), 用于日志聚合输出端，如果未部署时日志聚合会将日志收集到控制台
8. rocketmq （可选）, 用于日志聚合的中间管道
9. canal（可选），在使用elasticsearch来查询登录日志和操作日志时，开启canal来同步数据到es中；如果未开启canal则不要使用elasticsearch，避免数据不一致

说明：
1. 所有组件的默认用户名：sophie， 密码：12345678

### mysql 部署
mysql 开放端口3306（默认），mysql 初始化时指定字符集 --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
部署后创建远程登录账户 sophie，默认密码: 12345678
在登录mysql时指定--default-character-set=utf8mb4 ，执行configs 中的sophie.sql，data.sql 分别创建项目运行需要的库表，并插入初始数据

### etcd 部署
etcd 开放端口2379 和 2380
指定环境变量 ALLOW_NONE_AUTHENTICATION=yes，开启无密码访问

### redis部署
redis 开放端口6379
redis 初始化时的配置文件使用deployments/docker/redis/conf/redis.conf，访问密码12345678

### jaeger部署
jaeger 需要开启环境变量：COLLECTOR_OTLP_ENABLED=true、COLLECTOR_ZIPKIN_HOST_PORT=:9411
jaeger 用到的端口：5775 6831 6832 5778 16686 14250 14268 14269 4317 4318 9411

### minio部署
minio 的密码长度至少要8个字符，使用端口9000 9001
minio 初始化后需要创建存储桶 sophie，并指定读写策略：
```yaml
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "AWS": [
                    "*"
                ]
            },
            "Action": [
                "s3:GetObject"
            ],
            "Resource": [
                "arn:aws:s3:::sophie/*"
            ]
        },
        {
            "Effect": "Allow",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::123456789012:user/sophie"
                ]
            },
            "Action": [
                "s3:PutObject",
                "s3:DeleteObject",
                "s3:GetObject",
                "s3:ListBucket"
            ],
            "Resource": [
                "arn:aws:s3:::sophie",
                "arn:aws:s3:::sophie/*"
            ]
        }
    ]
}
```

### elasticsearch部署
在使用高版本的rocketmq时出了需要部署nameserver 和 broker 以外还需要部署proxy， 但是proxy和broker可以属于同一个节点
elasticsearch 部署后执行configs/es-init.sh 初始化es用到的索引

### canal部署
详情见[tools/canal/README.md](../../tools/canal/README.md)

### 公钥和私匙
它们用于配置sophie-gateway 的https服务，这里介绍自签名证书的配置方式
1. 安装 cfssl 工具集

```
cfssl: 证书签发工具
cfssljson: 将cfssl生成的证书（json格式）变为文件承载式证书
```

执行命令:

```
cd ${SOPHIE_ROOT}/configs
mkdir cert
cfssl gencert -initca cfssl/ca-csr.json | cfssljson -bare ca
ls ca*
mv ca* cert
```

生成 ca-key.pem(私钥) 和 ca.pem(公钥)

2. 修改 hosts

```
sudo tee -a /etc/hosts << EOF
127.0.0.1 sophie.gateway.com
127.0.0.1 sophie.system.com
EOF
```

3. 配置 gateway

```
cd $SOPHIE_ROOT/configs
sudo mkdir /var/run/sophie
cfssl gencert -ca=cert/ca.pem -ca-key=cert/ca-key.pem -config=cfssl/ca-config.json \
-profile=sophie cfssl/sophie-gateway-csr.json | cfssljson -bare sophie-gateway
sudo cp sophie-gateway*pem /var/run/sophie/
mv sophie-gateway*pem cert
mv sophie-gateway*csr cert
```

## 组件部署
本后台管理系统可部署的组件包括：
sophie-gateway（网关）
sophie-system（系统核心服务）
sophie-file (文件服务)
sophie-gen (代码生成服务)
sophie-schedule（定时任务-管理节点）
sophie-schedule-worker（定时任务-负载节点）
sophie-logstash（日志聚合系统组件）
它们之间的关系见[sophie 系统架构说明](../devel/architecture.md)

部署主要使用scripts中的脚本, 所有脚本都在项目根目录下执行
```
# 编译所有组件
bash scripts/build.sh 
```

environment.sh 统一管理设置所有组件的环境变量，并提供了默认值，可以直接在命令行中设置对应的环境变量后执行脚本覆盖默认值
```
bash scripts/environment.sh
```
所有组件的启动都需要读取配置文件（组件有三种配置方式，见[app 启动流程](./app.md), 配置文件使用configs/templates中的yml模板生成
```
# 以生成sophie-gateway.yml 文件为例
bash scripts/genconfig.sh scripts/environment.sh configs/templates/sophie-gateway.yml > configs/sophie-gateway.yml
```

