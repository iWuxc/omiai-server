# Admin 后台部署文档

## 概述

Admin 后台已集成到一键部署流程中，与后端服务和 H5 前端一起部署。

## 仓库结构

```
├── omiai-server (主仓库)
│   ├── .github/workflows/deploy.yml  # 部署配置
│   ├── Dockerfile                    # 后端 Dockerfile
│   └── ...
│
├── omiai-miniapp (独立仓库)          # H5 前端
│   ├── Dockerfile
│   └── ...
│
└── omiai-admin (独立仓库)            # Admin 后台
    ├── Dockerfile                    # Admin Dockerfile
    ├── apps/web-antd/nginx.conf      # Nginx 配置
    └── ...
```

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

## GitHub Actions 配置

### 必需的 Secrets

在 `omiai-server` 仓库的 Settings → Secrets and variables → Actions 中配置：

| Secret 名称 | 说明 |
|------------|------|
| `GH_PAT` | GitHub Personal Access Token (访问 omiai-admin 仓库) |
| `SERVER_HOST` | 服务器地址 |
| `SERVER_USER` | 服务器用户名 |
| `SERVER_SSH_KEY` | SSH 私钥 |
| `DB_PASSWORD` | 数据库密码 |
| `REDIS_PASSWORD` | Redis 密码 |
| 其他 Secrets... | 参考原有配置 |

### 创建 GH_PAT

1. 访问 GitHub → Settings → Developer settings → Personal access tokens → Tokens (classic)
2. 点击 "Generate new token (classic)"
3. 勾选权限：
   - `repo` (完整仓库访问权限)
   - `read:packages`
   - `write:packages`
4. 生成并复制 Token
5. 在 `omiai-server` 仓库的 Secrets 中添加 `GH_PAT`

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

## 部署触发

### 自动部署

推送到以下仓库的 `main` 分支会触发部署：

1. `omiai-server` - 触发完整部署
2. `omiai-miniapp` - 通过 `repository_dispatch` 触发
3. `omiai-admin` - 通过 `repository_dispatch` 触发

### 手动部署

在 GitHub Actions 页面点击 "Run workflow"

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

## 本地测试 Docker 构建

```bash
cd omiai-admin

# 构建 Docker 镜像
docker build -t omiai-admin:test .

# 本地运行测试
docker run -d -p 10081:80 --name omiai-admin-test omiai-admin:test

# 访问测试
open http://localhost:10081
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

### 常见问题

#### 1. `path "./omiai-admin" not found`

**原因**: `omiai-admin` 仓库访问权限问题

**解决**: 确保 `GH_PAT` Secret 已正确配置，且 Token 有 `repo` 权限

#### 2. API 请求 404

**原因**: Nginx 代理配置问题

**解决**: 检查 `apps/web-antd/nginx.conf` 中的 API 代理配置

#### 3. 页面空白

**原因**: 路由模式或构建问题

**解决**: 检查 `.env.production` 中的 `VITE_ROUTER_HISTORY=history`
