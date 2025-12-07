package usecases

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/ekalinin/terago/pkg/core"
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
func getMovedValue(tech core.Technology) int {
	if tech.IsNew {
		return 0 // New technology
	}
	if tech.IsMoved {
		// Determine direction based on ring movement
		// This is a simplified logic - you might want to enhance it
		return 1 // Moved outward
	}
	return 0 // No movement
}

// convertTechnologiesToEntries converts Technology structs to RadarEntry structs
func convertTechnologiesToEntries(technologies []core.Technology, meta core.Meta) []core.RadarEntry {
	var entries []core.RadarEntry

	for _, tech := range technologies {
		quadrantIndex := getQuadrantIndex(tech.Quadrant, meta.Quadrants)
		ringIndex := getRingIndex(tech.Ring, meta.Rings)
		moved := getMovedValue(tech)

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

func GenerateRadar(outputDir, templatePath string, files []core.TechnologiesFile, meta core.Meta) error {
	// Create output directory if it doesn't exist
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return err
		}
	}

	// read template file
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return err
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
		data := core.RadarData{
			Title:   meta.Title,
			Date:    file.Date,
			Entries: entries,
		}
		data.EntriesJSON = data.ToJSON()

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
