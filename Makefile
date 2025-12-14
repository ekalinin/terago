# Variables
BINARY_NAME=terago
BUILD_DIR=build
DIST_DIR=dist
TEST_INPUT_DIR=test/test_input
TEST_OUTPUT_DIR=test/test_output
TEMPLATE_PATH=pkg/radar/radar.html
META_PATH=$(TEST_INPUT_DIR)/meta.yaml
RADAR_PKG_DIR=pkg/radar

# JavaScript library versions (single source of truth)
D3_VERSION=7.9.0          # Full D3.js version (for documentation)
D3_MAJOR_VERSION=7        # Major version used in D3.js CDN URL
RADAR_VERSION=0.12        # Zalando Tech Radar version

# JavaScript library URLs and filenames
D3_URL=https://d3js.org/d3.v$(D3_MAJOR_VERSION).min.js
RADAR_URL=https://zalando.github.io/tech-radar/release/radar-$(RADAR_VERSION).js
D3_FILE=$(RADAR_PKG_DIR)/d3.min.js
RADAR_FILE=$(RADAR_PKG_DIR)/radar.min.js

# Main commands
.PHONY: all build clean test run-test help update-libs minify-libs build-minifier

# Default target - build the project
all: build

# Build the project
build:
	@echo "Building the project..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/terago
	@echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# Clean up
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(TEST_OUTPUT_DIR)
	@rm -rf $(DIST_DIR)
	@echo "Cleanup completed"

# Run tests
test: build
	@echo "Running on test data..."
	@mkdir -p $(TEST_OUTPUT_DIR)
	./$(BUILD_DIR)/$(BINARY_NAME) --input $(TEST_INPUT_DIR) --output $(TEST_OUTPUT_DIR) --template $(TEMPLATE_PATH)
	@echo "Test completed. Results in directory: $(TEST_OUTPUT_DIR)"
	@echo "Open result: open $(TEST_OUTPUT_DIR)/*.html"

# Quick run (without rebuilding)
run:
	@echo "Running on test data..."
	@mkdir -p $(TEST_OUTPUT_DIR)
	./$(BUILD_DIR)/$(BINARY_NAME) --input $(TEST_INPUT_DIR) --output $(TEST_OUTPUT_DIR) --template $(TEMPLATE_PATH) --add-changes --force
	@echo "Test completed. Results in directory: $(TEST_OUTPUT_DIR)"
	@echo "Open result: open $(TEST_OUTPUT_DIR)/*.html"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed"

# Build minifier tool
build-minifier:
	@echo "Building minifier tool..."
	@go build -o $(BUILD_DIR)/minifyjs ./cmd/minifyjs
	@echo "Minifier tool built: $(BUILD_DIR)/minifyjs"

# Download JavaScript libraries
download-libs:
	@echo "Downloading JavaScript libraries..."
	@echo "  D3.js version: $(D3_VERSION)"
	@echo "  Radar version: $(RADAR_VERSION)"
	@mkdir -p $(RADAR_PKG_DIR)
	@echo "Downloading D3.js..."
	@curl -sS -L -o $(D3_FILE) $(D3_URL)
	@echo "Downloading Zalando Tech Radar..."
	@curl -sS -L -o $(RADAR_FILE).tmp.js $(RADAR_URL)
	@echo "Libraries downloaded"

# Minify JavaScript libraries
minify-libs: build-minifier
	@echo "Minifying radar.min.js..."
	@if [ -f "$(RADAR_FILE).tmp.js" ]; then \
		./$(BUILD_DIR)/minifyjs -input $(RADAR_FILE).tmp.js -output $(RADAR_FILE); \
	else \
		echo "Error: $(RADAR_FILE).tmp.js not found. Run 'make download-libs' first."; \
		exit 1; \
	fi

# Update JavaScript libraries (download + minify)
update-libs: download-libs minify-libs
	@echo "JavaScript libraries updated successfully"
	@rm -f $(RADAR_FILE).tmp.js
	@echo "Cleanup completed"
	@echo ""
	@echo "Library versions:"
	@echo "  D3.js: $(D3_VERSION)"
	@echo "  Radar: $(RADAR_VERSION)"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted"

# Check code
vet:
	@echo "Checking code..."
	go vet ./...
	@echo "Check completed"

# Run unit tests
test-unit:
	@echo "Running unit tests..."
	go test ./... -v
	@echo "Unit tests completed"

# Run tests and checks
check: fmt vet test-unit test
	@echo "All checks completed"

