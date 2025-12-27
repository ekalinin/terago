package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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
	generator := GenerateRadar{
		OutputDir:             tempDir,
		TemplatePath:          "",
		Files:                 files,
		Meta:                  meta,
		Force:                 false,
		Verbose:               false,
		IncludeLinks:          false,
		AddChanges:            false,
		SkipFirstRadarChanges: false,
		EmbedLibs:             false,
	}
	err = generator.Do()
	if err != nil {
		t.Fatalf("GenerateRadar failed: %v", err)
	}

	// Check that files were created
	file1 := filepath.Join(tempDir, "20231201.html")
	file2 := filepath.Join(tempDir, "20231202.html")

	if _, err := os.Stat(file1); err != nil {
		t.Fatalf("File 20231201.html should have been created: %v", err)
	}
	if _, err := os.Stat(file2); err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("File 20231202.html should have been created: %v", err)
		} else {
			t.Fatalf("Unexpected error when checking 20231202.html: %v", err)
		}
	}

	// Get file modification times
	info1, err := os.Stat(file1)
	if err != nil {
		t.Fatalf("Failed to stat %s: %v", file1, err)
	}
	info2, err := os.Stat(file2)
	if err != nil {
		t.Fatalf("Failed to stat %s: %v", file2, err)
	}
	modTime1 := info1.ModTime()
	modTime2 := info2.ModTime()

	// Test 2: Generate without force again (should not modify existing files)
	generator = GenerateRadar{
		OutputDir:             tempDir,
		TemplatePath:          "",
		Files:                 files,
		Meta:                  meta,
		Force:                 false,
		Verbose:               false,
		IncludeLinks:          false,
		AddChanges:            false,
		SkipFirstRadarChanges: false,
		EmbedLibs:             false,
	}
	err = generator.Do()
	if err != nil {
		t.Fatalf("GenerateRadar failed: %v", err)
	}

	// Check that files were not modified
	info1After, err := os.Stat(file1)
	if err != nil {
		t.Fatalf("Failed to stat %s: %v", file1, err)
	}
	info2After, err := os.Stat(file2)
	if err != nil {
		t.Fatalf("Failed to stat %s: %v", file2, err)
	}
	if !info1After.ModTime().Equal(modTime1) {
		t.Error("File 20231201.html should not have been modified when force=false")
	}
	if !info2After.ModTime().Equal(modTime2) {
		t.Error("File 20231202.html should not have been modified when force=false")
	}

	// Test 3: Generate with force (should modify existing files)
	// Add a small delay to ensure modification time will be different (fails in Github Actions)
	time.Sleep(1100 * time.Millisecond)
	generator = GenerateRadar{
		OutputDir:             tempDir,
		TemplatePath:          "",
		Files:                 files,
		Meta:                  meta,
		Force:                 true,
		Verbose:               false,
		IncludeLinks:          false,
		AddChanges:            false,
		SkipFirstRadarChanges: false,
		EmbedLibs:             false,
	}
	err = generator.Do()
	if err != nil {
		t.Fatalf("GenerateRadar failed: %v", err)
	}

	// Check that files were modified
	info1AfterForce, err := os.Stat(file1)
	if err != nil {
		t.Fatalf("Failed to stat %s: %v", file1, err)
	}
	info2AfterForce, err := os.Stat(file2)
	if err != nil {
		t.Fatalf("Failed to stat %s: %v", file2, err)
	}
	if !info1AfterForce.ModTime().After(modTime1) {
		t.Error("File 20231201.html should have been modified when force=true")
	}
	if !info2AfterForce.ModTime().After(modTime2) {
		t.Error("File 20231202.html should have been modified when force=true")
	}
}

