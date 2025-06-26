#!/bin/bash

# 获取脚本所在目录的上级目录（即安装目录）
BASE_DIR="$(dirname "$(dirname "$(readlink -f "$0")")")"

# 检测操作系统和架构
OS_TYPE="$(uname)"
ARCH_TYPE="$(uname -m)"

if [[ "$OS_TYPE" == "Darwin" ]]; then
  if [[ "$ARCH_TYPE" == "arm64" ]]; then
    # macOS M1/M2/M3
    PYTHON_DIR="tools/python-embed-mac-arm64"
  else
    # macOS Intel
    PYTHON_DIR="tools/python-embed-mac-x86_64"
  fi
elif [[ "$OS_TYPE" == "Linux" ]]; then
  # Linux 平台
  PYTHON_DIR="tools/python-embed-linux"
else
  # window 平台
  PYTHON_DIR="tools/python-embed"
fi

# 设置 PYTHONPATH
export PYTHONPATH="${BASE_DIR}/${PYTHON_DIR}/lib/python3.x/site-packages"
export PYTHONHOME="${BASE_DIR}/${PYTHON_DIR}"

echo "Start AiPy..."
# 启动 Python 应用
"${BASE_DIR}/${PYTHON_DIR}/bin/python3" -c "from aipyapp.__main__ import main; main()"
