package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ekalinin/terago/pkg/core"
	"github.com/ekalinin/terago/pkg/usecases"
)

func main() {
	// Parse command-line arguments
	inputDir := flag.String("input", "", "Directory path containing YAML files")
	outputDir := flag.String("output", "output", "Directory path for HTML output")
	templatePath := flag.String("template", "", "path to template file (if empty, uses default template)")
	exportTemplate := flag.String("export-template", "", "Export embedded (default) template to file (for customization)")
	metaPath := flag.String("meta", "meta.yaml", "path to meta file")
	showVersion := flag.Bool("version", false, "print version")
	// debugMode := flag.Bool("debug", false, "enable debug mode")

	flag.Parse()

	// Export template if requested
	if *exportTemplate != "" {
		if err := usecases.ExportEmbeddedTemplate(*exportTemplate); err != nil {
			log.Fatalf("Failed to export template: %v", err)
		}
		log.Printf("Template exported to %s\n", *exportTemplate)
		os.Exit(0)
	}

	// Print version if requested
	if *showVersion {
		fmt.Println(core.Version)
		os.Exit(0)
	}

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
