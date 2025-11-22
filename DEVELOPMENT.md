# 开发指南

## 项目概览

Live Channels 是一个 Glance 扩展组件，用于展示多个直播平台的主播开播状态。项目采用 Go 后端 + 前端网页的架构。

## 开发环境设置

### 前置条件

-   Go 1.21+
-   Git
-   （可选）Docker & Docker Compose

### 快速开始

1. **克隆或初始化仓库**

```bash
cd f:\code\LiveChannelsCN
```

2. **下载依赖**

```bash
go mod download
```

3. **复制配置文件**

```bash
# Windows
copy config.example.json config.json

# Linux/Mac
cp config.example.json config.json
```

4. **编辑配置文件**

在 `config.json` 中添加你要监视的主播信息。

5. **运行开发服务器**

```bash
go run main.go
```

6. **访问应用**

打开浏览器访问 `http://localhost:8080`

## 项目结构详解

```
live-channels/
├── main.go                 # 程序入口点
├── go.mod                  # Go 模块定义
├── config.json             # 应用配置
├── config.example.json     # 配置示例
├── Dockerfile              # Docker 构建文件
├── docker-compose.yml      # Docker Compose 编排
├── Makefile               # 构建脚本
├── README.md              # 项目说明
├── .gitignore             # Git 忽略文件
│
├── internal/              # 内部包（不能被外部导入）
│   ├── models/            # 数据模型
│   │   ├── models.go      # 核心数据结构定义
│   │   └── models_test.go # 模型测试
│   │
│   ├── config/            # 配置管理
│   │   └── config.go      # 配置加载器
│   │
│   ├── platform/          # 直播平台 API 接口
│   │   ├── factory.go     # 平台工厂（策略模式）
│   │   ├── bilibili.go    # B 站 API 客户端
│   │   ├── douyu.go       # 斗鱼 API 客户端
│   │   └── huya.go        # 虎牙 API 客户端
│   │
│   ├── service/           # 业务逻辑层
│   │   ├── stream_service.go      # 直播服务
│   │   └── stream_service_test.go # 服务测试
│   │
│   └── api/               # HTTP API 层
│       └── router.go      # 路由定义
│
└── web/                   # 前端静态文件
    └── index.html         # 前端 UI 页面
```

## 核心模块说明

### Models (`internal/models/models.go`)

定义数据结构：

-   `Platform` - 直播平台枚举
-   `ChannelConfig` - 频道配置
-   `StreamStatus` - 直播状态
-   `APIResponse` - 统一 API 响应格式

### Platform 模块 (`internal/platform/`)

为每个直播平台实现 API 客户端，实现 `StreamProvider` 接口：

```go
type StreamProvider interface {
    GetStreamStatus(channelID string) (*StreamStatus, error)
}
```

**工厂模式**：

-   `Factory.go` 根据平台类型创建对应的客户端
-   各平台客户端实现相同的接口
-   易于扩展新平台

### Service 层 (`internal/service/stream_service.go`)

业务逻辑处理：

-   `GetAllStreamStatus()` - 并发获取所有频道状态
-   `GetStreamStatusByPlatform()` - 获取特定平台的状态

### API 层 (`internal/api/router.go`)

HTTP 路由定义和请求处理：

-   `/api/streams` - 获取所有直播状态
-   `/api/streams/:platform` - 获取特定平台的状态
-   `/health` - 健康检查

## 开发流程

### 添加新的直播平台

1. **创建新平台文件** (`internal/platform/newplatform.go`)

```go
package platform

import (
    "live-channels/internal/models"
    "github.com/go-resty/resty/v2"
)

type NewPlatformClient struct {
    client *resty.Client
}

func NewNewPlatformClient() *NewPlatformClient {
    return &NewPlatformClient{
        client: resty.New(),
    }
}

func (c *NewPlatformClient) GetStreamStatus(channelID string) (*models.StreamStatus, error) {
    // 实现 API 调用逻辑
    // 解析响应
    // 返回 StreamStatus
}
```

