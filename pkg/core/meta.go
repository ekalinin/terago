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

// MetaFile represents the metadata of the radar file.
type MetaFile struct {
	Title           string     `yaml:"title"`
	Description     string     `yaml:"description"`
	Quadrants       []Quadrant `yaml:"quadrants"`
	Rings           []Ring     `yaml:"rings"`
	FileNamePattern string     `yaml:"fileNamePattern"`
}

// Meta represents the metadata of the radar data used in main logic.
type Meta struct {
	Title           string      `yaml:"title"`
	Description     string      `yaml:"description"`
	Quadrants       []Quadrant  `yaml:"quadrants"`
	Rings           []Ring      `yaml:"rings"`
	FileNamePattern string      `yaml:"fileNamePattern"`
	ringSet         Set[string] `yaml:"-"`
	quadrantSet     Set[string] `yaml:"-"`
}

var defaultMeta = Meta{
	Title:           "My Radar",
	Description:     "Technology Radar",
	Quadrants:       DefaultQuadrants,
	Rings:           DefaultRings,
	FileNamePattern: `^\d{8}\.yaml$`, // default YYYYMMDD.yaml pattern
}

// NewMeta creates a new Meta with initialized ringSet and quadrantSet
func NewMeta(title, description string, quadrants []Quadrant, rings []Ring) Meta {
	m := defaultMeta

	if title != "" {
		m.Title = title
	}
	if description != "" {
		m.Description = description
	}
	if quadrants != nil {
		m.Quadrants = quadrants
	}
	if rings != nil {
		m.Rings = rings
	}

	m.PopulateSets()

	return m
}

// NewMetaFromFile creates a new Meta from a MetaFile
func NewMetaFromFile(metaFile MetaFile) Meta {
	m := NewMeta(
		metaFile.Title,
		metaFile.Description,
		metaFile.Quadrants,
		metaFile.Rings,
	)

	// Override FileNamePattern if provided
	if metaFile.FileNamePattern != "" {
		m.FileNamePattern = metaFile.FileNamePattern
	}

	return m
}

// PopulateSets fills the ringSet and quadrantSet with values from Rings and Quadrants
func (m *Meta) PopulateSets() {
	// Populate ringSet
	m.ringSet = make(Set[string])
	for _, r := range m.Rings {
		m.ringSet[r.Name] = struct{}{}
		m.ringSet[r.Alias] = struct{}{}
	}

	// Populate quadrantSet
	m.quadrantSet = make(Set[string])
	for _, q := range m.Quadrants {
		m.quadrantSet[q.Name] = struct{}{}
		m.quadrantSet[q.Alias] = struct{}{}
	}
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

// defaultMeta returns the default metadata of the radar.
func DefaultMeta() Meta {
	return NewMeta("", "", nil, nil)
}
