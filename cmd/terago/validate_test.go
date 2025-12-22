package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ekalinin/terago/pkg/core"
	"github.com/ekalinin/terago/pkg/usecases"
)

func TestValidateTechnologiesFile(t *testing.T) {
	meta := core.DefaultMeta()

	t.Run("valid file", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.yaml")

		content := `technologies:
  - name: "Go"
    ring: "Adopt"
    quadrant: "Languages"
    description: "Programming language"
`
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		err := usecases.ValidateTechnologiesFile(filePath, meta)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("invalid ring", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.yaml")

		content := `technologies:
  - name: "Go"
    ring: "InvalidRing"
    quadrant: "Languages"
    description: "Programming language"
`
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		err := usecases.ValidateTechnologiesFile(filePath, meta)
		if err == nil {
			t.Error("Expected error for invalid ring, got nil")
		}
	})

	t.Run("invalid quadrant", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.yaml")

		content := `technologies:
  - name: "Go"
    ring: "Adopt"
    quadrant: "InvalidQuadrant"
    description: "Programming language"
`
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		err := usecases.ValidateTechnologiesFile(filePath, meta)
		if err == nil {
			t.Error("Expected error for invalid quadrant, got nil")
		}
	})

	t.Run("missing name", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.yaml")

		content := `technologies:
  - ring: "Adopt"
    quadrant: "Languages"
    description: "Programming language"
`
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		err := usecases.ValidateTechnologiesFile(filePath, meta)
		if err == nil {
			t.Error("Expected error for missing name, got nil")
		}
	})

	t.Run("missing ring", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.yaml")

		content := `technologies:
  - name: "Go"
    quadrant: "Languages"
    description: "Programming language"
`
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		err := usecases.ValidateTechnologiesFile(filePath, meta)
		if err == nil {
			t.Error("Expected error for missing ring, got nil")
		}
	})

	t.Run("missing quadrant", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.yaml")

		content := `technologies:
  - name: "Go"
    ring: "Adopt"
    description: "Programming language"
`
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		err := usecases.ValidateTechnologiesFile(filePath, meta)
		if err == nil {
			t.Error("Expected error for missing quadrant, got nil")
		}
	})

	t.Run("missing description", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.yaml")

		content := `technologies:
  - name: "Go"
    ring: "Adopt"
    quadrant: "Languages"
`
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		err := usecases.ValidateTechnologiesFile(filePath, meta)
		if err == nil {
			t.Error("Expected error for missing description, got nil")
		}
	})

	t.Run("empty technologies list", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.yaml")

		content := `technologies: []
`
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		err := usecases.ValidateTechnologiesFile(filePath, meta)
		if err == nil {
			t.Error("Expected error for empty technologies list, got nil")
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.yaml")

		content := `invalid yaml content
  - this is: not properly formatted
    - nested: wrong
`
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		err := usecases.ValidateTechnologiesFile(filePath, meta)
		if err == nil {
			t.Error("Expected error for invalid YAML, got nil")
		}
	})

	t.Run("file not found", func(t *testing.T) {
		err := usecases.ValidateTechnologiesFile("/nonexistent/file.yaml", meta)
		if err == nil {
			t.Error("Expected error for nonexistent file, got nil")
		}
	})
}
