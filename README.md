# TeraGo - Technology Radar in Go

[![Build Status][build-badge]][build-url] [![Download][download-badge]][download-url] [![Test Status][test-badge]][test-url]

## Table of Contents

- [Key Features](#key-features)
- [Installation](#installation)
  - [From Source](#from-source)
  - [From GitHub Releases](#from-github-releases)
- [Usage](#usage)
  - [Available Commands](#available-commands)
  - [Generate Command](#generate-command)
  - [List Command](#list-command)
  - [Validate Command](#validate-command)
  - [Export Template Command](#export-template-command)
  - [Customizing the Radar Template](#customizing-the-radar-template)
  - [Input Data Format](#input-data-format)
    - [Metadata File (meta.yaml)](#metadata-file-metayaml)
    - [Technology Files (YYYYMMDD.yaml)](#technology-files-yyyymmddyaml)
- [Project Structure](#project-structure)
- [Visualization](#visualization)
- [License](#license)

TeraGo (**Te**chnology **Ra**dar in **Go**) is a tool for creating and visualizing
technology radars.
It allows you to track changes in your company's technology stack over time,
visualizing technologies by categories and adoption status.

## Key Features

- Generate technology radars in HTML format from YAML files
- Track technology changes between periods (new, moved technologies)
- Customizable technology categories (quadrants) and statuses (rings)
- Integration with templates for flexible appearance customization
- Multilingual support

## Installation

TeraGo can be installed from pre-built binaries or built from source.

### From Source

To build TeraGo from source, you need Go 1.20 or higher.

```bash
# Clone the repository
$ git clone https://github.com/ekalinin/terago.git
$ cd terago

# Build the project
$ make build
```

### From GitHub Releases

You can download the latest pre-built binary from [GitHub Releases][download-url]:

1. Go to the [releases page][download-url]
2. Download the appropriate version for your operating system
3. Extract the archive
4. Run the `terago` binary

## Usage

### Available Commands

TeraGo uses a command-based interface. Each command has its own set of options:

```bash
terago <command> [options]
```

**Available commands:**
- `generate` (or `g`) - Generate HTML radars from YAML files
- `export-template` (or `e`) - Export embedded template to file for customization
- `list` (or `l`) - List available radars and their render status
- `validate` (or `val`) - Validate YAML files structure and data
- `version` (or `v`) - Show version information
- `help` (or `h`) - Show help message

You can use full command names or short aliases (shown in parentheses) for convenience.

### Generate Command

Generate HTML technology radars from YAML files.

**Basic usage:**

```bash
./terago generate --input ./test/test_input --output ./output --meta ./test/test_input/test_meta.yaml
```

Or using the short alias:

```bash
./terago g --input ./test/test_input --output ./output
```

The `--meta` parameter is optional. If not specified, default values will be
used for quadrants and rings. These default values can be found in the
[source code](pkg/core/meta.go#L101-L115). The meta.yml file can partially
override these default values.

**Incremental Generation**: By default, TeraGo only generates HTML files for YAML files that don't have corresponding HTML files yet. This makes incremental updates efficient when adding new technology files.

**Force Regeneration**: To regenerate all HTML files (useful when updating templates or metadata), use the `--force` flag:

```bash
./terago generate --input ./test/test_input --output ./output --force
```

**Verbose Logging**: To see detailed information about file processing (which files are being generated or skipped), use the `--verbose` flag:

```bash
./terago generate --input ./test/test_input --output ./output --verbose
```

You can combine both flags:

```bash
./terago generate --input ./test/test_input --output ./output --force --verbose
```

#### Generate Command Options

- `--input` - path to directory with technology YAML files (required)
- `--output` - path to directory for saving HTML files (default: "output")
- `--template` - path to HTML template (if empty, uses default embedded template)
- `--meta` - path to metadata file (default: "meta.yaml")
- `--force` - force regeneration of all HTML files (ignore existing files)
- `--verbose` - enable verbose logging (show file processing details)
- `--include-links` - include links in radar entries (based on quadrant and technology name)
- `--add-changes` - add table with description of changed or new technologies
- `--skip-first-radar-changes` - skip changes table for the first (earliest) radar (default: true)
- `--embed-libs` - embed JavaScript libraries (D3.js and tech-radar) in HTML instead of loading from CDN

### List Command

List all available radar files in the input directory and check their render status.

**Basic usage:**

```bash
./terago list --input ./test/test_input --output ./output
```

Or using the short alias:

```bash
./terago l --input ./test/test_input
```

**Example output:**

```
Found 2 radar(s) in test/test_input:

  20231201 ✓ (rendered: 2025-12-17 18:36:01)
  20231202 ✓ (rendered: 2025-12-17 18:36:05)
  20231203 ✗ (not rendered)
```

The command shows:
- Total number of radar files found
- Each radar file with its date (YYYYMMDD format)
- Render status: ✓ (rendered with timestamp) or ✗ (not rendered)

#### List Command Options

- `--input` - path to directory with technology YAML files (required)
- `--output` - path to directory for HTML output (default: "output")

### Validate Command

Validate YAML files structure and data before generating radars. This command checks:
- YAML syntax correctness
- Presence of required fields (name, ring, quadrant, description)
- Validity of ring and quadrant values according to metadata
- Non-empty technologies list

**Basic usage:**

```bash
./terago validate --input ./test/test_input
```

Or using the short alias:

```bash
./terago val --input ./test/test_input
```

**Example output (normal mode):**

```
OK: 2 file(s) processed
```

**Example output (verbose mode):**

```bash
./terago validate --input ./test/test_input --verbose
```

```
OK: 20231201.yaml
OK: 20231202.yaml

OK: 2 file(s) processed
```

**Example with errors:**

```
ERROR: 20231203.yaml - invalid ring 'InvalidRing' in technology 'Python'
ERROR: 20231204.yaml - technology 'Go' is missing 'description' field

Validation completed with errors: 2 file(s) failed, 1 file(s) passed
```

#### Validate Command Options

- `--input` - path to directory with technology YAML files (required)
- `--meta` - path to metadata file (optional, default: searches for meta.yaml in input directory)
- `--verbose` - verbose output showing status for each file

The validate command is useful for:
- Checking data before generating radars
- CI/CD pipeline integration
- Quick validation of manually edited YAML files
- Finding errors in technology files

**Note about `--add-changes` and `--skip-first-radar-changes`**: When using the `--add-changes` flag, a table showing new and moved technologies is added to each radar. By default, this table is skipped for the first (earliest) radar because all technologies would be marked as "NEW" in the initial radar. You can control this behavior with the `--skip-first-radar-changes` flag:

```bash
# Default behavior: changes table shown for all radars except the first one
./terago generate --input ./test/test_input --output ./output --add-changes

# Explicitly skip changes table for the first radar
./terago generate --input ./test/test_input --output ./output --add-changes --skip-first-radar-changes=true

# Show changes table for ALL radars, including the first one
./terago generate --input ./test/test_input --output ./output --add-changes --skip-first-radar-changes=false
```

**Note about `--embed-libs`**: By default, the generated HTML files load D3.js and Zalando Tech Radar libraries from CDN (Content Delivery Network). This keeps the HTML files small (~11KB) but requires internet connection to view them. When you use the `--embed-libs` flag, the libraries (which are bundled with terago at compile time from `pkg/radar/`) are embedded directly into each HTML file (~304KB each). This makes the files self-contained and viewable offline, but significantly increases their size.

Example with embedded libraries:

```bash
./terago generate --input ./test/test_input --output ./output --embed-libs
```

### Export Template Command

Export the embedded HTML template to a file for customization.

**Basic usage:**

```bash
./terago export-template --output ./my-template.html
```

Or using the short alias:

```bash
./terago e --output ./my-template.html
```

This command exports the default embedded template that TeraGo uses for generating radar visualizations. Once exported, you can modify the template to customize the appearance of your radars.

#### Export Template Command Options

- `--output` - output file path for the template (required)

### Customizing the Radar Template

TeraGo uses an embedded HTML template for radar visualization. To customize
the appearance of your radar, first export the template using the `export-template` command
(see [Export Template Command](#export-template-command) section above):

```bash
./terago export-template --output ./my-template.html
```

This will create a file `my-template.html` with the default template content.
You can then modify this file according to your needs and use it with the `--template` parameter:

```bash
./terago generate --input ./test/test_input --output ./output --template ./my-template.html
```

The template uses Go's [text/template](https://pkg.go.dev/text/template) package
and has access to the following data:

- `.Title` - Radar title from metadata
- `.Date` - Current date
- `.Version` - Application version (see [version.go](pkg/core/version.go#L4))
- `.GeneratedAt` - Timestamp when the radar was generated (see [template.go](pkg/core/template.go#L18-L26))
- `.EntriesJSON` - Technologies data in JSON format
- `.Quadrants` - Array of quadrants from metadata
- `.Rings` - Array of rings from metadata
- `.QuadrantsJSON` - Quadrants data in JSON format
- `.RingsJSON` - Rings data in JSON format

The `.EntriesJSON` contains an array of technology entries with the following structure:

```json
[
  {
    "quadrant": 0,
    "ring": 0,
    "moved": 0,
    "label": "Technology Name",
    "link": "/quadrant/technology/",
    "active": false
  }
]
```

Where:
- `quadrant` - Quadrant index (0-3)
- `ring` - Ring index (0-3)
- `moved` - Movement indicator (-1 for deprecated, 0 for unchanged, 1 for improved, 2 for new)
- `label` - Technology name
- `link` - Technology link
- `active` - Active status (always false in current implementation)

The structure is defined in the [RadarEntry](pkg/core/template.go#L9-L16) struct,
and the conversion from Technology to RadarEntry is done in the
[convertTechnologiesToEntries](pkg/usecases/generateradar.go#L74-L98) function.

The `.QuadrantsJSON` contains an array of quadrants with the following structure:

```json
[
  {
    "name": "Quadrant Name",
    "id": "q1"
  }
]
```

The `.RingsJSON` contains an array of rings with the following structure:

```json
[
  {
    "name": "RING NAME",
    "color": "#93c47d",
    "id": "ring-alias"
  }
]
```

### Input Data Format

#### Metadata File (meta.yaml)

```yaml
title: "Technology Radar"
description: "Radar description"
# Optional: custom file name pattern (regex) for technology files
# Default pattern is ^\d{8}\.yaml$ (YYYYMMDD.yaml)
# You can override it to use custom naming convention:
# fileNamePattern: "^radar-\\d{4}-\\d{2}-\\d{2}\\.yaml$"  # radar-YYYY-MM-DD.yaml
# fileNamePattern: "^tech-\\d{8}\\.yaml$"                 # tech-YYYYMMDD.yaml
quadrants:
  - name: "Languages"
    alias: "languages"
  - name: "Frameworks"
    alias: "frameworks"
  - name: "Platforms"
    alias: "platforms"
  - name: "Techniques"
    alias: "techniques"
rings:
  - name: "Adopt"
    alias: "adopt"
  - name: "Trial"
    alias: "trial"
  - name: "Assess"
    alias: "assess"
  - name: "Hold"
    alias: "hold"
```

**Custom File Name Pattern**: By default, TeraGo looks for technology files with names in `YYYYMMDD.yaml` format (e.g., `20231201.yaml`). You can customize this behavior by specifying a `fileNamePattern` in your `meta.yaml` file using regular expression syntax. This allows you to use alternative naming conventions for your technology files, such as:

- `radar-2023-12-01.yaml` with pattern `^radar-\d{4}-\d{2}-\d{2}\.yaml$`
- `tech-20231201.yaml` with pattern `^tech-\d{8}\.yaml$`
- Any other pattern that matches your naming convention

Note: The file name (without the `.yaml` extension) will be used as the date identifier for the radar.

#### Technology Files (YYYYMMDD.yaml)

```yaml
technologies:
  - name: "Go"
    ring: "Adopt"
    quadrant: "Languages"
    description: "An efficient programming language"
  - name: "React"
    ring: "Trial"
    quadrant: "Frameworks"
    description: "A library for building user interfaces"
```

## Project Structure

```
terago/
├── cmd/
│   └── terago/
│       ├── main.go          # Application entry point
│       └── list.go          # List command implementation
├── pkg/
│   ├── core/                # Core data structures
│   ├── radar/               # Embedded HTML template
│   └── usecases/            # Business logic
├── test/
│   └── test_input/          # Test data
└── go.mod                   # Project dependencies
```

## Visualization

This project uses [Zalando Tech Radar](https://github.com/zalando/tech-radar) for visualization.
The embedded template can be found in [pkg/radar/radar.html](pkg/radar/radar.html).

## License

MIT License - see [LICENSE](LICENSE) file for details.

<!-- Badge links -->
[build-badge]: https://github.com/ekalinin/terago/actions/workflows/release.yml/badge.svg
[build-url]: https://github.com/ekalinin/terago/actions/workflows/release.yml
[test-badge]: https://github.com/ekalinin/terago/actions/workflows/test.yml/badge.svg
[test-url]: https://github.com/ekalinin/terago/actions/workflows/test.yml
[download-badge]: https://img.shields.io/github/v/release/ekalinin/terago
[download-url]: https://github.com/ekalinin/terago/releases/latest

