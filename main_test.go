package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProject_IsGem(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(dir string) error
		expected bool
	}{
		{
			name: "with gemspec file",
			setup: func(dir string) error {
				f, err := os.Create(filepath.Join(dir, "test.gemspec"))
				if err != nil {
					return err
				}
				return f.Close()
			},
			expected: true,
		},
		{
			name:     "without gemspec file",
			setup:    func(dir string) error { return nil },
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := tt.setup(dir); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			project := NewProject(dir)
			if got := project.IsGem(); got != tt.expected {
				t.Errorf("IsGem() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestProject_IsSpec(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(dir string) error
		expected bool
	}{
		{
			name: "with spec/spec_helper.rb",
			setup: func(dir string) error {
				specDir := filepath.Join(dir, "spec")
				if err := os.MkdirAll(specDir, 0755); err != nil {
					return err
				}
				f, err := os.Create(filepath.Join(specDir, "spec_helper.rb"))
				if err != nil {
					return err
				}
				return f.Close()
			},
			expected: true,
		},
		{
			name: "with .rspec file",
			setup: func(dir string) error {
				f, err := os.Create(filepath.Join(dir, ".rspec"))
				if err != nil {
					return err
				}
				return f.Close()
			},
			expected: true,
		},
		{
			name:     "without rspec indicators",
			setup:    func(dir string) error { return nil },
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := tt.setup(dir); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			project := NewProject(dir)
			if got := project.IsSpec(); got != tt.expected {
				t.Errorf("IsSpec() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestProject_SrcPaths(t *testing.T) {
	tests := []struct {
		name     string
		isGem    bool
		expected []string
	}{
		{
			name:     "gem project",
			isGem:    true,
			expected: []string{"lib", ""},
		},
		{
			name:     "regular project",
			isGem:    false,
			expected: []string{"app", "lib"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if tt.isGem {
				f, _ := os.Create(filepath.Join(dir, "test.gemspec"))
				f.Close()
			}

			project := NewProject(dir)
			got := project.SrcPaths()

			if len(got) != len(tt.expected) {
				t.Errorf("SrcPaths() returned %d paths, want %d", len(got), len(tt.expected))
				return
			}

			for i, path := range got {
				if path != tt.expected[i] {
					t.Errorf("SrcPaths()[%d] = %q, want %q", i, path, tt.expected[i])
				}
			}
		})
	}
}

func TestProject_TestAnchor(t *testing.T) {
	tests := []struct {
		name     string
		isSpec   bool
		expected string
	}{
		{
			name:     "rspec project",
			isSpec:   true,
			expected: "spec",
		},
		{
			name:     "minitest project",
			isSpec:   false,
			expected: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if tt.isSpec {
				specDir := filepath.Join(dir, "spec")
				os.MkdirAll(specDir, 0755)
				f, _ := os.Create(filepath.Join(specDir, "spec_helper.rb"))
				f.Close()
			}

			project := NewProject(dir)
			if got := project.TestAnchor(); got != tt.expected {
				t.Errorf("TestAnchor() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestProject_Testify(t *testing.T) {
	tests := []struct {
		name     string
		isSpec   bool
		input    string
		expected string
	}{
		{
			name:     "rspec project",
			isSpec:   true,
			input:    "user.rb",
			expected: "user_spec.rb",
		},
		{
			name:     "minitest project",
			isSpec:   false,
			input:    "user.rb",
			expected: "user_test.rb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if tt.isSpec {
				specDir := filepath.Join(dir, "spec")
				os.MkdirAll(specDir, 0755)
				f, _ := os.Create(filepath.Join(specDir, "spec_helper.rb"))
				f.Close()
			}

			project := NewProject(dir)
			if got := project.Testify(tt.input); got != tt.expected {
				t.Errorf("Testify(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestSourceFile_IsTestFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		isSpec   bool
		expected bool
	}{
		{
			name:     "rspec test file",
			filename: "user_spec.rb",
			isSpec:   true,
			expected: true,
		},
		{
			name:     "minitest test file with _test suffix",
			filename: "user_test.rb",
			isSpec:   false,
			expected: true,
		},
		{
			name:     "minitest test file with test_ prefix",
			filename: "test_user.rb",
			isSpec:   false,
			expected: true,
		},
		{
			name:     "source file in rspec project",
			filename: "user.rb",
			isSpec:   true,
			expected: false,
		},
		{
			name:     "source file in minitest project",
			filename: "user.rb",
			isSpec:   false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if tt.isSpec {
				specDir := filepath.Join(dir, "spec")
				os.MkdirAll(specDir, 0755)
				f, _ := os.Create(filepath.Join(specDir, "spec_helper.rb"))
				f.Close()
			}

			project := NewProject(dir)
			sourceFile := NewSourceFile(tt.filename, project)

			if got := sourceFile.IsTestFile(); got != tt.expected {
				t.Errorf("IsTestFile() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSourceFile_AlternateFile(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(dir string) error
		inputFile string
		expected  string
		isSpec    bool
	}{
		{
			name: "find test file for source file in rspec project",
			setupFunc: func(dir string) error {
				// Create lib/user.rb
				libDir := filepath.Join(dir, "lib")
				os.MkdirAll(libDir, 0755)
				f1, _ := os.Create(filepath.Join(libDir, "user.rb"))
				f1.Close()

				// Create spec/lib/user_spec.rb
				specLibDir := filepath.Join(dir, "spec", "lib")
				os.MkdirAll(specLibDir, 0755)
				f2, _ := os.Create(filepath.Join(specLibDir, "user_spec.rb"))
				f2.Close()

				// Create .rspec
				f3, _ := os.Create(filepath.Join(dir, ".rspec"))
				f3.Close()

				return nil
			},
			inputFile: "lib/user.rb",
			expected:  "spec/lib/user_spec.rb",
			isSpec:    true,
		},
		{
			name: "find source file for test file in rspec project",
			setupFunc: func(dir string) error {
				// Create lib/user.rb
				libDir := filepath.Join(dir, "lib")
				os.MkdirAll(libDir, 0755)
				f1, _ := os.Create(filepath.Join(libDir, "user.rb"))
				f1.Close()

				// Create spec/lib/user_spec.rb
				specLibDir := filepath.Join(dir, "spec", "lib")
				os.MkdirAll(specLibDir, 0755)
				f2, _ := os.Create(filepath.Join(specLibDir, "user_spec.rb"))
				f2.Close()

				// Create .rspec
				f3, _ := os.Create(filepath.Join(dir, ".rspec"))
				f3.Close()

				return nil
			},
			inputFile: "spec/lib/user_spec.rb",
			expected:  "lib/user.rb",
			isSpec:    true,
		},
		{
			name: "find test file for source file in minitest project",
			setupFunc: func(dir string) error {
				// Create lib/user.rb
				libDir := filepath.Join(dir, "lib")
				os.MkdirAll(libDir, 0755)
				f1, _ := os.Create(filepath.Join(libDir, "user.rb"))
				f1.Close()

				// Create test/lib/user_test.rb
				testLibDir := filepath.Join(dir, "test", "lib")
				os.MkdirAll(testLibDir, 0755)
				f2, _ := os.Create(filepath.Join(testLibDir, "user_test.rb"))
				f2.Close()

				return nil
			},
			inputFile: "lib/user.rb",
			expected:  "test/lib/user_test.rb",
			isSpec:    false,
		},
		{
			name: "find source file for test file in minitest project",
			setupFunc: func(dir string) error {
				// Create lib/user.rb
				libDir := filepath.Join(dir, "lib")
				os.MkdirAll(libDir, 0755)
				f1, _ := os.Create(filepath.Join(libDir, "user.rb"))
				f1.Close()

				// Create test/lib/user_test.rb
				testLibDir := filepath.Join(dir, "test", "lib")
				os.MkdirAll(testLibDir, 0755)
				f2, _ := os.Create(filepath.Join(testLibDir, "user_test.rb"))
				f2.Close()

				return nil
			},
			inputFile: "test/lib/user_test.rb",
			expected:  "lib/user.rb",
			isSpec:    false,
		},
		{
			name: "no alternate file found",
			setupFunc: func(dir string) error {
				// Create only the source file
				libDir := filepath.Join(dir, "lib")
				os.MkdirAll(libDir, 0755)
				f1, _ := os.Create(filepath.Join(libDir, "user.rb"))
				f1.Close()
				return nil
			},
			inputFile: "lib/user.rb",
			expected:  "",
			isSpec:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := tt.setupFunc(dir); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			project := NewProject(dir)
			sourceFile := NewSourceFile(tt.inputFile, project)

			got := sourceFile.AlternateFile()

			if tt.expected == "" {
				if got != "" {
					t.Errorf("AlternateFile() = %q, want empty string", got)
				}
			} else {
				expectedPath := filepath.Join(dir, tt.expected)
				if got != expectedPath {
					t.Errorf("AlternateFile() = %q, want %q", got, expectedPath)
				}
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	dir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(dir, "test.txt")
	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	f.Close()

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "existing file",
			path:     testFile,
			expected: true,
		},
		{
			name:     "non-existing file",
			path:     filepath.Join(dir, "nonexistent.txt"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fileExists(tt.path); got != tt.expected {
				t.Errorf("fileExists(%q) = %v, want %v", tt.path, got, tt.expected)
			}
		})
	}
}
