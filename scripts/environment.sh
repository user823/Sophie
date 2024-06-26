#!/bin/bash

# Sophie 项目源码根目录
SOPHIE_ROOT=$(pwd)

# 生成文件存放目录
LOCAL_OUTPUT_ROOT="${SOPHIE_ROOT}/${OUT_DIR:-_output}"

# 所有组件采用统一密码
readonly PASSWORD=${PASSWORD:-'12345678'}
# linux系统 sudo 用户和密码
readonly LINUX_USERNAME=${LINUX_USERNAME:-sophie}
readonly LINUX_PASSWORD=${LINUX_PASSWORD:-'12345678'}

# 设置安装目录
readonly INSTALL_DIR=${INSTALL_DIR:-/tmp/installation}
mkdir -p ${INSTALL_DIR}
readonly ENV_FILE=${SOPHIE_ROOT}/scripts/environment.sh

# mysql 配置信息
readonly MYSQL_ADMIN_USERNAME=${MYSQL_ADMIN_USERNAME:-root}
readonly MYSQL_ADMIN_PASSWORD=${MYSQL_ADMIN_PASSWORD:-${PASSWORD}}
readonly MYSQL_HOST=${MYSQL_HOST:-127.0.0.1:3306}
readonly MYSQL_DATABASE=${MYSQL_DATABASE:-sophie}
readonly MYSQL_USERNAME=${MYSQL_USERNAME:-sophie}
readonly MYSQL_PASSWORD=${MYSQL_PASSWORD:-${PASSWORD}}

# redis 配置信息
readonly REDIS_HOST=${REDIS_HOST:-127.0.0.1:6379}
readonly REDIS_USERNAME=${REDIS_USERNAME} # Redis 用户名
readonly REDIS_PASSWORD=${REDIS_PASSWORD:-${PASSWORD}} # Redis 密码

