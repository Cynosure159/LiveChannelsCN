package platform

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	httpClient *resty.Client
	once       sync.Once
)

// GetHTTPClient 返回共享的 HTTP 客户端单例
func GetHTTPClient() *resty.Client {
	once.Do(func() {
		// 配置自定义 Transport 以优化连接池
		// 禁用长连接(或者减少空闲时间)可以减少 wsarecv 错误，但会降低一点性能
		// 这里选择保守配置：保留长连接，但减少空闲时间和数量
		transport := &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 5,
			IdleConnTimeout:     30 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 5 * time.Second,
		}

		httpClient = resty.New().
			SetTransport(transport).
			SetTimeout(5 * time.Second).              // 降低总超时，避免长时间卡住
			SetRetryCount(2).                         // 减少重试次数
			SetRetryWaitTime(500 * time.Millisecond). // 重试前等待
			SetRetryMaxWaitTime(2 * time.Second).     // 最大等待时间
			// 添加重试条件：只在网络错误或 5xx 错误时重试
			AddRetryCondition(func(r *resty.Response, err error) bool {
				return err != nil || r.StatusCode() >= 500
			})
	})
	return httpClient
}
