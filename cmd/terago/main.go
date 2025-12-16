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
	metaPath := flag.String("meta", "", "path to meta file (if empty, searches for meta.yaml in input directory)")
	showVersion := flag.Bool("version", false, "print version")
	forceRegenerate := flag.Bool("force", false, "force regeneration of all HTML files (ignore existing files)")
	verbose := flag.Bool("verbose", false, "enable verbose logging (show file processing details)")
	includeLinks := flag.Bool("include-links", false, "include links in radar entries (based on quadrant and technology name)")
	addChanges := flag.Bool("add-changes", false, "add table with description of changed or new technologies")
	embedLibs := flag.Bool("embed-libs", false, "embed JavaScript libraries in HTML instead of loading from CDN")

	// Custom help message with version
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "terago version %s\n\n", core.Version)
		fmt.Fprintf(os.Stderr, "Technology Radar Generator - generates interactive HTML radar visualizations from YAML files\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  %s -input <directory> [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example:\n")
		fmt.Fprintf(os.Stderr, "  %s -input ./data -output ./public -meta ./data/meta.yaml\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

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

	// Read input directory (with yaml files)
	if *inputDir == "" {
		log.Fatalln("Error: Directory path is required (--input)")
	}

	// Read meta file
	meta, err := usecases.ReadMeta(*metaPath, *inputDir, *verbose)
	if err != nil {
		log.Fatalf("Failed to read meta file: %v", err)
	}

	files, err := usecases.ReadTechnologiesFiles(*inputDir, meta)
	if err != nil {
		log.Fatalf("Failed to read input directory: %v", err)
	}

	// Generate radar (html files)
	generator := usecases.GenerateRadar{
		OutputDir:    *outputDir,
		TemplatePath: *templatePath,
		Files:        files,
		Meta:         meta,
		Force:        *forceRegenerate,
		Verbose:      *verbose,
		IncludeLinks: *includeLinks,
		AddChanges:   *addChanges,
		EmbedLibs:    *embedLibs,
	}
	if err := generator.Do(); err != nil {
		log.Fatalf("Failed to generate radar: %v", err)
	}

	if *verbose {
		log.Println("Done.")
	}
}