func TestGenerateRadarWithChanges(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "terago_test_changes")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test data with new and moved technologies
	meta := core.Meta{
		Title: "Test Radar with Changes",
		Quadrants: []core.Quadrant{
			{Name: "Languages", Alias: "languages"},
			{Name: "Frameworks", Alias: "frameworks"},
		},
		Rings: []core.Ring{
			{Name: "Adopt", Alias: "adopt"},
			{Name: "Trial", Alias: "trial"},
			{Name: "Assess", Alias: "assess"},
		},
	}

	files := []core.TechnologiesFile{
		{
			Date: "20231201",
			Technologies: []core.Technology{
				{Name: "Python", Ring: "Adopt", Quadrant: "Languages", Description: "Base technology"},
			},
		},
		{
			Date: "20231202",
			Technologies: []core.Technology{
				{Name: "Go", Ring: "Adopt", Quadrant: "Languages", Description: "Fast and efficient", IsNew: true},
				{Name: "React", Ring: "Trial", Quadrant: "Frameworks", Description: "UI library", IsMoved: true, PreviousRing: "Assess"},
			},
		},
	}

	// Test 1: Generate with addChanges=false (should not include changes table)
	generator := GenerateRadar{
		OutputDir:             tempDir,
		TemplatePath:          "",
		Files:                 files,
		Meta:                  meta,
		Force:                 false,
		Verbose:               false,
		IncludeLinks:          false,
		AddChanges:            false,
		SkipFirstRadarChanges: false,
		EmbedLibs:             false,
	}
	err = generator.Do()
	if err != nil {
		t.Fatalf("GenerateRadar failed: %v", err)
	}

	// Read the generated file for first radar
	outputFile1 := filepath.Join(tempDir, "20231201.html")
	content1, err := os.ReadFile(outputFile1)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Should not contain changes table
	if strings.Contains(string(content1), "Changes in this Radar") {
		t.Error("Output should not contain changes table when addChanges=false")
	}

	// Test 2: Generate with addChanges=true and SkipFirstRadarChanges=true
	// First radar (earliest date) should NOT have changes table
	// Second radar should have changes table
	generator = GenerateRadar{
		OutputDir:             tempDir,
		TemplatePath:          "",
		Files:                 files,
		Meta:                  meta,
		Force:                 true,
		Verbose:               false,
		IncludeLinks:          false,
		AddChanges:            true,
		SkipFirstRadarChanges: true,
		EmbedLibs:             false,
	}
	err = generator.Do()
	if err != nil {
		t.Fatalf("GenerateRadar failed: %v", err)
	}

	// Read the first radar file
	content1, err = os.ReadFile(outputFile1)
	if err != nil {
		t.Fatalf("Failed to read first output file: %v", err)
	}

	// First radar should NOT contain changes table (even with addChanges=true)
	content1Str := string(content1)
	if strings.Contains(content1Str, "Changes in this Radar") {
		t.Error("First radar should not contain changes table")
	}

	// Read the second radar file
	outputFile2 := filepath.Join(tempDir, "20231202.html")
	content2, err := os.ReadFile(outputFile2)
	if err != nil {
		t.Fatalf("Failed to read second output file: %v", err)
	}

	// Second radar should contain changes table
	content2Str := string(content2)
	if !strings.Contains(content2Str, "changes-section") {
		t.Error("Output should contain changes section when addChanges=true")
	}
	if !strings.Contains(content2Str, "Changes in this Radar") {
		t.Error("Output should contain changes title")
	}
	if !strings.Contains(content2Str, "Go") {
		t.Error("Output should contain Go technology")
	}
	if !strings.Contains(content2Str, "React") {
		t.Error("Output should contain React technology")
	}
	if !strings.Contains(content2Str, "NEW") {
		t.Error("Output should contain NEW status for Go")
	}
	if !strings.Contains(content2Str, "MOVED") {
		t.Error("Output should contain MOVED status for React")
	}

	// Test 3: Generate with addChanges=true and SkipFirstRadarChanges=false
	// Both radars should have changes table
	generator = GenerateRadar{
		OutputDir:             tempDir,
		TemplatePath:          "",
		Files:                 files,
		Meta:                  meta,
		Force:                 true,
		Verbose:               false,
		IncludeLinks:          false,
		AddChanges:            true,
		SkipFirstRadarChanges: false,
		EmbedLibs:             false,
	}
	err = generator.Do()
	if err != nil {
		t.Fatalf("GenerateRadar failed: %v", err)
	}

	// Read the first radar file
	content1, err = os.ReadFile(outputFile1)
	if err != nil {
		t.Fatalf("Failed to read first output file: %v", err)
	}

	// First radar should contain changes table when SkipFirstRadarChanges=false
	content1Str = string(content1)
	if !strings.Contains(content1Str, "changes-section") {
		t.Error("First radar should contain changes section when SkipFirstRadarChanges=false")
	}
}

