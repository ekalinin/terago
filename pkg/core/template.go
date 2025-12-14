package core

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
)

// RadarEntry represents an entry in the radar visualization
type RadarEntry struct {
	Quadrant    int    `json:"quadrant"`
	Ring        int    `json:"ring"`
	Moved       int    `json:"moved"`
	Label       string `json:"label"`
	Link        string `json:"link"`
	Active      bool   `json:"active"`
	Description string `json:"description"`
}

// RadarData represents the data needed for the HTML template
type RadarData struct {
	Title       string
	Date        string
	Version     string
	GeneratedAt string
	Entries     []RadarEntry
	Quadrants   []Quadrant // Adding field for quadrants
	Rings       []Ring     // Adding field for rings
	// for JSON representation in the template
	EntriesJSON   template.JS
	QuadrantsJSON template.JS
	RingsJSON     template.JS
	DescriptionJS template.JS   // JavaScript for description modal
	ChangesTable  template.HTML // HTML table with changes
	// Embedded JavaScript libraries (empty if using CDN)
	D3JS    template.JS
	RadarJS template.JS
}

// UpdateJSON updates all JSON fields in the RadarData struct
func (rd *RadarData) UpdateJSON() error {
	// Update EntriesJSON
	entriesJSON, err := json.Marshal(rd.Entries)
	if err != nil {
		return err
	}
	rd.EntriesJSON = template.JS(entriesJSON)

	// Update QuadrantsJSON
	type QuadrantData struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}

	quadrantData := make([]QuadrantData, len(rd.Quadrants))
	for i, q := range rd.Quadrants {
		quadrantData[i] = QuadrantData{
			Name: q.Name,
			ID:   "q" + fmt.Sprintf("%d", i+1),
		}
	}

	quadrantsJSON, err := json.Marshal(quadrantData)
	if err != nil {
		return err
	}
	rd.QuadrantsJSON = template.JS(quadrantsJSON)

	// Update RingsJSON
	type RingData struct {
		Name  string `json:"name"`
		Color string `json:"color"`
		ID    string `json:"id"`
	}

	// Define colors for rings (same as in the original template)
	colors := []string{"#93c47d", "#93d2c2", "#fbdb84", "#efafa9"}

	ringData := make([]RingData, len(rd.Rings))
	for i, r := range rd.Rings {
		color := "#ddd" // default color
		if i < len(colors) {
			color = colors[i]
		}

		ringData[i] = RingData{
			Name:  strings.ToUpper(r.Name),
			Color: color,
			ID:    r.Alias,
		}
	}

	ringsJSON, err := json.Marshal(ringData)
	if err != nil {
		return err
	}
	rd.RingsJSON = template.JS(ringsJSON)

	return nil
}

// SetDescriptionJS sets the JavaScript code for description modal
func (rd *RadarData) SetDescriptionJS(js string) {
	rd.DescriptionJS = template.JS(js)
}

// SetChangesTable sets the HTML table with changes
func (rd *RadarData) SetChangesTable(html string) {
	rd.ChangesTable = template.HTML(html)
}

// SetEmbeddedLibs sets the embedded JavaScript libraries
func (rd *RadarData) SetEmbeddedLibs(d3JS, radarJS string) {
	rd.D3JS = template.JS(d3JS)
	rd.RadarJS = template.JS(radarJS)
}
