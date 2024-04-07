#!/bin/bash

if [ -z "$1" ]; then
  echo "Must provide target app"
  exit 1
fi

# 如果是直接可以执行的命令则执行
if which "$1" > /dev/null 2>&1; then
  exec "$@"
  exit 0
fi

# 检测templates下是否存在对应的yml文件
template_file="configs/templates/$1.yml"
if [ ! -f "$template_file" ]; then
    echo "Template file not found: $template_file"
    exit 1
fi

# 检测cmd下是否有对应的目录
cmd_file="bin/$1"
if [ ! -f "$cmd_file" ]; then
    echo "Target bin not found: $cmd_file"
    exit 1
fi

# 生成配置文件
bash scripts/genconfig.sh scripts/environment.sh "$template_file" > "configs/$1.yml"

shift

# 运行应用程序
./${cmd_file} "$@"