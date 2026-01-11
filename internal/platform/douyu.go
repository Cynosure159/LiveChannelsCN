package platform

import (
	"encoding/json"
	"fmt"
	"live-channels/internal/models"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

// DouyuResponse 斗鱼直播间信息响应
type DouyuResponse struct {
	Room struct {
		ShowStatus int    `json:"show_status"` // 直播状态：1 为在直播，0 为离线
		OwnerName  string `json:"owner_name"`  // 主播名字
		AvatarMid  string `json:"avatar_mid"`  // 主播头像（中等尺寸）
		RoomName   string `json:"room_name"`   // 直播间名称
		RoomPic    string `json:"room_pic"`    // 直播间封面
		VideoLoop  int    `json:"videoLoop"`   // 视频循环标志：0 为正常直播，1 为循环播放
		RoomBizAll struct {
			Hot string `json:"hot"` // 在线观众数
		} `json:"room_biz_all"`
	} `json:"room"`
}

// DouyuClient 斗鱼平台客户端
type DouyuClient struct {
	client *resty.Client
}

// NewDouyuClient 创建斗鱼客户端
func NewDouyuClient() *DouyuClient {
	return &DouyuClient{
		client: GetHTTPClient(),
	}
}

// GetStreamStatus 获取斗鱼直播状态
// API: https://www.douyu.com/betard/{roomId}
func (d *DouyuClient) GetStreamStatus(channelID string) (*models.StreamStatus, error) {
	url := fmt.Sprintf("https://www.douyu.com/betard/%s", channelID)

	// 获取直播间信息
	resp, err := d.client.R().
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch douyu room info: %w", err)
	}

	var douyuResp DouyuResponse
	if err := json.Unmarshal(resp.Body(), &douyuResp); err != nil {
		return nil, fmt.Errorf("failed to parse douyu response: %w", err)
	}

	// 判断是否直播：show_status == 1 且 videoLoop == 0
	isLive := douyuResp.Room.ShowStatus == 1 && douyuResp.Room.VideoLoop == 0
	viewers, _ := strconv.Atoi(douyuResp.Room.RoomBizAll.Hot)

	status := &models.StreamStatus{
		ChannelID:    channelID,
		Name:         douyuResp.Room.OwnerName,
		Platform:     "douyu",
		IsLive:       isLive,
		Title:        douyuResp.Room.RoomName,
		Viewers:      viewers,
		ThumbnailURL: douyuResp.Room.RoomPic,
		AvatarURL:    douyuResp.Room.AvatarMid,
		ProfileURL:   fmt.Sprintf("https://www.douyu.com/%s", channelID),
		UpdatedAt:    time.Now().Unix(),
	}

	return status, nil
}
