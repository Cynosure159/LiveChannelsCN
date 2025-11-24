# 第一阶段：构建阶段
FROM golang:1.25-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制源码
COPY . .

# 复制 go.mod 和 go.sum（如果存在）
# COPY go.mod go.sum ./
# RUN go env -w GO111MODULE=on 
# RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download



# 编译二进制文件，启用静态链接以兼容 Alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o live-channels .

# 第二阶段：运行阶段（基于 Alpine，轻量且可调试）
FROM alpine:3.20

# 安装必要工具（可选，用于调试）
RUN apk --no-cache add ca-certificates && \
    mkdir -p /config

# 复制编译好的二进制文件
COPY --from=builder /app/live-channels /live-channels

# 设置挂载点（仅保留 /config）
VOLUME ["/config"]

# 设置环境变量
ENV CONFIG_PATH=/config/channel.json

# 暴露端口
EXPOSE 8080

# 设置启动命令
CMD ["/live-channels"]