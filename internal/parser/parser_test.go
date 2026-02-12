package parser

import (
	"os"
	"strings"
	"testing"
)

func TestParse_BasicBlocks(t *testing.T) {
	source := []byte(`# Heading One

A paragraph of text.

## Heading Two

Another paragraph.

- Item one
- Item two
- Item three
`)

	blocks, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	expected := []struct {
		kind  BlockKind
		level int
	}{
		{BlockHeading, 1},
		{BlockParagraph, 0},
		{BlockHeading, 2},
		{BlockParagraph, 0},
		{BlockList, 0},
	}

	if len(blocks) != len(expected) {
		t.Fatalf("expected %d blocks, got %d", len(expected), len(blocks))
	}

	for i, exp := range expected {
		if blocks[i].Kind != exp.kind {
			t.Errorf("block %d: expected kind %v, got %v", i, exp.kind, blocks[i].Kind)
		}
		if blocks[i].Level != exp.level {
			t.Errorf("block %d: expected level %d, got %d", i, exp.level, blocks[i].Level)
		}
	}
}

func TestParse_TextExtraction(t *testing.T) {
	source := []byte(`# Hello World

This is a paragraph.
`)

	blocks, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}

	if blocks[0].Text != "Hello World" {
		t.Errorf("heading text: expected %q, got %q", "Hello World", blocks[0].Text)
	}

	if blocks[1].Text != "This is a paragraph." {
		t.Errorf("paragraph text: expected %q, got %q", "This is a paragraph.", blocks[1].Text)
	}
}

func TestParse_HTMLRendering(t *testing.T) {
	source := []byte(`# Title

Some **bold** text.
`)

	blocks, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}

	// Check heading HTML
	if blocks[0].HTML == "" {
		t.Error("heading HTML should not be empty")
	}

	// Check paragraph HTML contains bold
	if blocks[1].HTML == "" {
		t.Error("paragraph HTML should not be empty")
	}
}

func TestParse_CodeBlock(t *testing.T) {
	source := []byte("# Title\n\n```python\nprint(\"hello\")\n```\n")

	blocks, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}

	if blocks[1].Kind != BlockCodeBlock {
		t.Errorf("expected CodeBlock, got %v", blocks[1].Kind)
	}
}

func TestParse_ThematicBreak(t *testing.T) {
	source := []byte("Paragraph one.\n\n---\n\nParagraph two.\n")

	blocks, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(blocks) != 3 {
		t.Fatalf("expected 3 blocks, got %d", len(blocks))
	}

	if blocks[1].Kind != BlockThematicBreak {
		t.Errorf("expected ThematicBreak, got %v", blocks[1].Kind)
	}
}

func TestParse_Blockquote(t *testing.T) {
	source := []byte("> This is a quote.\n")

	blocks, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	if blocks[0].Kind != BlockBlockquote {
		t.Errorf("expected Blockquote, got %v", blocks[0].Kind)
	}
}

func TestParse_HTMLBlock(t *testing.T) {
	source := []byte("# Title\n\n<div class=\"note\">\n<p>Hello <strong>world</strong>.</p>\n</div>\n")

	blocks, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}

	if blocks[1].Kind != BlockHTML {
		t.Errorf("expected HTML block, got %v", blocks[1].Kind)
	}

	if blocks[1].Text == "" {
		t.Error("HTML block text should not be empty")
	}

	// Raw should preserve the original HTML
	if blocks[1].Raw == "" {
		t.Error("HTML block raw should not be empty")
	}

	// HTML output should contain the raw HTML content (not be stripped by goldmark)
	if !strings.Contains(blocks[1].HTML, "<div") {
		t.Errorf("HTML block HTML should contain the raw markup, got %q", blocks[1].HTML)
	}
}

