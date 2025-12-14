package radar

import _ "embed"

//go:embed radar.html
var HTML string

//go:embed showDescription.js
var DescriptionJS string

//go:embed d3.min.js
var D3JS string

//go:embed radar.min.js
var RadarJS string
