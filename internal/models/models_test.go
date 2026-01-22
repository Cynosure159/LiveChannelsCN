package models

import (
	"testing"
)

func TestPlatformValidation(t *testing.T) {
	tests := []struct {
		platform Platform
		expected bool
	}{
		{PlatformBilibili, true},
		{PlatformDouyu, true},
		{PlatformHuya, true},
		{Platform("invalid"), false},
		{Platform(""), false},
	}

	for _, tt := range tests {
		if got := tt.platform.IsValid(); got != tt.expected {
			t.Errorf("Platform %q.IsValid() = %v; want %v", tt.platform, got, tt.expected)
		}
	}
}
