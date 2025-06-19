#!/bin/bash
# 获取脚本所在目录的上级目录（即安装目录）
BASE_DIR="$(dirname "$(dirname "$(readlink -f "$0")")")"
export PYTHONPATH="${BASE_DIR}:${BASE_DIR}/aipyapp:${BASE_DIR}/tools/python-embed-linux/lib/python3.x/site-packages"
"${BASE_DIR}/tools/python-embed-linux/bin/python3" -c "from aipyapp.__main__ import main; main()"