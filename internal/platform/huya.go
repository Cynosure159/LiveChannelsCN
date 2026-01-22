package platform

import (
	"fmt"
	"live-channels/internal/models"
	"regexp"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

// HuyaClient 虎牙平台客户端
type HuyaClient struct {
	client *resty.Client
}

// NewHuyaClient 创建虎牙客户端
func NewHuyaClient() *HuyaClient {
	return &HuyaClient{
		client: GetHTTPClient(),
	}
}

// GetStreamStatus 获取虎牙直播状态
// 直接访问直播间页面，从 HTML 中解析 JSON 数据
// API: https://www.huya.com/{roomId}
func (h *HuyaClient) GetStreamStatus(channelID string) (*models.StreamStatus, error) {
	url := fmt.Sprintf("https://www.huya.com/%s", channelID)

	// 获取直播间页面
	resp, err := h.client.R().
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch huya room page: %w", err)
	}

	body := string(resp.Body())

	// 从 HTML 中解析出所需字段
	name, err := extractField(body, `"nick":"([^"]*)"`)
	if err != nil {
		return nil, fmt.Errorf("failed to extract nick: %w", err)
	}

	title, err := extractField(body, `"introduction":"([^"]*)"`)
	if err != nil {
		return nil, fmt.Errorf("failed to extract introduction: %w", err)
	}

	isOnStr, err := extractField(body, `"isOn":(\w+)`)
	if err != nil {
		return nil, fmt.Errorf("failed to extract isOn: %w", err)
	}

	viewersStr, err := extractField(body, `"attendeeCount":(\d+)`)
	if err != nil {
		return nil, fmt.Errorf("failed to extract attendeeCount: %w", err)
	}

	avatar, err := extractField(body, `"avatar180":"([^"]*)"`)
	if err != nil {
		// 头像可能不存在，使用空字符串
		avatar = ""
	}

	thumbnail, err := extractField(body, `"screenshot":"([^"]*)"`)
	if err != nil {
		// 缩略图可能不存在，使用空字符串
		thumbnail = ""
	}

	// 转换类型
	isLive := isOnStr == "true"
	viewers, _ := strconv.Atoi(viewersStr)

	status := &models.StreamStatus{
		ChannelID:    channelID,
		Name:         name,
		Platform:     "huya",
		IsLive:       isLive,
		Title:        title,
		Viewers:      viewers,
		ThumbnailURL: thumbnail,
		AvatarURL:    avatar,
		ProfileURL:   fmt.Sprintf("https://www.huya.com/%s", channelID),
		UpdatedAt:    time.Now().Unix(),
	}

	return status, nil
}

// extractField 从 HTML 中使用正则表达式提取字段
func extractField(html, pattern string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("invalid regex pattern: %w", err)
	}

	matches := re.FindStringSubmatch(html)
	if len(matches) < 2 {
		return "", fmt.Errorf("field not found with pattern: %s", pattern)
	}

	return matches[1], nil
}
