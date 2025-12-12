package usecases

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/ekalinin/terago/pkg/core"
)

// ReadMeta reads meta data from file.
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

	metaFile := core.MetaFile{}
	if err := yaml.Unmarshal(data, &metaFile); err != nil {
		log.Printf("Using default meta: failed to parse meta file '%s': %v", filePath, err)
		return core.DefaultMeta(), nil
	}

	meta := core.NewMetaFromFile(metaFile)
	return meta, nil
}
