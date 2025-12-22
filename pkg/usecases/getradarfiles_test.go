package usecases

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ekalinin/terago/pkg/core"
)

func TestGetRadarFiles(t *testing.T) {
	t.Run("valid files matching pattern", func(t *testing.T) {
		tmpDir := t.TempDir()
		meta := core.DefaultMeta()

		// Create test files
		testFiles := []string{"20231201.yaml", "20231202.yaml", "20231203.yaml"}
		for _, filename := range testFiles {
			content := `technologies:
  - name: "Test"
    ring: "Adopt"
    quadrant: "Languages"
    description: "Test technology"
`
			filePath := filepath.Join(tmpDir, filename)
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
		}

		// Create meta.yaml (should be excluded)
		metaContent := `title: "Test Radar"
description: "Test"
`
		metaPath := filepath.Join(tmpDir, "meta.yaml")
		if err := os.WriteFile(metaPath, []byte(metaContent), 0644); err != nil {
			t.Fatalf("Failed to create meta file: %v", err)
		}

		// Get radar files
		files, err := GetRadarFiles(tmpDir, meta)
		if err != nil {
			t.Fatalf("GetRadarFiles failed: %v", err)
		}

		// Verify count
		if len(files) != 3 {
			t.Errorf("Expected 3 files, got %d", len(files))
		}

		// Verify sorted order
		for i, file := range files {
			expected := filepath.Join(tmpDir, testFiles[i])
			if file != expected {
				t.Errorf("File %d: expected %s, got %s", i, expected, file)
			}
		}
	})

	t.Run("no matching files", func(t *testing.T) {
		tmpDir := t.TempDir()
		meta := core.DefaultMeta()

		// Create files that don't match the pattern
		invalidFiles := []string{"test.yaml", "data.yaml", "config.yaml"}
		for _, filename := range invalidFiles {
			filePath := filepath.Join(tmpDir, filename)
			if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
		}

		// Get radar files
		files, err := GetRadarFiles(tmpDir, meta)
		if err != nil {
			t.Fatalf("GetRadarFiles failed: %v", err)
		}

		// Verify empty result
		if len(files) != 0 {
			t.Errorf("Expected 0 files, got %d", len(files))
		}
	})

	t.Run("empty directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		meta := core.DefaultMeta()

		// Get radar files from empty directory
		files, err := GetRadarFiles(tmpDir, meta)
		if err != nil {
			t.Fatalf("GetRadarFiles failed: %v", err)
		}

		// Verify empty result
		if len(files) != 0 {
			t.Errorf("Expected 0 files, got %d", len(files))
		}
	})

	t.Run("files are sorted", func(t *testing.T) {
		tmpDir := t.TempDir()
		meta := core.DefaultMeta()

		// Create test files in random order
		testFiles := []string{"20231203.yaml", "20231201.yaml", "20231202.yaml"}
		for _, filename := range testFiles {
			content := `technologies:
  - name: "Test"
    ring: "Adopt"
    quadrant: "Languages"
    description: "Test technology"
`
			filePath := filepath.Join(tmpDir, filename)
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
		}

		// Get radar files
		files, err := GetRadarFiles(tmpDir, meta)
		if err != nil {
			t.Fatalf("GetRadarFiles failed: %v", err)
		}

		// Verify sorted order
		expectedOrder := []string{"20231201.yaml", "20231202.yaml", "20231203.yaml"}
		for i, file := range files {
			expected := filepath.Join(tmpDir, expectedOrder[i])
			if file != expected {
				t.Errorf("File %d: expected %s, got %s", i, expected, file)
			}
		}
	})
}
