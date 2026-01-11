<p align="center">
  <h1 align="center">📺 LiveChannelsCN</h1>
  <p align="center">
    A <a href="https://github.com/glanceapp/glance">Glance</a> extension widget for monitoring Chinese live streaming platforms
    <br />
    <a href="./README-ZH.md">中文文档</a> · <a href="#quick-start">Quick Start</a> · <a href="https://github.com/glanceapp/glance">Glance</a>
  </p>
</p>

<p align="center">
  <a href="https://github.com/Cynosure159/LiveChannelsCN/actions/workflows/ci.yml">
    <img src="https://github.com/Cynosure159/LiveChannelsCN/actions/workflows/ci.yml/badge.svg" alt="CI Status" />
  </a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go" alt="Go Version" />
  <img src="https://img.shields.io/badge/License-MIT-green?style=flat-square" alt="License" />
  <img src="https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker" alt="Docker" />
</p>

---

## ✨ Features

- 🎮 **Multi-Platform Support** - Bilibili, Douyu, Huya
- 🔴 **Real-time Status** - Live/offline indicators with viewer counts
- 🎨 **Glance Native Style** - Matches Twitch Channels widget design
- ⚡ **Concurrent Requests** - Fast parallel API calls
- 🐳 **Docker Ready** - Easy deployment with Docker Compose

## 📸 Preview

The widget displays streamers in a familiar Twitch-style layout:
- Avatar with live indicator
- Streamer name and game category
- Current viewer count
- Hover preview with stream thumbnail

## 🚀 Quick Start

### Prerequisites

- Go 1.21+ or Docker
- A running [Glance](https://github.com/glanceapp/glance) instance

### Option 1: Docker (Recommended)

```bash
# Clone the repository
git clone https://github.com/yourusername/LiveChannelsCN.git
cd LiveChannelsCN

# Copy and edit config
cp config.example.json config.json

# Start the service
docker-compose up -d
```

### Option 2: Build from Source

```bash
# Clone and build
git clone https://github.com/yourusername/LiveChannelsCN.git
cd LiveChannelsCN
go build -o live-channels

# Configure and run
cp config.example.json config.json
./live-channels
```

## ⚙️ Configuration

Edit `config.json` to add streamers:

```json
{
  "channels": [
    {
      "platform": "bilibili",
      "channel_id": "21013446",
      "name": "Streamer Name"
    },
    {
      "platform": "douyu",
      "channel_id": "5279",
      "name": "Another Streamer"
    },
    {
      "platform": "huya",
      "channel_id": "11336",
      "name": "Huya Streamer"
    }
  ]
}
```

### Supported Platforms

| Platform | `platform` value | How to get `channel_id` |
|----------|------------------|-------------------------|
| Bilibili | `bilibili` | Room ID from `live.bilibili.com/{id}` |
| Douyu | `douyu` | Room ID from `douyu.com/{id}` |
| Huya | `huya` | Room ID from `huya.com/{id}` |

## 🔗 Glance Integration

Add to your `glance.yml`:

```yaml
- type: extension
  url: http://localhost:8081
  allow-potentially-dangerous-html: true
  cache: 5m
  title: Live Channels
```

## 📡 API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | HTML widget for Glance <br> Params: `?cache=60` (cache TTL in sec), `?collapse=10` (max items before collapse) |
| `/api/streams` | GET | All stream statuses (JSON) <br> Params: `?cache=60` |
| `/api/streams/:platform` | GET | Filter by platform <br> Params: `?cache=60` |
| `/health` | GET | Health check |

## 🛠️ Development

```bash
# Run tests
go test -v ./...

# Format code
go fmt ./...

# Build
make build
```

## ⚙️ Advanced Configuration

You can adjust service behavior via environment variables or command-line flags:

| Env Variable | CLI Flag | Default | Description |
|--------------|----------|---------|-------------|
| `LOG_LEVEL` | `-level` | `info` | Log level (`debug`, `info`, `warn`, `error`) |
| `GIN_MODE` | `-mode` | `debug` | Set to `release` for production mode (JSON logs) |
| `CONFIG_PATH` | `-config` | `./config/config.json` | Path to your configuration file |
| `PORT` | `-port` | `8081` | Server listening port |
| `USER_AGENT` | `-ua` | (Built-in Default) | Custom HTTP User-Agent |

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- [Glance](https://github.com/glanceapp/glance) - The amazing self-hosted dashboard
- Inspired by Glance's built-in Twitch Channels widget
