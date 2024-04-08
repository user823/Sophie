Canel 用于将mysql中的binlog 数据同步到elasticsearch中

安装参考：[官方文档](https://www.alibabacloud.com/help/zh/es/use-cases/use-canal-to-synchronize-mysql-data-to-alibaba-cloud-elasticsearch)

安装canal-server 和 canal-adapter 在tools/canal 目录下执行命令：
```bash
chmod +x install.sh
./install.sh
```

## 安装教程
1. 安装admin, server, adapter
进入canal 目录
```yaml
chmod +x install.sh
./install.sh
```
es、mysql 安装参考 Sophie/deployments/docker 部分
添加es http_ca.cert 到java证书链中
```
mkdir ~/.certs
docker cp sophie-es:/usr/share/elasticsearch/config/certs/http_ca.crt ~/.certs
sudo keytool -import -trustcacerts -keystore cacerts -storepass changeit -noprompt -alias docker_es_ca -file ~/.cert/http_ca.crt 
```
2. 准备测试环境
测试数据库: sophie
mysql、es登陆账号(所有权限):sophie 123456
测试数据表：
```sql
CREATE TABLE `test_book` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '题名',
  `isbn` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'isbn',
  `author` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '作者',
  `publisher_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '出版社名',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC
```
es 中建立index 索引
```
PUT test_book
{
  "settings": {
    "index": {
      "number_of_shards": 1,
      "number_of_replicas": 1
    }
  },
  "mappings": {
    "properties": {
      "id": {
        "type": "long"
      },
      "title": {
        "type": "text"
      },
      "isbn": {
        "type": "text"
      },
      "author": {
        "type": "text"
      },
      "publisherName": {
        "type": "text"
      }
    }
  }
}
```

3. 修改canal-admin配置
首先在数据库中执行canal-admin/conf/canal_manager.sql 创建用于管理的表，初始化管理员账号
```yaml
server:
  port: 8089
spring:
  jackson:
    date-format: yyyy-MM-dd HH:mm:ss
    time-zone: GMT+8

spring.datasource:
  # 数据库ip:port
  address: 127.0.0.1:3306 
  database: canal_manager
  username: sophie
  password: 12345678
  driver-class-name: com.mysql.jdbc.Driver
  url: jdbc:mysql://${spring.datasource.address}/${spring.datasource.database}?useUnicode=true&characterEncoding=UTF-8&useSSL=false
  hikari:
    maximum-pool-size: 30
    minimum-idle: 1

canal:
  # 用于登陆admin的账号和密码
  adminUser: admin
  # 123456 密文为 6BB4837EB74329105EE4568DDA7DC67ED2CA2AD9
  adminPasswd: 123456
```

4. 修改canal-server配置
全局配置：conf/canal.properties
局部配置（会覆盖部分全局配置）:conf/canal_local.properties
对应于数据库的实例配置 conf/example/instance.properties

修改conf/canal.properties
```yaml
# canal admin config
canal.admin.manager = 127.0.0.1:8089
canal.admin.port = 11110
canal.admin.user = admin
# 这里要对应第3步中的密文
canal.admin.passwd = 6BB4837EB74329105EE4568DDA7DC67ED2CA2AD9

# admin auto register
canal.admin.register.auto = true
canal.admin.register.cluster =
canal.admin.register.name =
```

修改 conf/canal_local.properties
```yaml
# register ip
canal.register.ip =

# canal admin config
canal.admin.manager = 127.0.0.1:8089
canal.admin.port = 11110
# canal.admin.user = admin
# canal.admin.passwd = 6bb4837eb74329105ee4568dda7dc67ed2ca2ad9
# admin auto register
canal.admin.register.auto = true
canal.admin.register.cluster =
canal.admin.register.name =
```

修改 conf/example/instance.properties
```yaml
# username/password
canal.instance.dbUsername=sophie
canal.instance.dbPassword=12345678
canal.instance.connectionCharset = UTF-8
# enable druid Decrypt database password
canal.instance.enableDruid=false
#canal.instance.pwdPublicKey=MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBALK4BUxdDltRRE5/zXpVEVPUgunvscYFtEip3pmLlhrWpacX7y7GCMo2/JM6LeHmiiNdH1FWgGCpUfircSwlWKUCAwEAAQ==
```
注意：由于要配置admin，这里修改此文件是无效的，不过可以利用此处的修改用于后续步骤创建实例

5. 修改 canal-adapter配置
目录下的es6, es7, es8, hbase 等都是用来配置sql 到目标源的映射
修改application.yml
```yaml
  srcDataSources:
    defaultDS:
      url: jdbc:mysql://127.0.0.1:3306/sophie?useUnicode=true
      username: sophie
      password: 12345678
    canalAdapters:
    - instance: example # canal instance Name or mq topic name
      groups:
        - groupId: g1
          outerAdapters:
            - name: logger
            - name: es8
              key: es
              hosts: https://127.0.0.1:9200 # 127.0.0.1:9200 for rest mode
              properties:
                mode: rest # or rest
                security.auth: sophie:12345678 #  only used for rest mode
                cluster.name: docker-cluster
```

es 集群名称获取方式: 在kibana 输入
`GET _nodes/stats`

添加es8/test_book.yml 映射
```yaml
dataSourceKey: defaultDS
destination: example
outerAdapterKey: es
groupId: g1
esMapping:
  _index: test_book
  _id: _id
  sql: "select t.id as _id, t.title, t.isbn, t.author, t.publisher_name as publisherName from test_book t"
  commitBatch: 3000
```

6. 依次启动admin, server
```yaml
./canal-admin/bin/restart.sh
./canal-server/bin/restart.sh

# 检查启动日志
tail canal-admin/logs/admin.log
tail canal-server/logs/canal/canal.log
```

7. 创建实例
在admin instance管理中选择新建 instance，然后将canal/server/conf/example/instance.properties 内容贴过来
instance 名称填example，所属集群选择自动注册的主机
实例创建好后可以在 canal-server/logs/example 下看见日志

8. 运行adapter
```yaml
./canal-adapter/bin/restart.sh

# 日志查看
tail canal-adapter/logs/adapter/adapter.log
```

9. 插入验证数据，测试
```
INSERT INTO `sophie`.`test_book`( `title`, `isbn`, `author`, `publisher_name`) VALUES (  '三体', '98741254125', '刘慈欣', '工业出版社');

# 查看adapter日志
tail canal-adapter/logs/adapter/adapter.log
```

此时应该输出:
```
2024-01-29 15:35:33.433 [pool-3-thread-1] INFO  c.a.o.canal.client.adapter.logger.LoggerAdapterExample - DML: {"data":[{"id":3,"title":"三体","isbn":"98741254125","author":"刘慈欣","publisher_name":"工业出版社"}],"database":"sophie","destination":"example","es":1706513733000,"groupId":"g1","isDdl":false,"old":null,"pkNames":["id"],"sql":"","table":"test_book","ts":1706513733430,"type":"INSERT"}
2024-01-29 15:35:33.439 [pool-3-thread-1] DEBUG c.a.o.canal.client.adapter.es.core.service.ESSyncService - DML: {"data":[{"id":3,"title":"三体","isbn":"98741254125","author":"刘慈欣","publisher_name":"工业出版社"}],"database":"sophie","destination":"example","es":1706513733000,"groupId":"g1","isDdl":false,"old":null,"pkNames":["id"],"sql":"","table":"test_book","ts":1706513733430,"type":"INSERT"} 
Affected indexes: test_book 
```

es 端查看：
`GET test_book/_search`

## 报错说明
1. 
```
Caused by: com.alibaba.otter.canal.common.CanalException: com.alibaba.otter.canal.common.CanalException: instance : example config is not found
Caused by: com.alibaba.otter.canal.common.CanalException: instance : example config is not found
```
需要在admin中配置instance

2. adapter 高级用法参考：
[对象型数组同步](https://juejin.cn/post/7093385171168133157)