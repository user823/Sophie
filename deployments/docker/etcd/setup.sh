#!/bin/bash

docker run -d --name sophie-etcd -p 2379:2379 -p 2380:2380 --env ALLOW_NONE_AUTHENTICATION=yes bitnami/etcd:latest

# 测试
# docker run -it --rm bitnami/etcd:latest etcdctl --endpoints http://10.211.55.3:2379 put /test hello

docker run -it --rm bitnami/etcd:latest etcdctl --endpoints http://10.211.55.3:2379 put /config/test.json {"test": "hh"}
docker run -it --rm bitnami/etcd:latest etcdctl --endpoints http://10.211.55.3:2379 get /config/test.json
docker run -it --rm bitnami/etcd:latest etcdctl --endpoints http://10.211.55.3:2379 delete /config/test.json