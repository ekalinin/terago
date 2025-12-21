package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ekalinin/terago/pkg/usecases"
)

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
