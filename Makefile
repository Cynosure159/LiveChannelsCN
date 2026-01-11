.PHONY: build run clean test dev docker-build docker-run help

help:
	@echo "Live Channels - Makefile 命令"
	@echo ""
	@echo "可用命令:"
	@echo "  make dev          - 运行开发服务器"
	@echo "  make build        - 构建可执行文件"
	@echo "  make run          - 运行应用"
	@echo "  make clean        - 清理构建文件"
	@echo "  make test         - 运行测试"
	@echo "  make docker-build - 构建 Docker 镜像"
	@echo "  make docker-run   - 运行 Docker 容器"
	@echo "  make help         - 显示此帮助信息"

dev:
	@echo "启动开发服务器..."
	go run main.go

build:
	@echo "构建应用..."
	go build -o live-channels .

run: build
	@echo "运行应用..."
	./live-channels

clean:
	@echo "清理构建文件..."
	rm -f live-channels live-channels.exe

test:
	@echo "运行测试..."
	go test -v ./...

docker-build:
	@echo "构建 Docker 镜像..."
	docker build -t live-channels:latest .

docker-run: docker-build
	@echo "运行 Docker 容器..."
	docker-compose up -d

deps:
	@echo "下载依赖..."
	go mod download
	go mod tidy

fmt:
	@echo "格式化代码..."
	go fmt ./...

lint:
	@echo "检查代码..."
	go vet ./...
