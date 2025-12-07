package core

import (
	"encoding/json"
	"html/template"
)

// RadarEntry represents an entry in the radar visualization
type RadarEntry struct {
	Quadrant int    `json:"quadrant"`
	Ring     int    `json:"ring"`
	Moved    int    `json:"moved"`
	Label    string `json:"label"`
	Link     string `json:"link"`
	Active   bool   `json:"active"`
}

// RadarData represents the data needed for the HTML template
type RadarData struct {
	Title       string
	Date        string
	Entries     []RadarEntry
	EntriesJSON template.JS // Adding field for JSON representation
}

// ToJSON converts the entries to JSON format for use in templates
func (rd *RadarData) ToJSON() template.JS {
	jsonBytes, _ := json.Marshal(rd.Entries)
	return template.JS(jsonBytes)
}
