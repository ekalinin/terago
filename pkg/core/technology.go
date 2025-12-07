package core

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
