#!/bin/bash

# ========================================
# Let's Encrypt SSL 证书自动配置脚本
# ========================================

set -e

DOMAIN="www.omiai.cn"
EMAIL="your-email@example.com"  # 请修改为你的邮箱

echo "=== 开始配置 Let's Encrypt SSL 证书 ==="

# 1. 安装 Certbot
echo ">>> 安装 Certbot..."
if command -v yum &> /dev/null; then
    yum install -y certbot python3-certbot-nginx
elif command -v apt &> /dev/null; then
    apt update
    apt install -y certbot python3-certbot-nginx
fi

# 2. 创建 ACME 验证目录
echo ">>> 创建 ACME 验证目录..."
mkdir -p /var/www/certbot
chown -R nginx:nginx /var/www/certbot 2>/dev/null || chown -R www-data:www-data /var/www/certbot 2>/dev/null || true

# 3. 临时修改 Nginx 配置以支持 ACME 验证
echo ">>> 配置 Nginx 支持 ACME 验证..."
cat > /etc/nginx/conf.d/omiai.conf << 'EOF'
server {
    listen 80;
    server_name www.omiai.cn;
    
    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }
    
    location / {
        return 301 https://$host$request_uri;
    }
}
EOF

# 4. 测试并重载 Nginx
echo ">>> 重载 Nginx..."
nginx -t && nginx -s reload

# 5. 申请证书
echo ">>> 申请 SSL 证书..."
certbot certonly --webroot \
    -w /var/www/certbot \
    -d $DOMAIN \
    --email $EMAIL \
    --agree-tos \
    --no-eff-email

# 6. 设置自动续期
echo ">>> 配置自动续期..."
# Certbot 会自动创建 systemd timer，检查一下
systemctl enable certbot.timer
systemctl start certbot.timer

# 7. 测试续期
echo ">>> 测试自动续期..."
certbot renew --dry-run

echo ""
echo "=== SSL 证书配置完成 ==="
echo ""
echo "证书文件位置："
echo "  - 证书: /etc/letsencrypt/live/$DOMAIN/fullchain.pem"
echo "  - 私钥: /etc/letsencrypt/live/$DOMAIN/privkey.pem"
echo ""
echo "自动续期状态："
systemctl status certbot.timer --no-pager
echo ""
echo "请更新 Nginx 配置文件以使用新证书，然后执行："
echo "  nginx -t && nginx -s reload"
