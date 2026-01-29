#!/bin/bash

code_dir=$(pwd)

# 判断文件执行位置
if [ ! -d "${code_dir}/scripts" ]; then
  echo "请在项目目录下执行此文件"
  exit
fi

bash ./scripts/run.sh stop
