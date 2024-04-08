# Sophie 部署指南

docker 部署方式
deployments 中提供了项目运行依赖组件的docker部署、项目本身运行的镜像构建docker file 以及 开启服务的docker compose

## sophie 镜像使用方式
首先build 项目, 然后docker build
```
# 根目录下
bash scripts/build.sh
docker build -t sophie -f deployments/docker/dockerfile .
```

该镜像包含所有组件，运行单个组件:
```
docker run -d sophie:latest sophie-gateway(这里是组件名）
```
组件名支持：
sophie-gateway（网关）
sophie-system（系统核心服务）
sophie-file (文件服务)
sophie-gen (代码生成服务)
sophie-schedule（定时任务-管理节点）
sophie-schedule-worker（定时任务-负载节点）
sophie-logstash（日志聚合系统组件）

## 不使用docker-compose, 单容器部署
### mysql部署
直接进入docker/mysql 执行setup.sh。容器启动后将configs/sophie.sql、configs/data.sql 拷贝到mysql/data下，创建sophie用户
```
# 进入容器
docker exec -it sophie-mysql bash
cd /var/lib/mysql
# 创建用户...
# 初始化库表
source sophie.sql;
source data.sql;
```

### redis 部署
直接进入redis执行setup.sh 脚本

### minio 部署
进入minio执行setup.sh 脚本后，创建存储桶sophie，并配置读写策略
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

### jaeger 部署
直接进入jaeger 目录，执行setup.sh

### etcd 部署
直接进入etcd 目录，执行setup.sh

### elasticsearch 部署
直接进入etcd 目录，执行setup.sh
在两个容器启动后执行命令， 初始化elastic 和 kibana的密码
```
docker exec -it sophie-es /usr/share/elasticsearch/bin/elasticsearch-reset-password -u elastic
# 它会生成kibana 的token
docker exec -it sophie-es /usr/share/elasticsearch/bin/elasticsearch-create-enrollment-token -s kibana --url https://localhost:9200
```
然后去浏览器端访问 [kibana](http://localhost:5601)，粘贴token，进入后创建用户sophie, 密码12345678
执行configs/es-init.sh 初始化索引

### rocketmq 部署
首先创建rocketmq镜像，进入image-build， 里面包含rocketmq 和 dashboard（用于管理rmq）的镜像制作方式，dashboard 目前存在问题后续会修复
执行：
```
bash build-image.sh 5.1.4
```
执行后调用setup.sh 启用容器， setup.sh 提供两种部署方式：
```
# 1. 单nameserver， 单broker + proxy
# 2. nameserver、broker、proxy都在一个容器里面
bash setup.sh -q
```

说明：在arm 机器上apache/rocketmq没有提供官方镜像，因项目提供了rocketmq的镜像制作方式


### docker compose 部署方式
进入到deployments/docker 目录下，首先启动基础组件
```
docker compose up -d mysql etcd redis jaeger minio elasticsearch rocketmq
```
基础组件完成后初始化，参考单容器部署方式

然后启动项目组件, 其中sophie-gateway sophie-system 必不可少，其他可选
```
docker compose up -d sophie-gateway sophie-system sophie-file sophie-gen sophie-schedule sophie-sched schedule-worker sophie-logstash
```

## 腾讯云部署实践
将所有服务部署到 4核 8g的轻量级服务器上, 系统: ubuntu 20.04
1. 首先编译[Sophie-ui](https://github.com/user823/Sophie-ui), 得到静态资源文件夹dist
```
npm run build:prod
```
放到$HOME下
2. 安装nginx
```
sudo apt-get upgrade
sudo apt-get install nginx 
```
nginx 默认安装1.8 版本，nginx 配置文件在/etc/nginx下，修改nginx.conf:
```
user www-data;
worker_processes auto;
pid /run/nginx.pid;
include /etc/nginx/modules-enabled/*.conf;

events {
        worker_connections 768;
        # multi_accept on;
}

http {
        
        # 这里配置http向https自动跳转
        server {
            listen 80;
            server_name www.my-sophie.love;
            return 301 https://$host$request_uri;
        }

        server {
            listen 443 ssl;
            server_name www.my-sophie.love;
            ssl_certificate /home/sophie/.cert/my-sophie.love_bundle.crt;
            ssl_certificate_key /home/sophie/.cert/my-sophie.love.key;
            ssl_session_timeout 5m;
            ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE:ECDH:AES:HIGH:!NULL:!aNULL:!MD5:!ADH:!RC4;
            ssl_protocols TLSv1.2 TLSv1.3;
            ssl_prefer_server_ciphers on;

            
            # 这里设置为得到的编译文件 作为web站点
            location / {
                root ${HOME}/dist;
                index index.html index.htm;
                try_files $uri $uri/ /prod-api/$uri /index.html;
            }
            
            # 这里配置反向代理
            location /prod-api {
                rewrite ^/prod-api(.*)$ $1 break;
                proxy_pass http://localhost:8082;
                proxy_set_header Host $host;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            }
        }
        ##
        # Basic Settings
        ##

        sendfile on;
        tcp_nopush on;
        tcp_nodelay on;
        keepalive_timeout 65;
        types_hash_max_size 2048;
        # server_tokens off;

        # server_names_hash_bucket_size 64;
        # server_name_in_redirect off;

        include /etc/nginx/mime.types;
        default_type application/octet-stream;

        ##
        # SSL Settings
        ##

        ssl_protocols TLSv1 TLSv1.1 TLSv1.2 TLSv1.3; # Dropping SSLv3, ref: POODLE
        ssl_prefer_server_ciphers on;

        ##
        # Logging Settings
        ##

        access_log /var/log/nginx/access.log;
        error_log /var/log/nginx/error.log;

        ##
        # Gzip Settings
        ##

        gzip on;

        # gzip_vary on;
        # gzip_proxied any;
        # gzip_comp_level 6;
        # gzip_buffers 16 8k;
        # gzip_http_version 1.1;
        # gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;

        ##
        # Virtual Host Configs
        ##

        include /etc/nginx/conf.d/*.conf;
        include /etc/nginx/sites-enabled/*;
}
```

3. 参考docker-compose 部署方式部署sophie 应用
4. 浏览器访问 https://xxxx 正常情况下可看到sophie登录页面