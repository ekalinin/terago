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

func TestMetaPopulateSets(t *testing.T) {
	// Create a meta with some quadrants and rings
	quadrants := []Quadrant{
		{Name: "Languages", Alias: "languages"},
		{Name: "Frameworks", Alias: "frameworks"},
	}

	rings := []Ring{
		{Name: "Adopt", Alias: "adopt"},
		{Name: "Trial", Alias: "trial"},
	}

	meta := Meta{
		Title:       "Test Radar",
		Description: "Test Description",
		Quadrants:   quadrants,
		Rings:       rings,
		ringSet:     make(Set[string]),
		quadrantSet: make(Set[string]),
	}

	// Before populating, sets should be empty
	if len(meta.ringSet) != 0 {
		t.Errorf("Expected ringSet to be empty before PopulateSets, got %d elements", len(meta.ringSet))
	}

	if len(meta.quadrantSet) != 0 {
		t.Errorf("Expected quadrantSet to be empty before PopulateSets, got %d elements", len(meta.quadrantSet))
	}

	// Populate sets
	meta.PopulateSets()

	// After populating, sets should contain all names and aliases
	expectedRingSet := Set[string]{
		"Adopt": struct{}{},
		"adopt": struct{}{},
		"Trial": struct{}{},
		"trial": struct{}{},
	}

	expectedQuadrantSet := Set[string]{
		"Languages":  struct{}{},
		"languages":  struct{}{},
		"Frameworks": struct{}{},
		"frameworks": struct{}{},
	}

	// Check ringSet
	if len(meta.ringSet) != len(expectedRingSet) {
		t.Errorf("Expected ringSet to have %d elements, got %d", len(expectedRingSet), len(meta.ringSet))
	}

	for key := range expectedRingSet {
		if _, exists := meta.ringSet[key]; !exists {
			t.Errorf("Expected ringSet to contain %s", key)
		}
	}

	// Check quadrantSet
	if len(meta.quadrantSet) != len(expectedQuadrantSet) {
		t.Errorf("Expected quadrantSet to have %d elements, got %d", len(expectedQuadrantSet), len(meta.quadrantSet))
	}

	for key := range expectedQuadrantSet {
		if _, exists := meta.quadrantSet[key]; !exists {
			t.Errorf("Expected quadrantSet to contain %s", key)
		}
	}
}

func TestNewMeta(t *testing.T) {
	// Test with default values
	meta1 := NewMeta("", "", nil, nil)

	// Check that default values are set
	if meta1.Title != "My Radar" {
		t.Errorf("Expected default title 'My Radar', got %s", meta1.Title)
	}

	if meta1.Description != "Technology Radar" {
		t.Errorf("Expected default description 'Technology Radar', got %s", meta1.Description)
	}

	// Check that default quadrants are set
	if len(meta1.Quadrants) != len(DefaultQuadrants) {
		t.Errorf("Expected %d default quadrants, got %d", len(DefaultQuadrants), len(meta1.Quadrants))
	}

	// Check that default rings are set
	if len(meta1.Rings) != len(DefaultRings) {
		t.Errorf("Expected %d default rings, got %d", len(DefaultRings), len(meta1.Rings))
	}

	// Check that sets are populated
	if len(meta1.ringSet) == 0 {
		t.Error("Expected ringSet to be populated")
	}

	if len(meta1.quadrantSet) == 0 {
		t.Error("Expected quadrantSet to be populated")
	}

	// Test with custom values
	customTitle := "Custom Radar"
	customDescription := "Custom Description"
	customQuadrants := []Quadrant{
		{Name: "Tools", Alias: "tools"},
		{Name: "Services", Alias: "services"},
	}
	customRings := []Ring{
		{Name: "Use", Alias: "use"},
		{Name: "Evaluate", Alias: "evaluate"},
	}

	meta2 := NewMeta(customTitle, customDescription, customQuadrants, customRings)

	// Check that custom values are set
	if meta2.Title != customTitle {
		t.Errorf("Expected title '%s', got %s", customTitle, meta2.Title)
	}

	if meta2.Description != customDescription {
		t.Errorf("Expected description '%s', got %s", customDescription, meta2.Description)
	}

	// Check that custom quadrants are set
	if len(meta2.Quadrants) != len(customQuadrants) {
		t.Errorf("Expected %d custom quadrants, got %d", len(customQuadrants), len(meta2.Quadrants))
	}

	// Check that custom rings are set
	if len(meta2.Rings) != len(customRings) {
		t.Errorf("Expected %d custom rings, got %d", len(customRings), len(meta2.Rings))
	}

	// Check that sets are populated with custom values
	if len(meta2.ringSet) == 0 {
		t.Error("Expected ringSet to be populated with custom values")
	}

	if len(meta2.quadrantSet) == 0 {
		t.Error("Expected quadrantSet to be populated with custom values")
	}

	// Verify specific custom values in sets
	if !meta2.IsValidRing("Use") || !meta2.IsValidRing("use") {
		t.Error("Expected ringSet to contain custom ring 'Use' and 'use'")
	}

	if !meta2.IsValidQuadrant("Tools") || !meta2.IsValidQuadrant("tools") {
		t.Error("Expected quadrantSet to contain custom quadrant 'Tools' and 'tools'")
	}
}

func TestNewMetaFromFile(t *testing.T) {
	// Test with default FileNamePattern
	metaFile1 := MetaFile{
		Title:       "Test Radar",
		Description: "Test Description",
		Quadrants: []Quadrant{
			{Name: "Languages", Alias: "languages"},
		},
		Rings: []Ring{
			{Name: "Adopt", Alias: "adopt"},
		},
	}

	meta1 := NewMetaFromFile(metaFile1)

	// Check that default FileNamePattern is set
	expectedDefaultPattern := `^\d{8}\.yaml$`
	if meta1.FileNamePattern != expectedDefaultPattern {
		t.Errorf("Expected default FileNamePattern '%s', got '%s'", expectedDefaultPattern, meta1.FileNamePattern)
	}

	// Test with custom FileNamePattern
	customPattern := `^radar-\d{4}-\d{2}-\d{2}\.yaml$`
	metaFile2 := MetaFile{
		Title:           "Test Radar",
		Description:     "Test Description",
		FileNamePattern: customPattern,
		Quadrants: []Quadrant{
			{Name: "Languages", Alias: "languages"},
		},
		Rings: []Ring{
			{Name: "Adopt", Alias: "adopt"},
		},
	}

	meta2 := NewMetaFromFile(metaFile2)

	// Check that custom FileNamePattern is set
	if meta2.FileNamePattern != customPattern {
		t.Errorf("Expected custom FileNamePattern '%s', got '%s'", customPattern, meta2.FileNamePattern)
	}

	// Verify other fields are properly set
	if meta2.Title != metaFile2.Title {
		t.Errorf("Expected title '%s', got '%s'", metaFile2.Title, meta2.Title)
	}

	if meta2.Description != metaFile2.Description {
		t.Errorf("Expected description '%s', got '%s'", metaFile2.Description, meta2.Description)
	}
}
