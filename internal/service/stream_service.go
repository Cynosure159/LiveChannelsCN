package service

import (
	"live-channels/internal/models"
	"live-channels/internal/platform"
	"sort"
	"sync"
)

// StreamService 直播服务
type StreamService struct {
	config *models.Config
}

// NewStreamService 创建直播服务
func NewStreamService(config *models.Config) *StreamService {
	return &StreamService{
		config: config,
	}
}

// GetAllStreamStatus 获取所有直播状态
func (s *StreamService) GetAllStreamStatus() ([]models.StreamStatus, error) {
	return s.fetchStreamStatuses(s.config.Channels), nil
}

// GetStreamStatusByPlatform 获取指定平台的直播状态
func (s *StreamService) GetStreamStatusByPlatform(platformType models.Platform) ([]models.StreamStatus, error) {
	var targetChannels []models.ChannelConfig
	for _, channel := range s.config.Channels {
		if channel.Platform == platformType {
			targetChannels = append(targetChannels, channel)
		}
	}
	return s.fetchStreamStatuses(targetChannels), nil
}

// fetchStreamStatuses 并发获取频道列表的直播状态
func (s *StreamService) fetchStreamStatuses(channels []models.ChannelConfig) []models.StreamStatus {
	var statuses []models.StreamStatus
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, channel := range channels {
		wg.Add(1)
		go func(ch models.ChannelConfig) {
			defer wg.Done()

			provider := platform.CreateProvider(ch.Platform)
			if provider == nil {
				return
			}

			status, err := provider.GetStreamStatus(ch.ChannelID)
			if err != nil {
				// TODO: Add logging here
				return
			}

			if status != nil {
				// Apply config name override if present
				if ch.Name != "" {
					status.Name = ch.Name
				}

				mu.Lock()
				statuses = append(statuses, *status)
				mu.Unlock()
			}
		}(channel)
	}

	wg.Wait()

	// 排序：先按直播状态（在直播的在前），再按观众数量（多的在前）
	s.sortStreamStatus(statuses)

	return statuses
}

// sortStreamStatus 对直播状态进行排序
// 排序规则：1. 在直播的在前，不在直播的在后；2. 同状态下按观众数量多的在前
func (s *StreamService) sortStreamStatus(statuses []models.StreamStatus) {
	sort.Slice(statuses, func(i, j int) bool {
		// 首先按直播状态排序：IsLive=true 的排在前面
		if statuses[i].IsLive != statuses[j].IsLive {
			return statuses[i].IsLive // true > false
		}

		// 如果直播状态相同，按观众数量排序：数量多的在前
		return statuses[i].Viewers > statuses[j].Viewers
	})
}
