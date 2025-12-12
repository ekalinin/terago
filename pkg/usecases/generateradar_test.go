package usecases

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ekalinin/terago/pkg/core"
)

func TestGetMovedValue(t *testing.T) {
	// Define default rings for testing
	rings := []core.Ring{
		{Name: "Adopt", Alias: "adopt"},
		{Name: "Trial", Alias: "trial"},
		{Name: "Assess", Alias: "assess"},
		{Name: "Hold", Alias: "hold"},
	}

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
			expected: MovedValueNew,
		},
		{
			name: "Unchanged technology",
			tech: core.Technology{
				IsNew:   false,
				IsMoved: false,
			},
			expected: MovedValueUnchanged,
		},
		{
			name: "Moved to inner ring (improved)",
			tech: core.Technology{
				IsNew:        false,
				IsMoved:      true,
				Ring:         "Adopt",
				PreviousRing: "Trial",
			},
			expected: MovedValueImproved,
		},
		{
			name: "Moved to outer ring (deprecated)",
			tech: core.Technology{
				IsNew:        false,
				IsMoved:      true,
				Ring:         "Hold",
				PreviousRing: "Adopt",
			},
			expected: MovedValueDeprecated,
		},
		{
			name: "Same ring but marked as moved",
			tech: core.Technology{
				IsNew:        false,
				IsMoved:      true,
				Ring:         "Adopt",
				PreviousRing: "Adopt",
			},
			expected: MovedValueDeprecated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getMovedValue(tt.tech, rings)
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

func TestGenerateRadarWithForce(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "terago_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test data
	meta := core.Meta{
		Title: "Test Radar",
		Quadrants: []core.Quadrant{
			{Name: "Languages", Alias: "languages"},
			{Name: "Frameworks", Alias: "frameworks"},
		},
		Rings: []core.Ring{
			{Name: "Adopt", Alias: "adopt"},
			{Name: "Trial", Alias: "trial"},
		},
	}

	files := []core.TechnologiesFile{
		{
			Date: "20231201",
			Technologies: []core.Technology{
				{Name: "Go", Ring: "Adopt", Quadrant: "Languages"},
			},
		},
		{
			Date: "20231202",
			Technologies: []core.Technology{
				{Name: "React", Ring: "Trial", Quadrant: "Frameworks"},
			},
		},
	}

	// Test 1: Generate without force (should create files)
	err = GenerateRadar(tempDir, "", files, meta, false, false, false)
	if err != nil {
		t.Fatalf("GenerateRadar failed: %v", err)
	}

	// Check that files were created
	file1 := filepath.Join(tempDir, "20231201.html")
	file2 := filepath.Join(tempDir, "20231202.html")

	if _, err := os.Stat(file1); os.IsNotExist(err) {
		t.Error("File 20231201.html should have been created")
	}
	if _, err := os.Stat(file2); os.IsNotExist(err) {
		t.Error("File 20231202.html should have been created")
	}

	// Get file modification times
	info1, _ := os.Stat(file1)
	info2, _ := os.Stat(file2)
	modTime1 := info1.ModTime()
	modTime2 := info2.ModTime()

	// Test 2: Generate without force again (should not modify existing files)
	err = GenerateRadar(tempDir, "", files, meta, false, false, false)
	if err != nil {
		t.Fatalf("GenerateRadar failed: %v", err)
	}

	// Check that files were not modified
	info1After, _ := os.Stat(file1)
	info2After, _ := os.Stat(file2)
	if !info1After.ModTime().Equal(modTime1) {
		t.Error("File 20231201.html should not have been modified when force=false")
	}
	if !info2After.ModTime().Equal(modTime2) {
		t.Error("File 20231202.html should not have been modified when force=false")
	}

	// Test 3: Generate with force (should modify existing files)
	err = GenerateRadar(tempDir, "", files, meta, true, false, false)
	if err != nil {
		t.Fatalf("GenerateRadar failed: %v", err)
	}

	// Check that files were modified
	info1AfterForce, _ := os.Stat(file1)
	info2AfterForce, _ := os.Stat(file2)
	if !info1AfterForce.ModTime().After(modTime1) {
		t.Error("File 20231201.html should have been modified when force=true")
	}
	if !info2AfterForce.ModTime().After(modTime2) {
		t.Error("File 20231202.html should have been modified when force=true")
	}
}
