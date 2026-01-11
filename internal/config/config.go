package config

import (
	"encoding/json"
	"live-channels/internal/models"
	"os"
)

// LoadConfig 从 JSON 文件加载配置
func LoadConfig(filePath string) (*models.Config, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg models.Config
	err = json.Unmarshal(content, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
