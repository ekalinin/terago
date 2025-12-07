package usecases

import (
	"os"
	"testing"

	"github.com/ekalinin/terago/pkg/core"
)

func TestReadTechnologiesFiles(t *testing.T) {
	// Use default meta
	meta := core.DefaultMeta()

	// Create temporary directory for tests
	tempDir, err := os.MkdirTemp("", "terago-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test with valid technologies file
	validInputDir := tempDir + "/valid"
	err = os.MkdirAll(validInputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create valid input directory: %v", err)
	}

	validTechFile := validInputDir + "/20231201.yaml"
	// Create valid YAML file
	err = os.WriteFile(validTechFile, []byte("date: \"20231201\"\ntechnologies:\n  - name: \"Go\"\n    ring: \"Adopt\"\n    quadrant: \"Languages\"\n    description: \"Go programming language\"\n  - name: \"React\"\n    ring: \"Trial\"\n    quadrant: \"Frameworks\"\n    description: \"React framework\""), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test with valid files - should succeed
	files, err := ReadTechnologiesFiles(validInputDir, meta)
	if err != nil {
		t.Errorf("ReadTechnologiesFiles with valid files returned error: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}

	if len(files[0].Technologies) != 2 {
		t.Errorf("Expected 2 technologies, got %d", len(files[0].Technologies))
	}

	// Test with invalid ring
	invalidRingDir := tempDir + "/invalid_ring"
	err = os.MkdirAll(invalidRingDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create invalid ring directory: %v", err)
	}

	invalidRingFile := invalidRingDir + "/20231201.yaml"
	// Create YAML file with invalid ring
	err = os.WriteFile(invalidRingFile, []byte("date: \"20231201\"\ntechnologies:\n  - name: \"Invalid Technology\"\n    ring: \"InvalidRing\"\n    quadrant: \"Languages\"\n    description: \"Technology with invalid ring\""), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test with invalid ring - should fail
	_, err = ReadTechnologiesFiles(invalidRingDir, meta)
	if err == nil {
		t.Error("ReadTechnologiesFiles with invalid ring should return error")
	}

	// Test with invalid quadrant
	invalidQuadrantDir := tempDir + "/invalid_quadrant"
	err = os.MkdirAll(invalidQuadrantDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create invalid quadrant directory: %v", err)
	}

	invalidQuadrantFile := invalidQuadrantDir + "/20231201.yaml"
	// Create YAML file with invalid quadrant
	err = os.WriteFile(invalidQuadrantFile, []byte("date: \"20231201\"\ntechnologies:\n  - name: \"Invalid Technology\"\n    ring: \"Adopt\"\n    quadrant: \"InvalidQuadrant\"\n    description: \"Technology with invalid quadrant\""), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test with invalid quadrant - should fail
	_, err = ReadTechnologiesFiles(invalidQuadrantDir, meta)
	if err == nil {
		t.Error("ReadTechnologiesFiles with invalid quadrant should return error")
	}
}
