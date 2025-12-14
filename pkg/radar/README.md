# Radar Package

This package contains embedded resources for technology radar visualization.

## Files

- `radar.html` - HTML template for radar visualization
- `showDescription.js` - JavaScript for showing technology descriptions in modal
- `d3.min.js` - D3.js library for data visualization (minified)
- `radar.min.js` - Zalando Tech Radar library for radar visualization (minified)

## JavaScript Libraries

The JavaScript libraries are embedded into the Go binary at compile time using `//go:embed` directives.
This allows the application to work without internet access and ensures version consistency.

**Library versions are managed in the Makefile** (single source of truth):
- D3.js version: controlled by `D3_VERSION` (full) and `D3_MAJOR_VERSION` (for URL)
- Radar version: controlled by `RADAR_VERSION` variable

Note: D3.js uses major version in CDN URLs (e.g., `d3.v7.min.js`), so both versions are needed.

### Updating Libraries

If you need to update the JavaScript libraries, use the Makefile commands:

```bash
# Download and minify libraries in one command
make update-libs

# Or step by step:
make download-libs  # Downloads D3.js and Zalando Tech Radar
make minify-libs    # Minifies radar.min.js (reduces size by ~50%)
```

The minification process:
- Downloads the original radar library (~20KB)
- Minifies it using built-in Go minifier (~10KB, 51% reduction)
- Reduces the final HTML size by ~10KB when using --embed-libs

**Requirements:**
- No external dependencies required!
- Minification is done using pure Go (`github.com/tdewolff/minify/v2`)

**To update library versions:**
1. Edit `Makefile` and change version variables:
   - For D3.js: update both `D3_VERSION` (e.g., 8.0.0) and `D3_MAJOR_VERSION` (e.g., 8)
   - For Radar: update `RADAR_VERSION` (e.g., 0.13)
2. Run `make update-libs`
3. Rebuild the project

Manual minification (if needed):

```bash
# Check current versions in Makefile
grep -E "D3_VERSION|RADAR_VERSION" Makefile

# Download libraries manually
curl -L -o pkg/radar/d3.min.js https://d3js.org/d3.v7.min.js
curl -L -o pkg/radar/radar.min.js.tmp.js https://zalando.github.io/tech-radar/release/radar-0.12.js

# Build minifier and minify
make build-minifier
./build/minifyjs -input pkg/radar/radar.min.js.tmp.js -output pkg/radar/radar.min.js
```

After updating, rebuild the project:

```bash
go build ./cmd/terago
```

## Usage

These files are automatically embedded and used by the `generateradar.go` use case:

- When `--embed-libs` flag is NOT used: HTML files reference CDN URLs
- When `--embed-libs` flag IS used: JavaScript libraries are embedded inline in the HTML
