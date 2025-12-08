package core

import (
	"encoding/json"
	"html/template"
	"testing"
)

func TestRadarDataUpdateJSON(t *testing.T) {
	// Test data
	entries := []RadarEntry{
		{Quadrant: 0, Ring: 0, Moved: 0, Label: "Go", Link: "/languages/go/", Active: false},
		{Quadrant: 1, Ring: 1, Moved: 1, Label: "React", Link: "/frameworks/react/", Active: false},
	}

	quadrants := []Quadrant{
		{Name: "Languages", Alias: "languages"},
		{Name: "Frameworks", Alias: "frameworks"},
		{Name: "Platforms", Alias: "platforms"},
		{Name: "Techniques", Alias: "techniques"},
	}

	rings := []Ring{
		{Name: "Adopt", Alias: "adopt"},
		{Name: "Trial", Alias: "trial"},
		{Name: "Assess", Alias: "assess"},
		{Name: "Hold", Alias: "hold"},
	}

	// Create RadarData
	data := RadarData{
		Title:       "Test Radar",
		Date:        "2023-12-01",
		Version:     "1.0.0",
		GeneratedAt: "2023-12-01 12:00:00",
		Entries:     entries,
		Quadrants:   quadrants,
		Rings:       rings,
	}

	// Call UpdateJSON
	err := data.UpdateJSON()
	if err != nil {
		t.Fatalf("UpdateJSON() error = %v", err)
	}

	// Test EntriesJSON
	if data.EntriesJSON == "" {
		t.Error("EntriesJSON is empty")
	}

	// Verify EntriesJSON content
	var parsedEntries []RadarEntry
	err = json.Unmarshal([]byte(data.EntriesJSON), &parsedEntries)
	if err != nil {
		t.Fatalf("Failed to parse EntriesJSON: %v", err)
	}

	if len(parsedEntries) != len(entries) {
		t.Errorf("EntriesJSON length = %v, want %v", len(parsedEntries), len(entries))
	}

	// Test QuadrantsJSON
	if data.QuadrantsJSON == "" {
		t.Error("QuadrantsJSON is empty")
	}

	// Verify QuadrantsJSON content
	var parsedQuadrants []map[string]interface{}
	err = json.Unmarshal([]byte(data.QuadrantsJSON), &parsedQuadrants)
	if err != nil {
		t.Fatalf("Failed to parse QuadrantsJSON: %v", err)
	}

	if len(parsedQuadrants) != len(quadrants) {
		t.Errorf("QuadrantsJSON length = %v, want %v", len(parsedQuadrants), len(quadrants))
	}

	// Check first quadrant
	if parsedQuadrants[0]["name"] != "Languages" {
		t.Errorf("First quadrant name = %v, want %v", parsedQuadrants[0]["name"], "Languages")
	}

	if parsedQuadrants[0]["id"] != "q1" {
		t.Errorf("First quadrant id = %v, want %v", parsedQuadrants[0]["id"], "q1")
	}

	// Test RingsJSON
	if data.RingsJSON == "" {
		t.Error("RingsJSON is empty")
	}

	// Verify RingsJSON content
	var parsedRings []map[string]interface{}
	err = json.Unmarshal([]byte(data.RingsJSON), &parsedRings)
	if err != nil {
		t.Fatalf("Failed to parse RingsJSON: %v", err)
	}

	if len(parsedRings) != len(rings) {
		t.Errorf("RingsJSON length = %v, want %v", len(parsedRings), len(rings))
	}

	// Check first ring
	if parsedRings[0]["name"] != "ADOPT" {
		t.Errorf("First ring name = %v, want %v", parsedRings[0]["name"], "ADOPT")
	}

	if parsedRings[0]["id"] != "adopt" {
		t.Errorf("First ring id = %v, want %v", parsedRings[0]["id"], "adopt")
	}

	if parsedRings[0]["color"] != "#93c47d" {
		t.Errorf("First ring color = %v, want %v", parsedRings[0]["color"], "#93c47d")
	}
}

func TestRadarDataUpdateJSONEmpty(t *testing.T) {
	// Test with empty data
	data := RadarData{
		Title:       "Empty Radar",
		Date:        "2023-12-01",
		Version:     "1.0.0",
		GeneratedAt: "2023-12-01 12:00:00",
		Entries:     []RadarEntry{},
		Quadrants:   []Quadrant{},
		Rings:       []Ring{},
	}

	// Call UpdateJSON
	err := data.UpdateJSON()
	if err != nil {
		t.Fatalf("UpdateJSON() error = %v", err)
	}

	// Test EntriesJSON
	if data.EntriesJSON != template.JS("[]") {
		t.Errorf("EntriesJSON = %v, want []", data.EntriesJSON)
	}

	// Test QuadrantsJSON
	if data.QuadrantsJSON != template.JS("[]") {
		t.Errorf("QuadrantsJSON = %v, want []", data.QuadrantsJSON)
	}

	// Test RingsJSON
	if data.RingsJSON != template.JS("[]") {
		t.Errorf("RingsJSON = %v, want []", data.RingsJSON)
	}
}
