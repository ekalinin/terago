package usecases

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ekalinin/terago/pkg/core"
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
  - name: "Docker"
    ring: "Trial"
    quadrant: "Platforms"
    description: "Container platform"
`
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		err := ValidateTechnologiesFile(filePath, meta)
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

		err := ValidateTechnologiesFile(filePath, meta)
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

		err := ValidateTechnologiesFile(filePath, meta)
		if err == nil {
			t.Error("Expected error for invalid quadrant, got nil")
		}
	})

	t.Run("missing required fields", func(t *testing.T) {
		testCases := []struct {
			name    string
			content string
		}{
			{
				name: "missing name",
				content: `technologies:
  - ring: "Adopt"
    quadrant: "Languages"
    description: "Programming language"
`,
			},
			{
				name: "missing ring",
				content: `technologies:
  - name: "Go"
    quadrant: "Languages"
    description: "Programming language"
`,
			},
			{
				name: "missing quadrant",
				content: `technologies:
  - name: "Go"
    ring: "Adopt"
    description: "Programming language"
`,
			},
			{
				name: "missing description",
				content: `technologies:
  - name: "Go"
    ring: "Adopt"
    quadrant: "Languages"
`,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "test.yaml")

				if err := os.WriteFile(filePath, []byte(tc.content), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}

				err := ValidateTechnologiesFile(filePath, meta)
				if err == nil {
					t.Errorf("Expected error for %s, got nil", tc.name)
				}
			})
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

		err := ValidateTechnologiesFile(filePath, meta)
		if err == nil {
			t.Error("Expected error for empty technologies list, got nil")
		}
	})

	t.Run("invalid yaml syntax", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.yaml")

		content := `invalid yaml content
  - this is: not properly formatted
    - nested: wrong
`
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		err := ValidateTechnologiesFile(filePath, meta)
		if err == nil {
			t.Error("Expected error for invalid YAML syntax, got nil")
		}
	})

	t.Run("file not found", func(t *testing.T) {
		err := ValidateTechnologiesFile("/nonexistent/file.yaml", meta)
		if err == nil {
			t.Error("Expected error for nonexistent file, got nil")
		}
	})

	t.Run("multiple technologies with mixed validity", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.yaml")

		content := `technologies:
  - name: "Go"
    ring: "Adopt"
    quadrant: "Languages"
    description: "Programming language"
  - name: "InvalidTech"
    ring: "InvalidRing"
    quadrant: "Languages"
    description: "Test"
`
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		err := ValidateTechnologiesFile(filePath, meta)
		if err == nil {
			t.Error("Expected error for invalid ring in second technology, got nil")
		}
	})
}
