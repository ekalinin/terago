package core

import (
	"testing"
)

func TestTechnologiesFileValidateRingsAndQuadrants(t *testing.T) {
	// Create test meta with default values
	meta := DefaultMeta()

	// Test with valid technologies
	validFile := TechnologiesFile{
		Date: "20231201",
		Technologies: []Technology{
			{
				Name:     "Go",
				Ring:     "Adopt",
				Quadrant: "Languages",
			},
			{
				Name:     "React",
				Ring:     "Trial",
				Quadrant: "Frameworks",
			},
		},
	}

	// Valid file should not return error
	err := validFile.ValidateRingsAndQuadrants(meta)
	if err != nil {
		t.Errorf("ValidateRingsAndQuadrants with valid file returned error: %v", err)
	}

	// Test with invalid ring
	invalidRingFile := TechnologiesFile{
		Date: "20231201",
		Technologies: []Technology{
			{
				Name:     "Invalid Technology",
				Ring:     "InvalidRing",
				Quadrant: "Languages",
			},
		},
	}

	// Invalid ring should return error
	err = invalidRingFile.ValidateRingsAndQuadrants(meta)
	if err == nil {
		t.Error("ValidateRingsAndQuadrants with invalid ring should return error")
	}

	// Test with invalid quadrant
	invalidQuadrantFile := TechnologiesFile{
		Date: "20231201",
		Technologies: []Technology{
			{
				Name:     "Invalid Technology",
				Ring:     "Adopt",
				Quadrant: "InvalidQuadrant",
			},
		},
	}

	// Invalid quadrant should return error
	err = invalidQuadrantFile.ValidateRingsAndQuadrants(meta)
	if err == nil {
		t.Error("ValidateRingsAndQuadrants with invalid quadrant should return error")
	}
}
