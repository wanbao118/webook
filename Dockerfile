# 使用官方的Go镜像作为构建环境
FROM golang:1.23 AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件以缓存依赖
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制项目文件
COPY . .

# 构建项目
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app .

# 使用轻量级的Alpine镜像作为运行环境
FROM alpine:latest

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /bin/app .

# 暴露应用端口
EXPOSE 8080

# 运行应用
CMD ["./app"]