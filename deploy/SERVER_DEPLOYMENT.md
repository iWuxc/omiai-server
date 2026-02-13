# Omiai 服务器完整部署文档

## 目录

1. [服务器环境要求](#1-服务器环境要求)
2. [域名与DNS配置](#2-域名与dns配置)
3. [SSL证书配置](#3-ssl证书配置)
4. [Docker环境安装](#4-docker环境安装)
5. [项目部署](#5-项目部署)
6. [Nginx配置](#6-nginx配置)
7. [GitHub Actions CI/CD](#7-github-actions-cicd)
8. [常用运维命令](#8-常用运维命令)
9. [故障排查](#9-故障排查)

---

## 1. 服务器环境要求

### 1.1 硬件要求

| 项目 | 最低配置 | 推荐配置 |
|------|---------|---------|
| CPU | 2核 | 4核+ |
| 内存 | 4GB | 8GB+ |
| 硬盘 | 40GB | 100GB+ SSD |

### 1.2 软件要求

| 软件 | 版本 |
|------|------|
| 操作系统 | CentOS 7+ / Ubuntu 20.04+ |
| Docker | 24.0+ |
| Docker Compose | 2.0+ |
| Nginx | 1.20+ |
| Git | 2.0+ |

### 1.3 端口开放

| 端口 | 用途 |
|------|------|
| 22 | SSH |
| 80 | HTTP |
| 443 | HTTPS |
| 10131-10133 | 后端服务（容器内部） |
| 10080 | H5前端（容器内部） |
| 10081 | Admin后台（容器内部） |
| 13018 | MySQL（可选外部访问） |
| 6379 | Redis（可选外部访问） |

---

## 2. 域名与DNS配置

### 2.1 域名解析

在域名服务商处添加 A 记录：

| 主机记录 | 记录类型 | 记录值 |
|---------|---------|--------|
| www | A | 服务器IP |
| @ | A | 服务器IP |

### 2.2 访问地址规划

| URL | 服务 |
|-----|------|
| `https://www.omiai.cn/` | 重定向到 /h5 |
| `https://www.omiai.cn/h5` | H5 前端 |
| `https://www.omiai.cn/web` | Admin 后台 |
| `https://www.omiai.cn/api` | 后端 API |

---

## 3. SSL证书配置

### 3.1 使用 Let's Encrypt（推荐）

Let's Encrypt 提供免费SSL证书，支持自动续期。

#### 安装 Certbot

```bash
# CentOS
yum install -y certbot python3-certbot-nginx

# Ubuntu/Debian
apt update
apt install -y certbot python3-certbot-nginx
```

#### 创建验证目录

```bash
mkdir -p /var/www/certbot
chown -R nginx:nginx /var/www/certbot
```

#### 申请证书

```bash
# 临时配置 Nginx 支持 ACME 验证
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

nginx -t && nginx -s reload

# 申请证书
certbot certonly --webroot \
    -w /var/www/certbot \
    -d www.omiai.cn \
    --email your-email@example.com \
    --agree-tos \
    --no-eff-email
```

#### 配置自动续期

```bash
# 启用自动续期定时器
systemctl enable certbot.timer
systemctl start certbot.timer

# 测试续期
certbot renew --dry-run
```

#### 证书文件位置

```
/etc/letsencrypt/live/www.omiai.cn/
├── fullchain.pem  # 证书链
├── privkey.pem    # 私钥
├── cert.pem       # 证书
└── chain.pem      # 中间证书
```

### 3.2 使用自定义证书（备选）

如果使用阿里云/腾讯云等免费证书：

```bash
# 创建证书目录
mkdir -p /etc/nginx/ssl

# 上传证书文件
# scp www.omiai.cn_bundle.crt root@server:/etc/nginx/ssl/
# scp www.omiai.cn.key root@server:/etc/nginx/ssl/

# 设置权限
chmod 600 /etc/nginx/ssl/www.omiai.cn.key
```

---

## 4. Docker环境安装

### 4.1 安装 Docker

```bash
# CentOS
yum install -y yum-utils
yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
yum install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Ubuntu
apt update
apt install -y ca-certificates curl gnupg
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
chmod a+r /etc/apt/keyrings/docker.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
apt update
apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# 启动 Docker
systemctl enable docker
systemctl start docker
```

### 4.2 配置 Docker 镜像加速（腾讯云内网）

```bash
mkdir -p /etc/docker
cat > /etc/docker/daemon.json << 'EOF'
{
    "registry-mirrors": ["https://mirror.ccs.tencentyun.com"],
    "log-driver": "json-file",
    "log-opts": {
        "max-size": "10m",
        "max-file": "3"
    }
}
EOF

systemctl daemon-reload
systemctl restart docker
```

### 4.3 登录 GitHub Container Registry

```bash
# 创建 Personal Access Token (需要 read:packages 权限)
# GitHub -> Settings -> Developer settings -> Personal access tokens

echo YOUR_GITHUB_TOKEN | docker login ghcr.io -u YOUR_GITHUB_USERNAME --password-stdin
```

---

## 5. 项目部署

### 5.1 目录结构

```
/data/omiai-server/
├── deploy/
│   ├── docker-compose.prod.yml
│   └── .env
├── doc/sql/
│   └── omiai.sql
├── runtime/
│   └── logs/
└── nginx/
    └── ssl/
```

### 5.2 创建目录

```bash
mkdir -p /data/omiai-server/{deploy,runtime/logs,nginx/ssl}
cd /data/omiai-server
```

### 5.3 创建环境变量文件

```bash
cat > deploy/.env << 'EOF'
# 数据库配置
DB_HOST=mysql
DB_PORT=3306
DB_USER=root
DB_PASSWORD=YOUR_DB_PASSWORD
DB_NAME=omiai

# Redis配置
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=YOUR_REDIS_PASSWORD

# AI配置
ZHIPUAI_API_KEY=YOUR_ZHIPUAI_API_KEY

# 域名配置
DOMAIN_H5=https://www.omiai.cn/h5

# 存储配置
STORAGE_DRIVER=cos

# 腾讯云COS配置
COS_BUCKET_URL=https://your-bucket.cos.ap-guangzhou.myqcloud.com
COS_REGION=ap-guangzhou
COS_SECRET_ID=YOUR_COS_SECRET_ID
COS_SECRET_KEY=YOUR_COS_SECRET_KEY
EOF

chmod 600 deploy/.env
```

### 5.4 创建 Docker Compose 配置

```bash
cat > deploy/docker-compose.prod.yml << 'EOF'
version: '3.8'

services:
  server:
    image: ghcr.io/iwuxc/omiai-server:latest
    container_name: omiai-server
    restart: always
    ports:
      - "10131:10131"
      - "10132:10132"
      - "10133:10133"
    environment:
      - APP_ENV=prod
      - DEBUG=false
      - ENABLE_CRON=true
      - LOG_LEVEL=info
      - STORAGE_DRIVER=${STORAGE_DRIVER}
      - TZ=Asia/Shanghai
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=root
      - DB_DEBUG=false
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=omiai
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - ZHIPUAI_API_KEY=${ZHIPUAI_API_KEY}
      - ZHIPUAI_MODEL=glm-4-flash
      - DOMAIN_H5=${DOMAIN_H5}
      - COS_BUCKET_URL=${COS_BUCKET_URL}
      - COS_REGION=${COS_REGION}
      - COS_SECRET_ID=${COS_SECRET_ID}
      - COS_SECRET_KEY=${COS_SECRET_KEY}
    depends_on:
      - mysql
      - redis
    volumes:
      - ../runtime/logs:/app/runtime/logs
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  frontend:
    image: ghcr.io/iwuxc/omiai-server-web:latest
    container_name: omiai-frontend
    restart: always
    ports:
      - "10080:80"
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  admin:
    image: ghcr.io/iwuxc/omiai-server-admin:latest
    container_name: omiai-admin
    restart: always
    ports:
      - "10081:80"
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  mysql:
    image: mysql:8.0
    container_name: omiai-mysql
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_PASSWORD}
      MYSQL_DATABASE: omiai
    ports:
      - "13018:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ../doc/sql/omiai.sql:/docker-entrypoint-initdb.d/init.sql
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  redis:
    image: redis:6.2
    container_name: omiai-redis
    restart: always
    command: redis-server --requirepass ${REDIS_PASSWORD}
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

volumes:
  mysql_data:
  redis_data:
EOF
```

### 5.5 拉取镜像并启动

```bash
cd /data/omiai-server/deploy

# 拉取最新镜像
docker pull ghcr.io/iwuxc/omiai-server:latest
docker pull ghcr.io/iwuxc/omiai-server-web:latest
docker pull ghcr.io/iwuxc/omiai-server-admin:latest

# 启动服务
docker compose -f docker-compose.prod.yml up -d

# 查看状态
docker compose -f docker-compose.prod.yml ps
```

---

## 6. Nginx配置

### 6.1 安装 Nginx

```bash
# CentOS
yum install -y nginx

# Ubuntu
apt install -y nginx

# 启动
systemctl enable nginx
systemctl start nginx
```

### 6.2 配置文件

```bash
cat > /etc/nginx/conf.d/omiai.conf << 'EOF'
# HTTP Server - 处理 ACME 验证和重定向
server {
    listen 80;
    server_name www.omiai.cn;
    client_max_body_size 50m;
    
    # Let's Encrypt ACME 验证路径
    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    # 其他请求重定向到 HTTPS
    location / {
        return 301 https://$host$request_uri;
    }
}

# HTTPS Server
server {
    listen 443 ssl http2;
    server_name www.omiai.cn;
    client_max_body_size 50m;

    # Let's Encrypt SSL 证书
    ssl_certificate /etc/letsencrypt/live/www.omiai.cn/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/www.omiai.cn/privkey.pem;
    
    # SSL 安全配置
    ssl_session_timeout 1d;
    ssl_session_cache shared:SSL:50m;
    ssl_session_tickets off;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    
    # HSTS
    add_header Strict-Transport-Security "max-age=63072000" always;

    # API 请求
    location /api/ {
        add_header 'Access-Control-Allow-Origin' '*' always;
        add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS, PUT, DELETE' always;
        add_header 'Access-Control-Allow-Headers' 'DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Authorization' always;
        
        if ($request_method = 'OPTIONS') {
            add_header 'Access-Control-Max-Age' 1728000;
            add_header 'Content-Type' 'text/plain; charset=utf-8';
            add_header 'Content-Length' 0;
            return 204;
        }

        proxy_pass http://localhost:10131;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # H5 前端
    location /h5 {
        proxy_pass http://localhost:10080/h5;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Admin 后台
    location /web {
        proxy_pass http://localhost:10081/web;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 根路径重定向
    location = / {
        return 301 /h5;
    }

    location / {
        return 404;
    }
}
EOF
```

### 6.3 测试并重载

```bash
nginx -t
nginx -s reload
```

---

## 7. GitHub Actions CI/CD

### 7.1 仓库结构

```
omiiai-server (主仓库)
├── .github/workflows/deploy.yml
├── Dockerfile
└── deploy/

omiiai-miniapp (H5前端仓库)
├── .github/workflows/trigger-deploy.yml
├── Dockerfile
└── nginx.conf

omiiai-admin (Admin后台仓库)
├── .github/workflows/trigger-deploy.yml
├── Dockerfile
└── apps/web-antd/nginx.conf
```

### 7.2 必需的 GitHub Secrets

在 `omiiai-server` 仓库配置：

| Secret 名称 | 说明 |
|------------|------|
| `SERVER_HOST` | 服务器地址 |
| `SERVER_USER` | 服务器用户名 |
| `SERVER_SSH_KEY` | SSH 私钥 |
| `DB_PASSWORD` | 数据库密码 |
| `REDIS_PASSWORD` | Redis 密码 |
| `ZHIPUAI_API_KEY` | 智谱 AI API Key |
| `DOMAIN_H5` | H5 域名 |
| `STORAGE_DRIVER` | 存储驱动 (cos) |
| `COS_BUCKET_URL` | 腾讯云 COS Bucket URL |
| `COS_REGION` | 腾讯云 COS Region |
| `COS_SECRET_ID` | 腾讯云 COS Secret ID |
| `COS_SECRET_KEY` | 腾讯云 COS Secret Key |

在 `omiiai-miniapp` 和 `omiiai-admin` 仓库配置：

| Secret 名称 | 说明 |
|------------|------|
| `PAT_TOKEN` | GitHub Personal Access Token (repo, workflow 权限) |

### 7.3 部署触发

| 事件 | 触发方式 |
|------|---------|
| 推送到 omiai-server main 分支 | 自动触发完整部署 |
| 推送到 omiai-miniapp main 分支 | 通过 repository_dispatch 触发 |
| 推送到 omiai-admin main 分支 | 通过 repository_dispatch 触发 |
| 手动触发 | GitHub Actions 页面点击 "Run workflow" |

---

## 8. 常用运维命令

### 8.1 Docker 命令

```bash
# 查看容器状态
docker ps

# 查看容器日志
docker logs omiai-server
docker logs omiai-frontend
docker logs omiai-admin

# 实时查看日志
docker logs -f omiai-server

# 进入容器
docker exec -it omiai-server sh

# 重启容器
docker restart omiai-server

# 拉取最新镜像并重启
cd /data/omiai-server/deploy
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d

# 清理未使用的镜像
docker image prune -f
```

### 8.2 Nginx 命令

```bash
# 测试配置
nginx -t

# 重载配置
nginx -s reload

# 查看状态
systemctl status nginx

# 重启
systemctl restart nginx
```

### 8.3 SSL 证书命令

```bash
# 查看证书状态
certbot certificates

# 手动续期
certbot renew

# 测试续期
certbot renew --dry-run

# 查看自动续期状态
systemctl status certbot.timer
```

### 8.4 数据库命令

```bash
# 连接 MySQL
docker exec -it omiai-mysql mysql -uroot -p

# 备份数据库
docker exec omiai-mysql mysqldump -uroot -p omiai > backup.sql

# 恢复数据库
docker exec -i omiai-mysql mysql -uroot -p omiai < backup.sql
```

---

## 9. 故障排查

### 9.1 服务无法启动

```bash
# 检查容器日志
docker logs omiai-server

# 检查端口占用
netstat -tlnp | grep -E '10131|10080|10081'

# 检查磁盘空间
df -h

# 检查内存
free -m
```

### 9.2 页面 404 错误

```bash
# 检查容器内文件
docker exec omiai-frontend ls -la /usr/share/nginx/html/
docker exec omiai-admin ls -la /usr/share/nginx/html/

# 检查 Nginx 配置
docker exec omiai-frontend cat /etc/nginx/conf.d/default.conf
docker exec omiai-admin cat /etc/nginx/conf.d/default.conf
```

### 9.3 API 请求失败

```bash
# 检查后端服务
curl http://localhost:10131/api/health

# 检查 Nginx 代理
curl http://localhost/api/

# 检查 CORS 配置
curl -I -X OPTIONS http://localhost/api/
```

### 9.4 SSL 证书问题

```bash
# 检查证书文件
ls -la /etc/letsencrypt/live/www.omiai.cn/

# 检查证书有效期
openssl x509 -in /etc/letsencrypt/live/www.omiai.cn/fullchain.pem -noout -dates

# 检查 Nginx SSL 配置
nginx -t
```

---

## 附录：快速部署清单

```bash
# 1. 安装基础软件
yum install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin nginx certbot python3-certbot-nginx

# 2. 配置 Docker
mkdir -p /etc/docker
cat > /etc/docker/daemon.json << 'EOF'
{"registry-mirrors": ["https://mirror.ccs.tencentyun.com"], "log-driver": "json-file", "log-opts": {"max-size": "10m", "max-file": "3"}}
EOF
systemctl enable docker && systemctl start docker

# 3. 申请 SSL 证书
mkdir -p /var/www/certbot
certbot certonly --webroot -w /var/www/certbot -d www.omiai.cn --email your@email.com --agree-tos
systemctl enable certbot.timer && systemctl start certbot.timer

# 4. 创建项目目录
mkdir -p /data/omiai-server/{deploy,runtime/logs}

# 5. 配置环境变量（修改 .env 文件）

# 6. 拉取镜像并启动
docker login ghcr.io
docker compose -f docker-compose.prod.yml up -d

# 7. 配置 Nginx
nginx -t && nginx -s reload

# 8. 验证服务
curl https://www.omiai.cn/api/
```

---

**文档版本**: v1.0  
**更新日期**: 2025-02-13  
**维护者**: Omiai Team
