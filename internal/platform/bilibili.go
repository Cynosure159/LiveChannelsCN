package platform

import (
	"encoding/json"
	"fmt"
	"live-channels/internal/models"
	"time"

	"github.com/go-resty/resty/v2"
)

// BilibiliResponse 直播间信息响应
type BilibiliResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		LiveStatus  int    `json:"live_status"`
		Title       string `json:"title"`
		RoomID      int    `json:"room_id"`
		OnlineCount int    `json:"online"`
		KeyFrame    string `json:"keyframe"`
		UserInfo    struct {
			Info struct {
				Uname string `json:"uname"`
			} `json:"info"`
		} `json:"user_info"`
	} `json:"data"`
}

// BilibiliAnchorResponse 主播详细信息响应
type BilibiliAnchorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Info struct {
			Uname string `json:"uname"`
			Face  string `json:"face"`
		} `json:"info"`
	} `json:"data"`
}

// AnchorInfo 主播信息
type AnchorInfo struct {
	Uname string // 主播名字
	Face  string // 主播头像 URL
}

// BilibiliClient Bilibili 平台客户端
type BilibiliClient struct {
	client *resty.Client
}

// NewBilibiliClient 创建 Bilibili 客户端
func NewBilibiliClient() *BilibiliClient {
	return &BilibiliClient{
		client: GetHTTPClient(),
	}
}

// GetStreamStatus 获取 B 站直播状态
// API: https://api.live.bilibili.com/room/v1/Room/get_info?room_id={roomId}
func (b *BilibiliClient) GetStreamStatus(channelID string) (*models.StreamStatus, error) {
	url := fmt.Sprintf("https://api.live.bilibili.com/room/v1/Room/get_info?room_id=%s", channelID)

	// 获取直播间基本信息
	resp, err := b.client.R().
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch bilibili room info: %w", err)
	}

	var biliResp BilibiliResponse
	if err := json.Unmarshal(resp.Body(), &biliResp); err != nil {
		return nil, fmt.Errorf("failed to parse bilibili response: %w", err)
	}

	if biliResp.Code != 0 {
		return nil, fmt.Errorf("bilibili api error: %s", biliResp.Message)
	}

	isLive := biliResp.Data.LiveStatus == 1
	roomID := biliResp.Data.RoomID

	// 获取主播详细信息（名字和头像）
	anchorInfo, err := b.getAnchorInfo(roomID)
	if err != nil {
		// 如果获取主播信息失败，使用channelID作为降级方案
		anchorInfo = &AnchorInfo{
			Uname: channelID,
			Face:  "",
		}
	}

	status := &models.StreamStatus{
		ChannelID:    channelID,
		Name:         anchorInfo.Uname,
		Platform:     "bilibili",
		IsLive:       isLive,
		Title:        biliResp.Data.Title,
		Viewers:      biliResp.Data.OnlineCount,
		ThumbnailURL: biliResp.Data.KeyFrame,
		AvatarURL:    anchorInfo.Face,
		ProfileURL:   fmt.Sprintf("https://live.bilibili.com/%s", channelID),
		UpdatedAt:    time.Now().Unix(),
	}

	return status, nil
}

// getAnchorInfo 获取主播详细信息
// API: https://api.live.bilibili.com/live_user/v1/UserInfo/get_anchor_in_room
func (b *BilibiliClient) getAnchorInfo(roomID int) (*AnchorInfo, error) {
	url := fmt.Sprintf("https://api.live.bilibili.com/live_user/v1/UserInfo/get_anchor_in_room?roomid=%d", roomID)

	resp, err := b.client.R().
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch anchor info: %w", err)
	}

	var anchorResp BilibiliAnchorResponse
	if err := json.Unmarshal(resp.Body(), &anchorResp); err != nil {
		return nil, fmt.Errorf("failed to parse anchor response: %w", err)
	}

	if anchorResp.Code != 0 {
		return nil, fmt.Errorf("bilibili anchor api error: %s", anchorResp.Message)
	}

	return &AnchorInfo{
		Uname: anchorResp.Data.Info.Uname,
		Face:  anchorResp.Data.Info.Face,
	}, nil
}
