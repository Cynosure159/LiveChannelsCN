package models

import (
	"testing"
)

func TestPlatformValidation(t *testing.T) {
	tests := []struct {
		platform Platform
		valid    bool
	}{
		{PlatformBilibili, true},
		{PlatformDouyu, true},
		{PlatformHuya, true},
		{Platform("invalid"), false},
	}

	validPlatforms := []Platform{PlatformBilibili, PlatformDouyu, PlatformHuya}

	for _, tt := range tests {
		found := false
		for _, vp := range validPlatforms {
			if vp == tt.platform {
				found = true
				break
			}
		}

		if found != tt.valid {
			t.Errorf("Platform %s validation failed", tt.platform)
		}
	}
}
