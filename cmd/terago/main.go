package main

import (
	"flag"
	"log"

	"github.com/ekalinin/terago/pkg/usecases"
)

func main() {
	// Parse command-line arguments
	inputDir := flag.String("input", "", "Directory path containing YAML files")
	outputDir := flag.String("output", "output", "Directory path for HTML output")
	templatePath := flag.String("template", "./templates/index.html", "path to template file")
	metaPath := flag.String("meta", "meta.yaml", "path to meta file")
	// debugMode := flag.Bool("debug", false, "enable debug mode")

	flag.Parse()
	log.Println("Start, input=", *inputDir, ", output=", *outputDir, ", template=", *templatePath, ", meta=", *metaPath)

	// Try to read meta file, use defaults if not available
	log.Printf("Reading meta file: %s", *metaPath)
	meta, err := usecases.ReadMeta(*metaPath)
	if err != nil {
		log.Fatalf("Failed to read meta file: %v", err)
	}

	if *inputDir == "" {
		log.Fatalln("Error: Directory path is required (--input)")
	}
	files, err := usecases.ReadTechnologiesFiles(*inputDir, meta)
	if err != nil {
		log.Fatalf("Failed to read input directory: %v", err)
	}

	if err := usecases.GenerateRadar(*outputDir, *templatePath, files, meta); err != nil {
		log.Fatalf("Failed to generate radar: %v", err)
	}
	log.Println("Done.")
}
