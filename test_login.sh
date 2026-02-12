#!/bin/bash

# 测试登录流程的脚本

echo "================================"
echo "测试登录流程"
echo "================================"

# 1. 测试登录接口
echo ""
echo "1. 测试登录接口..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:10131/api/auth/login/h5 \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800138000","password":"123456"}')

echo "登录响应:"
echo "$LOGIN_RESPONSE" | python3 -m json.tool

# 提取 accessToken
ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['accessToken'])" 2>/dev/null)

if [ -z "$ACCESS_TOKEN" ]; then
  echo "❌ 登录失败，无法获取 accessToken"
  echo "请确保："
  echo "1. 后端服务正在运行 (端口 10131)"
  echo "2. 数据库中存在测试用户"
  echo "3. 密码使用 MD5 加密"
  exit 1
fi

echo ""
echo "✅ 登录成功，AccessToken: ${ACCESS_TOKEN:0:20}..."

# 2. 测试获取用户信息接口
echo ""
echo "2. 测试获取用户信息接口..."
USER_INFO_RESPONSE=$(curl -s -X GET http://localhost:10131/api/user/info \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "用户信息响应:"
echo "$USER_INFO_RESPONSE" | python3 -m json.tool

# 3. 测试获取权限码接口
echo ""
echo "3. 测试获取权限码接口..."
ACCESS_CODES_RESPONSE=$(curl -s -X GET http://localhost:10131/api/auth/codes \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "权限码响应:"
echo "$ACCESS_CODES_RESPONSE" | python3 -m json.tool

echo ""
echo "================================"
echo "测试完成"
echo "================================"
