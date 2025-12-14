package main

import (
	"os"
	"strings"
	"testing"
)

func TestMinifyEmptyFile(t *testing.T) {
	// Create temporary empty file
	tmpInput, err := os.CreateTemp("", "empty-*.js")
	if err != nil {
		t.Fatalf("Failed to create temp input file: %v", err)
	}
	defer os.Remove(tmpInput.Name())
	tmpInput.Close()

	tmpOutput, err := os.CreateTemp("", "empty-out-*.js")
	if err != nil {
		t.Fatalf("Failed to create temp output file: %v", err)
	}
	defer os.Remove(tmpOutput.Name())
	tmpOutput.Close()

	// This should not panic with division by zero
	err = minifyFile(tmpInput.Name(), tmpOutput.Name())
	if err != nil {
		t.Fatalf("minifyFile failed: %v", err)
	}

	// Verify output file exists and is empty
	content, err := os.ReadFile(tmpOutput.Name())
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if len(content) != 0 {
		t.Errorf("Expected empty output, got %d bytes", len(content))
	}
}

func TestMinifyNonEmptyFile(t *testing.T) {
	// Create temporary file with JavaScript content
	tmpInput, err := os.CreateTemp("", "test-*.js")
	if err != nil {
		t.Fatalf("Failed to create temp input file: %v", err)
	}
	defer os.Remove(tmpInput.Name())

	inputContent := "function test() { return 42; }"
	if _, err := tmpInput.WriteString(inputContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpInput.Close()

	tmpOutput, err := os.CreateTemp("", "test-out-*.js")
	if err != nil {
		t.Fatalf("Failed to create temp output file: %v", err)
	}
	defer os.Remove(tmpOutput.Name())
	tmpOutput.Close()

	// Run minifier
	err = minifyFile(tmpInput.Name(), tmpOutput.Name())
	if err != nil {
		t.Fatalf("minifyFile failed: %v", err)
	}

	// Verify output file exists and is smaller or equal to input
	content, err := os.ReadFile(tmpOutput.Name())
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if len(content) == 0 {
		t.Error("Expected non-empty output")
	}

	if len(content) > len(inputContent) {
		t.Errorf("Minified size (%d) is larger than original (%d)", len(content), len(inputContent))
	}

	// Verify it's valid JavaScript (basic check - contains "function")
	if !strings.Contains(string(content), "function") {
		t.Error("Minified output doesn't contain expected JavaScript keyword")
	}
}