# elasticsearch 配置信息
readonly ES_ENDPOINT=${ES_ENDPOINT:-https://127.0.0.1:9200}
readonly ES_USERNAME=${ES_USERNAME:-sophie}
readonly ES_PASSWORD=${ES_PASSWORD:-${PASSWORD}}
readonly ES_APIKEY=${ES_APIKEY}
readonly ES_CLOUDID=${ES_CLOUDID}

# Etcd 配置信息
readonly ETCD_HOST=${ETCD_HOST:-127.0.0.1:2379}
readonly ETCD_USERNAME=${ETCD_USERNAME:-""}
readonly ETCD_PASSWORD=${ETCD_PASSWORD:-""}

# jaeger 配置信息
readonly JAEGER_HOST=${JAEGER_HOST:-127.0.0.1:4317}

# rocketmq 配置信息
readonly ROCKETMQ_HOST=${ROCKETMQ_HOST:-127.0.0.1:8081}
readonly ROCKETMQ_ACCESSKEY=${ROCKETMQ_ACCESSKEY:-sophie}
readonly ROCKETMQ_ACCESSSECET=${ROCKETMQ_ACCESSSECET:-${PASSWORD}}

# minio 配置信息
readonly MINIO_HOST=${MINIO_HOST:-127.0.0.1:9000}
readonly MINIO_ACCESSKEY=${MINIO_ACCESSKEY:-sophie}
readonly MINIO_ACCESSSECRET=${MINIO_ACCESSSECRET:-${PASSWORD}}

# RPC 客户端配置（保持默认即可）
readonly RPC_CLIENT_MIN_IDLE_PER_ADDRESS=${RPC_CLIENT_MIN_IDLE_PER_ADDRESS:-2} # rpc每个地址最小空闲连接（提升性能）
readonly RPC_CLIENT_MAX_IDLE_PER_ADDRESS=${RPC_CLIENT_MAX_IDLE_PER_ADDRESS:-10} # rpc每个地址的最大空闲连接
readonly RPC_CLIENT_MAX_IDLE_GLOBAL=${RPC_CLIENT_MAX_IDLE_GLOBAL:-100} # rpc全局最大空闲连接
readonly RPC_CLIENT_MAX_IDLE_TIMEOUT=${RPC_CLIENT_MAX_IDLE_TIMEOUT:-1m} # rpc最大空闲时间
readonly RPC_CLIENT_CONN_TIMEOUT=${RPC_CLIENT_CONN_TIMEOUT:-3s} # rpc连接超时时间
readonly RPC_CLIENT_RPC_TIMEOUT=${RPC_CLIENT_RPC_TIMEOUT:-3s} # rpc调用超时时间
readonly RPC_CLIENT_MAX_RETRY_TIME=${RPC_CLIENT_MAX_RETRY_TIME:-2} # rpc调用最大重试次数
readonly RPC_CLIENT_MAX_DURATION_MS=${RPC_CLIENT_MAX_DURATION_MS:-0} # 复用熔断器
readonly RPC_CLIENT_CIRCUITBREAK=${RPC_CLIENT_CIRCUITBREAK:-0.8} # 熔断率
readonly RPC_CLIENT_MINSAMPLE=${RPC_CLIENT_MINSAMPLE:-200} # 熔断恢复采样数

# Log 配置
readonly LOG_OUTPUT_PATH=${LOG_OUTPUT_PATH:-stdout} # 日志输出端
readonly LOG_ERR_OUTPUT_PATH=${LOG_ERR_OUTPUT_PATH:-stderr} # 日志错误输出
readonly LOG_LEVEL=${LOG_LEVEL:-INFO} # 日志输出级别
readonly LOG_AGGREGATION=${LOG_AGGREGATION:-false} #是否开启日志聚合

# Sophie 配置
readonly SOPHIE_DATA_DIR=${SOPHIE_DATA_DIR:-/data/sophie}
readonly SOPHIE_INSTALL_DIR=${SOPHIE_INSTALL_DIR:-/opt/sophie} # 安装文件存放目录
readonly SOPHIE_CONFIG_DIR=${SOPHIE_CONFIG_DIR:-${SOPHIE_INSTALL_DIR}/configs} # 配置文件存放目录
readonly SOPHIE_LOG_DIR=${SOPHIE_LOG_DIR:-/var/log/sophie}
readonly CA_FILE=${CA_FILE:-${SOPHIE_CONFIG_DIR}/cert/ca.pem} #CA 证书

# Sophie Gateway 配置 (ip 从8082开始分配)
readonly SOPHIE_GATEWAY_INSECURE_ADDRESS=${SOPHIE_GATEWAY_INSECURE_ADDRESS:-127.0.0.1}
readonly SOPHIE_GATEWAY_INSECURE_PORT=${SOPHIE_GATEWAY_INSECURE_PORT:-8082}
readonly SOPHIE_GATEWAY_SECURE_ADDRESS=${SOPHIE_GATEWAY_SECURE_ADDRESS:-127.0.0.1}
readonly SOPHIE_GATEWAY_SECURE_PORT=${SOPHIE_GATEWAY_SECURE_PORT:-8083}
readonly SOPHIE_GATEWAY_SECURE_CERT_FILE=${SOPHIE_GATEWAY_SECURE_CERT_FILE:-${SOPHIE_CONFIG_DIR}/cert/sophie-gateway.pem}
readonly SOPHIE_GATEWAY_SECURE_PK_FILE=${SOPHIE_GATEWAY_SECURE_PK_FILE:-${SOPHIE_CONFIG_DIR}/cert/sophie-gateway-key.pem}

# Sophie System 配置
readonly SOPHIE_SYSTEM_ADDRESS=${SOPHIE_SYSTEM_ADDRESS:-127.0.0.1}
readonly SOPHIE_SYSTEM_PORT=${SOPHIE_SYSTEM_PORT:-8084}
readonly SOPHIE_SYSTEM_MYSQL_DATABASE=${SOPHIE_SYSTEM_MYSQL_DATABASE:-sophie}
readonly SOPHIE_SYSTEM_QPS=${SOPHIE_SYSTEM_QPS:-200}
readonly SOPHIE_SYSTEM_CONN_LIMIT=${SOPHIE_SYSTEM_CONN_LIMIT:-1000}

# Sophie Schedule 配置
readonly SOPHIE_SCHEDULE_ADDRESS=${SOPHIE_SCHEDULE_ADDRESS:-127.0.0.1}
readonly SOPHIE_SCHEDULE_PORT=${SOPHIE_SCHEDULE_PORT:-8085}
readonly SOPHIE_SCHEDULE_QPS=${SOPHIE_SCHEDULE_QPS:-100}
readonly SOPHIE_SCHEDULE_CONN_LIMIT=${SOPHIE_SCHEDULE_CONN_LIMIT:-20}

# Sophie Worker 配置
readonly SOPHIE_SCHEDULE_WORKER_ADDRESS=${SOPHIE_SCHEDULE_WORKER_ADDRESS:-127.0.0.1}
readonly SOPHIE_SCHEDULE_WORKER_PORT=${SOPHIE_SCHEDULE_WORKER_PORT:-8091}
readonly SOPHIE_SCHEDULE_WORKER_QPS=${SOPHIE_SCHEDULE_WORKER_QPS:-100}
readonly SOPHIE_SCHEDULE_WORKER_CONN_LIMIT=${SOPHIE_SCHEDULE_WORKER_CONN_LIMIT:-20}

# Sophie Gen 配置
readonly SOPHIE_GEN_ADDRESS=${SOPHIE_GEN_ADDRESS:-127.0.0.1}
readonly SOPHIE_GEN_PORT=${SOPHIE_GEN_PORT:-8086}
readonly SOPHIE_GEN_QPS=${SOPHIE_GEN_QPS:-200}
readonly SOPHIE_GEN_CONN_LIMIT=${SOPHIE_GEN_CONN_LIMIT:-1000}
readonly SOPHIE_GEN_TEMPLATE_PATH=${SOPHIE_GEN_TEMPLATE_PATH:-$SOPHIE_ROOT/templates}}
readonly SOPHIE_GEN_AUTHOR=${SOPHIE_GEN_AUTHOR:-sophie}

# Sophie File 配置
readonly SOPHIE_FILE_ADDRESS=${SOPHIE_FILE_ADDRESS:-127.0.0.1}
readonly SOPHIE_FILE_PORT=${SOPHIE_FILE_PORT:-8087}
readonly SOPHIE_FILE_QPS=${SOPHIE_FILE_QPS:-200}
readonly SOPHIE_FILE_CONN_LIMIT=${SOPHIE_FILE_CONN_LIMIT:-1000}

# Sophie Logstash 配置
readonly SOPHIE_LOGSTASH_INSECURE_ADDRESS=${SOPHIE_LOGSTASH_INSECURE_ADDRESS:-127.0.0.1}
readonly SOPHIE_LOGSTASH_INSECURE_PORT=${SOPHIE_LOGSTASH_INSECURE_PORT:-8088}