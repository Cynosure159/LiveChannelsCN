<p align="center">
  <h1 align="center">📺 LiveChannelsCN</h1>
  <p align="center">
    <a href="https://github.com/glanceapp/glance">Glance</a> 看板的中国直播平台扩展组件
    <br />
    <a href="./README.md">English</a> · <a href="#快速开始">快速开始</a> · <a href="https://github.com/glanceapp/glance">Glance</a>
  </p>
</p>

<p align="center">
  <a href="https://github.com/Cynosure159/LiveChannelsCN/actions/workflows/ci.yml">
    <img src="https://github.com/Cynosure159/LiveChannelsCN/actions/workflows/ci.yml/badge.svg" alt="CI 状态" />
  </a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go" alt="Go Version" />
  <img src="https://img.shields.io/badge/License-MIT-green?style=flat-square" alt="License" />
  <a href="https://hub.docker.com/r/cynosure159/live-channels-cn">
    <img src="https://img.shields.io/docker/v/cynosure159/live-channels-cn?sort=semver&style=flat-square&logo=docker&start=latest" alt="Docker Image" />
  </a>
</p>

---

## 📚 项目介绍

由于Glance原生组件只支持twitch直播状态查看，因此开发了这个插件，用于查看国内直播平台的主播状态
本人也不太懂go，因此绝大部分代码由AI生成
如果需要支持更多的平台请提issue，或者直接PR（这当然最好，可以直接在 ./platform 目录下实现新的平台类）

## ✨ 功能特性

- 🎮 **多平台支持** - B站、斗鱼、虎牙
- 🔴 **实时状态** - 开播/离线指示，显示观看人数
- 🎨 **原生风格** - 与 Glance 内置 Twitch 组件保持一致
- ⚡ **并发请求** - 快速并行获取数据
- 🐳 **容器部署** - Docker Compose 一键启动

## 📸 效果预览

组件采用 Twitch 风格布局展示：
- 带直播状态指示的头像
- 主播名称与游戏分类
- 实时观看人数
- 悬浮预览直播封面

## 🚀 快速开始

### 环境要求

- Go 1.21+ 或 Docker
- 运行中的 [Glance](https://github.com/glanceapp/glance) 实例

### 方式一：Docker Hub 镜像（推荐）

```bash
# 1. 创建配置文件
mkdir -p config
cat > config/config.json << 'EOF'
{
  "channels": [
    {
      "platform": "bilibili",
      "channel_id": "21013446",
      "name": "主播名称"
    }
  ]
}
EOF

# 2. 运行容器
docker run -d \
  --name live-channels \
  -p 8081:8081 \
  -v $(pwd)/config/config.json:/config/config.json:ro \
  -e GIN_MODE=release \
  -e LOG_LEVEL=info \
  --restart unless-stopped \
  cynosure159/live-channels-cn:latest
```

### 方式二：Docker Compose

```bash
# 克隆仓库
git clone https://github.com/Cynosure159/LiveChannelsCN.git
cd LiveChannelsCN

# 编辑配置文件
vim config/config.json

# 启动服务
docker-compose up -d
```

### 方式三：源码编译

```bash
# 克隆并编译
git clone https://github.com/Cynosure159/LiveChannelsCN.git
cd LiveChannelsCN
go build -o live-channels

# 配置并运行
cp config.example.json config/config.json
./live-channels
```

## ⚙️ 配置说明

编辑 `config.json` 添加主播：

```json
{
  "channels": [
    {
      "platform": "bilibili",
      "channel_id": "21013446",
      "name": "主播名称"
    },
    {
      "platform": "douyu",
      "channel_id": "5279",
      "name": "斗鱼主播"
    },
    {
      "platform": "huya",
      "channel_id": "11336",
      "name": "虎牙主播"
    }
  ]
}
```

### 支持的平台

| 平台 | `platform` 值 | 如何获取 `channel_id` |
|------|---------------|----------------------|
| B站 | `bilibili` | 直播间链接 `live.bilibili.com/{房间号}` |
| 斗鱼 | `douyu` | 直播间链接 `douyu.com/{房间号}` |
| 虎牙 | `huya` | 直播间链接 `huya.com/{房间号}` |

## 🔗 Glance 集成

在 `glance.yml` 中添加：

```yaml
- type: extension
  url: http://localhost:8081
  allow-potentially-dangerous-html: true
  cache: 5m
  title: 直播状态
```

## 📡 API 接口

| 端点 | 方法 | 描述 |
|------|------|------|
| `/` | GET | HTML 组件（供 Glance 嵌入） <br> 参数：`?cache=60` (缓存时间秒), `?collapse=10` (折叠数量) |
| `/api/streams` | GET | 所有主播状态 (JSON) <br> 参数：`?cache=60` |
| `/api/streams/:platform` | GET | 按平台筛选 <br> 参数：`?cache=60` |
| `/health` | GET | 健康检查 |

## 🛠️ 开发指南

```bash
# 运行测试
go test -v ./...

# 格式化代码
go fmt ./...

# 构建
make build
```

## ⚙️ 进阶配置

可以通过环境变量或命令行参数调整服务行为：

| 环境变量 | 命令行参数 | 默认值 | 说明 |
|----------|------------|--------|------|
| `LOG_LEVEL` | `-level` | `info` | 日志等级 (`debug`, `info`, `warn`, `error`) |
| `GIN_MODE` | `-mode` | `debug` | 设置为 `release` 可切换到生产模式（JSON 日志） |
| `CONFIG_PATH` | `-config` | `./config/config.json` | 配置文件路径 |
| `PORT` | `-port` | `8081` | 服务监听端口 |
| `USER_AGENT` | `-ua` | (内置默认值) | 自定义 HTTP User-Agent |

**Docker 部署示例**：
```bash
docker run -d \
  -p 8081:8081 \
  -v ./config/config.json:/config/config.json:ro \
  -e LOG_LEVEL=info \
  -e GIN_MODE=release \
  cynosure159/live-channels-cn:latest
```

## 📄 开源许可

MIT License - 详见 [LICENSE](LICENSE)

## 🙏 致谢

- [Glance](https://github.com/glanceapp/glance) - 优秀的自托管看板项目
- 设计灵感来自 Glance 内置的 Twitch Channels 组件
