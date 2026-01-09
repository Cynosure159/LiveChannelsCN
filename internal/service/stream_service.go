package service

import (
	"live-channels/internal/logger"
	"live-channels/internal/models"
	"live-channels/internal/platform"
	"sort"
	"sync"
	"time"

	"go.uber.org/zap"
)

// StreamService 直播服务
type StreamService struct {
	config  *models.Config
	cache   map[string]cacheItem
	cacheMu sync.RWMutex
}

type cacheItem struct {
	status    *models.StreamStatus
	timestamp time.Time
}

// NewStreamService 创建直播服务
func NewStreamService(config *models.Config) *StreamService {
	return &StreamService{
		config: config,
		cache:  make(map[string]cacheItem),
	}
}

// GetAllStreamStatus 获取所有直播状态
func (s *StreamService) GetAllStreamStatus(cacheDuration time.Duration) ([]models.StreamStatus, error) {
	return s.fetchStreamStatuses(s.config.Channels, cacheDuration), nil
}

// GetStreamStatusByPlatform 获取指定平台的直播状态
func (s *StreamService) GetStreamStatusByPlatform(platformType models.Platform, cacheDuration time.Duration) ([]models.StreamStatus, error) {
	var targetChannels []models.ChannelConfig
	for _, channel := range s.config.Channels {
		if channel.Platform == platformType {
			targetChannels = append(targetChannels, channel)
		}
	}
	return s.fetchStreamStatuses(targetChannels, cacheDuration), nil
}

// 默认 Worker 数量
const DefaultWorkerCount = 10

// fetchStreamStatuses 使用 Worker Pool 并发获取频道列表的直播状态
func (s *StreamService) fetchStreamStatuses(channels []models.ChannelConfig, cacheDuration time.Duration) []models.StreamStatus {
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
		go s.worker(jobs, results, cacheDuration)
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
func (s *StreamService) worker(jobs <-chan models.ChannelConfig, results chan<- *models.StreamStatus, cacheDuration time.Duration) {
	for ch := range jobs {
		// 1. 尝试从缓存获取
		cacheKey := string(ch.Platform) + ":" + ch.ChannelID
		s.cacheMu.RLock()
		item, found := s.cache[cacheKey]
		s.cacheMu.RUnlock()

		if found && time.Since(item.timestamp) < cacheDuration {
			// 缓存命中且未过期
			logger.Debug("Cache Hit",
				zap.String("platform", string(ch.Platform)),
				zap.String("channel_id", ch.ChannelID),
			)
			if item.status != nil {
				// 返回副本以防止外部修改影响缓存
				copiedStatus := *item.status
				results <- &copiedStatus
			} else {
				results <- nil
			}
			continue
		}

		// 2. 缓存未命中或过期，从 Provider 获取
		logger.Debug("Fetching API",
			zap.String("platform", string(ch.Platform)),
			zap.String("channel_id", ch.ChannelID),
		)
		provider := platform.CreateProvider(ch.Platform)
		if provider == nil {
			results <- nil
			continue
		}

		status, err := provider.GetStreamStatus(ch.ChannelID)
		if err != nil {
			// 发生错误时，如果缓存中还有（即使过期），优先返回旧缓存作为容错
			if found && item.status != nil {
				logger.Warn("Using stale cache due to error",
					zap.String("platform", string(ch.Platform)),
					zap.String("channel_id", ch.ChannelID),
					zap.Error(err),
				)
				copiedStatus := *item.status
				results <- &copiedStatus
				continue
			}
			logger.Error("Failed to fetch stream status",
				zap.String("platform", string(ch.Platform)),
				zap.String("channel_id", ch.ChannelID),
				zap.Error(err),
			)
			results <- nil
			continue
		}

		if status != nil {
			// Apply config name override if present
			if ch.Name != "" {
				status.Name = ch.Name
			}

			// 更新缓存
			s.cacheMu.Lock()
			s.cache[cacheKey] = cacheItem{
				status:    status,
				timestamp: time.Now(),
			}
			s.cacheMu.Unlock()
			logger.Debug("Cache Updated",
				zap.String("platform", string(ch.Platform)),
				zap.String("channel_id", ch.ChannelID),
			)
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
