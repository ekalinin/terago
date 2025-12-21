package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ekalinin/terago/pkg/usecases"
)

func generateCommand(args []string) {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)

	inputDir := fs.String("input", "", "Directory path containing YAML files")
	outputDir := fs.String("output", "output", "Directory path for HTML output")
	templatePath := fs.String("template", "", "path to template file (if empty, uses default template)")
	metaPath := fs.String("meta", "", "path to meta file (if empty, searches for meta.yaml in input directory)")
	forceRegenerate := fs.Bool("force", false, "force regeneration of all HTML files (ignore existing files)")
	verbose := fs.Bool("verbose", false, "enable verbose logging (show file processing details)")
	includeLinks := fs.Bool("include-links", false, "include links in radar entries (based on quadrant and technology name)")
	addChanges := fs.Bool("add-changes", false, "add table with description of changed or new technologies")
	skipFirstRadarChanges := fs.Bool("skip-first-radar-changes", true, "skip changes table for the first (earliest) radar (default: true)")
	embedLibs := fs.Bool("embed-libs", false, "embed JavaScript libraries in HTML instead of loading from CDN")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  terago generate -input <directory> [options]\n\n")
		fmt.Fprintf(os.Stderr, "Example:\n")
		fmt.Fprintf(os.Stderr, "  terago generate -input ./data -output ./public -meta ./data/meta.yaml\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
	}

	fs.Parse(args)

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
		OutputDir:             *outputDir,
		TemplatePath:          *templatePath,
		Files:                 files,
		Meta:                  meta,
		Force:                 *forceRegenerate,
		Verbose:               *verbose,
		IncludeLinks:          *includeLinks,
		AddChanges:            *addChanges,
		SkipFirstRadarChanges: *skipFirstRadarChanges,
		EmbedLibs:             *embedLibs,
	}
	if err := generator.Do(); err != nil {
		log.Fatalf("Failed to generate radar: %v", err)
	}

	if *verbose {
		log.Println("Done.")
	}
}

