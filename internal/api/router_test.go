package api

import (
	"live-channels/internal/models"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	// 切换到项目根目录运行测试，以便 SetupRouter 能找到 ./web 目录
	// 或者手动创建模拟目录
	gin.SetMode(gin.TestMode)

	// 这里假设我们在项目根目录运行，如果不是，尝试向上找
	originalWd, _ := os.Getwd()
	for i := 0; i < 3; i++ {
		if _, err := os.Stat("web"); err == nil {
			break
		}
		os.Chdir("..")
	}

	code := m.Run()

	os.Chdir(originalWd)
	os.Exit(code)
}

func TestHealthCheck(t *testing.T) {
	cfg := &models.Config{}
	router := SetupRouter(cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", w.Code)
	}

	if !strings.Contains(w.Body.String(), `"status":"ok"`) {
		t.Errorf("Body does not contain status ok: %v", w.Body.String())
	}
}

func TestInvalidPlatformAPI(t *testing.T) {
	cfg := &models.Config{}
	router := SetupRouter(cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/streams/invalid_platform", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest, got %v", w.Code)
	}

	if !strings.Contains(w.Body.String(), "invalid platform") {
		t.Errorf("Body content mismatch: %v", w.Body.String())
	}
}

func TestGetCacheDuration(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected int // seconds
	}{
		{"Default", "", 60},
		{"Valid", "?cache=120", 120},
		{"Invalid", "?cache=abc", 60},
		{"Negative", "?cache=-10", 60},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			url := "/" + tt.query
			c.Request, _ = http.NewRequest("GET", url, nil)

			duration := getCacheDuration(c)
			if int(duration.Seconds()) != tt.expected {
				t.Errorf("getCacheDuration() = %v, want %v", duration.Seconds(), tt.expected)
			}
		})
	}
}
