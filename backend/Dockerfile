# 构建阶段
FROM golang:1.23-alpine AS builder

# 安装必要的工具
RUN apk add --no-cache git ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 设置Go环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download && go mod verify

# 复制源代码
COPY . .

# 构建参数
ARG BUILD_ENV=production
ARG BUILD_TIME
ARG GIT_COMMIT

# 编译应用
RUN go build -ldflags="-w -s -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -o mcprapi ./cmd/main.go

# 开发阶段（用于开发环境）
FROM golang:1.23-alpine AS development

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# 设置Go环境变量（使用国内镜像源）
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOPROXY=https://goproxy.cn,direct \
    GOSUMDB=sum.golang.google.cn

# 安装air用于热重载（使用兼容版本）
RUN go install github.com/cosmtrek/air@v1.49.0

COPY go.mod go.sum ./

# 设置Go代理并下载依赖
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]

# 生产阶段
FROM alpine:3.18 AS production

# 安装必要的包
RUN apk --no-cache add ca-certificates tzdata wget && \
    addgroup -g 1001 -S mcprapi && \
adduser -u 1001 -S mcprapi -G mcprapi

# 设置时区
ENV TZ=Asia/Shanghai

# 创建必要的目录
RUN mkdir -p /app/logs /app/configs && \
    chown -R mcprapi:mcprapi /app

# 切换到非root用户
USER mcprapi

# 设置工作目录
WORKDIR /app

# 从builder阶段复制编译好的二进制文件
COPY --from=builder --chown=mcprapi:mcprapi /app/mcprapi /app/

# 复制配置文件
COPY --from=builder --chown=mcprapi:mcprapi /app/configs /app/configs/

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 暴露端口
EXPOSE 8080

# 设置启动命令
CMD ["/app/mcprapi", "--config", "configs/prod.yaml"]