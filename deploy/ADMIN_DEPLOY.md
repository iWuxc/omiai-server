# Admin 后台部署文档

## 概述

Admin 后台已集成到一键部署流程中，与后端服务和 H5 前端一起部署。

## 仓库结构

```
├── omiai-server (主仓库)
│   ├── .github/workflows/deploy.yml  # 部署配置
│   ├── Dockerfile                    # 后端 Dockerfile
│   └── deploy/nginx/omiai.conf       # 主 Nginx 配置
│
├── omiai-miniapp (公开仓库)          # H5 前端
│   ├── Dockerfile
│   ├── nginx.conf                    # H5 容器 Nginx 配置
│   └── src/manifest.json             # H5 router base 配置
│
└── omiai-admin (公开仓库)            # Admin 后台
    ├── Dockerfile                    # Admin Dockerfile
    ├── apps/web-antd/nginx.conf      # Admin 容器 Nginx 配置
    └── apps/web-antd/.env.production # Admin 构建配置
```

## 架构

```
┌─────────────────────────────────────────────────────────────┐
│                     主 Nginx (80/443)                        │
├─────────────────────────────────────────────────────────────┤
│  /api/*     →  omiai-server:10131 (后端服务)                 │
│  /h5/*      →  omiai-frontend:10080/h5 (H5前端)              │
│  /web/*     →  omiai-admin:10081/web (Admin后台)             │
│  /          →  重定向到 /h5                                  │
└─────────────────────────────────────────────────────────────┘
```

## 访问地址

| 服务 | 访问地址 | 说明 |
|------|---------|------|
| H5 前端 | http://www.omiai.cn/h5 | 移动端 H5 页面 |
| Admin 后台 | http://www.omiai.cn/web | 管理后台 |
| API 接口 | http://www.omiai.cn/api | 后端 API |
| 根路径 | http://www.omiai.cn/ | 重定向到 /h5 |

## 配置说明

### 1. 主服务器 Nginx 配置

文件: `deploy/nginx/omiai.conf`

```nginx
# H5 前端
location /h5 {
    proxy_pass http://localhost:10080/h5;
}

# Admin 后台
location /web {
    proxy_pass http://localhost:10081/web;
}

# 根路径重定向
location = / {
    return 301 /h5;
}
```

### 2. H5 前端配置

**manifest.json** - 设置 H5 router base:
```json
{
    "h5": {
        "router": {
            "base": "/h5/"
        }
    }
}
```

**nginx.conf** - 容器内 Nginx:
```nginx
location /h5 {
    alias /usr/share/nginx/html;
    try_files $uri $uri/ /h5/index.html;
}
```

### 3. Admin 后台配置

**.env.production** - 设置构建 base:
```
VITE_BASE=/web/
```

**nginx.conf** - 容器内 Nginx:
```nginx
location /web {
    alias /usr/share/nginx/html;
    try_files $uri $uri/ /web/index.html;
}
```

## GitHub Actions 配置

### 必需的 Secrets

在 `omiai-server` 仓库的 Settings → Secrets and variables → Actions 中配置：

| Secret 名称 | 说明 |
|------------|------|
| `SERVER_HOST` | 服务器地址 |
| `SERVER_USER` | 服务器用户名 |
| `SERVER_SSH_KEY` | SSH 私钥 |
| `DB_PASSWORD` | 数据库密码 |
| `REDIS_PASSWORD` | Redis 密码 |
| `ZHIPUAI_API_KEY` | 智谱 AI API Key |
| `DOMAIN_H5` | H5 域名 |
| `STORAGE_DRIVER` | 存储驱动 |
| `COS_BUCKET_URL` | 腾讯云 COS Bucket URL |
| `COS_REGION` | 腾讯云 COS Region |
| `COS_SECRET_ID` | 腾讯云 COS Secret ID |
| `COS_SECRET_KEY` | 腾讯云 COS Secret Key |

## 服务端口

| 服务 | 容器端口 | 主机端口 | 说明 |
|------|---------|---------|------|
| omiai-server | 10131-10133 | 10131-10133 | 后端服务 |
| omiai-frontend | 80 | 10080 | H5 前端 |
| omiai-admin | 80 | 10081 | Admin 后台 |
| mysql | 3306 | 13018 | 数据库 |
| redis | 6379 | 6379 | 缓存 |

## 部署触发

### 自动部署

推送到 `omiai-server` 仓库的 `main` 分支会触发完整部署。

### 手动部署

在 GitHub Actions 页面点击 "Run workflow"

## 本地开发

### Admin 后台

```bash
cd omiai-admin

# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev:antd

# 构建生产版本
pnpm build:antd
```

### H5 前端

```bash
cd omiai-miniapp

# 安装依赖
npm install

# 启动开发服务器
npm run dev:h5

# 构建生产版本
npm run build:h5
```

## 故障排查

### 查看容器日志

```bash
docker logs omiai-admin
docker logs omiai-frontend
```

### 进入容器调试

```bash
docker exec -it omiai-admin sh
docker exec -it omiai-frontend sh
```

### 检查 Nginx 配置

```bash
docker exec omiai-admin nginx -t
docker exec omiai-frontend nginx -t
```

### 常见问题

#### 1. 页面空白或 404

**原因**: base path 配置不正确

**解决**: 
- Admin: 检查 `.env.production` 中 `VITE_BASE=/web/`
- H5: 检查 `manifest.json` 中 `h5.router.base: "/h5/"`

#### 2. 静态资源加载失败

**原因**: Nginx alias 配置问题

**解决**: 确保容器内 Nginx 使用 `alias` 而非 `root`

#### 3. API 请求 404

**原因**: API 代理配置问题

**解决**: 检查主 Nginx 配置中的 `/api/` location

## 需要提交的文件

### omiai-admin 仓库

- `Dockerfile`
- `apps/web-antd/nginx.conf`
- `apps/web-antd/.env.production`
- `.dockerignore`

### omiai-miniapp 仓库

- `nginx.conf`
- `src/manifest.json` (更新 h5.router.base)

### omiai-server 仓库

- `.github/workflows/deploy.yml`
- `deploy/nginx/omiai.conf`
