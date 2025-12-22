package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ekalinin/terago/pkg/usecases"
)

func validateCommand(args []string) {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)

	inputDir := fs.String("input", "", "Directory path containing YAML files")
	metaPath := fs.String("meta", "", "Path to meta.yaml file (optional)")
	verbose := fs.Bool("verbose", false, "Verbose output - show status for each file")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  terago validate -input <directory> [options]\n\n")
		fmt.Fprintf(os.Stderr, "Example:\n")
		fmt.Fprintf(os.Stderr, "  terago validate -input ./data\n")
		fmt.Fprintf(os.Stderr, "  terago validate -input ./data --verbose\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
	}

	fs.Parse(args)

	if *inputDir == "" {
		log.Fatalln("Error: Directory path is required (--input)")
	}

	// Read meta configuration
	meta, err := usecases.ReadMeta(*metaPath, *inputDir, false)
	if err != nil {
		log.Fatalf("Failed to read meta: %v", err)
	}

	// Get all valid YAML files
	validFiles, err := usecases.GetRadarFiles(*inputDir, meta)
	if err != nil {
		log.Fatalf("Failed to read input directory: %v", err)
	}

	if len(validFiles) == 0 {
		fmt.Println("No radar files found in", *inputDir)
		return
	}

	// Validate each file
	errorCount := 0
	for _, file := range validFiles {
		fileName := filepath.Base(file)
		err := usecases.ValidateTechnologiesFile(file, meta)

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s - %v\n", fileName, err)
			errorCount++
		} else {
			if *verbose {
				fmt.Printf("OK: %s\n", fileName)
			}
		}
	}

	// Print summary
	if errorCount > 0 {
		fmt.Fprintf(os.Stderr, "\nValidation completed with errors: %d file(s) failed, %d file(s) passed\n",
			errorCount, len(validFiles)-errorCount)
		os.Exit(1)
	} else {
		if !*verbose {
			fmt.Printf("OK: %d file(s) processed\n", len(validFiles))
		} else {
			fmt.Printf("\nOK: %d file(s) processed\n", len(validFiles))
		}
	}
}
