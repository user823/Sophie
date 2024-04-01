项目运行环境Docker部署

注意：
1. mysql先运行sophie.sql 脚本建立数据库和表，然后运行data.sql 初始化数据。 mysql 服务器默认使用utfmb4字符集，导入数据前确保客户端连接字符集也是utfmb4
连接时通过 --default-character-set=utf8mb4设置客户端连接字符集
查看方式：
```
SHOW VARIABLES LIKE 'character_set_client';
```

2. 截止到(2024.03.28) apache 官方的rocketmq 只有x86架构linux 镜像。为了使用docker在arm架构的mac上部署rocketmq，在rocketmq中提供了自定义的rocketmq
镜像进行build，但是dashboard镜像没有通过测试，后续会进行修复。

3. 使用rocketmq 之前确保topic 建立好，topic 建立命令:
```
export ip=xxxxx
docker exec --rm sophie-rmq bash mqadmin updatetopic -n ${ip}:9876 -t sophie_record_aggregation -c DefaultCluster
```