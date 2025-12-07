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

	// Test with empty file path - should return defaults
	meta, err := ReadMeta("")
	if err != nil {
		t.Errorf("ReadMeta(\"\") returned error: %v", err)
	}

	// Check that we got the default values
	if meta.Title != core.DefaultMeta.Title {
		t.Errorf("Expected title %s, got %s", core.DefaultMeta.Title, meta.Title)
	}

	if meta.Description != core.DefaultMeta.Description {
		t.Errorf("Expected description %s, got %s", core.DefaultMeta.Description, meta.Description)
	}

	if len(meta.Quadrants) != len(core.DefaultMeta.Quadrants) {
		t.Errorf("Expected %d quadrants, got %d", len(core.DefaultMeta.Quadrants), len(meta.Quadrants))
	}

	if len(meta.Rings) != len(core.DefaultMeta.Rings) {
		t.Errorf("Expected %d rings, got %d", len(core.DefaultMeta.Rings), len(meta.Rings))
	}

	// Test with non-existent file - should return defaults
	meta, err = ReadMeta("non-existent-file.yaml")
	if err != nil {
		t.Errorf("ReadMeta(\"non-existent-file.yaml\") returned error: %v", err)
	}

	// Check that we got the default values
	if meta.Title != core.DefaultMeta.Title {
		t.Errorf("Expected title %s, got %s", core.DefaultMeta.Title, meta.Title)
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
	if meta.Title != core.DefaultMeta.Title {
		t.Errorf("Expected title %s, got %s", core.DefaultMeta.Title, meta.Title)
	}
}
