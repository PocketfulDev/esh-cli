package utils

import (
	"fmt"
	"testing"
)

func TestParseSemanticVersion(t *testing.T) {
	tests := []struct {
		version   string
		wantMajor int
		wantMinor int
		wantPatch int
		wantPre   string
		shouldErr bool
	}{
		{"1.2.3", 1, 2, 3, "", false},
		{"0.0.1", 0, 0, 1, "", false},
		{"10.20.30", 10, 20, 30, "", false},
		{"1.2.3-alpha", 1, 2, 3, "alpha", false},
		{"1.2.3-beta.1", 1, 2, 3, "beta.1", false},
		{"v1.2.3", 1, 2, 3, "", false}, // with v prefix
		{"1.2", 0, 0, 0, "", true},     // missing patch
		{"1.2.3.4", 0, 0, 0, "", true}, // too many parts
		{"a.b.c", 0, 0, 0, "", true},   // non-numeric
		{"", 0, 0, 0, "", true},        // empty
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			sv, err := ParseSemanticVersion(tt.version)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("ParseSemanticVersion(%q) expected error, got nil", tt.version)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseSemanticVersion(%q) unexpected error: %v", tt.version, err)
				return
			}

			if sv.Major != tt.wantMajor {
				t.Errorf("ParseSemanticVersion(%q) major = %d, want %d", tt.version, sv.Major, tt.wantMajor)
			}
			if sv.Minor != tt.wantMinor {
				t.Errorf("ParseSemanticVersion(%q) minor = %d, want %d", tt.version, sv.Minor, tt.wantMinor)
			}
			if sv.Patch != tt.wantPatch {
				t.Errorf("ParseSemanticVersion(%q) patch = %d, want %d", tt.version, sv.Patch, tt.wantPatch)
			}
			if sv.Prerelease != tt.wantPre {
				t.Errorf("ParseSemanticVersion(%q) prerelease = %q, want %q", tt.version, sv.Prerelease, tt.wantPre)
			}
		})
	}
}

func TestSemanticVersionString(t *testing.T) {
	tests := []struct {
		sv   SemanticVersion
		want string
	}{
		{SemanticVersion{1, 2, 3, ""}, "1.2.3"},
		{SemanticVersion{0, 0, 1, ""}, "0.0.1"},
		{SemanticVersion{1, 2, 3, "alpha"}, "1.2.3-alpha"},
		{SemanticVersion{1, 2, 3, "beta.1"}, "1.2.3-beta.1"},
	}

	for _, tt := range tests {
		got := tt.sv.String()
		if got != tt.want {
			t.Errorf("SemanticVersion.String() = %q, want %q", got, tt.want)
		}
	}
}

func TestBumpSemanticVersion(t *testing.T) {
	tests := []struct {
		version   string
		bumpType  BumpType
		want      string
		shouldErr bool
	}{
		{"1.2.3", BumpMajor, "2.0.0", false},
		{"1.2.3", BumpMinor, "1.3.0", false},
		{"1.2.3", BumpPatch, "1.2.4", false},
		{"0.0.1", BumpMajor, "1.0.0", false},
		{"0.1.0", BumpMinor, "0.2.0", false},
		{"1.0.0", BumpPatch, "1.0.1", false},
		{"1.2.3-alpha", BumpPatch, "1.2.4", false}, // removes prerelease
		{"invalid", BumpPatch, "", true},
		{"1.2.3", "invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.version+"_"+string(tt.bumpType), func(t *testing.T) {
			got, err := BumpSemanticVersion(tt.version, tt.bumpType)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("BumpSemanticVersion(%q, %q) expected error, got nil", tt.version, tt.bumpType)
				}
				return
			}

			if err != nil {
				t.Errorf("BumpSemanticVersion(%q, %q) unexpected error: %v", tt.version, tt.bumpType, err)
				return
			}

			if got != tt.want {
				t.Errorf("BumpSemanticVersion(%q, %q) = %q, want %q", tt.version, tt.bumpType, got, tt.want)
			}
		})
	}
}

func TestCompareSemanticVersions(t *testing.T) {
	tests := []struct {
		v1        string
		v2        string
		want      int
		shouldErr bool
	}{
		{"1.2.3", "1.2.3", 0, false},  // equal
		{"1.2.3", "1.2.4", -1, false}, // v1 < v2 (patch)
		{"1.2.4", "1.2.3", 1, false},  // v1 > v2 (patch)
		{"1.2.3", "1.3.0", -1, false}, // v1 < v2 (minor)
		{"1.3.0", "1.2.3", 1, false},  // v1 > v2 (minor)
		{"1.2.3", "2.0.0", -1, false}, // v1 < v2 (major)
		{"2.0.0", "1.2.3", 1, false},  // v1 > v2 (major)
		{"0.0.1", "0.0.2", -1, false}, // small versions
		{"invalid", "1.2.3", 0, true}, // invalid v1
		{"1.2.3", "invalid", 0, true}, // invalid v2
	}

	for _, tt := range tests {
		t.Run(tt.v1+"_vs_"+tt.v2, func(t *testing.T) {
			got, err := CompareSemanticVersions(tt.v1, tt.v2)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("CompareSemanticVersions(%q, %q) expected error, got nil", tt.v1, tt.v2)
				}
				return
			}

			if err != nil {
				t.Errorf("CompareSemanticVersions(%q, %q) unexpected error: %v", tt.v1, tt.v2, err)
				return
			}

			if got != tt.want {
				t.Errorf("CompareSemanticVersions(%q, %q) = %d, want %d", tt.v1, tt.v2, got, tt.want)
			}
		})
	}
}

