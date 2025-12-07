# Variables
BINARY_NAME=terago
BUILD_DIR=build
TEST_INPUT_DIR=test/test_input
TEST_OUTPUT_DIR=test_output
TEMPLATE_PATH=template/radar.html
META_PATH=$(TEST_INPUT_DIR)/test_meta.yaml

# Main commands
.PHONY: all build clean test run-test help

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
	@echo "Cleanup completed"

# Run tests
test: build
	@echo "Running on test data..."
	@mkdir -p $(TEST_OUTPUT_DIR)
	./$(BUILD_DIR)/$(BINARY_NAME) --input $(TEST_INPUT_DIR) --output $(TEST_OUTPUT_DIR) --template $(TEMPLATE_PATH) --meta $(META_PATH)
	@echo "Test completed. Results in directory: $(TEST_OUTPUT_DIR)"
	@echo "Open result: open $(TEST_OUTPUT_DIR)/*.html"

# Quick run (without rebuilding)
run: 
	@echo "Running on test data..."
	@mkdir -p $(TEST_OUTPUT_DIR)
	./$(BUILD_DIR)/$(BINARY_NAME) --input $(TEST_INPUT_DIR) --output $(TEST_OUTPUT_DIR) --template $(TEMPLATE_PATH) --meta $(META_PATH)
	@echo "Test completed. Results in directory: $(TEST_OUTPUT_DIR)"
	@echo "Open result: open $(TEST_OUTPUT_DIR)/*.html"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed"

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

# Run tests and checks
check: fmt vet test
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

# Build release
release: clean fmt vet test
	@echo "Building release..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/terago
	@echo "Release built: $(BUILD_DIR)/$(BINARY_NAME)"

# Help
help:
	@echo "Available commands:"
	@echo "  build      - Build the project"
	@echo "  clean      - Clean up temporary files"
	@echo "  test       - Build and run on test data"
	@echo "  run        - Run on test data (without rebuilding)"
	@echo "  deps       - Install dependencies"
	@echo "  fmt        - Format code"
	@echo "  vet        - Check code"
	@echo "  check      - Format, check and test"
	@echo "  setup-test - Create test data"
	@echo "  demo       - Full setup and demo run"
	@echo "  release    - Build release version"
	@echo "  help       - Show this help"