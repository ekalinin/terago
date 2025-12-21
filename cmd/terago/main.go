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
	if len(os.Args) < 2 {
		printMainUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "generate", "g":
		generateCommand(os.Args[2:])
	case "export-template", "e":
		exportTemplateCommand(os.Args[2:])
	case "list", "l":
		listCommand(os.Args[2:])
	case "version", "v", "-version", "--version":
		fmt.Println(core.Version)
		os.Exit(0)
	case "help", "h", "-help", "--help", "-h":
		printMainUsage()
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printMainUsage()
		os.Exit(1)
	}
}

func printMainUsage() {
	fmt.Fprintf(os.Stderr, "terago version %s\n\n", core.Version)
	fmt.Fprintf(os.Stderr, "Technology Radar Generator - generates interactive HTML radar visualizations from YAML files\n\n")
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  terago <command> [options]\n\n")
	fmt.Fprintf(os.Stderr, "Available Commands:\n")
	fmt.Fprintf(os.Stderr, "  generate, g         Generate HTML radars from YAML files\n")
	fmt.Fprintf(os.Stderr, "  export-template, e  Export embedded template to file for customization\n")
	fmt.Fprintf(os.Stderr, "  list, l             List available radars and their render status\n")
	fmt.Fprintf(os.Stderr, "  version, v          Show version information\n")
	fmt.Fprintf(os.Stderr, "  help, h             Show this help message\n\n")
	fmt.Fprintf(os.Stderr, "Use \"terago <command> -h\" for more information about a command.\n")
}

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

func exportTemplateCommand(args []string) {
	fs := flag.NewFlagSet("export-template", flag.ExitOnError)

	outputPath := fs.String("output", "", "Output file path for the template")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  terago export-template -output <file>\n\n")
		fmt.Fprintf(os.Stderr, "Example:\n")
		fmt.Fprintf(os.Stderr, "  terago export-template -output ./my-template.html\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
	}

	fs.Parse(args)

	if *outputPath == "" {
		log.Fatalln("Error: Output file path is required (--output)")
	}

	if err := usecases.ExportEmbeddedTemplate(*outputPath); err != nil {
		log.Fatalf("Failed to export template: %v", err)
	}

	log.Printf("Template exported to %s\n", *outputPath)
}
