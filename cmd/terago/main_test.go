package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ekalinin/terago/pkg/core"
)

// buildBinary builds the terago binary for testing
func buildBinary(t *testing.T) string {
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "terago")

	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	return binaryPath
}

// runCommand runs the terago command with given arguments
func runCommand(t *testing.T, binary string, args ...string) (stdout string, stderr string, exitCode int) {
	cmd := exec.Command(binary, args...)

	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("Failed to run command: %v", err)
		}
	}

	return outBuf.String(), errBuf.String(), exitCode
}

func TestVersionFlag(t *testing.T) {
	binary := buildBinary(t)

	stdout, _, exitCode := runCommand(t, binary, "-version")

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	expectedVersion := core.Version + "\n"
	if stdout != expectedVersion {
		t.Errorf("Expected version %q, got %q", expectedVersion, stdout)
	}
}

func TestHelpFlag(t *testing.T) {
	binary := buildBinary(t)

	_, stderr, exitCode := runCommand(t, binary, "-help")

	// -help flag exits with code 0 in our implementation
	if exitCode != 0 {
		t.Errorf("Expected exit code 0 for -help, got %d", exitCode)
	}

	// Check that help output contains version
	if !strings.Contains(stderr, "terago version "+core.Version) {
		t.Errorf("Help output should contain version, got: %s", stderr)
	}

	// Check that help output contains description
	if !strings.Contains(stderr, "Technology Radar Generator") {
		t.Errorf("Help output should contain description, got: %s", stderr)
	}

	// Check that help output contains usage section
	if !strings.Contains(stderr, "Usage:") {
		t.Errorf("Help output should contain Usage section, got: %s", stderr)
	}

	// Check that help output contains example
	if !strings.Contains(stderr, "Example:") {
		t.Errorf("Help output should contain Example section, got: %s", stderr)
	}

	if !strings.Contains(stderr, "-input ./data -output ./public") {
		t.Errorf("Help output should contain example command, got: %s", stderr)
	}

	// Check that help output contains options
	if !strings.Contains(stderr, "Options:") {
		t.Errorf("Help output should contain Options section, got: %s", stderr)
	}

	// Check that some flags are listed
	if !strings.Contains(stderr, "-input") {
		t.Errorf("Help output should contain -input flag, got: %s", stderr)
	}

	if !strings.Contains(stderr, "-output") {
		t.Errorf("Help output should contain -output flag, got: %s", stderr)
	}
}

func TestExportTemplate(t *testing.T) {
	binary := buildBinary(t)
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "template.html")

	stdout, stderr, exitCode := runCommand(t, binary, "-export-template", templatePath)

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	// Check that template file was created
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		t.Errorf("Template file was not created at %s", templatePath)
	}

	// Check output message (log output goes to stderr)
	output := stdout + stderr
	expectedMsg := "Template exported to " + templatePath
	if !strings.Contains(output, expectedMsg) {
		t.Errorf("Expected output to contain %q, got stdout: %q, stderr: %q", expectedMsg, stdout, stderr)
	}
}

func TestMissingInputDirectory(t *testing.T) {
	binary := buildBinary(t)

	_, stderr, exitCode := runCommand(t, binary)

	if exitCode == 0 {
		t.Error("Expected non-zero exit code when input directory is missing")
	}

	expectedError := "Error: Directory path is required (--input)"
	if !strings.Contains(stderr, expectedError) {
		t.Errorf("Expected error message %q, got %q", expectedError, stderr)
	}
}

func TestWithTestData(t *testing.T) {
	binary := buildBinary(t)

	// Check if test input directory exists
	testInputDir := "../../test/test_input"
	if _, err := os.Stat(testInputDir); os.IsNotExist(err) {
		t.Skip("Test input directory not found, skipping integration test")
	}

	tmpOutputDir := t.TempDir()

	_, stderr, exitCode := runCommand(t, binary,
		"-input", testInputDir,
		"-output", tmpOutputDir,
		"-verbose")

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
	}

	// Check that output files were created
	entries, err := os.ReadDir(tmpOutputDir)
	if err != nil {
		t.Fatalf("Failed to read output directory: %v", err)
	}

	if len(entries) == 0 {
		t.Error("No output files were generated")
	}

	// Check that HTML files were created
	htmlFound := false
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".html") {
			htmlFound = true
			break
		}
	}

	if !htmlFound {
		t.Error("No HTML files were generated")
	}
}

