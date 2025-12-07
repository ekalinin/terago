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
	Version     string
	GeneratedAt string
	Entries     []RadarEntry
	EntriesJSON template.JS // Adding field for JSON representation
}

// ToJSON converts the entries to JSON format for use in templates
func (rd *RadarData) ToJSON() (template.JS, error) {
	jsonBytes, err := json.Marshal(rd.Entries)
	if err != nil {
		return template.JS(""), err
	}
	return template.JS(jsonBytes), nil
}
