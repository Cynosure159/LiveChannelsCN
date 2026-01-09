package platform

import (
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
		httpClient = resty.New().
			SetTimeout(10 * time.Second).
			SetRetryCount(3)
	})
	return httpClient
}
