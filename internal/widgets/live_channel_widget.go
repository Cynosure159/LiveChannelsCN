package widgets

import (
	"live-channels/internal/models"
	"time"
)

// LiveChannelWidget 直播频道组件
type LiveChannelWidget struct {
	Channels     []ChannelItem `json:"channels"`
	CollapseAfter int          `json:"collapse_after"`
}

// ChannelItem 频道信息项
type ChannelItem struct {
	// 基础信息
	Name        string `json:"name"`
	Platform    string `json:"platform"`
	Exists      bool   `json:"exists"`
	ProfileURL  string `json:"profile_url"`

	// 直播状态
	IsLive     bool   `json:"is_live"`
	Title      string `json:"title"`
	Category   string `json:"category"`
	CategoryURL string `json:"category_url"`

	// 媒体信息
	AvatarURL    string `json:"avatar_url"`
	ThumbnailURL string `json:"thumbnail_url"`

	// 实时数据
	ViewersCount int       `json:"viewers_count"`
	LiveSince    string    `json:"live_since"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewLiveChannelWidget 创建直播频道组件
func NewLiveChannelWidget(statuses []models.StreamStatus) *LiveChannelWidget {
	channels := make([]ChannelItem, 0, len(statuses))

	for _, status := range statuses {
		item := ChannelItem{
			Name:         status.Name,
			Platform:     status.Platform,
			Exists:       true,
			ProfileURL:   status.ProfileURL,
			IsLive:       status.IsLive,
			Title:        status.Title,
			AvatarURL:    status.ProfileURL, // 可以设置为真实头像 URL
			ThumbnailURL: status.ThumbnailURL,
			ViewersCount: status.Viewers,
			UpdatedAt:    time.Unix(status.UpdatedAt, 0),
		}

		if status.IsLive {
			item.LiveSince = formatLiveTime(time.Now(), item.UpdatedAt)
		}

		channels = append(channels, item)
	}

	return &LiveChannelWidget{
		Channels:      channels,
		CollapseAfter: 5,
	}
}

// formatLiveTime 格式化直播时间
func formatLiveTime(now, liveTime time.Time) string {
	duration := now.Sub(liveTime)

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60

	if hours > 0 {
		return time.Now().Add(-duration).Format("15:04")
	}
	if minutes > 0 {
		return time.Now().Add(-duration).Format("15:04")
	}
	return "刚开播"
}
