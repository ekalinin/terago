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

// Set is a generic set implementation
type Set[T comparable] map[T]struct{}

// Meta represents the metadata of the radar.
type Meta struct {
	Title       string      `yaml:"title"`
	Description string      `yaml:"description"`
	Quadrants   []Quadrant  `yaml:"quadrants"`
	Rings       []Ring      `yaml:"rings"`
	ringSet     Set[string] `yaml:"-"`
	quadrantSet Set[string] `yaml:"-"`
}

// NewMeta creates a new Meta with initialized ringSet and quadrantSet
func NewMeta(title, description string, quadrants []Quadrant, rings []Ring) Meta {
	m := Meta{
		Title:       title,
		Description: description,
		Quadrants:   quadrants,
		Rings:       rings,
		ringSet:     make(Set[string]),
		quadrantSet: make(Set[string]),
	}

	// Populate ringSet
	for _, r := range m.Rings {
		m.ringSet[r.Name] = struct{}{}
		m.ringSet[r.Alias] = struct{}{}
	}

	// Populate quadrantSet
	for _, q := range m.Quadrants {
		m.quadrantSet[q.Name] = struct{}{}
		m.quadrantSet[q.Alias] = struct{}{}
	}

	return m
}

// IsValidRing checks if a ring is valid according to the meta configuration
func (m *Meta) IsValidRing(ring string) bool {
	_, exists := m.ringSet[ring]
	return exists
}

// IsValidQuadrant checks if a quadrant is valid according to the meta configuration
func (m *Meta) IsValidQuadrant(quadrant string) bool {
	_, exists := m.quadrantSet[quadrant]
	return exists
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

// DefaultMeta returns the default metadata of the radar.
func DefaultMeta() Meta {
	return NewMeta(
		"My Radar",
		"Technology Radar",
		DefaultQuadrants,
		DefaultRings,
	)
}