func TestBuildChangesTable(t *testing.T) {
	meta := core.Meta{
		Rings: []core.Ring{
			{Name: "Adopt", Alias: "adopt"},
			{Name: "Trial", Alias: "trial"},
		},
	}

	tests := []struct {
		name         string
		technologies []core.Technology
		wantContains []string
		wantEmpty    bool
	}{
		{
			name: "with new technology",
			technologies: []core.Technology{
				{Name: "Go", Ring: "Adopt", Quadrant: "Languages", Description: "Fast language", IsNew: true},
			},
			wantContains: []string{"Go", "NEW", "Fast language", "Languages"},
			wantEmpty:    false,
		},
		{
			name: "with moved technology",
			technologies: []core.Technology{
				{Name: "React", Ring: "Trial", Quadrant: "Frameworks", Description: "UI library", IsMoved: true, PreviousRing: "Adopt"},
			},
			wantContains: []string{"React", "MOVED", "Adopt", "Trial", "UI library"},
			wantEmpty:    false,
		},
		{
			name: "with deleted technology",
			technologies: []core.Technology{
				{Name: "Angular", Ring: "Trial", Quadrant: "Frameworks", Description: "Deprecated framework", IsDeleted: true},
			},
			wantContains: []string{"Angular", "DELETED from Trial", "Deprecated framework", "status-deleted"},
			wantEmpty:    false,
		},
		{
			name: "with no changes",
			technologies: []core.Technology{
				{Name: "Docker", Ring: "Adopt", Quadrant: "Infrastructure", Description: "Container platform", IsNew: false, IsMoved: false},
			},
			wantContains: []string{},
			wantEmpty:    true,
		},
		{
			name: "with mixed technologies including deleted",
			technologies: []core.Technology{
				{Name: "Go", Ring: "Adopt", Quadrant: "Languages", Description: "Fast", IsNew: true},
				{Name: "Docker", Ring: "Adopt", Quadrant: "Infrastructure", Description: "Containers", IsNew: false, IsMoved: false},
				{Name: "React", Ring: "Trial", Quadrant: "Frameworks", Description: "UI", IsMoved: true, PreviousRing: "Assess"},
				{Name: "Angular", Ring: "Adopt", Quadrant: "Frameworks", Description: "Old framework", IsDeleted: true},
			},
			wantContains: []string{"Go", "React", "Angular", "NEW", "MOVED", "DELETED from Adopt"},
			wantEmpty:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildChangesTable(tt.technologies, meta)

			if tt.wantEmpty {
				if result != "" {
					t.Errorf("Expected empty result, got: %s", result)
				}
				return
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(result, want) {
					t.Errorf("Expected result to contain %q, but it doesn't. Result: %s", want, result)
				}
			}
		})
	}
}

