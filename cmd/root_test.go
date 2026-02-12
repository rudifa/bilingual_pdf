package cmd

import (
	"strings"
	"testing"

	"bilingual_pdf/internal/parser"
)

func TestReconstructMarkdown_ListWithLink(t *testing.T) {
	// Simulate a translated list that preserves link syntax
	sourceBlock := parser.Block{
		Kind: parser.BlockList,
		Raw:  "1. download from the [Releases](https://example.com/releases) page\n2. unzip the file\n",
	}

	translatedText := "1. descargar del [Releases](https://example.com/releases) página\n2. descomprimir el archivo\n"
	result := reconstructMarkdown(sourceBlock, translatedText)

	// Should pass through the translated text as-is (preserving link syntax)
	if !strings.Contains(result, "[Releases](https://example.com/releases)") {
		t.Errorf("reconstructMarkdown should preserve link syntax, got %q", result)
	}
}

func TestBuildTranslatedBlocks_ListWithLink(t *testing.T) {
	sourceBlocks := []parser.Block{
		{
			Kind: parser.BlockList,
			Raw:  "1. download from the [Releases](https://example.com/releases) page\n2. unzip the file\n",
		},
	}

	// Simulate Google Translate returning the raw markdown with links preserved
	translatedTexts := []string{
		"1. descargar del [Releases](https://example.com/releases) página\n2. descomprimir el archivo\n",
	}

	result := buildTranslatedBlocks(sourceBlocks, translatedTexts)

	if len(result) != 1 {
		t.Fatalf("expected 1 block, got %d", len(result))
	}

	b := result[0]
	if b.Kind != parser.BlockList {
		t.Errorf("expected BlockList, got %v", b.Kind)
	}

	// HTML should contain the link as an <a> tag
	if !strings.Contains(b.HTML, `<a href="https://example.com/releases">Releases</a>`) {
		t.Errorf("translated list HTML should contain <a> tag with href, got %q", b.HTML)
	}
}
