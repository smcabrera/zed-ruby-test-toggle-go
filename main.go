package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// Version information (set at build time)
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

// Project represents a Ruby/Rails project structure
type Project struct {
	Root string
}

// NewProject creates a new Project instance
func NewProject(root string) *Project {
	// Remove trailing slash
	root = strings.TrimSuffix(root, "/")
	return &Project{Root: root}
}

// IsGem checks if the project is a gem
func (p *Project) IsGem() bool {
	pattern := filepath.Join(p.Root, "*.gemspec")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return false
	}
	return len(matches) > 0
}

// IsSpec checks if the project uses RSpec
func (p *Project) IsSpec() bool {
	specClues := []string{
		filepath.Join(p.Root, "spec", "spec_helper.rb"),
		filepath.Join(p.Root, ".rspec"),
		filepath.Join(p.Root, "**", "spec", "spec_helper.rb"),
	}

	for _, clue := range specClues {
		matches, err := filepath.Glob(clue)
		if err != nil {
			continue
		}
		if len(matches) > 0 {
			return true
		}
	}
	return false
}

// SrcPaths returns the source paths for the project
func (p *Project) SrcPaths() []string {
	if p.IsGem() {
		return []string{"lib", ""}
	}
	return []string{"app", "lib"}
}

// TestAnchor returns the test directory name
func (p *Project) TestAnchor() string {
	if p.IsSpec() {
		return "spec"
	}
	return "test"
}

// TestPaths returns the test paths for the project
func (p *Project) TestPaths() []string {
	anchor := p.TestAnchor()
	return []string{anchor, filepath.Join(anchor, "lib")}
}

// TestRegexes returns the regexes for matching test files
func (p *Project) TestRegexes() []*regexp.Regexp {
	if p.IsSpec() {
		return []*regexp.Regexp{
			regexp.MustCompile(`_spec\.rb$`),
		}
	}
	return []*regexp.Regexp{
		regexp.MustCompile(`_test\.rb$`),
		regexp.MustCompile(`test_[a-zA-Z0-9_]*\.rb$`),
	}
}

// TestSuffix returns the test file suffix
func (p *Project) TestSuffix() string {
	if p.IsSpec() {
		return "_spec.rb"
	}
	return "_test.rb"
}

// Testify converts a source file path to a test file path
func (p *Project) Testify(path string) string {
	return strings.Replace(path, ".rb", p.TestSuffix(), 1)
}

// SourceFile represents a source or test file
type SourceFile struct {
	Filename string
	Project  *Project
}

// NewSourceFile creates a new SourceFile instance
func NewSourceFile(filename string, project *Project) *SourceFile {
	return &SourceFile{
		Filename: filename,
		Project:  project,
	}
}

// IsTestFile checks if the file is a test file
func (s *SourceFile) IsTestFile() bool {
	for _, regex := range s.Project.TestRegexes() {
		if regex.MatchString(s.Filename) {
			return true
		}
	}
	return false
}

// IsController checks if the file is a Rails controller
func (s *SourceFile) IsController() bool {
	return strings.Contains(s.Filename, "app/controllers/") && strings.HasSuffix(s.Filename, "_controller.rb")
}

// IsRequestSpec checks if the file is a Rails request spec
func (s *SourceFile) IsRequestSpec() bool {
	return strings.Contains(s.Filename, "spec/requests/") && strings.HasSuffix(s.Filename, "_controller_spec.rb")
}

// AlternateFile finds the alternate file (test->source or source->test)
func (s *SourceFile) AlternateFile() string {
	if s.IsTestFile() {
		return s.findAlternateSrc()
	}
	return s.findAlternateTest()
}

// findAlternateSrc finds the source file for a test file
func (s *SourceFile) findAlternateSrc() string {
	// Special handling for request specs with _controller suffix
	if s.IsRequestSpec() {
		candidate := strings.Replace(s.Filename, "spec/requests/", "app/controllers/", 1)
		candidate = strings.Replace(candidate, "_controller_spec.rb", "_controller.rb", 1)
		target := filepath.Join(s.Project.Root, candidate)
		if fileExists(target) {
			return target
		}
	}

	srcPaths := s.Project.SrcPaths()
	testPaths := s.Project.TestPaths()
	testRegexes := s.Project.TestRegexes()

	for _, srcPath := range srcPaths {
		for _, testPath := range testPaths {
			for _, regex := range testRegexes {
				// Replace test path with src path
				candidate := strings.Replace(s.Filename, testPath, srcPath, 1)
				// Replace test suffix with .rb
				candidate = regex.ReplaceAllString(candidate, ".rb")

				target := filepath.Join(s.Project.Root, candidate)
				if fileExists(target) {
					return target
				}
			}
		}
	}
	return ""
}

