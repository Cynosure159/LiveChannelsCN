package main

import (
	"flag"
	"live-channels/internal/api"
	"live-channels/internal/config"
	"live-channels/internal/logger"
	"live-channels/internal/platform"
	"os"
	"strings"

	"go.uber.org/zap"
)

func main() {
	// 定义命令行参数
	flagLevel := flag.String("level", os.Getenv("LOG_LEVEL"), "日志级别 (debug, info, warn, error)")
	flagMode := flag.String("mode", os.Getenv("GIN_MODE"), "运行模式 (debug, release)")
	flagConfig := flag.String("config", os.Getenv("CONFIG_PATH"), "配置文件路径")
	flagPort := flag.String("port", os.Getenv("PORT"), "服务端口")
	flagUA := flag.String("ua", os.Getenv("USER_AGENT"), "自定义 User-Agent")
	flag.Parse()

	// 1. 确定基本参数与默认值
	logLevel := *flagLevel
	if logLevel == "" {
		// 如果是 go run 运行，默认 debug；否则默认 info
		if isGoRun() {
			logLevel = "debug"
		} else {
			logLevel = "info"
		}
	}

	runMode := *flagMode
	if runMode == "" {
		runMode = "debug"
	}

	configPath := *flagConfig
	if configPath == "" {
		configPath = "./config/config.json"
	}

	port := *flagPort
	if port == "" {
		port = "8081"
	}

	ua := *flagUA

	// 2. 初始化日志
	logMode := "dev"
	if runMode == "release" {
		logMode = "prod"
	}
	logger.Init(logMode, logLevel)
	defer logger.Sync()

	// 3. 加载配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logger.Fatal("Failed to load config", zap.String("path", configPath), zap.Error(err))
	}

	// 3.5 设置全局 User-Agent
	// 优先级：命令行参数/环境变量 > 配置文件 > 内置默认值
	if ua == "" {
		ua = cfg.UserAgent
	}
	platform.SetUserAgent(ua)

	// 4. 启动 API 服务器
	router := api.SetupRouter(cfg)

	logger.Info("Starting server",
		zap.String("port", port),
		zap.String("mode", runMode),
		zap.String("log_level", logLevel),
	)
	if err := router.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

// isGoRun 启发式检测是否是通过 go run 启动
func isGoRun() bool {
	exe, err := os.Executable()
	if err != nil {
		return false
	}
	exe = strings.ToLower(exe)
	// go-build 是 Go 编译器的通用临时标识
	// /tmp/ 是 Unix-like 系统的默认临时目录
	// \temp\ 是 Windows 的默认临时目录
	return strings.Contains(exe, "go-build") ||
		strings.Contains(exe, "/tmp/") ||
		strings.Contains(exe, "\\temp\\")
}
