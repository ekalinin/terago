package usecases

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

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
// If includeLinks is true, Link field will be populated based on technology name and quadrant
func convertTechnologiesToEntries(technologies []core.Technology, meta core.Meta, includeLinks bool) []core.RadarEntry {
	var entries []core.RadarEntry

	for _, tech := range technologies {
		quadrantIndex := getQuadrantIndex(tech.Quadrant, meta.Quadrants)
		ringIndex := getRingIndex(tech.Ring, meta.Rings)
		moved := getMovedValue(tech, meta.Rings)

		// Create link based on technology name and quadrant if includeLinks is true
		link := ""
		if includeLinks {
			link = "/" + tech.Quadrant + "/" + tech.Name + "/"
		}

		entry := core.RadarEntry{
			Quadrant:    quadrantIndex,
			Ring:        ringIndex,
			Moved:       moved,
			Label:       tech.Name,
			Link:        link,
			Active:      false,
			Description: tech.Description,
		}

		entries = append(entries, entry)
	}

	return entries
}

// buildChangesTable creates an HTML table with changed or new technologies
func buildChangesTable(technologies []core.Technology, meta core.Meta) string {
	var changedTechs []core.Technology

	// Filter technologies that are new or moved
	for _, tech := range technologies {
		if tech.IsNew || tech.IsMoved {
			changedTechs = append(changedTechs, tech)
		}
	}

	// If no changes, return empty string
	if len(changedTechs) == 0 {
		return ""
	}

	// Build HTML table
	html := `
	<div class="changes-section">
		<h3>Changes in this Radar</h3>
		<table class="changes-table">
			<thead>
				<tr>
					<th>Technology</th>
					<th>Quadrant</th>
					<th>Status</th>
					<th>Description</th>
				</tr>
			</thead>
			<tbody>`

	for _, tech := range changedTechs {
		status := ""
		if tech.IsNew {
			status = "NEW"
		} else if tech.IsMoved {
			status = "MOVED: " + tech.PreviousRing + " → " + tech.Ring
		}

		html += "\n\t\t\t\t<tr>"
		html += "\n\t\t\t\t\t<td><strong>" + tech.Name + "</strong></td>"
		html += "\n\t\t\t\t\t<td>" + tech.Quadrant + "</td>"
		html += "\n\t\t\t\t\t<td class=\"status-" + strings.ToLower(tech.Ring) + "\">" + status + "</td>"
		html += "\n\t\t\t\t\t<td>" + tech.Description + "</td>"
		html += "\n\t\t\t\t</tr>"
	}

	html += `
			</tbody>
		</table>
	</div>`

	return html
}

// GenerateRadar generates Radar (HTML file).
// If includeLinks is true, each radar entry will have a link based on its quadrant and name.
// If addChanges is true, a table with changed or new technologies will be included.
func GenerateRadar(outputDir, templatePath string, files []core.TechnologiesFile, meta core.Meta, force, verbose, includeLinks, addChanges bool) error {
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

	// Process each file and generate HTML only if it doesn't exist or force is true
	for _, file := range files {
		// Check if HTML file already exists
		outputFile := filepath.Join(outputDir, file.Date+".html")
		if !force {
			if _, err := os.Stat(outputFile); err == nil {
				// File exists and force is false, skip generation
				if verbose {
					log.Printf("Skipping %s.html (already exists, use --force to regenerate)", file.Date)
				}
				continue
			}
		}

		// Convert technologies to radar entries
		entries := convertTechnologiesToEntries(file.Technologies, meta, includeLinks)

		// Prepare data for template
		formattedDate := formatDate(file.Date)
		data := core.RadarData{
			Title:       meta.Title,
			Date:        formattedDate,
			Version:     core.Version,
			GeneratedAt: time.Now().Format("2006-01-02 15:04:05"),
			Entries:     entries,
			Quadrants:   meta.Quadrants,
			Rings:       meta.Rings,
		}
		if err := data.UpdateJSON(); err != nil {
			return err
		}
		// Set description JavaScript
		data.SetDescriptionJS(radar.DescriptionJS)

		// Build and set changes table if addChanges is true
		if addChanges {
			changesHTML := buildChangesTable(file.Technologies, meta)
			data.SetChangesTable(changesHTML)
		}

		// Create output file
		f, err := os.Create(outputFile)
		if err != nil {
			return err
		}
		defer f.Close()

		// Execute template
		if err := tmpl.Execute(f, data); err != nil {
			return err
		}

		if verbose {
			log.Printf("Generated %s.html", file.Date)
		}
	}

	return nil
}