// findAlternateTest finds the test file for a source file
func (s *SourceFile) findAlternateTest() string {
	// Special handling for controllers -> request specs
	if s.IsController() {
		candidate := strings.Replace(s.Filename, "app/controllers/", "spec/requests/", 1)
		candidate = strings.Replace(candidate, "_controller.rb", "_controller_spec.rb", 1)
		target := filepath.Join(s.Project.Root, candidate)
		if fileExists(target) {
			return target
		}
	}

	testPaths := s.Project.TestPaths()
	srcPaths := s.Project.SrcPaths()

	for _, testPath := range testPaths {
		for _, srcPath := range srcPaths {
			var candidate string
			if srcPath == "" {
				// For empty src path (gem root files), prepend test path
				candidate = filepath.Join(testPath, s.Filename)
			} else {
				// Replace src path with test path
				candidate = strings.Replace(s.Filename, srcPath, testPath, 1)
			}
			// Convert to test file name
			candidate = s.Project.Testify(candidate)

			target := filepath.Join(s.Project.Root, candidate)
			if fileExists(target) {
				return target
			}
		}
	}
	return ""
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// CLI handles command line interface
type CLI struct {
	Root string
	Path string
}

// NewCLI creates a new CLI instance from command line arguments
func NewCLI() *CLI {
	cli := &CLI{}

	// Define the lookup command
	lookupCmd := flag.NewFlagSet("lookup", flag.ExitOnError)
	lookupCmd.StringVar(&cli.Root, "r", "", "Project root directory")
	lookupCmd.StringVar(&cli.Root, "root", "", "Project root directory")
	lookupCmd.StringVar(&cli.Path, "p", "", "Path to file")
	lookupCmd.StringVar(&cli.Path, "path", "", "Path to file")

	// Check if we have a subcommand
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Handle version flag
	if os.Args[1] == "-v" || os.Args[1] == "--version" || os.Args[1] == "version" {
		printVersion()
		os.Exit(0)
	}

	// Parse the subcommand
	switch os.Args[1] {
	case "lookup":
		lookupCmd.Parse(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}

	// Set defaults
	if cli.Root == "" {
		cli.Root, _ = os.Getwd()
	}

	return cli
}

// Run executes the CLI logic
func (c *CLI) Run() error {
	if c.Path == "" {
		return fmt.Errorf("path is required")
	}

	project := NewProject(c.Root)
	sourceFile := NewSourceFile(c.Path, project)

	alternateFile := sourceFile.AlternateFile()
	if alternateFile == "" {
		// No alternate file found, exit silently
		return nil
	}

	// Execute zed command
	cmd := exec.Command("zed", alternateFile)
	return cmd.Run()
}

// printVersion prints version information
func printVersion() {
	fmt.Printf("go-zed-test-toggle %s\n", Version)
	fmt.Printf("Commit: %s\n", Commit)
	fmt.Printf("Built: %s\n", BuildTime)
}

// printUsage prints usage information
func printUsage() {
	fmt.Fprintln(os.Stderr, "go-zed-test-toggle - Toggle between source and test files in Zed editor")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  go-zed-test-toggle lookup [options]  Find and open the alternate file")
	fmt.Fprintln(os.Stderr, "  go-zed-test-toggle version           Show version information")
	fmt.Fprintln(os.Stderr, "  go-zed-test-toggle help              Show this help message")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Lookup options:")
	fmt.Fprintln(os.Stderr, "  -p, --path string    Path to file (required)")
	fmt.Fprintln(os.Stderr, "  -r, --root string    Project root directory (default: current directory)")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintln(os.Stderr, `  go-zed-test-toggle lookup -p "lib/user.rb" -r "/path/to/project"`)
	fmt.Fprintln(os.Stderr, `  go-zed-test-toggle lookup --path="$ZED_RELATIVE_FILE" --root="$ZED_WORKTREE_ROOT"`)
}

func main() {
	cli := NewCLI()
	if err := cli.Run(); err != nil {
		log.Fatal(err)
	}
}
