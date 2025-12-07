package core

import (
	"testing"
)

func TestMetaIsValidRing(t *testing.T) {
	meta := NewMeta("Test Radar", "Test Description", nil, []Ring{
		{Name: "Adopt", Alias: "adopt"},
		{Name: "Trial", Alias: "trial"},
	})

	tests := []struct {
		ring     string
		expected bool
	}{
		{"Adopt", true},
		{"adopt", true},
		{"Trial", true},
		{"trial", true},
		{"Hold", false},
		{"hold", false},
		{"", false},
	}

	for _, test := range tests {
		result := meta.IsValidRing(test.ring)
		if result != test.expected {
			t.Errorf("IsValidRing(%s) = %v; expected %v", test.ring, result, test.expected)
		}
	}
}

func TestMetaIsValidQuadrant(t *testing.T) {
	meta := NewMeta("Test Radar", "Test Description", []Quadrant{
		{Name: "Languages", Alias: "languages"},
		{Name: "Frameworks", Alias: "frameworks"},
	}, nil)

	tests := []struct {
		quadrant string
		expected bool
	}{
		{"Languages", true},
		{"languages", true},
		{"Frameworks", true},
		{"frameworks", true},
		{"Platforms", false},
		{"platforms", false},
		{"", false},
	}

	for _, test := range tests {
		result := meta.IsValidQuadrant(test.quadrant)
		if result != test.expected {
			t.Errorf("IsValidQuadrant(%s) = %v; expected %v", test.quadrant, result, test.expected)
		}
	}
}
