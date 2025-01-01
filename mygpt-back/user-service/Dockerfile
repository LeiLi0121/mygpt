# 使用官方 Golang 镜像作为基础镜像
FROM golang:1.20

# 设置工作目录
WORKDIR /app

# 将当前目录内容复制到容器内
COPY . .

# 下载依赖并构建项目
RUN go mod tidy
RUN go build -o main .

# 设置默认启动命令
CMD ["./main"]

# 暴露服务端口
EXPOSE 8080
