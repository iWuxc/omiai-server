#!/bin/bash

dir=$(pwd)
bin=bin/server
code_dir=$(pwd)
bak_time=$(date "+%m%d-%H.%M.%S")
go_bin=go

# 判断文件执行位置
if [ -f "${dir}/scripts/run.sh" ]; then
  code_dir=${dir}
elif [ -f "${dir}/run.sh" ]; then
  code_dir=$(dirname "${dir}")
else
  echo "请在项目目录下执行此文件"
  exit
fi

getpid() {
  # shellcheck disable=SC2009
  # shellcheck disable=SC2155
  local run_pid=$(ps -ef | grep "${code_dir}"/"${bin}" | grep -v "grep" | awk '{ print $2 }')
  if [ "${run_pid}" == "" ]; then
    echo "0"
  fi
  echo "${run_pid}"
}

status() {
  # 检测运行状态
  # shellcheck disable=SC2155
  local pid=$(getpid)
  if [ "${pid}" != "0" ]; then
    echo "程序正在运行中, PID：${pid}"
  else
    echo "程序已停止运行"
  fi
}

start() {
  # shellcheck disable=SC2155
  local pid=$(getpid)
  if [ "${pid}" != "0" ]; then
    echo "程序正在运行中, PID：$(getpid)"
  else
    # 删除旧的编译包
    if [ -e "${code_dir}/${bin}" ]; then
      rm -f "${code_dir}"/${bin}
    fi

    # 打包
    echo "正在打包程序..."
    ${go_bin} build -o ./bin/ "${code_dir}"/...

    # 检测是否打包失败
    if [ ! -f "${code_dir}/${bin}" ]; then
      echo "打包失败，请手动打包"
      exit
    fi

    # 后台运行
    echo "正在启动..."
    nohup "${code_dir}"/${bin} -conf="${code_dir}"/configs >"${code_dir}"/runtime/output.log 2>&1 &

    # 等待写入PID到文件
    while :; do
      if [ "$(getpid)" != "0" ]; then
        echo "Done，PID：$(getpid)"
        break
      fi
    done

  fi
}

restart() {
  # shellcheck disable=SC2155
  local pid=$(getpid)
  if [ "${pid}" != "0" ]; then
    # 删除旧的编译包
    if [ -e "${code_dir}/${bin}" ]; then
      rm -f "${code_dir}"/${bin}
    fi

    # 重新编译
    ${go_bin} build -o ./bin/ "${code_dir}"/...

    echo "正在重启..."
    bash -c "kill -TERM ${pid}"
    # 备份旧的运行日志
    if [ -s "${code_dir}"/runtime/output.log ]; then
      echo "正在备份运行日志..."
      cp "${code_dir}"/runtime/output.log "${code_dir}"/runtime/logs/output_"${bak_time}".log
      echo >"${code_dir}"/runtime/output.log
    fi
    start
    sleep 1
    # 等待执行完成
    while :; do
      # shellcheck disable=SC2155
      local new_pid=$(getpid)
      if [ "${new_pid}" != "${pid}" ]; then
        echo "Done，PID：${new_pid}"
        break
      fi
    done

  else
    echo "程序不在运行状态, 正在启动中..."
    start
  fi
}

stop() {
  # shellcheck disable=SC2155
  local pid=$(getpid)
  if [ "${pid}" != "0" ]; then
    bash -c "kill -TERM ${pid}"
  fi
  echo "Done."
}

run() {
  case $1 in
  pid)
    getpid
    ;;
  status)
    status
    ;;
  start)
    start
    ;;
  restart)
    restart
    ;;
  stop)
    stop
    ;;
  *)
    echo "status|start|restart|stop"
    ;;
  esac
  return 0
}
run "$1"
