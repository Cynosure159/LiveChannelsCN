package platform

import (
	"live-channels/internal/models"
)

// StreamProvider 直播平台接口
type StreamProvider interface {
	GetStreamStatus(channelID string) (*models.StreamStatus, error)
}

// Factory 工厂函数
func CreateProvider(platform models.Platform) StreamProvider {
	switch platform {
	case models.PlatformBilibili:
		return NewBilibiliClient()
	case models.PlatformDouyu:
		return NewDouyuClient()
	case models.PlatformHuya:
		return NewHuyaClient()
	default:
		return nil
	}
}
