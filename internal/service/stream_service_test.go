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

func TestApplyConfigOverrides(t *testing.T) {
	service := NewStreamService(&models.Config{})
	status := &models.StreamStatus{
		Name: "Original Name",
	}
	ch := models.ChannelConfig{
		Name: "Overridden Name",
	}

	service.applyConfigOverrides(status, ch)

	if status.Name != "Overridden Name" {
		t.Errorf("Expected Name to be 'Overridden Name', got '%s'", status.Name)
	}

	// 测试空覆盖
	status2 := &models.StreamStatus{Name: "Original"}
	ch2 := models.ChannelConfig{Name: ""}
	service.applyConfigOverrides(status2, ch2)
	if status2.Name != "Original" {
		t.Errorf("Expected Name to stay 'Original', got '%s'", status2.Name)
	}
}

func TestSortStreamStatus(t *testing.T) {
	service := NewStreamService(&models.Config{})
	statuses := []models.StreamStatus{
		{Name: "A", IsLive: false, Viewers: 100},
		{Name: "B", IsLive: true, Viewers: 50},
		{Name: "C", IsLive: true, Viewers: 200},
		{Name: "D", IsLive: false, Viewers: 300},
	}

	service.sortStreamStatus(statuses)

	// Expected order: C (live 200), B (live 50), D (not live 300), A (not live 100)
	expected := []string{"C", "B", "D", "A"}
	for i, name := range expected {
		if statuses[i].Name != name {
			t.Errorf("At index %d: expected %s, got %s", i, name, statuses[i].Name)
		}
	}
}
