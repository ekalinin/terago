package usecases

import (
	"os"

	"github.com/ekalinin/terago/pkg/radar"
)

// ExportEmbeddedTemplate exports the embedded template to a file
func ExportEmbeddedTemplate(filepath string) error {
	// Create the file
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the embedded template content to the file
	_, err = file.WriteString(radar.HTML)
	return err
}
