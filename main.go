package main

import (
	"live-channels/internal/api"
	"live-channels/internal/config"
	"live-channels/internal/logger"
	"os"

	"go.uber.org/zap"
)

func main() {
	// 初始化日志
	logger.Init("dev")
	defer logger.Sync()

	// 初始化配置
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config/config.json"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		// 使用 Zap Fatal
		logger.Fatal("Failed to load config", zap.String("path", configPath), zap.Error(err))
	}

	// 启动 API 服务器
	router := api.SetupRouter(cfg)

	if err := router.Run(":8081"); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