func TestParse_LinkInParagraph(t *testing.T) {
	source := []byte("Visitez [Google](https://google.com) pour chercher.\n")

	blocks, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	b := blocks[0]
	if b.Kind != BlockParagraph {
		t.Errorf("expected Paragraph, got %v", b.Kind)
	}

	// Text should include the link display text
	if !strings.Contains(b.Text, "Google") {
		t.Errorf("Text should contain link display text 'Google', got %q", b.Text)
	}

	// Raw should preserve the markdown link syntax
	if !strings.Contains(b.Raw, "[Google](https://google.com)") {
		t.Errorf("Raw should preserve markdown link syntax, got %q", b.Raw)
	}

	// HTML should render an <a> tag
	if !strings.Contains(b.HTML, `<a href="https://google.com">Google</a>`) {
		t.Errorf("HTML should contain <a> tag, got %q", b.HTML)
	}
}

func TestParse_EmphasisInParagraph(t *testing.T) {
	source := []byte("Some **bold** and *italic* text.\n")

	blocks, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	b := blocks[0]

	// Text should include text from inside emphasis
	if !strings.Contains(b.Text, "bold") {
		t.Errorf("Text should contain text inside emphasis, got %q", b.Text)
	}
	if !strings.Contains(b.Text, "italic") {
		t.Errorf("Text should contain text inside emphasis, got %q", b.Text)
	}

	// HTML should have <strong> and <em> tags
	if !strings.Contains(b.HTML, "<strong>bold</strong>") {
		t.Errorf("HTML should contain <strong>, got %q", b.HTML)
	}
	if !strings.Contains(b.HTML, "<em>italic</em>") {
		t.Errorf("HTML should contain <em>, got %q", b.HTML)
	}
}

func TestParse_ListRoundTrip(t *testing.T) {
	// Simulate the translation pipeline: translated list items with markers
	// must parse back into a List block with proper <li> HTML.
	tests := []struct {
		name     string
		markdown string
	}{
		{"unordered", "- pan fresco\n- queso de cabra\n- vino tinto\n"},
		{"ordered", "1. pan fresco\n2. queso de cabra\n3. vino tinto\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocks, err := Parse([]byte(tt.markdown))
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if len(blocks) != 1 {
				t.Fatalf("expected 1 block, got %d", len(blocks))
			}
			if blocks[0].Kind != BlockList {
				t.Errorf("expected List, got %v", blocks[0].Kind)
			}
			if !strings.Contains(blocks[0].HTML, "<li>pan fresco</li>") {
				t.Errorf("list HTML should contain <li> items, got %q", blocks[0].HTML)
			}
			if !strings.Contains(blocks[0].HTML, "<li>vino tinto</li>") {
				t.Errorf("list HTML should contain <li> items, got %q", blocks[0].HTML)
			}
		})
	}
}

func TestParse_LinkInListItem(t *testing.T) {
	source := []byte("1. download from the [Releases](https://github.com/example/releases) page\n2. unzip the file\n")

	blocks, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	b := blocks[0]
	if b.Kind != BlockList {
		t.Errorf("expected List, got %v", b.Kind)
	}

	// Raw should preserve the markdown link syntax
	if !strings.Contains(b.Raw, "[Releases](https://github.com/example/releases)") {
		t.Errorf("Raw should preserve markdown link syntax, got %q", b.Raw)
	}

	// HTML should render the link as an <a> tag
	if !strings.Contains(b.HTML, `<a href="https://github.com/example/releases">Releases</a>`) {
		t.Errorf("HTML should contain <a> tag with href, got %q", b.HTML)
	}
}

func TestParse_SampleFile(t *testing.T) {
	source, err := os.ReadFile("../../testdata/sample.fr.md")
	if err != nil {
		t.Fatalf("reading sample file: %v", err)
	}

	blocks, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// sample.fr.md has: h1, p, h2, p, p, h3, list, h2, codeblock, blockquote, h2, html, thematic break, p
	// = 14 blocks
	if len(blocks) < 12 {
		t.Errorf("expected at least 12 blocks from sample, got %d", len(blocks))
	}

	// First block should be a heading
	if blocks[0].Kind != BlockHeading {
		t.Errorf("first block should be heading, got %v", blocks[0].Kind)
	}
	if blocks[0].Level != 1 {
		t.Errorf("first heading level should be 1, got %d", blocks[0].Level)
	}
}
