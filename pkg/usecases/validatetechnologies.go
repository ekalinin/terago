package usecases

import (
	"fmt"
	"os"

	"github.com/ekalinin/terago/pkg/core"
	"gopkg.in/yaml.v3"
)

// ValidateTechnologiesFile validates a single technologies YAML file.
// It checks YAML syntax, required fields, and validates rings and quadrants.
func ValidateTechnologiesFile(filePath string, meta core.Meta) error {
	// Read file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Parse YAML content
	var technologiesFile core.TechnologiesFile
	if err := yaml.Unmarshal(data, &technologiesFile); err != nil {
		return fmt.Errorf("error parsing YAML: %v", err)
	}

	// Check if technologies list is empty
	if len(technologiesFile.Technologies) == 0 {
		return fmt.Errorf("no technologies found in file")
	}

	// Validate rings and quadrants
	if err := technologiesFile.ValidateRingsAndQuadrants(meta); err != nil {
		return err
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
