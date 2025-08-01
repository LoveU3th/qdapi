# 使用官方 Go 镜像作为构建环境
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 文件
COPY go.mod ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o qdapi ./cmd

# 使用轻量级的 Alpine 镜像作为运行环境
FROM alpine:latest

# 安装 ca-certificates 用于 HTTPS 请求
RUN apk --no-cache add ca-certificates tzdata

# 设置时区为上海
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone

# 创建非 root 用户
RUN addgroup -g 1001 -S qdapi && \
    adduser -u 1001 -S qdapi -G qdapi

# 设置工作目录
WORKDIR /app

# 从构建阶段复制可执行文件
COPY --from=builder /app/qdapi .

# 创建配置文件目录并设置权限
RUN mkdir -p /app/config && chown -R qdapi:qdapi /app

# 切换到非 root 用户
USER qdapi

# 暴露端口（如果需要）
# EXPOSE 8080

# 运行应用
CMD ["./qdapi"] 