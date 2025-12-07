package usecases

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/ekalinin/terago/pkg/core"
)

func ReadMeta(filePath string) (core.Meta, error) {
	var metaFile core.Meta

	data, err := os.ReadFile(filePath)
	if err != nil {
		return metaFile, fmt.Errorf("error reading file: %v", err)
	}

	if err := yaml.Unmarshal(data, &metaFile); err != nil {
		return metaFile, fmt.Errorf("error parsing YAML: %v", err)
	}

	return metaFile, nil
}
