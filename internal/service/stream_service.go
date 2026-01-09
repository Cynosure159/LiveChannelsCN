package service

import (
	"live-channels/internal/models"
	"live-channels/internal/platform"
	"sort"
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

// 默认 Worker 数量
const DefaultWorkerCount = 10

// fetchStreamStatuses 使用 Worker Pool 并发获取频道列表的直播状态
func (s *StreamService) fetchStreamStatuses(channels []models.ChannelConfig) []models.StreamStatus {
	// 如果频道数量少于 Worker 数量，就用频道数量，避免启动多余 Goroutine
	workerCount := DefaultWorkerCount
	if len(channels) < workerCount {
		workerCount = len(channels)
	}
	if workerCount == 0 {
		return []models.StreamStatus{}
	}

	jobs := make(chan models.ChannelConfig, len(channels))
	results := make(chan *models.StreamStatus, len(channels))

	// 启动 Workers
	for w := 0; w < workerCount; w++ {
		go s.worker(jobs, results)
	}

	// 发送任务
	for _, ch := range channels {
		jobs <- ch
	}
	close(jobs)

	// 收集结果
	var statuses []models.StreamStatus
	for i := 0; i < len(channels); i++ {
		status := <-results
		if status != nil {
			statuses = append(statuses, *status)
		}
	}

	// 排序
	s.sortStreamStatus(statuses)

	return statuses
}

// worker 处理具体的获取任务
func (s *StreamService) worker(jobs <-chan models.ChannelConfig, results chan<- *models.StreamStatus) {
	for ch := range jobs {
		provider := platform.CreateProvider(ch.Platform)
		if provider == nil {
			results <- nil
			continue
		}

		status, err := provider.GetStreamStatus(ch.ChannelID)
		if err != nil {
			// TODO: Add logging here
			results <- nil
			continue
		}

		if status != nil {
			// Apply config name override if present
			if ch.Name != "" {
				status.Name = ch.Name
			}
		}
		results <- status
	}
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
