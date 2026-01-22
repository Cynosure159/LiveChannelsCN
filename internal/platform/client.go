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
	userAgent  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
)

// SetUserAgent 设置全局 User-Agent
func SetUserAgent(ua string) {
	if ua != "" {
		userAgent = ua
	}
}

// GetHTTPClient 返回共享的 HTTP 客户端单例
func GetHTTPClient() *resty.Client {
	once.Do(func() {
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
			SetHeader("User-Agent", userAgent). // 设置全局 UA
			SetTimeout(5 * time.Second).
			SetRetryCount(2).
			SetRetryWaitTime(500 * time.Millisecond).
			SetRetryMaxWaitTime(2 * time.Second).
			AddRetryCondition(func(r *resty.Response, err error) bool {
				return err != nil || r.StatusCode() >= 500
			})
	})
	return httpClient
}
