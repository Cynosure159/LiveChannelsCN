package platform

import (
	"encoding/json"
	"fmt"
	"live-channels/internal/models"
	"time"

	"github.com/go-resty/resty/v2"
)

type HuyaResponse struct {
	Code    int `json:"code"`
	Message string `json:"message"`
	Data struct {
		IsLive      int `json:"isLive"`
		RoomName    string `json:"roomName"`
		NickName    string `json:"nickName"`
		TotalCount  int `json:"totalCount"`
		RoomPic     string `json:"roomPic"`
		ProfileRoom string `json:"profileRoom"`
	} `json:"data"`
}

type HuyaClient struct {
	client *resty.Client
}

func NewHuyaClient() *HuyaClient {
	return &HuyaClient{
		client: resty.New(),
	}
}

// GetStreamStatus 获取虎牙直播状态
func (h *HuyaClient) GetStreamStatus(channelID string) (*models.StreamStatus, error) {
	url := fmt.Sprintf("https://www.huya.com/cache.php?m=LiveRoom&do=getLiveRoomInfo&roomid=%s", channelID)

	resp, err := h.client.R().
		SetHeader("User-Agent", "Mozilla/5.0").
		Get(url)

	if err != nil {
		return nil, err
	}

	var huyaResp HuyaResponse
	err = json.Unmarshal(resp.Body(), &huyaResp)
	if err != nil {
		return nil, err
	}

	if huyaResp.Code != 200 {
		return nil, fmt.Errorf("huya api error: %s", huyaResp.Message)
	}

	isLive := huyaResp.Data.IsLive == 1

	status := &models.StreamStatus{
		ChannelID:    channelID,
		Name:         huyaResp.Data.NickName,
		Platform:     "huya",
		IsLive:       isLive,
		Title:        huyaResp.Data.RoomName,
		Viewers:      huyaResp.Data.TotalCount,
		ThumbnailURL: huyaResp.Data.RoomPic,
		ProfileURL:   fmt.Sprintf("https://www.huya.com/%s", channelID),
		UpdatedAt:    time.Now().Unix(),
	}

	return status, nil
}
