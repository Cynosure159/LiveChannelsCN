package service

import (
	"live-channels/internal/models"
	"testing"
)

func TestNewStreamService(t *testing.T) {
	cfg := &models.Config{
		Channels: []models.ChannelConfig{
			{
				Platform:  models.PlatformBilibili,
				ChannelID: "123",
				Name:      "test",
			},
		},
	}

	service := NewStreamService(cfg)
	if service == nil {
		t.Fatal("NewStreamService returned nil")
	}

	if service.config != cfg {
		t.Fatal("Config not set correctly")
	}
}