func TestConvertTechnologiesToEntries(t *testing.T) {
	meta := core.Meta{
		Quadrants: []core.Quadrant{
			{Name: "Languages", Alias: "languages"},
			{Name: "Frameworks", Alias: "frameworks"},
			{Name: "Infrastructure", Alias: "infrastructure"},
			{Name: "Techniques", Alias: "techniques"},
		},
		Rings: []core.Ring{
			{Name: "Adopt", Alias: "adopt"},
			{Name: "Trial", Alias: "trial"},
			{Name: "Assess", Alias: "assess"},
			{Name: "Hold", Alias: "hold"},
		},
	}

	tests := []struct {
		name               string
		technologies       []core.Technology
		includeLinks       bool
		expectedCount      int
		checkFirstEntry    bool
		expectedQuadrant   int
		expectedRing       int
		expectedMoved      int
		expectedLabel      string
		expectedLink       string
		expectedHasDeleted bool
	}{
		{
			name: "Basic conversion without links",
			technologies: []core.Technology{
				{Name: "Go", Ring: "Adopt", Quadrant: "Languages", Description: "Go language", IsNew: false, IsMoved: false},
				{Name: "React", Ring: "Trial", Quadrant: "Frameworks", Description: "React library", IsNew: true, IsMoved: false},
			},
			includeLinks:     false,
			expectedCount:    2,
			checkFirstEntry:  true,
			expectedQuadrant: 0, // Languages
			expectedRing:     0, // Adopt
			expectedMoved:    0, // Unchanged
			expectedLabel:    "Go",
			expectedLink:     "",
		},
		{
			name: "Basic conversion with links",
			technologies: []core.Technology{
				{Name: "Go", Ring: "Adopt", Quadrant: "Languages", Description: "Go language", IsNew: false, IsMoved: false},
			},
			includeLinks:     true,
			expectedCount:    1,
			checkFirstEntry:  true,
			expectedQuadrant: 0,
			expectedRing:     0,
			expectedMoved:    0,
			expectedLabel:    "Go",
			expectedLink:     "/Languages/Go/",
		},
		{
			name: "New technology conversion",
			technologies: []core.Technology{
				{Name: "Kubernetes", Ring: "Trial", Quadrant: "Infrastructure", Description: "K8s", IsNew: true, IsMoved: false},
			},
			includeLinks:     false,
			expectedCount:    1,
			checkFirstEntry:  true,
			expectedQuadrant: 2, // Infrastructure
			expectedRing:     1, // Trial
			expectedMoved:    2, // New
			expectedLabel:    "Kubernetes",
			expectedLink:     "",
		},
		{
			name: "Moved technology conversion",
			technologies: []core.Technology{
				{Name: "React", Ring: "Adopt", Quadrant: "Frameworks", Description: "React", IsNew: false, IsMoved: true, PreviousRing: "Trial"},
			},
			includeLinks:     false,
			expectedCount:    1,
			checkFirstEntry:  true,
			expectedQuadrant: 1, // Frameworks
			expectedRing:     0, // Adopt
			expectedMoved:    1, // Improved (moved to inner ring)
			expectedLabel:    "React",
			expectedLink:     "",
		},
		{
			name: "Deleted technology should be skipped",
			technologies: []core.Technology{
				{Name: "Go", Ring: "Adopt", Quadrant: "Languages", Description: "Go", IsNew: false, IsMoved: false, IsDeleted: false},
				{Name: "Angular", Ring: "Hold", Quadrant: "Frameworks", Description: "Old framework", IsNew: false, IsMoved: false, IsDeleted: true},
				{Name: "React", Ring: "Trial", Quadrant: "Frameworks", Description: "React", IsNew: false, IsMoved: false, IsDeleted: false},
			},
			includeLinks:       false,
			expectedCount:      2, // Angular should be skipped
			checkFirstEntry:    true,
			expectedQuadrant:   0, // Languages (Go)
			expectedRing:       0, // Adopt
			expectedMoved:      0,
			expectedLabel:      "Go",
			expectedLink:       "",
			expectedHasDeleted: false,
		},
		{
			name: "Only deleted technologies",
			technologies: []core.Technology{
				{Name: "Angular", Ring: "Hold", Quadrant: "Frameworks", Description: "Old", IsDeleted: true},
				{Name: "Backbone", Ring: "Hold", Quadrant: "Frameworks", Description: "Ancient", IsDeleted: true},
			},
			includeLinks:    false,
			expectedCount:   0, // All should be skipped
			checkFirstEntry: false,
		},
		{
			name:            "Empty technologies list",
			technologies:    []core.Technology{},
			includeLinks:    false,
			expectedCount:   0,
			checkFirstEntry: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertTechnologiesToEntries(tt.technologies, meta, tt.includeLinks)

			// Check count
			if len(result) != tt.expectedCount {
				t.Errorf("Expected %d entries, got %d", tt.expectedCount, len(result))
			}

			// Check that no deleted technologies are in result
			if tt.expectedHasDeleted == false && len(result) > 0 {
				for _, entry := range result {
					for _, tech := range tt.technologies {
						if entry.Label == tech.Name && tech.IsDeleted {
							t.Errorf("Deleted technology %s should not be in result", tech.Name)
						}
					}
				}
			}

			// Check first entry details if specified
			if tt.checkFirstEntry && len(result) > 0 {
				entry := result[0]
				if entry.Quadrant != tt.expectedQuadrant {
					t.Errorf("Expected quadrant %d, got %d", tt.expectedQuadrant, entry.Quadrant)
				}
				if entry.Ring != tt.expectedRing {
					t.Errorf("Expected ring %d, got %d", tt.expectedRing, entry.Ring)
				}
				if entry.Moved != tt.expectedMoved {
					t.Errorf("Expected moved %d, got %d", tt.expectedMoved, entry.Moved)
				}
				if entry.Label != tt.expectedLabel {
					t.Errorf("Expected label %q, got %q", tt.expectedLabel, entry.Label)
				}
				if entry.Link != tt.expectedLink {
					t.Errorf("Expected link %q, got %q", tt.expectedLink, entry.Link)
				}
			}
		})
	}
}
