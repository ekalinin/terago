package usecases

import (
	"fmt"

	"github.com/ekalinin/terago/pkg/core"
)

// ValidateTechnologiesFile validates a single technologies YAML file.
// It checks YAML syntax, required fields, and validates rings and quadrants.
func ValidateTechnologiesFile(filePath string, meta core.Meta) error {
	// Read and parse file
	technologiesFile, err := readTechnologiesFile(filePath, meta)
	if err != nil {
		return err
	}

	// Check if technologies list is empty
	if len(technologiesFile.Technologies) == 0 {
		return fmt.Errorf("no technologies found in file")
	}

	// Additional validation: check for required fields
	for i, tech := range technologiesFile.Technologies {
		if tech.Name == "" {
			return fmt.Errorf("technology #%d is missing 'name' field", i+1)
		}
		if tech.Ring == "" {
			return fmt.Errorf("technology '%s' is missing 'ring' field", tech.Name)
		}
		if tech.Quadrant == "" {
			return fmt.Errorf("technology '%s' is missing 'quadrant' field", tech.Name)
		}
		if tech.Description == "" {
			return fmt.Errorf("technology '%s' is missing 'description' field", tech.Name)
		}
	}

	return nil
}
