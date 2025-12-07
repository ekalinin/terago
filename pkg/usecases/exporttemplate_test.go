package usecases

import (
	"os"
	"testing"
)

func TestExportEmbeddedTemplate(t *testing.T) {
	// Test file path
	testFile := "test_template.html"

	// Ensure cleanup after test
	defer os.Remove(testFile)

	// Test the ExportEmbeddedTemplate function
	err := ExportEmbeddedTemplate(testFile)
	if err != nil {
		t.Fatalf("ExportEmbeddedTemplate failed: %v", err)
	}

	// Check that file was created
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("ExportEmbeddedTemplate did not create the file")
	}

	// Read the file content
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read exported file: %v", err)
	}

	// Check that content is not empty
	if len(content) == 0 {
		t.Fatal("Exported file is empty")
	}

	// Check that content contains expected HTML elements
	contentStr := string(content)
	if !contains(contentStr, "<!DOCTYPE html>") {
		t.Error("Exported template does not contain DOCTYPE declaration")
	}

	if !contains(contentStr, "<title>{{ .Title }}</title>") {
		t.Error("Exported template does not contain title placeholder")
	}

	if !contains(contentStr, "radar_visualization") {
		t.Error("Exported template does not contain radar_visualization function")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && indexOf(s, substr) != -1
}

// Helper function to find index of substring
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
