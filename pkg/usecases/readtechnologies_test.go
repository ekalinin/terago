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

func TestReadTechnologiesFilesWithCustomPattern(t *testing.T) {
	// Create temporary directory for tests
	tempDir, err := os.MkdirTemp("", "terago-test-pattern")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create custom meta with custom file pattern
	meta := core.NewMeta("Test Radar", "Test Description", nil, nil)
	meta.FileNamePattern = `^radar-\d{4}-\d{2}-\d{2}\.yaml$` // radar-YYYY-MM-DD.yaml format

	// Create test files with custom pattern
	customDir := tempDir + "/custom_pattern"
	err = os.MkdirAll(customDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create custom pattern directory: %v", err)
	}

	// Create valid files with custom pattern
	validFile1 := customDir + "/radar-2023-12-01.yaml"
	err = os.WriteFile(validFile1, []byte("date: \"2023-12-01\"\ntechnologies:\n  - name: \"Go\"\n    ring: \"Adopt\"\n    quadrant: \"Languages\"\n    description: \"Go programming language\""), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	validFile2 := customDir + "/radar-2023-12-15.yaml"
	err = os.WriteFile(validFile2, []byte("date: \"2023-12-15\"\ntechnologies:\n  - name: \"React\"\n    ring: \"Trial\"\n    quadrant: \"Frameworks\"\n    description: \"React framework\""), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create file that doesn't match the pattern (should be ignored)
	invalidPatternFile := customDir + "/20231201.yaml"
	err = os.WriteFile(invalidPatternFile, []byte("date: \"20231201\"\ntechnologies:\n  - name: \"Python\"\n    ring: \"Adopt\"\n    quadrant: \"Languages\"\n    description: \"Python programming language\""), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test with custom pattern
	files, err := ReadTechnologiesFiles(customDir, meta)
	if err != nil {
		t.Errorf("ReadTechnologiesFiles with custom pattern returned error: %v", err)
	}

	// Should only find files matching the custom pattern
	if len(files) != 2 {
		t.Errorf("Expected 2 files matching custom pattern, got %d", len(files))
	}

	// Verify that files are sorted correctly
	if len(files) == 2 {
		// Date is extracted as filename without .yaml extension
		if files[0].Date != "radar-2023-12-01" {
			t.Errorf("Expected first file date to be 'radar-2023-12-01', got '%s'", files[0].Date)
		}
		if files[1].Date != "radar-2023-12-15" {
			t.Errorf("Expected second file date to be 'radar-2023-12-15', got '%s'", files[1].Date)
		}
	}
}

func TestReadTechnologiesFilesWithDefaultPattern(t *testing.T) {
	// Create temporary directory for tests
	tempDir, err := os.MkdirTemp("", "terago-test-default")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Use meta with default FileNamePattern
	meta := core.NewMeta("Test Radar", "Test Description", nil, nil)
	// meta.FileNamePattern is already set to default by NewMeta

	// Create test directory
	testDir := tempDir + "/default_pattern"
	err = os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create files with default YYYYMMDD pattern
	file1 := testDir + "/20231201.yaml"
	err = os.WriteFile(file1, []byte("date: \"20231201\"\ntechnologies:\n  - name: \"Go\"\n    ring: \"Adopt\"\n    quadrant: \"Languages\"\n    description: \"Go programming language\""), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	file2 := testDir + "/20231215.yaml"
	err = os.WriteFile(file2, []byte("date: \"20231215\"\ntechnologies:\n  - name: \"React\"\n    ring: \"Trial\"\n    quadrant: \"Frameworks\"\n    description: \"React framework\""), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create file with invalid pattern (should be ignored)
	invalidFile := testDir + "/radar-2023.yaml"
	err = os.WriteFile(invalidFile, []byte("date: \"2023\"\ntechnologies:\n  - name: \"Python\"\n    ring: \"Adopt\"\n    quadrant: \"Languages\"\n    description: \"Python\""), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test with default pattern
	files, err := ReadTechnologiesFiles(testDir, meta)
	if err != nil {
		t.Errorf("ReadTechnologiesFiles with default pattern returned error: %v", err)
	}

	// Should only find files matching the default YYYYMMDD pattern
	if len(files) != 2 {
		t.Errorf("Expected 2 files matching default pattern, got %d", len(files))
	}
}
