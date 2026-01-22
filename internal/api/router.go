package api

import (
	"live-channels/internal/models"
	"live-channels/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter(cfg *models.Config) *gin.Engine {
	router := gin.Default()

	// 添加 CORS 中间件
	router.Use(corsMiddleware())

	// 创建服务
	streamService := service.NewStreamService(cfg)

	// 加载 HTML 模板
	router.LoadHTMLGlob("./web/*.html")

	// 提供静态文件
	router.Static("/web", "./web")

	// 提供 index.html，并带上主播数据
	router.GET("/", func(c *gin.Context) {
		cacheDuration := getCacheDuration(c)
		statuses, err := streamService.GetAllStreamStatus(cacheDuration)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		if statuses == nil {
			statuses = []models.StreamStatus{}
		}

		c.Header("Widget-Content-Type", "html")
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Channels":      statuses,
			"CollapseAfter": c.DefaultQuery("collapse", "10"),
		})
	})

	// 获取所有直播状态
	router.GET("/api/streams", func(c *gin.Context) {
		cacheDuration := getCacheDuration(c)
		statuses, err := streamService.GetAllStreamStatus(cacheDuration)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		if statuses == nil {
			statuses = []models.StreamStatus{}
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Status: "success",
			Data:   statuses,
		})
	})

	// 获取指定平台的直播状态
	router.GET("/api/streams/:platform", func(c *gin.Context) {
		platformStr := c.Param("platform")
		platformType := models.Platform(platformStr)

		// 验证平台
		if !platformType.IsValid() {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Status:  "error",
				Message: "invalid platform",
			})
			return
		}

		cacheDuration := getCacheDuration(c)
		statuses, err := streamService.GetStreamStatusByPlatform(platformType, cacheDuration)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		if statuses == nil {
			statuses = []models.StreamStatus{}
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Status: "success",
			Data:   statuses,
		})
	})

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	return router
}

// getCacheDuration 从请求参数获取缓存时间，默认 60s
func getCacheDuration(c *gin.Context) time.Duration {
	cacheSecondsStr := c.DefaultQuery("cache", "60")
	cacheSeconds, err := strconv.Atoi(cacheSecondsStr)
	if err != nil || cacheSeconds < 0 {
		cacheSeconds = 60
	}
	return time.Duration(cacheSeconds) * time.Second
}

// corsMiddleware CORS 中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