# Create test data (if not exists)
setup-test:
	@echo "Creating test data..."
	@mkdir -p $(TEST_INPUT_DIR)
	@if [ ! -f $(TEST_INPUT_DIR)/20231201.yaml ]; then \
		echo "date: \"20231201\"" > $(TEST_INPUT_DIR)/20231201.yaml; \
		echo "technologies:" >> $(TEST_INPUT_DIR)/20231201.yaml; \
		echo "  - name: \"Go\"" >> $(TEST_INPUT_DIR)/20231201.yaml; \
		echo "    ring: \"Adopt\"" >> $(TEST_INPUT_DIR)/20231201.yaml; \
		echo "    quadrant: \"Languages\"" >> $(TEST_INPUT_DIR)/20231201.yaml; \
		echo "    description: \"Efficient programming language\"" >> $(TEST_INPUT_DIR)/20231201.yaml; \
		echo "  - name: \"React\"" >> $(TEST_INPUT_DIR)/20231201.yaml; \
		echo "    ring: \"Trial\"" >> $(TEST_INPUT_DIR)/20231201.yaml; \
		echo "    quadrant: \"Frameworks\"" >> $(TEST_INPUT_DIR)/20231201.yaml; \
		echo "    description: \"Library for creating user interfaces\"" >> $(TEST_INPUT_DIR)/20231201.yaml; \
		echo "  - name: \"Docker\"" >> $(TEST_INPUT_DIR)/20231201.yaml; \
		echo "    ring: \"Adopt\"" >> $(TEST_INPUT_DIR)/20231201.yaml; \
		echo "    quadrant: \"Infrastructure\"" >> $(TEST_INPUT_DIR)/20231201.yaml; \
		echo "    description: \"Platform for containerizing applications\"" >> $(TEST_INPUT_DIR)/20231201.yaml; \
		echo "  - name: \"Microservices\"" >> $(TEST_INPUT_DIR)/20231201.yaml; \
		echo "    ring: \"Assess\"" >> $(TEST_INPUT_DIR)/20231201.yaml; \
		echo "    quadrant: \"Architecture\"" >> $(TEST_INPUT_DIR)/20231201.yaml; \
		echo "    description: \"Approach to software architecture\"" >> $(TEST_INPUT_DIR)/20231201.yaml; \
	fi
	@if [ ! -f $(META_PATH) ]; then \
		echo "title: \"Test Technology Radar\"" > $(META_PATH); \
		echo "description: \"Example technology radar for demonstration\"" >> $(META_PATH); \
		echo "quadrants:" >> $(META_PATH); \
		echo "  - name: \"Languages\"" >> $(META_PATH); \
		echo "    alias: \"languages\"" >> $(META_PATH); \
		echo "  - name: \"Frameworks\"" >> $(META_PATH); \
		echo "    alias: \"frameworks\"" >> $(META_PATH); \
		echo "  - name: \"Infrastructure\"" >> $(META_PATH); \
		echo "    alias: \"infrastructure\"" >> $(META_PATH); \
		echo "  - name: \"Architecture\"" >> $(META_PATH); \
		echo "    alias: \"architecture\"" >> $(META_PATH); \
		echo "rings:" >> $(META_PATH); \
		echo "  - name: \"Adopt\"" >> $(META_PATH); \
		echo "    alias: \"adopt\"" >> $(META_PATH); \
		echo "  - name: \"Trial\"" >> $(META_PATH); \
		echo "    alias: \"trial\"" >> $(META_PATH); \
		echo "  - name: \"Assess\"" >> $(META_PATH); \
		echo "    alias: \"assess\"" >> $(META_PATH); \
		echo "  - name: \"Hold\"" >> $(META_PATH); \
		echo "    alias: \"hold\"" >> $(META_PATH); \
	fi
	@echo "Test data created"

# Full setup and run
demo: setup-test build test
	@echo "Demo completed"

# Get version from source
VERSION := $(shell grep -E 'const Version = "([^"]+)"' pkg/core/version.go | cut -d '"' -f 2)

# Create git tag
tag:
	@echo "Creating git tag v$(VERSION)"
	git tag -a "v$(VERSION)" -m "Release version $(VERSION)"
	git push origin "v$(VERSION)"
	@echo "Git tag v$(VERSION) created and pushed"

# Goreleaser release
goreleaser:
	@echo "Creating release with Goreleaser..."
	goreleaser release --clean
	@echo "Release created with Goreleaser"

# Goreleaser snapshot
goreleaser-snapshot:
	@echo "Creating snapshot with Goreleaser..."
	goreleaser release --snapshot --clean
	@echo "Snapshot created with Goreleaser"

# Goreleaser check
goreleaser-check:
	@echo "Checking Goreleaser configuration..."
	goreleaser check
	@echo "Goreleaser configuration is valid"

# Build release
release: clean fmt vet test tag goreleaser


# Help
help:
	@echo "Available commands:"
	@echo "  build              - Build the project"
	@echo "  clean              - Clean up temporary files"
	@echo "  test               - Build and run on test data"
	@echo "  test-unit          - Run unit tests"
	@echo "  run                - Run on test data (without rebuilding)"
	@echo "  deps               - Install dependencies"
	@echo "  fmt                - Format code"
	@echo "  vet                - Check code"
	@echo "  check              - Format, check and test"
	@echo "  setup-test         - Create test data"
	@echo "  demo               - Full setup and demo run"
	@echo "  release            - Build release version"
	@echo "  goreleaser         - Create release with Goreleaser"
	@echo "  goreleaser-snapshot - Create snapshot with Goreleaser"
	@echo "  goreleaser-check   - Check Goreleaser configuration"
	@echo "  build-minifier     - Build JavaScript minifier tool"
	@echo "  download-libs      - Download JavaScript libraries"
	@echo "  minify-libs        - Minify JavaScript libraries"
	@echo "  update-libs        - Download and minify libraries (download-libs + minify-libs)"
	@echo "  help               - Show this help"
