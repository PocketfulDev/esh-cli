package utils

import (
	"testing"
)

func TestIsVersionValid(t *testing.T) {
	tests := []struct {
		version string
		hotFix  bool
		want    bool
	}{
		{"1.2", false, true},
		{"10.25", false, true},
		{"1.2.3", false, false},
		{"1.2-1.0", true, true},
		{"1.2-1", false, false},
		{"invalid", false, false},
	}

	for _, tt := range tests {
		got := IsVersionValid(tt.version, tt.hotFix)
		if got != tt.want {
			t.Errorf("IsVersionValid(%q, %t) = %t, want %t", tt.version, tt.hotFix, got, tt.want)
		}
	}
}

func TestIsTagValid(t *testing.T) {
	tests := []struct {
		tag  string
		want bool
	}{
		{"stg6_1.2-0", true},
		{"stg6_1.2-1.0", true},
		{"service_stg6_1.2-0", true},
		{"invalid_env_1.2-0", false},
		{"stg6_1.2.3-0", false},
		{"stg6_1.2", false},
		{"invalid", false},
	}

	for _, tt := range tests {
		got := IsTagValid(tt.tag)
		if got != tt.want {
			t.Errorf("IsTagValid(%q) = %t, want %t", tt.tag, got, tt.want)
		}
	}
}

func TestIncrementTag(t *testing.T) {
	tests := []struct {
		tag    string
		hotFix bool
		want   string
	}{
		{"stg6_1.2-0", false, "stg6_1.2-1"},
		{"stg6_1.2-0", true, "stg6_1.2-0.1"},
		{"stg6_1.2-1.2", true, "stg6_1.2-1.3"},
		{"service_stg6_1.2-0", false, "service_stg6_1.2-1"},
		{"", false, ""},            // empty tag
		{"invalid_tag", false, ""}, // invalid format
		{"dev_0.0.1", false, ""},   // invalid format (no dash)
	}

	for _, tt := range tests {
		got := IncrementTag(tt.tag, tt.hotFix)
		if got != tt.want {
			t.Errorf("IncrementTag(%q, %t) = %q, want %q", tt.tag, tt.hotFix, got, tt.want)
		}
	}
}

func TestGetEnvFromTag(t *testing.T) {
	tests := []struct {
		tag     string
		want    string
		wantErr bool
	}{
		{"stg6_1.2-0", "stg6", false},
		{"service_stg6_1.2-0", "stg6", false},
		{"invalid", "", true},
		{"too_many_parts_here_test", "", true},
	}

	for _, tt := range tests {
		got, err := GetEnvFromTag(tt.tag)
		if (err != nil) != tt.wantErr {
			t.Errorf("GetEnvFromTag(%q) error = %v, wantErr %v", tt.tag, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("GetEnvFromTag(%q) = %q, want %q", tt.tag, got, tt.want)
		}
	}
}

func TestIsReleaseBranch(t *testing.T) {
	tests := []struct {
		branch string
		want   bool
	}{
		{"release_1.2", true},
		{"release_10.25", true},
		{"master", false},
		{"main", false},
		{"feature/test", false},
		{"release_1", false},
	}

	for _, tt := range tests {
		got := IsReleaseBranch(tt.branch)
		if got != tt.want {
			t.Errorf("IsReleaseBranch(%q) = %t, want %t", tt.branch, got, tt.want)
		}
	}
}

func TestTagPrefix(t *testing.T) {
	tests := []struct {
		env     string
		version string
		service string
		want    string
	}{
		{"stg6", "1.2", "", "stg6_1.2"},
		{"stg6", "1.2", "myservice", "myservice_stg6_1.2"},
		{"production2", "2.1", "api", "api_production2_2.1"},
	}

	for _, tt := range tests {
		got := TagPrefix(tt.env, tt.version, tt.service)
		if got != tt.want {
			t.Errorf("TagPrefix(%q, %q, %q) = %q, want %q", tt.env, tt.version, tt.service, got, tt.want)
		}
	}
}
