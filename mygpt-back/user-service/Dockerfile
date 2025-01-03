# 使用官方 Golang 镜像
FROM golang:1.22.4

# 设置工作目录
WORKDIR /app

# 安装 MySQL 客户端
RUN apt-get update && apt-get install -y mysql-client iputils-ping

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制整个项目文件
COPY . .

# 构建服务
RUN go build -o main .

# 设置启动命令
CMD ["./main"]

# 暴露服务端口
EXPOSE 8080
