package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// 创建临时测试文件
	tmpDir := t.TempDir()
	validConfigFile := filepath.Join(tmpDir, "valid_config.json")
	invalidJsonFile := filepath.Join(tmpDir, "invalid_json.json")

	validContent := `{
		"channels": [
			{"platform": "bilibili", "channel_id": "123", "name": "test1"}
		],
		"user_agent": "test-ua"
	}`
	err := os.WriteFile(validConfigFile, []byte(validContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(invalidJsonFile, []byte(`{"channels": [`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "Valid Config",
			path:    validConfigFile,
			wantErr: false,
		},
		{
			name:    "File Not Found",
			path:    filepath.Join(tmpDir, "not_exists.json"),
			wantErr: true,
		},
		{
			name:    "Invalid JSON",
			path:    invalidJsonFile,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := LoadConfig(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && cfg == nil {
				t.Errorf("LoadConfig() returned nil config for valid file")
			}
			if !tt.wantErr && cfg.UserAgent != "test-ua" {
				t.Errorf("LoadConfig() UserAgent = %v, want %v", cfg.UserAgent, "test-ua")
			}
		})
	}
}
