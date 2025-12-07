package core

// Quadrant represents a quadrant of the radar.
type Quadrant struct {
	Name  string `yaml:"name"`
	Alias string `yaml:"alias"`
}

// Ring represents a ring of the radar.
type Ring struct {
	Name  string `yaml:"name"`
	Alias string `yaml:"alias"`
}

// Meta represents the metadata of the radar.
type Meta struct {
	Title       string     `yaml:"title"`
	Description string     `yaml:"description"`
	Quadrants   []Quadrant `yaml:"quadrants"`
	Rings       []Ring     `yaml:"rings"`
}

// DefaultRings is the default rings of the radar.
var DefaultRings = []Ring{
	{Name: "Adopt", Alias: "adopt"},
	{Name: "Trial", Alias: "trial"},
	{Name: "Assess", Alias: "assess"},
	{Name: "Hold", Alias: "hold"},
}

// DefaultQuadrants is the default quadrants of the radar.
var DefaultQuadrants = []Quadrant{
	{Name: "Languages", Alias: "languages"},
	{Name: "Frameworks", Alias: "frameworks"},
	{Name: "Platforms", Alias: "platforms"},
	{Name: "Techniques", Alias: "techniques"},
}

// DefaultMeta is the default metadata of the radar.
var DefaultMeta = Meta{
	Title:       "My Radar",
	Description: "Technology Radar",
	Quadrants:   DefaultQuadrants,
	Rings:       DefaultRings,
}