func TestVerboseFlag(t *testing.T) {
	binary := buildBinary(t)

	testInputDir := "../../test/test_input"
	if _, err := os.Stat(testInputDir); os.IsNotExist(err) {
		t.Skip("Test input directory not found, skipping verbose flag test")
	}

	tmpOutputDir := t.TempDir()

	_, stderr, exitCode := runCommand(t, binary,
		"-input", testInputDir,
		"-output", tmpOutputDir,
		"-verbose")

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	// Check that verbose output contains expected messages
	if !strings.Contains(stderr, "Start, input=") {
		t.Errorf("Verbose output should contain start message, got: %s", stderr)
	}

	if !strings.Contains(stderr, "Done.") {
		t.Errorf("Verbose output should contain done message, got: %s", stderr)
	}
}

func TestForceFlag(t *testing.T) {
	binary := buildBinary(t)

	testInputDir := "../../test/test_input"
	if _, err := os.Stat(testInputDir); os.IsNotExist(err) {
		t.Skip("Test input directory not found, skipping force flag test")
	}

	tmpOutputDir := t.TempDir()

	// First run
	_, _, exitCode := runCommand(t, binary,
		"-input", testInputDir,
		"-output", tmpOutputDir)

	if exitCode != 0 {
		t.Errorf("Expected exit code 0 for first run, got %d", exitCode)
	}

	// Second run with -force flag
	_, _, exitCode = runCommand(t, binary,
		"-input", testInputDir,
		"-output", tmpOutputDir,
		"-force")

	if exitCode != 0 {
		t.Errorf("Expected exit code 0 for second run with -force, got %d", exitCode)
	}
}

func TestEmbedLibsFlag(t *testing.T) {
	binary := buildBinary(t)

	testInputDir := "../../test/test_input"
	if _, err := os.Stat(testInputDir); os.IsNotExist(err) {
		t.Skip("Test input directory not found, skipping embed-libs flag test")
	}

	tmpOutputDir := t.TempDir()

	_, _, exitCode := runCommand(t, binary,
		"-input", testInputDir,
		"-output", tmpOutputDir,
		"-embed-libs")

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	// Check that HTML files were created
	entries, err := os.ReadDir(tmpOutputDir)
	if err != nil {
		t.Fatalf("Failed to read output directory: %v", err)
	}

	// Find an HTML file and check that it contains embedded libraries
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".html") {
			htmlPath := filepath.Join(tmpOutputDir, entry.Name())
			content, err := os.ReadFile(htmlPath)
			if err != nil {
				t.Fatalf("Failed to read HTML file: %v", err)
			}

			// When libs are embedded, the HTML should not contain CDN links
			// and should have inline JavaScript
			htmlStr := string(content)
			if strings.Contains(htmlStr, "cdn.jsdelivr.net") {
				t.Error("HTML should not contain CDN links when -embed-libs is used")
			}
			break
		}
	}
}

func TestIncludeLinksFlag(t *testing.T) {
	binary := buildBinary(t)

	testInputDir := "../../test/test_input"
	if _, err := os.Stat(testInputDir); os.IsNotExist(err) {
		t.Skip("Test input directory not found, skipping include-links flag test")
	}

	tmpOutputDir := t.TempDir()

	_, _, exitCode := runCommand(t, binary,
		"-input", testInputDir,
		"-output", tmpOutputDir,
		"-include-links")

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	// Check that HTML files were created
	entries, err := os.ReadDir(tmpOutputDir)
	if err != nil {
		t.Fatalf("Failed to read output directory: %v", err)
	}

	if len(entries) == 0 {
		t.Error("No output files were generated")
	}
}

func TestAddChangesFlag(t *testing.T) {
	binary := buildBinary(t)

	testInputDir := "../../test/test_input"
	if _, err := os.Stat(testInputDir); os.IsNotExist(err) {
		t.Skip("Test input directory not found, skipping add-changes flag test")
	}

	tmpOutputDir := t.TempDir()

	_, _, exitCode := runCommand(t, binary,
		"-input", testInputDir,
		"-output", tmpOutputDir,
		"-add-changes")

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	// Check that HTML files were created
	entries, err := os.ReadDir(tmpOutputDir)
	if err != nil {
		t.Fatalf("Failed to read output directory: %v", err)
	}

	if len(entries) == 0 {
		t.Error("No output files were generated")
	}
}

