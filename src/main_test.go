package main

import (
	"os"
	"path/filepath"
	"testing"
)

func number(value float64) *float64 {
	return &value
}

func TestValidateAPIResponse(t *testing.T) {
	response := JevResponse{Answers: map[string]Answer{
		"impact_score": {
			Score:      number(3),
			Confidence: number(0.8),
		},
		"has_breaking_changes": {
			Noul: number(0.2),
		},
	}}

	score, breaking, err := validateAPIResponse(response)
	if err != nil {
		t.Fatalf("expected valid response, got %v", err)
	}
	if *score.Score != 3 || *breaking.Noul != 0.2 {
		t.Fatalf("unexpected validated answers: %#v, %#v", score, breaking)
	}
}

func TestValidateAPIResponseRejectsMissingFields(t *testing.T) {
	tests := []JevResponse{
		{},
		{Answers: map[string]Answer{"impact_score": {Score: number(3)}}},
		{Answers: map[string]Answer{
			"impact_score":         {Score: number(3), Confidence: number(0.8)},
			"has_breaking_changes": {},
		}},
	}

	for _, response := range tests {
		if _, _, err := validateAPIResponse(response); err == nil {
			t.Fatal("expected incomplete API response to be rejected")
		}
	}
}

func TestReadFilesLimitsRetainedContent(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "large-diff.txt")
	if err := os.WriteFile(filePath, []byte("abcdefghij"), 0600); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	content, err := readFiles([]string{filePath}, 5)
	if err != nil {
		t.Fatalf("read files: %v", err)
	}
	if content != "abcdef" {
		t.Fatalf("expected max-chars plus one byte, got %q", content)
	}
	if got := truncateContent(content, 5); got != "abcde\n... [truncated]" {
		t.Fatalf("unexpected truncation result: %q", got)
	}
}
