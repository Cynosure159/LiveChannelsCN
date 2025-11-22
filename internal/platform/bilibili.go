package platform

import (
	"encoding/json"
	"fmt"
	"live-channels/internal/models"
	"time"

	"github.com/go-resty/resty/v2"
)

type BilibiliResponse struct {
	Code    int `json:"code"`
	Message string `json:"message"`
	Data struct {
		LiveStatus int `json:"live_status"`
		Title      string `json:"title"`
		RoomID     int `json:"room_id"`
		UserInfo struct {
			Info struct {
				Uname string `json:"uname"`
			} `json:"info"`
		} `json:"user_info"`
		OnlineCount int `json:"online"`
		KeyFrame    string `json:"keyframe"`
	} `json:"data"`
}

type BilibiliAnchorResponse struct {
	Code    int `json:"code"`
	Message string `json:"message"`
	Data struct {
		Info struct {
			Uname string `json:"uname"`
			Face  string `json:"face"`
		} `json:"info"`
	} `json:"data"`
}

type BilibiliClient struct {
	client *resty.Client
}

func NewBilibiliClient() *BilibiliClient {
	return &BilibiliClient{
		client: resty.New(),
	}
}

// GetStreamStatus 获取 B 站直播状态
func (b *BilibiliClient) GetStreamStatus(channelID string) (*models.StreamStatus, error) {
	url := fmt.Sprintf("https://api.live.bilibili.com/room/v1/Room/get_info?room_id=%s", channelID)

	resp, err := b.client.R().
		SetHeader("User-Agent", "Mozilla/5.0").
		Get(url)

	if err != nil {
		return nil, err
	}

	var biliResp BilibiliResponse
	err = json.Unmarshal(resp.Body(), &biliResp)
	if err != nil {
		return nil, err
	}

	if biliResp.Code != 0 {
		return nil, fmt.Errorf("bilibili api error: %s", biliResp.Message)
	}

	isLive := biliResp.Data.LiveStatus == 1
	roomID := biliResp.Data.RoomID

	// 获取主播的详细信息（名字和头像）
	anchorInfo, err := b.getAnchorInfo(roomID)
	if err != nil {
		// 如果获取主播信息失败，使用原始数据
		anchorInfo = &AnchorInfo{
			Uname: biliResp.Data.UserInfo.Info.Uname,
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

	// 如果有头像，添加到 ThumbnailURL 的注释中（或创建新字段）
	if anchorInfo.Face != "" {
		// 如果需要，可以在模型中添加 AvatarURL 字段
		// 这里暂时保留，可以后续扩展
	}

	return status, nil
}

// AnchorInfo 主播信息
type AnchorInfo struct {
	Uname string
	Face  string
}

// getAnchorInfo 获取主播的详细信息
func (b *BilibiliClient) getAnchorInfo(roomID int) (*AnchorInfo, error) {
	url := fmt.Sprintf("https://api.live.bilibili.com/live_user/v1/UserInfo/get_anchor_in_room?roomid=%d", roomID)

	resp, err := b.client.R().
		SetHeader("User-Agent", "Mozilla/5.0").
		Get(url)

	if err != nil {
		return nil, err
	}

	var anchorResp BilibiliAnchorResponse
	err = json.Unmarshal(resp.Body(), &anchorResp)
	if err != nil {
		return nil, err
	}

	if anchorResp.Code != 0 {
		return nil, fmt.Errorf("bilibili anchor api error: %s", anchorResp.Message)
	}

	return &AnchorInfo{
		Uname: anchorResp.Data.Info.Uname,
		Face:  anchorResp.Data.Info.Face,
	}, nil
}
