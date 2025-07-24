# 多阶段构建 - 构建阶段
FROM golang:1.23-alpine AS builder

# 设置Go代理
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git ca-certificates tzdata

# 复制go mod文件，利用缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制全部源代码和资源目录
COPY . .

# 构建Go应用
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o main \
    ./main.go

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 复制二进制和资源目录
COPY --from=builder /app/main .
COPY --from=builder /app/manifest ./manifest
COPY --from=builder /app/hack ./hack
COPY --from=builder /app/resource ./resource

# 设置权限
RUN chmod +x /app/main

# 暴露端口 8000
EXPOSE 8000

# 说明：以下为推荐挂载
# -v /path/to/config.yaml:/app/manifest/config.yaml:ro
# -v /path/to/resource/log:/app/resource/log
# -v /tmp:/tmp

CMD ["./main"]