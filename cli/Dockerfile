# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git ca-certificates tzdata

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X github.com/ciclebyte/template_starter/cli/cmd.Version=docker -X github.com/ciclebyte/template_starter/cli/cmd.GitCommit=$(git rev-parse HEAD 2>/dev/null || echo 'unknown') -X github.com/ciclebyte/template_starter/cli/cmd.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%S)" \
    -a -installsuffix cgo \
    -o template-cli main.go

# 运行阶段
FROM alpine:latest

# 安装ca证书和时区数据
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/template-cli .

# 创建配置目录
RUN mkdir -p /root/.template-cli

# 设置环境变量
ENV PATH="/root:${PATH}"

# 暴露配置目录作为volume
VOLUME ["/root/.template-cli"]

# 设置入口点
ENTRYPOINT ["./template-cli"]

# 默认命令
CMD ["--help"]