#!/bin/sh
set -e

# 确保 configs 目录存在
mkdir -p /app/configs

# 如果存在 config.template.yaml，则使用 envsubst 生成 config.yaml
if [ -f /app/configs/config.template.yaml ]; then
    echo "Found template, generating config.yaml..."
    # 只替换在模板中显式定义的变量，避免误伤其他 ${}
    # 这里为了简单，我们假设模板里只有环境变量占位符
    envsubst < /app/configs/config.template.yaml > /app/configs/config.yaml
    
    if [ -s /app/configs/config.yaml ]; then
        echo "Config generated successfully."
        # cat /app/configs/config.yaml # 调试时可以打开，但注意不要泄露密码
    else
        echo "Error: Generated config.yaml is empty!"
        exit 1
    fi
else
    echo "Warning: config.template.yaml not found at /app/configs/config.template.yaml"
fi

# 执行传入的命令
echo "Starting application..."
exec "$@"