func TestGetVersionFromTag(t *testing.T) {
	tests := []struct {
		tag       string
		want      string
		shouldErr bool
	}{
		{"stg6_1.2.3-1", "1.2.3", false},
		{"demo_0.1.0", "0.1.0", false},
		{"service_stg6_1.2.3-0", "1.2.3", false},
		{"production2_2.1.0-5", "2.1.0", false},
		{"invalid_tag", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			got, err := GetVersionFromTag(tt.tag)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("GetVersionFromTag(%q) expected error, got nil", tt.tag)
				}
				return
			}

			if err != nil {
				t.Errorf("GetVersionFromTag(%q) unexpected error: %v", tt.tag, err)
				return
			}

			if got != tt.want {
				t.Errorf("GetVersionFromTag(%q) = %q, want %q", tt.tag, got, tt.want)
			}
		})
	}
}

func TestBumpTagVersion(t *testing.T) {
	tests := []struct {
		tag       string
		bumpType  BumpType
		env       string
		service   string
		want      string
		shouldErr bool
	}{
		{"stg6_1.2.3-1", BumpMajor, "stg6", "", "stg6_2.0.0-1", false},
		{"stg6_1.2.3-1", BumpMinor, "stg6", "", "stg6_1.3.0-1", false},
		{"stg6_1.2.3-1", BumpPatch, "stg6", "", "stg6_1.2.4-1", false},
		{"service_stg6_1.2.3-1", BumpPatch, "stg6", "service", "service_stg6_1.2.4-1", false},
		{"invalid_tag", BumpPatch, "stg6", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.tag+"_"+string(tt.bumpType), func(t *testing.T) {
			got, err := BumpTagVersion(tt.tag, tt.bumpType, tt.env, tt.service)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("BumpTagVersion(%q, %q, %q, %q) expected error, got nil",
						tt.tag, tt.bumpType, tt.env, tt.service)
				}
				return
			}

			if err != nil {
				t.Errorf("BumpTagVersion(%q, %q, %q, %q) unexpected error: %v",
					tt.tag, tt.bumpType, tt.env, tt.service, err)
				return
			}

			if got != tt.want {
				t.Errorf("BumpTagVersion(%q, %q, %q, %q) = %q, want %q",
					tt.tag, tt.bumpType, tt.env, tt.service, got, tt.want)
			}
		})
	}
}

func TestDetectBumpType(t *testing.T) {
	tests := []struct {
		commits []string
		want    BumpType
	}{
		{[]string{"feat: add new feature"}, BumpMinor},
		{[]string{"fix: resolve bug"}, BumpPatch},
		{[]string{"feat!: breaking change"}, BumpMajor},
		{[]string{"BREAKING CHANGE: remove API"}, BumpMajor},
		{[]string{"feat: new feature", "fix: bug fix"}, BumpMinor},
		{[]string{"feat!: breaking", "feat: feature"}, BumpMajor},
		{[]string{"docs: update readme"}, BumpPatch},
		{[]string{"chore: update dependencies"}, BumpPatch},
		{[]string{}, BumpPatch}, // default
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			got := DetectBumpType(tt.commits)
			if got != tt.want {
				t.Errorf("DetectBumpType(%v) = %q, want %q", tt.commits, got, tt.want)
			}
		})
	}
}

func TestValidateSemanticVersionBump(t *testing.T) {
	tests := []struct {
		currentTag      string
		proposedVersion string
		bumpType        BumpType
		shouldErr       bool
	}{
		{"stg6_1.2.3-1", "1.3.0", BumpMinor, false},
		{"stg6_1.2.3-1", "2.0.0", BumpMajor, false},
		{"stg6_1.2.3-1", "1.2.4", BumpPatch, false},
		{"stg6_1.2.3-1", "1.4.0", BumpMinor, true}, // wrong version
		{"stg6_1.2.3-1", "1.2.3", BumpPatch, true}, // no increment
		{"invalid_tag", "1.2.4", BumpPatch, true},  // invalid tag
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			err := ValidateSemanticVersionBump(tt.currentTag, tt.proposedVersion, tt.bumpType)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("ValidateSemanticVersionBump expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ValidateSemanticVersionBump unexpected error: %v", err)
				}
			}
		})
	}
}
