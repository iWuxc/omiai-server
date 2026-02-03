# Build Stage
FROM golang:alpine AS builder

# 设置环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /build

# 复制 go.mod 和 go.sum 并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
# 使用包路径 ./cmd/server 以包含目录下所有文件(main.go, wire_gen.go, unix_panic.go等)
# 添加 -mod=mod 以忽略可能存在的 vendor 目录，强制使用 module 模式
RUN go build -mod=mod -ldflags="-s -w" -o server ./cmd/server

# Run Stage
FROM alpine:latest

# 安装基础工具和 timezone
RUN apk --no-cache add ca-certificates tzdata gettext \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /build/server /app/bin/server

# 复制配置文件模板和启动脚本
COPY configs/config.template.yaml /app/configs/config.template.yaml
COPY scripts/entrypoint.sh /app/bin/entrypoint.sh

# 创建日志和运行时目录
RUN mkdir -p /app/runtime/logs && chmod +x /app/bin/entrypoint.sh

# 暴露端口
EXPOSE 10131 10132 10133

# 设置入口点
ENTRYPOINT ["/app/bin/entrypoint.sh"]

# 默认命令
CMD ["/app/bin/server", "-conf", "/app/configs"]
