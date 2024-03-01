#!/bin/bash

# 设置环境变量脚本
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SOPHIE_ROOT="$(dirname "$SCRIPT_DIR")"
SOPHIE_CONFIG_DIR="${SOPHIE_ROOT}/configs"

export SOPHIE_ROOT
export SOPHIE_CONFIG_DIR