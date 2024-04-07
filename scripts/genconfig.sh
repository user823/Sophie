#!/usr/bin/env bash

# 脚本根据environment.sh 配置，生成sophie 组件 yaml 配置文件
# 示例: genconfig.sh scripts/environment.sh configs/sophie-gateway.yml

log::error() {
  timestamp=$(date +"[%m%d %H:%M:%S]")
  echo "!!! ${timestamp} ${1-}" >&2
  shift
  for message; do
    echo "    ${message}" >&2
  done
}

env_file="$1"
template_file="$2"

SOPHIE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

set -o errexit
set +o nounset
set -o pipefail

unset CDPATH

export GO111MODULE=on

ALLOW_EMPTY_VARS=("ETCD_USERNAME" "ETCD_PASSWORD" "REDIS_USERNAME")

if [ $# -ne 2 ];then
    log::error "Usage: genconfig.sh scripts/environment.sh configs/sophie-gateway.yml"
    exit 1
fi

source "${env_file}"

declare -A envs

set +u
for env in $(sed -n 's/^[^#].*${\(.*\)}.*/\1/p' ${template_file})
do
    # 检查是否为允许为空的环境变量，如果是，则跳过检查
    if [[ " ${ALLOW_EMPTY_VARS[@]} " =~ " ${env} " ]]; then
        continue
    fi

    if [ -z "$(eval echo \$${env})" ];then
        log::error "environment variable '${env}' not set"
        missing=true
    fi
done

if [ "${missing}" ];then
    log::error 'You may run `source scripts/environment.sh` to set these environment'
    exit 1
fi

eval "cat << EOF
$(cat ${template_file})
EOF"