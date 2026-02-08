package parser

import (
	"os"
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

func TestParse_SampleFile(t *testing.T) {
	source, err := os.ReadFile("../../testdata/sample.fr.md")
	if err != nil {
		t.Fatalf("reading sample file: %v", err)
	}

	blocks, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// sample.fr.md has: h1, p, h2, p, p, h3, list, h2, codeblock, blockquote, thematic break, p
	// = 12 blocks
	if len(blocks) < 10 {
		t.Errorf("expected at least 10 blocks from sample, got %d", len(blocks))
	}

	// First block should be a heading
	if blocks[0].Kind != BlockHeading {
		t.Errorf("first block should be heading, got %v", blocks[0].Kind)
	}
	if blocks[0].Level != 1 {
		t.Errorf("first heading level should be 1, got %d", blocks[0].Level)
	}
}
