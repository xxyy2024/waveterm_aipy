#!/bin/bash

# 获取脚本所在目录的上级目录（即安装目录）
BASE_DIR="$(dirname "$(dirname "$(readlink -f "$0")")")"

# 检测操作系统和架构
OS_TYPE="$(uname)"
ARCH_TYPE="$(uname -m)"

if [[ "$OS_TYPE" == "Darwin" ]]; then
  APP_DIR="tools/prompt-optimizer/mac"
elif [[ "$OS_TYPE" == "Linux" ]]; then
  # Linux 平台
  APP_DIR="tools/prompt-optimizer/linux"
else
 # Windows 平台（默认分支）
  echo "当前不支持 Windows 平台，请使用 macOS 或 Linux"
  exit 1
fi

mkdir -p ~/.prompt-optimizer
if [ ! -f ~/.prompt-optimizer/config.yaml ]; then
    cp "${BASE_DIR}/tools/prompt_config.yaml" ~/.prompt-optimizer/config.yaml
fi


echo "Start Prompt Optimizer..."
# 启动 Python 应用
"${BASE_DIR}/${APP_DIR}/app"  interactive -c  ~/.prompt-optimizer/config.yaml
