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
	forceRegenerate := flag.Bool("force", false, "force regeneration of all HTML files (ignore existing files)")
	verbose := flag.Bool("verbose", false, "enable verbose logging (show file processing details)")
	includeLinks := flag.Bool("include-links", false, "include links in radar entries (based on quadrant and technology name)")

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

	if *verbose {
		log.Println("Start, input=", *inputDir, ", output=", *outputDir, ", template=", *templatePath, ", meta=", *metaPath)
	}

	// Try to read meta file, use defaults if not available
	if *verbose {
		log.Println("Reading meta file: ", *metaPath)
	}
	meta, err := usecases.ReadMeta(*metaPath)
	if err != nil {
		log.Fatalf("Failed to read meta file: %v", err)
	}

	// Read input directory (with yaml files)
	if *inputDir == "" {
		log.Fatalln("Error: Directory path is required (--input)")
	}
	files, err := usecases.ReadTechnologiesFiles(*inputDir, meta)
	if err != nil {
		log.Fatalf("Failed to read input directory: %v", err)
	}

	// Generate radar (html files)
	if err := usecases.GenerateRadar(*outputDir, *templatePath, files, meta, *forceRegenerate, *verbose, *includeLinks); err != nil {
		log.Fatalf("Failed to generate radar: %v", err)
	}

	if *verbose {
		log.Println("Done.")
	}
}
