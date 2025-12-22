package main

import (
	"fmt"
	"os"

	"github.com/ekalinin/terago/pkg/core"
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
	case "validate", "val":
		validateCommand(os.Args[2:])
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
	fmt.Fprintf(os.Stderr, "  validate, val       Validate YAML files structure and data\n")
	fmt.Fprintf(os.Stderr, "  version, v          Show version information\n")
	fmt.Fprintf(os.Stderr, "  help, h             Show this help message\n\n")
	fmt.Fprintf(os.Stderr, "Use \"terago <command> -h\" for more information about a command.\n")
}
