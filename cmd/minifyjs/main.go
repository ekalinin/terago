package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

// minifyFile reads, minifies, and writes a JavaScript file
func minifyFile(inputFile, outputFile string) error {
	// Read input file
	input, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Create minifier
	m := minify.New()
	m.AddFunc("text/javascript", js.Minify)

	// Minify
	minified, err := m.Bytes("text/javascript", input)
	if err != nil {
		return fmt.Errorf("failed to minify: %w", err)
	}

	// Write output file
	if err := os.WriteFile(outputFile, minified, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	// Print statistics
	fmt.Printf("Minification completed:\n")
	fmt.Printf("  Original size: %d bytes\n", len(input))
	fmt.Printf("  Minified size: %d bytes\n", len(minified))
	if len(input) > 0 {
		reduction := float64(len(input)-len(minified)) / float64(len(input)) * 100
		fmt.Printf("  Reduction: %.1f%%\n", reduction)
	} else {
		fmt.Printf("  Reduction: 0.0%% (empty file)\n")
	}

	return nil
}

func main() {
	inputFile := flag.String("input", "", "Input JavaScript file")
	outputFile := flag.String("output", "", "Output minified JavaScript file")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		log.Fatal("Both --input and --output flags are required")
	}

	if err := minifyFile(*inputFile, *outputFile); err != nil {
		log.Fatal(err)
	}
}
