# omiai

## 项目概览

本项目为服务端应用，采用 Go 开发。文中出现的代码路径均以项目根路径为基准。

## 目录结构

```text
.
├── Dockerfile
├── Makefile
├── README.md
├── cmd
│  └── server
│    ├── main.go            项目入口
│    ├── wire.go            依赖注入定义
│    └── wire_gen.go
├── configs
│  ├── config.yaml          实际配置
│  └── config.yaml.bak      配置示例
├── data                    数据文件
├── doc                     文档中心
│  ├── FAQ.md
│  ├── README.md
│  └── images
├── internal
│  ├── api                  接口定义
│  ├── biz                  业务接口定义
│  ├── conf                 配置与初始化
│  ├── controller           业务处理
│  ├── cron                 定时任务
│  ├── data                 数据访问
│  ├── middleware           中间件
│  ├── server               服务初始化
│  ├── service              业务服务
│  └── validate             数据校验
├── pkg
│  └── response             响应与状态码
├── runtime                 运行时日志
└── scripts                 启停脚本
```

## 版本要求

Go 1.25

## 本地开发

1. 配置 Go 代理（建议写入 `~/.bashrc` 或 `~/.zshrc`）

```bash
export GO111MODULE=on
export GOPROXY=https://goproxy.io
```

2. 准备配置文件

```bash
cp configs/config.yaml.bak configs/config.yaml
``

