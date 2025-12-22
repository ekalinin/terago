package usecases

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/ekalinin/terago/pkg/core"
	"gopkg.in/yaml.v3"
)

// GetRadarFiles returns a sorted list of YAML files that match the pattern from meta.
func GetRadarFiles(inputDir string, meta core.Meta) ([]string, error) {
	// Get all YAML files in the directory
	files, err := filepath.Glob(filepath.Join(inputDir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	// Get file name pattern from meta or use default YYYYMMDD.yaml pattern
	datePattern := regexp.MustCompile(meta.FileNamePattern)
	var validFiles []string
	for _, file := range files {
		baseName := filepath.Base(file)
		if datePattern.MatchString(baseName) {
			validFiles = append(validFiles, file)
		}
	}

	// Sort files by date (filename)
	sort.Slice(validFiles, func(i, j int) bool {
		fileNameI := filepath.Base(validFiles[i])
		fileNameJ := filepath.Base(validFiles[j])
		return strings.Compare(fileNameI, fileNameJ) < 0
	})

	return validFiles, nil
}

// ReadTechnologiesFiles reads all Technologies files in the specified directory.
func ReadTechnologiesFiles(inputDir string, meta core.Meta) ([]core.TechnologiesFile, error) {
	var technologiesFiles []core.TechnologiesFile

	// Get all valid YAML files
	validFiles, err := GetRadarFiles(inputDir, meta)
	if err != nil {
		return technologiesFiles, err
	}

	// Process each file in order
	var previousTechnologies []core.Technology
	for _, file := range validFiles {
		fileName := filepath.Base(file)
		// fmt.Printf("Processing file: %s\n", fileName)

		// Extract date from filename (without .yaml extension)
		dateStr := strings.TrimSuffix(fileName, ".yaml")

		// Parse YAML file
		technologiesFile, err := readTechnologiesFile(file, meta)
		if err != nil {
			return technologiesFiles, fmt.Errorf("error processing file %s: %v", file, err)
		}

		// Set the date in the file
		technologiesFile.Date = dateStr

		// Compare with previous period to identify changes
		if previousTechnologies != nil {
			markChanges(technologiesFile.Technologies, previousTechnologies)
		} else {
			// If this is the first file, mark all as new
			for i := range technologiesFile.Technologies {
				technologiesFile.Technologies[i].IsNew = true
			}
		}

		// Save current technologies for next comparison
		previousTechnologies = technologiesFile.Technologies

		// Add to result
		technologiesFiles = append(technologiesFiles, technologiesFile)
	}

	return technologiesFiles, nil
}

// readTechnologiesFile reads and validates a single technologies YAML file.
func readTechnologiesFile(filePath string, meta core.Meta) (core.TechnologiesFile, error) {
	var technologiesFile core.TechnologiesFile

	// Read file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return technologiesFile, fmt.Errorf("error reading file: %v", err)
	}

	// Parse YAML content
	if err := yaml.Unmarshal(data, &technologiesFile); err != nil {
		return technologiesFile, fmt.Errorf("error parsing YAML: %v", err)
	}

	// Validate rings and quadrants
	if err := technologiesFile.ValidateRingsAndQuadrants(meta); err != nil {
		return technologiesFile, fmt.Errorf("validation error in file %s: %v", filePath, err)
	}

	return technologiesFile, nil
}

// markChanges compares current technologies with previous ones and marks changes
func markChanges(current, previous []core.Technology) {
	// Create a map of previous technologies for quick lookup
	previousMap := make(map[string]core.Technology)
	for _, tech := range previous {
		previousMap[tech.Name] = tech
	}

	// Check each current technology
	for i := range current {
		prevTech, exists := previousMap[current[i].Name]
		if exists {
			// Technology existed before, check if ring changed
			if prevTech.Ring != current[i].Ring {
				current[i].IsMoved = true
				current[i].PreviousRing = prevTech.Ring
			}
			current[i].IsNew = false
		} else {
			// New technology
			current[i].IsNew = true
		}
	}
}
