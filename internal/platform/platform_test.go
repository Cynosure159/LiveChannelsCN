package platform

import (
	"live-channels/internal/models"
	"testing"
)

func TestCreateProvider(t *testing.T) {
	tests := []struct {
		platform models.Platform
		wantNil  bool
	}{
		{models.PlatformBilibili, false},
		{models.PlatformDouyu, false},
		{models.PlatformHuya, false},
		{models.Platform("invalid"), true},
	}

	for _, tt := range tests {
		t.Run(string(tt.platform), func(t *testing.T) {
			got := CreateProvider(tt.platform)
			if (got == nil) != tt.wantNil {
				t.Errorf("CreateProvider() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}
