package usecases

import (
	"testing"

	"github.com/ekalinin/terago/pkg/core"
)

func TestGetMovedValue(t *testing.T) {
	tests := []struct {
		name     string
		tech     core.Technology
		expected int
	}{
		{
			name: "New technology",
			tech: core.Technology{
				IsNew:   true,
				IsMoved: false,
			},
			expected: 0,
		},
		{
			name: "Moved technology",
			tech: core.Technology{
				IsNew:   false,
				IsMoved: true,
			},
			expected: 1,
		},
		{
			name: "Unchanged technology",
			tech: core.Technology{
				IsNew:   false,
				IsMoved: false,
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getMovedValue(tt.tech)
			if result != tt.expected {
				t.Errorf("getMovedValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetRingIndex(t *testing.T) {
	rings := []core.Ring{
		{Name: "Adopt", Alias: "adopt"},
		{Name: "Trial", Alias: "trial"},
		{Name: "Assess", Alias: "assess"},
		{Name: "Hold", Alias: "hold"},
	}

	tests := []struct {
		name     string
		ring     string
		expected int
	}{
		{
			name:     "Exact match with name",
			ring:     "Adopt",
			expected: 0,
		},
		{
			name:     "Exact match with alias",
			ring:     "trial",
			expected: 1,
		},
		{
			name:     "Case insensitive match",
			ring:     "ASSESS",
			expected: 2,
		},
		{
			name:     "No match - default to first",
			ring:     "Unknown",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getRingIndex(tt.ring, rings)
			if result != tt.expected {
				t.Errorf("getRingIndex() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetQuadrantIndex(t *testing.T) {
	quadrants := []core.Quadrant{
		{Name: "Languages", Alias: "languages"},
		{Name: "Frameworks", Alias: "frameworks"},
		{Name: "Platforms", Alias: "platforms"},
		{Name: "Techniques", Alias: "techniques"},
	}

	tests := []struct {
		name     string
		quadrant string
		expected int
	}{
		{
			name:     "Exact match with name",
			quadrant: "Languages",
			expected: 0,
		},
		{
			name:     "Exact match with alias",
			quadrant: "frameworks",
			expected: 1,
		},
		{
			name:     "Case insensitive match",
			quadrant: "PLATFORMS",
			expected: 2,
		},
		{
			name:     "No match - default to first",
			quadrant: "Unknown",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getQuadrantIndex(tt.quadrant, quadrants)
			if result != tt.expected {
				t.Errorf("getQuadrantIndex() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid YYYYMMDD format",
			input:    "20231201",
			expected: "2023-12-01",
		},
		{
			name:     "Another valid date",
			input:    "20240115",
			expected: "2024-01-15",
		},
		{
			name:     "Invalid length - too short",
			input:    "202312",
			expected: "202312",
		},
		{
			name:     "Invalid length - too long",
			input:    "202312011",
			expected: "202312011",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Non-numeric string",
			input:    "abcdabcd",
			expected: "abcd-ab-cd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDate(tt.input)
			if result != tt.expected {
				t.Errorf("formatDate(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
