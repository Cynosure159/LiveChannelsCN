package models

// Platform 直播平台类型
type Platform string

const (
	PlatformBilibili Platform = "bilibili"
	PlatformDouyu    Platform = "douyu"
	PlatformHuya     Platform = "huya"
)

// ChannelConfig 频道配置
type ChannelConfig struct {
	Platform  Platform `json:"platform"`
	ChannelID string   `json:"channel_id"`
	Name      string   `json:"name"`
}

// Config 应用配置
type Config struct {
	Channels []ChannelConfig `json:"channels"`
}

// StreamStatus 直播状态
type StreamStatus struct {
	ChannelID   string `json:"channel_id"`
	Name        string `json:"name"`
	Platform    string `json:"platform"`
	IsLive      bool   `json:"is_live"`
	Title       string `json:"title"`
	Game        string `json:"game"`
	Viewers     int    `json:"viewers"`
	ThumbnailURL string `json:"thumbnail_url"`
	AvatarURL   string `json:"avatar_url"`
	ProfileURL  string `json:"profile_url"`
	UpdatedAt   int64  `json:"updated_at"`
}

// APIResponse API 响应
type APIResponse struct {
	Status  string         `json:"status"`
	Data    []StreamStatus `json:"data"`
	Message string         `json:"message,omitempty"`
}
