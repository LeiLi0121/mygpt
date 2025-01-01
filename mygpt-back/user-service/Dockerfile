# 使用官方 Golang 镜像
FROM golang:1.20

# 设置工作目录
WORKDIR /app

# 复制整个 mygpt-back 目录
COPY ../../mygpt-back /app

# 切换到 user-service 子目录
WORKDIR /app/user-service

# 下载依赖
RUN go mod download

# 构建项目
RUN go build -o main .

# 设置启动命令
CMD ["./main"]

# 暴露服务端口
EXPOSE 8080
