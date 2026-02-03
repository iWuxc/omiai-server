# 腾讯云轻量应用服务器 (OpenCloudOS 8) 自动化部署文档

本方案适用于基于 VPS/云服务器 (如腾讯云 Lighthouse) 的自动化部署场景，相比 Serverless 方案，它提供了对操作系统底层的完全控制能力，成本更低，适合中小型项目。

## 1. 架构设计

*   **OS**: OpenCloudOS 8 (TencentOS 衍生版，兼容 CentOS 8)
*   **容器引擎**: Docker 26 + Docker Compose
*   **网关**: Nginx (宿主机反向代理 -> 容器)
*   **CI/CD**: GitHub Actions (前后端一体化构建部署)
*   **镜像仓库**: GitHub Container Registry (ghcr.io)
*   **前端服务**: Nginx 容器 (端口 10080)
*   **后端服务**: Go Server 容器 (端口 10131)

## 2. 服务器初始化

登录服务器，执行初始化脚本安装 Docker 和 Nginx：

```bash
# 上传 setup_lighthouse.sh 到服务器
chmod +x deploy/setup_lighthouse.sh
sudo ./deploy/setup_lighthouse.sh
```

该脚本会自动：
1.  配置 Docker yum 源并安装最新版 Docker CE。
2.  安装 Nginx。
3.  创建项目目录 `/data/omiai-server`。
4.  配置 Docker 日志轮转 (防止磁盘爆满)。

## 3. GitHub Actions 配置

### 3.1 密钥配置

在 GitHub 仓库 -> Settings -> Secrets and variables -> Actions 中添加以下 Secrets：

| Secret Name | 说明 |
| :--- | :--- |
| `SERVER_HOST` | 服务器公网 IP |
| `SERVER_USER` | SSH 用户名 (通常为 root 或 lighthouse) |
| `SERVER_SSH_KEY` | SSH 私钥内容 (确保已将公钥添加到服务器 `~/.ssh/authorized_keys`) |
| `DB_HOST` | 数据库 Host (内网IP) |
| `DB_PASSWORD` | 数据库密码 |
| `REDIS_PASSWORD`| Redis 密码 |
| `ZHIPUAI_API_KEY`| AI Key |
| `DOMAIN_H5` | 前端域名 |

### 3.2 流程说明

每次向 `main` 分支推送代码时，`.github/workflows/deploy.yml` 会自动执行：
1.  **构建**: 编译 Go 代码并构建 Docker 镜像。
2.  **推送**: 将镜像推送到 ghcr.io (私有仓库)。
3.  **部署**: 通过 SSH 连接服务器，生成 `docker-compose.prod.yml` 和 `.env`，拉取新镜像并重启服务。

## 4. Nginx 与 SSL 配置

### 4.1 配置反向代理

将 `deploy/nginx/omiai.conf` 复制到服务器 `/etc/nginx/conf.d/`：

```bash
cp deploy/nginx/omiai.conf /etc/nginx/conf.d/omiai.conf
# 修改 server_name 为您的实际域名
vim /etc/nginx/conf.d/omiai.conf
# 重载配置
nginx -s reload
```

此配置已包含：
- `/api/` -> 转发给后端容器 (10131)
- `/` -> 转发给前端容器 (10080)

### 4.2 申请 SSL 证书 (Certbot)

使用 Certbot 自动申请免费证书：

```bash
yum install -y certbot python3-certbot-nginx
certbot --nginx -d your_domain.com
```

Certbot 会自动修改 Nginx 配置以启用 HTTPS。

## 5. 监控与日志

### 5.1 日志查看

*   **应用日志**: 挂载在 `/data/omiai-server/runtime/logs`。
*   **容器标准输出**:
    ```bash
    docker logs -f --tail 100 omiai-server
    ```

### 5.2 基础监控

对于轻量级场景，推荐使用 **Portainer** 进行可视化监控和管理：

```bash
docker volume create portainer_data
docker run -d -p 9000:9000 --name portainer --restart=always -v /var/run/docker.sock:/var/run/docker.sock -v portainer_data:/data portainer/portainer-ce:latest
```

访问 `http://ip:9000` 即可查看容器状态、资源占用和日志。

## 6. 故障排查手册

### Q1: GitHub Action 部署失败，提示 SSH 连接超时
*   检查 `SERVER_HOST` 是否正确。
*   检查服务器防火墙 (安全组) 是否放通 22 端口。

### Q2: 服务启动失败，提示数据库连接错误
*   检查 `.env` 文件中的 `DB_HOST` 是否正确。
*   如果是本机数据库，请使用宿主机 IP (通常是 `172.17.0.1`) 而非 `127.0.0.1`。
*   检查安全组是否放通 3306 端口。

### Q3: 访问 API 报 502 Bad Gateway
*   检查后端容器是否运行: `docker ps`
*   检查 Nginx 配置中 `proxy_pass` 端口是否与容器暴露端口 (10131) 一致。
*   查看 Nginx 错误日志: `tail -f /var/log/nginx/error.log`

### Q4: 磁盘空间报警
*   清理未使用的镜像: `docker image prune -a`
*   检查日志目录 `/data/omiai-server/runtime/logs` 是否过大，配置 Logrotate。
