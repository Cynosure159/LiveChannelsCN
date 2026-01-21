# AGENTS.md

> 项目知识库 - 为 AI Agents 和开发者提供完整项目上下文

## 📋 项目概述

**LiveChannelsCN** 是一个 [Glance](https://github.com/glanceapp/glance) 扩展组件，用于监控中国直播平台（B站、斗鱼、虎牙）的主播实时开播状态。设计完全参照 Glance 内置 **Twitch Channels** 组件风格。

**核心价值**：在自托管看板中集中查看多个直播平台主播状态

---

## 🏗️ 系统架构

### 分层设计

```
Glance 看板
    ↓ HTTP GET /
LiveChannelsCN (Gin Server)
    ├── API 层 (router.go) - HTTP 路由、CORS
    ├── 服务层 (stream_service.go) - 业务逻辑、并发处理、公共 Fetch 逻辑
    ├── 平台层 (factory.go) - 工厂模式 + 策略模式
    │   ├── client.go - 共享 HTTP 客户端单例 (Resty)
    │   ├── bilibili.go - B站 API 客户端
    │   ├── douyu.go - 斗鱼 API 客户端
    │   └── huya.go - 虎牙 API 客户端
    └── 模型层 (models.go) - 数据结构
    ↓ REST API
外部直播平台 API
```

### 核心数据模型

```go
// 直播状态 - internal/models/models.go
type StreamStatus struct {
    ChannelID    string  // 房间号
    Name         string  // 主播名
    Platform     string  // bilibili|douyu|huya
    IsLive       bool    // 是否在线
    Title        string  // 直播标题
    Viewers      int     // 观看人数
    ThumbnailURL string  // 封面
    AvatarURL    string  // 头像
    ProfileURL   string  // 主页链接
    UpdatedAt    int64   // 时间戳
}

// 平台接口 - internal/platform/factory.go
type StreamProvider interface {
    GetStreamStatus(channelID string) (*StreamStatus, error)
}
```

---

## ⚙️ 配置与部署

### 配置文件 (`config.json`)

```json
{
  "channels": [
    {
      "platform": "bilibili",
      "channel_id": "21013446",
      "name": "主播名称（可选）"
    }
  ]
}
```

**环境变量**：
- `CONFIG_PATH`: 配置文件路径（默认 `./config/config.json`，Docker 中为 `/config/config.json`）
- `PORT`: 服务端口（默认 `8081`）
- `LOG_LEVEL`: 日志等级 (`debug`, `info`, `warn`, `error`，默认 `info`)
- `GIN_MODE`: 运行模式 (`release` 会启用生产级 JSON 日志，默认为 `debug`)
- `USER_AGENT`: 自定义 HTTP User-Agent

**Channel ID 获取**：
- B站：`live.bilibili.com/{房间号}`
- 斗鱼：`douyu.com/{房间号}`
- 虎牙：`huya.com/{房间号}`

### Docker 部署

```yaml
# docker-compose.yml
services:
  live-channels:
    image: live-channels:latest
    ports:
      - "8081:8081"
    volumes:
      - ./config:/config
    environment:
      - CONFIG_PATH=/config/config.json
```

**镜像**：多阶段构建，基于 Alpine 3.20，约 15MB

### Glance 集成

```yaml
# glance.yml
- type: extension
  url: http://localhost:8081
  allow-potentially-dangerous-html: true
  cache: 5m
```

---

## 🔧 技术栈

| 组件 | 技术 | 版本 |
|------|------|------|
| 语言 | Go | 1.25+ |
| Web 框架 | Gin | v1.11.0 |
| HTTP 客户端 | Resty | v2.16.5 |
| 模板引擎 | Go Template | 内置 |
| 日志库 | Zap | v1.27.1 |
| 容器 | Docker + Alpine | 3.20 |

**前端**：复用 Glance CSS 类（`twitch-channel-live`, `list`, `collapsible-container` 等）

---

## 📡 API 端点

| 端点 | 方法 | 描述 |
|------|------|------|
| `/` | GET | HTML Widget（供 Glance 嵌入） <br> 参数：`?cache=60` (设置缓存时间) <br> 参数：`?collapse=10` (设置折叠数量) |
| `/api/streams` | GET | 所有主播状态 (JSON) <br> 参数：`?cache=60` |
| `/api/streams/:platform` | GET | 按平台筛选 <br> 参数：`?cache=60` |
| `/health` | GET | 健康检查 |

**响应格式**：
```json
{
  "status": "success|error",
  "data": [StreamStatus],
  "message": "错误信息（可选）"
}
```

---

## 🎯 设计模式 & 优化

### 工厂模式
`platform.CreateProvider` 根据平台类型创建对应的客户端实例。

### 集中校验
在 `models` 包中为 `Platform` 类型实现了 `IsValid()` 方法，统一了平台合法性校验逻辑，减少 API 层的硬编码。

### 统一配置覆盖逻辑
- **集中处理**：提取 `applyConfigOverrides` 私有方法，统一处理 `config.json` 中的字段覆盖（如主播名）。
- **解耦缓存**：缓存中存储原始 raw 数据，仅在返回给用户前实时应用覆盖逻辑，确保即便在运行时修改配置也能立即生效。
- **全路径覆盖**：无论通过 API 实时获取、内存缓存命中还是故障降级（Stale Cache），均通过同一套覆盖流程。

### 单例模式
使用 `platform.GetHTTPClient()` 获取全局共享的 Resty 客户端，复用 TCP 连接。

### 并发处理 (Worker Pool)
Service 层使用 **Worker Pool** 模式处理并发请求，默认 10 个 Worker，防止突发流量耗尽系统资源。

### HTTP 优化 (连接池)
- **全局单例**：所有平台共享同一个 Resty 客户端实例。
- **连接池配置**：自定义 `http.Transport`，配置 `MaxIdleConns` 和 `IdleConnTimeout`，复用 TCP 连接。
- **重试策略**：添加重试等待（Backoff）和基于状态码的重试条件，减少网络抖动影响。

### 缓存机制
- **内存缓存**：使用 `map` + `RWMutex` 实现简单的内存缓存，默认 TTL 为 **1 分钟**。
- **参数化配置**：支持通过 URL 参数 `?cache=SEC` 动态设置缓存时间（默认 60s）。
- **容错降级**：当 API 请求失败时，优先返回过期的缓存数据，避免前端显示为空。
- **并发安全**：Worker 读取缓存使用读锁，写入使用写锁，确保线程安全。

### 结构化日志 (Zap)
- **高性能**：使用 Uber Zap 替换标准库 log，提供极高性能的结构化日志记录。
- **环境隔离**：开发模式（高亮 Console）与生产模式（JSON）自动切换。
- **智能默认值**：使用 `go run` 启动时默认开启 `debug` 等级，编译后的二进制文件或 Docker 运行默认开启 `info` 等级。
- **上下文丰富**：日志自动携带 Platform、ChannelID 等关键字段，便于排查问题。

### 可配置 User-Agent
- **统一管理**：全局共享的 Resty 客户端统一配置 User-Agent，移除各平台代码中的硬编码。
- **灵活配置**：支持通过 `config.json`、环境变量 `USER_AGENT` 或命令行参数 `-ua` 进行动态设置。
- **安全增强**：默认使用更具辨识度的 UA，降低被直播平台作为爬虫拦截的风险。

### CI/CD 自动化
- **GitHub Actions**：配置了 `.github/workflows/ci.yml` 自动化流水线。
- **质量检查**：每次 Push 或 PR 自动运行 `go fmt` 和 `go vet`。
- **多平台构建**：在 Linux 和 Windows 上验证编译成功。
- **自动化测试**：运行 `go test` 确保代码正确性。
- **Docker 自动发布**：配置了 `.github/workflows/docker-publish.yml`，当推送 Tag (如 `v1.0.0`) 时自动构建并推送到 Docker Hub。

---

## 🛡️ 错误处理

| 场景 | 策略 |
|------|------|
| 网络请求失败 | 单个失败不影响其他频道，记录错误日志 |
| API 频率限制 | 建议 ≥30秒 间隔 |
| 主播信息缺失 | 降级使用 Channel ID |
| 头像加载失败 | 显示默认 SVG 图标 |

---

## 📁 目录结构

```
LiveChannelsCN/
├── main.go                    # 入口
├── config.json                # 配置（gitignore）
├── Dockerfile                 # 容器构建
├── docker-compose.yml         # 编排
├── internal/
│   ├── api/router.go          # HTTP 路由
│   ├── config/config.go       # 配置加载
│   ├── logger/
│   │   └── logger.go          # Zap 日志封装
│   ├── models/
│   │   ├── models.go          # 数据结构
│   │   └── models_test.go     # 单元测试
│   ├── platform/
│   │   ├── factory.go         # 平台工厂
│   │   ├── client.go          # HTTP 客户端单例
│   │   ├── bilibili.go        # B站客户端
│   │   ├── douyu.go           # 斗鱼客户端
│   │   └── huya.go            # 虎牙客户端
│   └── service/
│       ├── stream_service.go  # 业务逻辑
│       └── stream_service_test.go
└── web/
    └── index.html             # Go Template
```

---

## 🌐 外部 API

### Bilibili
- **房间信息**：`api.live.bilibili.com/room/v1/Room/get_info?room_id={id}`
- **主播信息**：`api.live.bilibili.com/live_user/v1/UserInfo/get_anchor_in_room?roomid={id}`

### Douyu
- **房间信息**：`open.douyu.com/api/RoomApi/room/{room_id}`

### Huya
- **房间信息**：`www.huya.com/cache.php?m=LiveList&do=getLiveListByPage&tagAll={room_id}`

---

## 🔮 开发规范

### 代码规范
```bash
# 格式化
go fmt ./...

# 测试
go test -v ./...

# 构建
go build -o live-channels.exe

# 运行 (支持参数)
./live-channels.exe -level debug -config ./my-config.json -mode release -port 8081
```

### 添加新功能
1. 阅读本文档理解架构
2. 遵循分层设计原则
3. 在对应层级添加代码
4. 编写单元测试
5. 更新本文档（如有架构变更）

---

## 📊 当前状态

**版本**：v0.9.3 (Quality)

**已完成优化**：
- ✅ 支持 `CONFIG_PATH` 环境变量
- ✅ 全局共享 HTTP 客户端 (Resty + 连接池)
- ✅ Service 层代码重构 (Worker Pool 并发)
- ✅ 内存缓存机制 (TTL + 智能降级 + 参数化)
- ✅ 结构化日志 (Zap)
- ✅ 移除无效测试目录
- ✅ **全面提升单元测试覆盖率**：
    - `internal/api`: 路由、健康检查、缓存参数解析测试
    - `internal/config`: 配置加载、错误处理测试
    - `internal/models`: 平台验证逻辑测试
    - `internal/platform`: 工厂模式测试
    - `internal/service`: 排序逻辑、配置覆盖逻辑测试
- ✅ GitHub Actions CI/CD (自动测试、Docker 自动发布)

**待办**：
- [ ] 更多平台支持（抖音、快手）
- [ ] WebSocket 实时推送
- [ ] 主播分组管理

---

## 📄 许可证

MIT License
