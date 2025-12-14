package usecases

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/ekalinin/terago/pkg/core"
)

// ReadMeta reads meta data from file.
// If metaPath is empty, searches for meta.yaml in inputDir.
func ReadMeta(metaPath, inputDir string, verbose bool) (core.Meta, error) {
	// Determine meta file path
	filePath := metaPath

	// Early return: no meta path and no input dir
	if filePath == "" && inputDir == "" {
		log.Println("Using default meta: no meta file specified")
		return core.DefaultMeta(), nil
	}

	// If no explicit meta path, search for meta.yaml in input directory
	if filePath == "" {
		filePath = filepath.Join(inputDir, "meta.yaml")
		if verbose {
			log.Println("No meta path specified, using default:", filePath)
		}
	}

	if verbose {
		log.Println("Reading meta file:", filePath)
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
