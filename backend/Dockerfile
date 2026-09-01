# 多阶段构建 - 编译阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装 git（用于获取依赖）
RUN apk add --no-cache git

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译应用（禁用 CGO，使用 linux 平台）
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/ai-wordbook .

# 运行阶段 - 使用最小镜像
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# 从编译阶段复制二进制文件
COPY --from=builder /app/ai-wordbook .
COPY --from=builder /app/.env.example .env

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动应用
CMD ["./ai-wordbook"]
