package usecases

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

	// Test with empty file path and no input dir - should return defaults
	meta, err := ReadMeta("", "", false)
	if err != nil {
		t.Errorf("ReadMeta(\"\", \"\", false) returned error: %v", err)
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
	meta, err = ReadMeta("non-existent-file.yaml", "", false)
	if err != nil {
		t.Errorf("ReadMeta(\"non-existent-file.yaml\", \"\", false) returned error: %v", err)
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

	meta, err = ReadMeta(invalidYamlFile, "", false)
	if err != nil {
		t.Errorf("ReadMeta with invalid YAML returned error: %v", err)
	}

	// Should still return defaults
	if meta.Title != defaultMeta.Title {
		t.Errorf("Expected title %s, got %s", defaultMeta.Title, meta.Title)
	}

	// Test with empty meta path but valid input dir with meta.yaml
	tmpDir, err := os.MkdirTemp("", "terago-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testMetaContent := `
title: "Test Radar"
description: "Test Description"
quadrants:
  - name: "Languages"
    alias: "languages"
rings:
  - name: "Adopt"
    alias: "adopt"
`
	metaPath := filepath.Join(tmpDir, "meta.yaml")
	err = os.WriteFile(metaPath, []byte(testMetaContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test meta file: %v", err)
	}

	meta, err = ReadMeta("", tmpDir, false)
	if err != nil {
		t.Errorf("ReadMeta with input dir returned error: %v", err)
	}

	// Check that we got the values from the file
	if meta.Title != "Test Radar" {
		t.Errorf("Expected title 'Test Radar', got %s", meta.Title)
	}

	if meta.Description != "Test Description" {
		t.Errorf("Expected description 'Test Description', got %s", meta.Description)
	}
}