2. **在 models.go 中添加平台常量**

```go
const (
    PlatformNewPlatform Platform = "newplatform"
)
```

3. **在 factory.go 中添加工厂逻辑**

```go
case models.PlatformNewPlatform:
    return NewNewPlatformClient()
```

### 修改 API 响应格式

1. 编辑 `internal/models/models.go` 中的 `StreamStatus` 结构体
2. 对应修改各平台的 API 响应解析逻辑
3. 前端调整对应的显示逻辑

### 修改前端 UI

编辑 `web/index.html` 中的 HTML 和 CSS。主要函数：

-   `loadStreams()` - 加载数据
-   `createChannelCard()` - 生成卡片 HTML
-   样式通过 CSS 变量管理，便于主题定制

## 常见任务

### 运行测试

```bash
go test -v ./...
```

### 代码格式化

```bash
go fmt ./...
```

### 代码检查

```bash
go vet ./...
```

### 构建可执行文件

```bash
go build -o live-channels.exe
```

或使用 Makefile：

```bash
make build
make run
```

### 使用 Docker 部署

```bash
docker-compose up -d
```

## API 频率限制和注意事项

### B 站 (Bilibili)

-   API 端点: `https://api.live.bilibili.com/room/v1/Room/get_info`
-   建议间隔: >= 30 秒
-   限制: 无严格限制，但避免高频访问

### 斗鱼 (Douyu)

-   API 端点: `https://open.douyu.com/source/query`
-   建议间隔: >= 30 秒
-   限制: 可能需要 API Key（当前使用公开接口）

### 虎牙 (Huya)

-   API 端点: `https://www.huya.com/cache.php`
-   建议间隔: >= 30 秒
-   限制: 可能触发反爬虫机制，需要合理的 User-Agent

## 错误处理

所有网络请求都应处理以下错误：

-   网络连接错误
-   API 返回错误码
-   响应解析错误
-   超时错误

示例：

```go
resp, err := client.R().Get(url)
if err != nil {
    // 处理网络错误
    return nil, err
}

var response Response
err = json.Unmarshal(resp.Body(), &response)
if err != nil {
    // 处理解析错误
    return nil, err
}

if response.Code != 0 {
    // 处理 API 错误
    return nil, fmt.Errorf("api error: %s", response.Message)
}
```

## 性能优化建议

1. **缓存层**

    - 在服务层实现缓存（内存 or Redis）
    - 减少 API 调用频率

2. **并发控制**

    - 当前使用 goroutine + WaitGroup
    - 可考虑添加连接池限制

3. **前端优化**
    - 缓存响应数据
    - 减少不必要的重新渲染
    - 使用 WebSocket 实现实时更新（可选）

## 问题排查

### 无法获取直播状态

1. 检查网络连接
2. 验证 `channel_id` 是否正确
3. 检查是否触发平台频率限制
4. 查看服务器日志获取详细错误信息

### CORS 错误

前端跨域请求被浏览器阻止。已在 `router.go` 中添加 CORS 中间件，如需调整：

```go
c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
```

### 缩略图无法加载

1. 检查图片 URL 是否有效
2. 前端已添加 fallback 处理（图片加载失败时隐藏）
3. 可能需要调整图片代理或缓存策略

## 贡献指南

1. 创建功能分支 (`git checkout -b feature/your-feature`)
2. 提交变更 (`git commit -am 'Add some feature'`)
3. 推送到分支 (`git push origin feature/your-feature`)
4. 创建 Pull Request

## 许可证

MIT License

## 相关资源

-   [Go 官方文档](https://golang.org/doc/)
-   [Gin Web 框架](https://github.com/gin-gonic/gin)
-   [Resty HTTP 客户端](https://github.com/go-resty/resty)
-   [Glance](https://github.com/glanceapp/glance)
