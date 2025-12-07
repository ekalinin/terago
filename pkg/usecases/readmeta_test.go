package usecases

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/ekalinin/terago/pkg/core"
)

func TestReadMeta(t *testing.T) {
	// Redirect log output to discard to avoid cluttering test output
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	// Get default meta once
	defaultMeta := core.DefaultMeta()

	// Test with empty file path - should return defaults
	meta, err := ReadMeta("")
	if err != nil {
		t.Errorf("ReadMeta(\"\") returned error: %v", err)
	}

	// Check that we got the default values
	if meta.Title != defaultMeta.Title {
		t.Errorf("Expected title %s, got %s", defaultMeta.Title, meta.Title)
	}

	if meta.Description != defaultMeta.Description {
		t.Errorf("Expected description %s, got %s", defaultMeta.Description, meta.Description)
	}

	if len(meta.Quadrants) != len(defaultMeta.Quadrants) {
		t.Errorf("Expected %d quadrants, got %d", len(defaultMeta.Quadrants), len(meta.Quadrants))
	}

	if len(meta.Rings) != len(defaultMeta.Rings) {
		t.Errorf("Expected %d rings, got %d", len(defaultMeta.Rings), len(meta.Rings))
	}

	// Test with non-existent file - should return defaults
	meta, err = ReadMeta("non-existent-file.yaml")
	if err != nil {
		t.Errorf("ReadMeta(\"non-existent-file.yaml\") returned error: %v", err)
	}

	// Check that we got the default values
	if meta.Title != defaultMeta.Title {
		t.Errorf("Expected title %s, got %s", defaultMeta.Title, meta.Title)
	}

	// Test with invalid YAML file - should return defaults
	invalidYamlFile := "test-invalid.yaml"
	invalidYamlContent := `
title: "Test Radar"
invalid: yaml: content:
  - this: is
  - not: valid
`

	err = os.WriteFile(invalidYamlFile, []byte(invalidYamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(invalidYamlFile)

	meta, err = ReadMeta(invalidYamlFile)
	if err != nil {
		t.Errorf("ReadMeta with invalid YAML returned error: %v", err)
	}

	// Should still return defaults
	if meta.Title != defaultMeta.Title {
		t.Errorf("Expected title %s, got %s", defaultMeta.Title, meta.Title)
	}
}
