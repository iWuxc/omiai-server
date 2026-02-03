#!/bin/sh
set -e

# 如果存在 config.template.yaml，则使用 envsubst 生成 config.yaml
if [ -f /app/configs/config.template.yaml ]; then
    echo "Generating config.yaml from template..."
    # 使用 envsubst 替换环境变量，只替换模板中存在的变量
    envsubst < /app/configs/config.template.yaml > /app/configs/config.yaml
    echo "Config generated."
fi

# 执行传入的命令
exec "$@"
