package core

import "fmt"

// Technology represents a single technology entry in the radar
type Technology struct {
	Name        string `yaml:"name"`
	Ring        string `yaml:"ring"`
	Quadrant    string `yaml:"quadrant"`
	Description string `yaml:"description"`
	Info        string `yaml:"info,omitempty"`
	// Used for tracking changes between periods
	IsNew        bool   `yaml:"-"`
	IsMoved      bool   `yaml:"-"`
	PreviousRing string `yaml:"-"`
}

// TechnologiesFile represents the structure of the YAML file
type TechnologiesFile struct {
	Date         string
	Technologies []Technology `yaml:"technologies"`
}

// ValidateRingsAndQuadrants validates that all technologies in the file
// have valid rings and quadrants according to the provided meta configuration
func (tf *TechnologiesFile) ValidateRingsAndQuadrants(meta Meta) error {
	// Validate each technology
	for _, tech := range tf.Technologies {
		if !meta.IsValidRing(tech.Ring) {
			return fmt.Errorf("invalid ring '%s' in technology '%s'", tech.Ring, tech.Name)
		}
		if !meta.IsValidQuadrant(tech.Quadrant) {
			return fmt.Errorf("invalid quadrant '%s' in technology '%s'", tech.Quadrant, tech.Name)
		}
	}

	return nil
}
