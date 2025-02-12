#!/bin/bash

# 创建主要目录
mkdir -p cmd/server
mkdir -p config
mkdir -p internal/{handler,middleware,router,service/{discovery,loadbalancer,proxy},model}
mkdir -p pkg/{cache,logger,utils}
mkdir -p api
mkdir -p scripts
mkdir -p test/integration

# 创建文件
# cmd
touch cmd/server/main.go

# config
touch config/config.go
touch config/config.yaml

# internal
touch internal/handler/{auth.go,proxy.go,response.go}
touch internal/middleware/{jwt.go,cors.go,ratelimit.go,logger.go}
touch internal/router/router.go
touch internal/service/discovery/nacos.go
touch internal/service/loadbalancer/balancer.go
touch internal/service/proxy/proxy.go
touch internal/model/{route.go,error.go}

# pkg
touch pkg/cache/redis.go
touch pkg/logger/logger.go
touch pkg/utils/{jwt.go,http.go}

# api
touch api/swagger.yaml

# scripts
touch scripts/{build.sh,deploy.sh}
chmod +x scripts/{build.sh,deploy.sh}

# root files
touch Dockerfile
touch Makefile
touch go.mod
touch README.md

echo "目录结构创建完成！" 