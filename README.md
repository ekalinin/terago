# TeraGo - Technology Radar

TeraGo (**Te**chnology **Ra**dar in **Go**) is a tool for creating and visualizing
technology radars based on [Zalando Tech Radar](https://github.com/zalando/tech-radar).
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

```bash
# Clone the repository
$ git clone https://github.com/ekalinin/terago.git
$ cd terago

# Build the project
$ make build
```

## Usage

### Basic Usage

```bash
./terago --input ./test/test_input --output ./output --meta ./test/test_input/test_meta.yaml
```

The `--meta` parameter is optional. If not specified, default values will be used for quadrants and rings. These default values can be found in the [source code](pkg/core/meta.go#L101-L115). The meta.yml file can partially override these default values.

### Command Line Parameters

- `--input` - path to directory with technology YAML files (required)
- `--output` - path to directory for saving HTML files (default: "output")
- `--template` - path to HTML template (default: "./templates/index.html")
- `--meta` - path to metadata file (default: "meta.yaml")

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
│   └── usecases/            # Business logic
├── template/
│   └── radar.html           # HTML template for visualization
├── test/
│   └── test_input/          # Test data
└── go.mod                   # Project dependencies
```

## Visualization

This project uses [Zalando Tech Radar](https://github.com/zalando/tech-radar) for visualization.

## License

MIT License - see [LICENSE](LICENSE) file for details.