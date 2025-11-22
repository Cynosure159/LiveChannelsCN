package main

import (
	"live-channels/internal/api"
	"live-channels/internal/config"
	"log"
)

func main() {
	// 初始化配置
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 启动 API 服务器
	router := api.SetupRouter(cfg)

	if err := router.Run(":12301"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
