package usecases

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/ekalinin/terago/pkg/core"
	"github.com/ekalinin/terago/pkg/radar"
)

const (
	// MovedValueUnchanged indicates no movement
	MovedValueUnchanged = 0
	// MovedValueDeprecated indicates movement to outer ring (deprecated)
	MovedValueDeprecated = -1
	// MovedValueImproved indicates movement to inner ring (improved)
	MovedValueImproved = 1
	// MovedValueNew indicates new technology (not in previous radar)
	MovedValueNew = 2
)

// getQuadrantIndex returns the index of the quadrant based on its name
func getQuadrantIndex(quadrant string, quadrants []core.Quadrant) int {
	for i, q := range quadrants {
		if strings.EqualFold(q.Name, quadrant) || strings.EqualFold(q.Alias, quadrant) {
			return i
		}
	}
	return 0 // Default to first quadrant if not found
}

// getRingIndex returns the index of the ring based on its name
func getRingIndex(ring string, rings []core.Ring) int {
	for i, r := range rings {
		if strings.EqualFold(r.Name, ring) || strings.EqualFold(r.Alias, ring) {
			return i
		}
	}
	return 0 // Default to first ring if not found
}

// getMovedValue determines the moved value based on technology changes
func getMovedValue(tech core.Technology, rings []core.Ring) int {
	if tech.IsNew {
		return MovedValueNew // New technology
	}
	if tech.IsMoved {
		// Determine direction based on ring movement
		currentRingIndex := getRingIndex(tech.Ring, rings)
		previousRingIndex := getRingIndex(tech.PreviousRing, rings)

		if currentRingIndex < previousRingIndex {
			return MovedValueImproved // Moved to inner ring (improved)
		} else if currentRingIndex > previousRingIndex {
			return MovedValueDeprecated // Moved to outer ring (deprecated)
		}
		// If ring indices are equal but IsMoved is true, return deprecated as default
		return MovedValueDeprecated
	}
	return MovedValueUnchanged // No movement
}

// formatDate converts date from YYYYMMDD format to YYYY-MM-DD format
func formatDate(dateStr string) string {
	if len(dateStr) != 8 {
		return dateStr // Return as is if not in expected format
	}
	return dateStr[:4] + "-" + dateStr[4:6] + "-" + dateStr[6:8]
}

// convertTechnologiesToEntries converts Technology structs to RadarEntry structs
func convertTechnologiesToEntries(technologies []core.Technology, meta core.Meta) []core.RadarEntry {
	var entries []core.RadarEntry

	for _, tech := range technologies {
		quadrantIndex := getQuadrantIndex(tech.Quadrant, meta.Quadrants)
		ringIndex := getRingIndex(tech.Ring, meta.Rings)
		moved := getMovedValue(tech, meta.Rings)

		// Create link based on technology name and quadrant
		link := "/" + tech.Quadrant + "/" + tech.Name + "/"

		entry := core.RadarEntry{
			Quadrant: quadrantIndex,
			Ring:     ringIndex,
			Moved:    moved,
			Label:    tech.Name,
			Link:     link,
			Active:   false,
		}

		entries = append(entries, entry)
	}

	return entries
}

// GenerateRadar generates Radar (HTML file).
func GenerateRadar(outputDir, templatePath string, files []core.TechnologiesFile, meta core.Meta) error {
	// Create output directory if it doesn't exist
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return err
		}
	}

	// read template file
	var templateContent []byte
	var err error
	if templatePath == "" {
		// Use embedded template
		templateContent = []byte(radar.HTML)
	} else {
		// Read template file from disk
		templateContent, err = os.ReadFile(templatePath)
		if err != nil {
			return err
		}
	}

	// Parse template
	tmpl, err := template.New("radar").Parse(string(templateContent))
	if err != nil {
		return err
	}

	// Process each file and generate HTML
	for _, file := range files {
		// Convert technologies to radar entries
		entries := convertTechnologiesToEntries(file.Technologies, meta)

		// Prepare data for template
		formattedDate := formatDate(file.Date)
		data := core.RadarData{
			Title:   meta.Title,
			Date:    formattedDate,
			Entries: entries,
		}
		entriesJSON, err := data.ToJSON()
		if err != nil {
			return err
		}
		data.EntriesJSON = entriesJSON

		// Create output file
		outputFile := filepath.Join(outputDir, file.Date+".html")
		f, err := os.Create(outputFile)
		if err != nil {
			return err
		}
		defer f.Close()

		// Execute template
		if err := tmpl.Execute(f, data); err != nil {
			return err
		}
	}

	return nil
}
