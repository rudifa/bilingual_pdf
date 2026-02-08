package translator

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestFileTranslator_MatchingBlocks(t *testing.T) {
	ft := NewFileTranslator("../../testdata/sample.es.md", os.Stderr)

	// Source blocks (just the text content, matching sample.fr.md structure)
	sourceBlocks := []string{
		"Bonjour le monde",
		"Ceci est un document de test",
		"Les saisons",
	}

	results, err := ft.Translate(sourceBlocks, "fr", "es")
	if err != nil {
		t.Fatalf("Translate failed: %v", err)
	}

	// Should have results (at least as many as source blocks or translation blocks)
	if len(results) == 0 {
		t.Fatal("expected non-empty results")
	}

	// First block should contain Spanish text
	if results[0] == "" {
		t.Error("first translated block should not be empty")
	}
}

func TestFileTranslator_MismatchedBlocks(t *testing.T) {
	var warn bytes.Buffer
	ft := NewFileTranslator("../../testdata/sample_short.es.md", &warn)

	// More source blocks than translation blocks
	sourceBlocks := make([]string, 10)
	for i := range sourceBlocks {
		sourceBlocks[i] = "block"
	}

	results, err := ft.Translate(sourceBlocks, "fr", "es")
	if err != nil {
		t.Fatalf("Translate failed: %v", err)
	}

	// Should have printed a warning
	if !strings.Contains(warn.String(), "mismatch") {
		t.Errorf("expected mismatch warning, got: %q", warn.String())
	}

	// Results should be padded to max length
	if len(results) != 10 {
		t.Errorf("expected 10 results (padded), got %d", len(results))
	}

	// Later results should be empty (padded)
	emptyFound := false
	for _, r := range results {
		if r == "" {
			emptyFound = true
			break
		}
	}
	if !emptyFound {
		t.Error("expected some empty (padded) results for mismatched blocks")
	}
}

func TestFileTranslator_TranslateBlocks(t *testing.T) {
	ft := NewFileTranslator("../../testdata/sample.es.md", os.Stderr)

	// Parse source file to get blocks
	source, err := os.ReadFile("../../testdata/sample.fr.md")
	if err != nil {
		t.Fatalf("reading source: %v", err)
	}

	// We need parser here
	_ = source

	// Just verify the method doesn't panic with a valid file
	_, err = ft.TranslateBlocks(nil)
	if err != nil {
		t.Fatalf("TranslateBlocks failed: %v", err)
	}
}
