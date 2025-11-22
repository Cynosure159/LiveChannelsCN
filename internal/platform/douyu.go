package platform

import (
	"fmt"
	"live-channels/internal/models"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

type DouyuClient struct {
	client *resty.Client
}

func NewDouyuClient() *DouyuClient {
	return &DouyuClient{
		client: resty.New(),
	}
}

// GetStreamStatus 获取斗鱼直播状态
func (d *DouyuClient) GetStreamStatus(channelID string) (*models.StreamStatus, error) {
	url := fmt.Sprintf("https://www.douyu.com/betard/%s", channelID)

	resp, err := d.client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		Get(url)

	if err != nil {
		return nil, err
	}

	body := resp.Body()

	// 使用 gjson 快速访问字段
	ownerName := gjson.GetBytes(body, "room.owner_name").String()
	ownerAvatar := gjson.GetBytes(body, "room.avatar_mid").String()
	showStatus := gjson.GetBytes(body, "room.show_status").Int()
	videoLoop := gjson.GetBytes(body, "room.videoLoop").Int()
	roomName := gjson.GetBytes(body, "room.room_name").String()
	roomThumb := gjson.GetBytes(body, "room.room_pic").String()
	onlineCount := gjson.GetBytes(body, "room.room_biz_all.hot").Int()

	log.Println("ownerName: %s", ownerName)
	log.Println("roomName: %s", roomName)
	log.Println("showStatus: %s", showStatus)

	// 判断是否直播：show_status == 1 且 videoLoop == 0
	isLive := showStatus == 1 && videoLoop == 0

	status := &models.StreamStatus{
		ChannelID:    channelID,
		Name:         ownerName,
		Platform:     "douyu",
		IsLive:       isLive,
		Title:        roomName,
		Viewers:      int(onlineCount),
		ThumbnailURL: roomThumb,
		AvatarURL:   ownerAvatar,
		ProfileURL:   fmt.Sprintf("https://www.douyu.com/%s", channelID),
		UpdatedAt:    time.Now().Unix(),
	}

	return status, nil
}
