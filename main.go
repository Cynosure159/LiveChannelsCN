package main

import (
	"live-channels/internal/api"
	"live-channels/internal/config"
	"log"
	"os"
)

func main() {
	// 初始化配置
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config/config.json"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config from %s: %v", configPath, err)
	}

	// 启动 API 服务器
	router := api.SetupRouter(cfg)

	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
