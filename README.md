# omiai - 婚恋匹配服务

## 项目简介

Go + Gin + GORM + Redis 开发的婚恋匹配服务端应用。

## 核心功能

- 客户管理：客户档案的增删改查、公海池管理
- 匹配推荐：AI 算法自动推荐合适对象
- 跟进提醒：自动生成跟进任务、提醒规则配置
- 数据统计：工作台数据概览
- H5 管理端：内置轻量级管理后台

## 快速开始

### 环境要求

- Go 1.25+

### 配置

```bash
cp configs/config.yaml.bak configs/config.yaml
```

### 启动服务

```bash
make build
./bin/server
```

## 命令

```bash
make build    # 构建
make lint     # 代码检查
make test     # 测试
make wire     # 重新生成依赖注入
```

## 目录结构

```
cmd/server/      # 入口
configs/         # 配置
internal/        # 业务代码 (controller/biz/data)
pkg/             # 公共包
```
