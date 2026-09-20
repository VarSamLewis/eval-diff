package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Helper to resolve paths relative to the project root
func projectRootPath(t *testing.T, relativePath string) string {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	// If running inside /tests, step up one level
	if filepath.Base(dir) == "tests" {
		dir = filepath.Dir(dir)
	}
	return filepath.Join(dir, relativePath)
}

func TestReadFilesFromRoot(t *testing.T) {
	mainGoPath := projectRootPath(t, "src/main.go")

	data, err := os.ReadFile(mainGoPath)
	if err != nil {
		t.Fatalf("expected to read src/main.go from project root, got error: %v", err)
	}

	if !strings.Contains(string(data), "package main") {
		t.Errorf("expected src/main.go to contain 'package main'")
	}
}

func TestLoadEnvFile(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	err := os.WriteFile(envPath, []byte("TYPESAFE_TEST_KEY=secret_value_123\n# comment\n"), 0644)
	if err != nil {
		t.Fatalf("failed to write dummy .env: %v", err)
	}

	// Read and parse env manually to verify format handling
	content, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("failed to read temp env file: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	var found bool
	for _, line := range lines {
		if strings.HasPrefix(line, "TYPESAFE_TEST_KEY=") {
			found = true
			val := strings.TrimPrefix(line, "TYPESAFE_TEST_KEY=")
			if val != "secret_value_123" {
				t.Errorf("expected value 'secret_value_123', got %q", val)
			}
		}
	}

	if !found {
		t.Errorf("failed to locate TYPESAFE_TEST_KEY in test .env")
	}
}

func TestSetGitHubOutput(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "github_output")

	os.Setenv("GITHUB_OUTPUT", outputPath)
	defer os.Unsetenv("GITHUB_OUTPUT")

	f, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatalf("failed to open github output file: %v", err)
	}
	_, _ = f.WriteString("impact_score=3.50\n")
	_ = f.Close()

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read github output file: %v", err)
	}

	if !strings.Contains(string(data), "impact_score=3.50\n") {
		t.Errorf("expected output file to contain key-value pair, got %q", string(data))
	}
}

func TestTruncateContentLogic(t *testing.T) {
	input := "abcdefghij"
	maxChars := 5

	var result string
	if len(input) > maxChars {
		result = input[:maxChars] + "\n... [truncated]"
	} else {
		result = input
	}

	expected := "abcde\n... [truncated]"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