func TestCustomTemplateFlag(t *testing.T) {
	binary := buildBinary(t)

	testInputDir := "../../test/test_input"
	if _, err := os.Stat(testInputDir); os.IsNotExist(err) {
		t.Skip("Test input directory not found, skipping custom template flag test")
	}

	// First export the template
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "custom_template.html")

	_, _, exitCode := runCommand(t, binary, "-export-template", templatePath)
	if exitCode != 0 {
		t.Fatalf("Failed to export template, exit code: %d", exitCode)
	}

	// Now use the custom template
	tmpOutputDir := t.TempDir()
	_, _, exitCode = runCommand(t, binary,
		"-input", testInputDir,
		"-output", tmpOutputDir,
		"-template", templatePath)

	if exitCode != 0 {
		t.Errorf("Expected exit code 0 when using custom template, got %d", exitCode)
	}

	// Check that HTML files were created
	entries, err := os.ReadDir(tmpOutputDir)
	if err != nil {
		t.Fatalf("Failed to read output directory: %v", err)
	}

	if len(entries) == 0 {
		t.Error("No output files were generated with custom template")
	}
}

func TestCustomMetaFlag(t *testing.T) {
	binary := buildBinary(t)

	testInputDir := "../../test/test_input"
	testMetaPath := "../../test/test_input/meta.yaml"

	if _, err := os.Stat(testInputDir); os.IsNotExist(err) {
		t.Skip("Test input directory not found, skipping custom meta flag test")
	}

	if _, err := os.Stat(testMetaPath); os.IsNotExist(err) {
		t.Skip("Test meta file not found, skipping custom meta flag test")
	}

	tmpOutputDir := t.TempDir()

	_, _, exitCode := runCommand(t, binary,
		"-input", testInputDir,
		"-output", tmpOutputDir,
		"-meta", testMetaPath)

	if exitCode != 0 {
		t.Errorf("Expected exit code 0 when using custom meta, got %d", exitCode)
	}

	// Check that HTML files were created
	entries, err := os.ReadDir(tmpOutputDir)
	if err != nil {
		t.Fatalf("Failed to read output directory: %v", err)
	}

	if len(entries) == 0 {
		t.Error("No output files were generated with custom meta")
	}
}

func TestMultipleFlags(t *testing.T) {
	binary := buildBinary(t)

	testInputDir := "../../test/test_input"
	if _, err := os.Stat(testInputDir); os.IsNotExist(err) {
		t.Skip("Test input directory not found, skipping multiple flags test")
	}

	tmpOutputDir := t.TempDir()

	_, stderr, exitCode := runCommand(t, binary,
		"-input", testInputDir,
		"-output", tmpOutputDir,
		"-verbose",
		"-force",
		"-include-links",
		"-add-changes",
		"-embed-libs")

	if exitCode != 0 {
		t.Errorf("Expected exit code 0 with multiple flags, got %d. Stderr: %s", exitCode, stderr)
	}

	// Check that HTML files were created
	entries, err := os.ReadDir(tmpOutputDir)
	if err != nil {
		t.Fatalf("Failed to read output directory: %v", err)
	}

	if len(entries) == 0 {
		t.Error("No output files were generated with multiple flags")
	}

	// Check verbose output
	if !strings.Contains(stderr, "Start, input=") {
		t.Error("Verbose output should be present with -verbose flag")
	}
}

func TestInvalidMetaPath(t *testing.T) {
	binary := buildBinary(t)

	testInputDir := "../../test/test_input"
	if _, err := os.Stat(testInputDir); os.IsNotExist(err) {
		t.Skip("Test input directory not found, skipping invalid meta path test")
	}

	invalidMetaPath := "/non/existent/meta.yaml"
	tmpOutputDir := t.TempDir()

	_, stderr, exitCode := runCommand(t, binary,
		"-input", testInputDir,
		"-output", tmpOutputDir,
		"-meta", invalidMetaPath)

	if exitCode == 0 {
		t.Error("Expected non-zero exit code for invalid meta path")
	}

	if stderr == "" {
		t.Error("Expected error message for invalid meta path")
	}
}

func TestInvalidTemplatePath(t *testing.T) {
	binary := buildBinary(t)

	testInputDir := "../../test/test_input"
	if _, err := os.Stat(testInputDir); os.IsNotExist(err) {
		t.Skip("Test input directory not found, skipping invalid template path test")
	}

	invalidTemplatePath := "/non/existent/template.html"
	tmpOutputDir := t.TempDir()

	_, stderr, exitCode := runCommand(t, binary,
		"-input", testInputDir,
		"-output", tmpOutputDir,
		"-template", invalidTemplatePath)

	if exitCode == 0 {
		t.Error("Expected non-zero exit code for invalid template path")
	}

	if stderr == "" {
		t.Error("Expected error message for invalid template path")
	}
}

