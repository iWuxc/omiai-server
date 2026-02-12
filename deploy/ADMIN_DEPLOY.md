# Admin 后台部署文档

## 概述

Admin 后台已集成到一键部署流程中，与后端服务和 H5 前端一起部署。

## 架构

```
┌─────────────────────────────────────────────────────────────┐
│                        Nginx (80/443)                        │
├─────────────────────────────────────────────────────────────┤
│  /api/*     →  omiai-server:10131 (后端服务)                 │
│  /admin/*   →  omiai-admin:80 (Admin后台)                    │
│  /*         →  omiai-frontend:80 (H5前端)                    │
└─────────────────────────────────────────────────────────────┘
```

## 服务端口

| 服务 | 容器端口 | 主机端口 | 说明 |
|------|---------|---------|------|
| omiai-server | 10131-10133 | 10131-10133 | 后端服务 |
| omiai-frontend | 80 | 10080 | H5 前端 |
| omiai-admin | 80 | 10081 | Admin 后台 |
| mysql | 3306 | 13018 | 数据库 |
| redis | 6379 | 6379 | 缓存 |

## 访问地址

- **H5 前端**: http://www.omiai.cn/
- **Admin 后台**: http://www.omiai.cn/admin/
- **API 接口**: http://www.omiai.cn/api/

## 部署方式

### 自动部署 (推荐)

推送到 `main` 分支后，GitHub Actions 会自动：

1. 构建后端服务镜像 (`ghcr.io/owner/omiai-server:latest`)
2. 构建 H5 前端镜像 (`ghcr.io/owner/omiai-server-web:latest`)
3. 构建 Admin 后台镜像 (`ghcr.io/owner/omiai-server-admin:latest`)
4. 部署到服务器

### 手动部署

```bash
# 1. 构建 Admin 镜像
cd omiai-admin
docker build -t omiai-admin:latest .

# 2. 运行容器
docker run -d \
  --name omiai-admin \
  -p 10081:80 \
  --network omiai-network \
  omiai-admin:latest

# 3. 更新 Nginx 配置
# 添加 /admin/ 路由到 omiai-admin:80
```

## 本地开发

```bash
cd omiai-admin

# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev:antd

# 构建生产版本
pnpm build:antd
```

## 环境变量

### 构建时环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| VITE_API_BASE_URL | API 基础地址 | /api |
| VITE_APP_TITLE | 应用标题 | Omiai Admin |

### 运行时配置

Admin 容器内的 Nginx 会自动代理 `/api/` 请求到后端服务。

## 文件结构

```
omiai-admin/
├── Dockerfile              # Docker 构建文件
├── .dockerignore          # Docker 忽略文件
├── apps/
│   └── web-antd/
│       ├── nginx.conf     # Nginx 配置
│       ├── .env.production # 生产环境变量
│       └── ...
└── ...
```

## 故障排查

### 查看容器日志

```bash
docker logs omiai-admin
```

### 进入容器调试

```bash
docker exec -it omiai-admin sh
```

### 检查 Nginx 配置

```bash
docker exec omiai-admin nginx -t
```

## 更新部署

```bash
# 拉取最新镜像
docker pull ghcr.io/owner/omiai-server-admin:latest

# 重启服务
cd /data/omiai-server/deploy
docker compose -f docker-compose.prod.yml up -d admin
```
