package usecases

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/ekalinin/terago/pkg/core"
)

func ReadMeta(filePath string) (core.Meta, error) {
	// If filePath is empty, return default meta
	if filePath == "" {
		log.Println("Using default meta: no meta file specified")
		return core.DefaultMeta(), nil
	}

	// Try to read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Using default meta: failed to read meta file '%s': %v", filePath, err)
		return core.DefaultMeta(), nil
	}

	// Try to parse YAML
	var metaFile core.Meta
	if err := yaml.Unmarshal(data, &metaFile); err != nil {
		log.Printf("Using default meta: failed to parse meta file '%s': %v", filePath, err)
		return core.DefaultMeta(), nil
	}

	// Apply defaults for missing fields
	needsDefaults := false
	defaultMeta := core.DefaultMeta()
	if metaFile.Title == "" {
		metaFile.Title = defaultMeta.Title
		needsDefaults = true
	}
	if metaFile.Description == "" {
		metaFile.Description = defaultMeta.Description
		needsDefaults = true
	}
	if len(metaFile.Quadrants) == 0 {
		metaFile.Quadrants = defaultMeta.Quadrants
		needsDefaults = true
	}
	if len(metaFile.Rings) == 0 {
		metaFile.Rings = defaultMeta.Rings
		needsDefaults = true
	}

	// If we applied any defaults, inform the user
	if needsDefaults {
		log.Printf("Applied default values for missing fields in meta file '%s'", filePath)
	}

	return metaFile, nil
}
