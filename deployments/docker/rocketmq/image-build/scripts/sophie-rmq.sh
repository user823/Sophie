#!/bin/bash
set -e

# 本脚本直接在一个容器中同时执行namesrv 和 broker、proxy
# 需要暴露9876、10911、8081等端口

nohup bash mqnamesrv &
bash mqbroker -n localhost:9876 --enable-proxy