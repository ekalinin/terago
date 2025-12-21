package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/ekalinin/terago/pkg/usecases"
)

func listCommand(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)

	inputDir := fs.String("input", "", "Directory path containing YAML files")
	outputDir := fs.String("output", "output", "Directory path for HTML output")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  terago list -input <directory> [options]\n\n")
		fmt.Fprintf(os.Stderr, "Example:\n")
		fmt.Fprintf(os.Stderr, "  terago list -input ./data -output ./public\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
	}

	fs.Parse(args)

	if *inputDir == "" {
		log.Fatalln("Error: Directory path is required (--input)")
	}

	// Read meta configuration to get file pattern
	meta, _ := usecases.ReadMeta("", *inputDir, false)

	// Get all YAML files in the input directory
	files, err := filepath.Glob(filepath.Join(*inputDir, "*.yaml"))
	if err != nil {
		log.Fatalf("Failed to read input directory: %v", err)
	}

	// Filter files by pattern from meta
	datePattern := regexp.MustCompile(meta.FileNamePattern)
	var validFiles []string
	for _, file := range files {
		baseName := filepath.Base(file)
		if datePattern.MatchString(baseName) {
			validFiles = append(validFiles, file)
		}
	}

	if len(validFiles) == 0 {
		fmt.Println("No radar files found in", *inputDir)
		return
	}

	// Sort files by date (filename)
	sort.Slice(validFiles, func(i, j int) bool {
		fileNameI := filepath.Base(validFiles[i])
		fileNameJ := filepath.Base(validFiles[j])
		return strings.Compare(fileNameI, fileNameJ) < 0
	})

	fmt.Printf("Found %d radar(s) in %s:\n\n", len(validFiles), *inputDir)

	// Check each radar file
	for _, file := range validFiles {
		fileName := filepath.Base(file)
		dateStr := strings.TrimSuffix(fileName, ".yaml")

		// Check if corresponding HTML file exists
		htmlFile := filepath.Join(*outputDir, dateStr+".html")
		stat, err := os.Stat(htmlFile)

		fmt.Printf("  %s", dateStr)

		if err == nil {
			// HTML file exists
			modTime := stat.ModTime()
			fmt.Printf(" ✓ (rendered: %s)\n", modTime.Format("2006-01-02 15:04:05"))
		} else if os.IsNotExist(err) {
			// HTML file doesn't exist
			fmt.Printf(" ✗ (not rendered)\n")
		} else {
			// Other error
			fmt.Printf(" ? (error checking: %v)\n", err)
		}
	}
	fmt.Println()
}
