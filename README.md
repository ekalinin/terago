# TeraGo - Technology Radar

[![Build Status][build-badge]][build-url] [![Download][download-badge]][download-url]

## Table of Contents

- [Key Features](#key-features)
- [Installation](#installation)
  - [From Source](#from-source)
  - [From GitHub Releases](#from-github-releases)
- [Usage](#usage)
  - [Basic Usage](#basic-usage)
  - [Command Line Parameters](#command-line-parameters)
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

To install TeraGo, you need Go 1.20 or higher.

### From Source

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

### Basic Usage

```bash
./terago --input ./test/test_input --output ./output --meta ./test/test_input/test_meta.yaml
```

The `--meta` parameter is optional. If not specified, default values will be
used for quadrants and rings. These default values can be found in the
[source code](pkg/core/meta.go#L101-L115). The meta.yml file can partially
override these default values.

### Command Line Parameters

- `--input` - path to directory with technology YAML files (required)
- `--output` - path to directory for saving HTML files (default: "output")
- `--template` - path to HTML template (if empty, uses default embedded template)
- `--export-template` - export embedded (default) template to file for customization
- `--meta` - path to metadata file (default: "meta.yaml")
- `--version` - print version and exit

### Customizing the Radar Template

TeraGo uses an embedded HTML template for radar visualization. If you want to customize
the appearance of your radar, you can export this template and modify it:

```bash
./terago --export-template ./my-template.html
```

This will create a file `my-template.html` with the default template content.
You can then modify this file according to your needs and use it with the `--template` parameter:

```bash
./terago --input ./test/test_input --output ./output --template ./my-template.html
```

The template uses Go's [text/template](https://pkg.go.dev/text/template) package
and has access to the following data:

- `.Title` - Radar title from metadata
- `.Date` - Current date
- `.Version` - Application version (see [version.go](pkg/core/version.go#L4))
- `.GeneratedAt` - Timestamp when the radar was generated (see [template.go](pkg/core/template.go#L18-L26))
- `.EntriesJSON` - Technologies data in JSON format

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

### Input Data Format

#### Metadata File (meta.yaml)

```yaml
title: "Technology Radar"
description: "Radar description"
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
│       └── main.go          # Application entry point
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
[download-badge]: https://img.shields.io/github/v/release/ekalinin/terago
[download-url]: https://github.com/ekalinin/terago/releases/latest