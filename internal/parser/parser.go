package parser

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// BlockKind identifies the structural type of a markdown block.
type BlockKind int

const (
	BlockHeading BlockKind = iota
	BlockParagraph
	BlockList
	BlockCodeBlock
	BlockBlockquote
	BlockThematicBreak
	BlockHTML
)

func (k BlockKind) String() string {
	switch k {
	case BlockHeading:
		return "Heading"
	case BlockParagraph:
		return "Paragraph"
	case BlockList:
		return "List"
	case BlockCodeBlock:
		return "CodeBlock"
	case BlockBlockquote:
		return "Blockquote"
	case BlockThematicBreak:
		return "ThematicBreak"
	case BlockHTML:
		return "HTML"
	default:
		return "Unknown"
	}
}

// Block represents a single structural unit extracted from the markdown.
type Block struct {
	Kind  BlockKind
	Level int    // heading level (1-6), 0 for non-headings
	Raw   string // reconstructed markdown text (includes syntax like #, -, >, etc.)
	HTML  string // rendered HTML fragment for this block
	Text  string // plain text content (for translation)
}

// Parse reads markdown source bytes and returns an ordered slice of Blocks.
func Parse(source []byte) ([]Block, error) {
	md := goldmark.New()
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	var blocks []Block

	for child := doc.FirstChild(); child != nil; child = child.NextSibling() {
		block := extractBlock(child, source)
		if block == nil {
			continue
		}

		// Render HTML: for HTML blocks the raw content is already HTML;
		// for everything else, convert the reconstructed markdown.
		if block.Kind == BlockHTML {
			block.HTML = block.Raw
		} else {
			var buf bytes.Buffer
			if err := md.Convert([]byte(block.Raw), &buf); err != nil {
				return nil, err
			}
			block.HTML = buf.String()
		}

		blocks = append(blocks, *block)
	}

	return blocks, nil
}

// extractBlock converts an AST node into a Block with proper Raw markdown.
func extractBlock(node ast.Node, source []byte) *Block {
	b := &Block{}

	switch n := node.(type) {
	case *ast.Heading:
		b.Kind = BlockHeading
		b.Level = n.Level
		b.Text = collectText(node, source)
		b.Raw = strings.Repeat("#", n.Level) + " " + b.Text

	case *ast.Paragraph:
		b.Kind = BlockParagraph
		b.Text = collectText(node, source)
		b.Raw = strings.TrimRight(extractLines(node, source), "\n")

	case *ast.List:
		b.Kind = BlockList
		b.Text, b.Raw = extractList(n, source)

	case *ast.FencedCodeBlock:
		b.Kind = BlockCodeBlock
		b.Text = extractCodeContent(n, source)
		lang := ""
		if n.Info != nil {
			lang = string(n.Info.Value(source))
		}
		b.Raw = "```" + lang + "\n" + b.Text + "```"

	case *ast.CodeBlock:
		b.Kind = BlockCodeBlock
		b.Text = extractCodeContent(n, source)
		// Indented code block
		lines := strings.Split(b.Text, "\n")
		var raw strings.Builder
		for _, line := range lines {
			if line != "" {
				raw.WriteString("    " + line + "\n")
			}
		}
		b.Raw = raw.String()

	case *ast.Blockquote:
		b.Kind = BlockBlockquote
		b.Text, b.Raw = extractBlockquote(n, source)

	case *ast.ThematicBreak:
		b.Kind = BlockThematicBreak
		b.Text = "---"
		b.Raw = "---"

	case *ast.HTMLBlock:
		b.Kind = BlockHTML
		b.Text = extractLines(node, source)
		b.Raw = b.Text

	default:
		return nil
	}

	return b
}

// collectText extracts plain text from inline children of a node.
func collectText(node ast.Node, source []byte) string {
	var buf bytes.Buffer
	collectInlineText(&buf, node, source)
	return strings.TrimRight(buf.String(), "\n")
}

// collectInlineText recursively collects text from inline children.
func collectInlineText(buf *bytes.Buffer, node ast.Node, source []byte) {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if child.Type() != ast.TypeInline {
			collectInlineText(buf, child, source)
			continue
		}
		collectInlineNode(buf, child, source)
	}
}

// collectInlineNode extracts text from a single inline node.
func collectInlineNode(buf *bytes.Buffer, node ast.Node, source []byte) {
	if t, ok := node.(*ast.Text); ok {
		buf.Write(t.Segment.Value(source))
		if t.SoftLineBreak() {
			buf.WriteByte('\n')
		}
		return
	}
	if _, ok := node.(*ast.CodeSpan); ok {
		for gc := node.FirstChild(); gc != nil; gc = gc.NextSibling() {
			if t, ok := gc.(*ast.Text); ok {
				buf.Write(t.Segment.Value(source))
			}
		}
		return
	}
	// For other inline nodes (Link, Emphasis, Image, etc.),
	// recurse into their children to collect text content.
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		collectInlineNode(buf, child, source)
	}
}

// extractLines gets line content directly from a node.
func extractLines(node ast.Node, source []byte) string {
	var buf bytes.Buffer
	lines := node.Lines()
	for i := 0; i < lines.Len(); i++ {
		seg := lines.At(i)
		buf.Write(seg.Value(source))
	}
	return buf.String()
}

// extractCodeContent gets the code content from a code block node.
func extractCodeContent(node ast.Node, source []byte) string {
	return extractLines(node, source)
}

// extractList extracts text and raw markdown from a list node.
func extractList(list *ast.List, source []byte) (text, raw string) {
	var textBuf, rawBuf strings.Builder
	idx := 1

	for child := list.FirstChild(); child != nil; child = child.NextSibling() {
		if _, ok := child.(*ast.ListItem); !ok {
			continue
		}
		itemText := collectText(child, source)
		itemContent := extractListItemContent(child, source)
		textBuf.WriteString(itemText)
		textBuf.WriteString("\n")

		if list.IsOrdered() {
			rawBuf.WriteString(fmt.Sprintf("%d. %s\n", idx, itemContent))
			idx++
		} else {
			rawBuf.WriteString("- " + itemContent + "\n")
		}
	}

	return strings.TrimRight(textBuf.String(), "\n"), rawBuf.String()
}

// extractListItemContent gets the original source content from a list item's
// children, preserving inline markdown syntax (links, emphasis, etc.).
func extractListItemContent(item ast.Node, source []byte) string {
	var buf bytes.Buffer
	for child := item.FirstChild(); child != nil; child = child.NextSibling() {
		lines := child.Lines()
		for i := 0; i < lines.Len(); i++ {
			seg := lines.At(i)
			buf.Write(seg.Value(source))
		}
	}
	return strings.TrimRight(buf.String(), "\n")
}

// extractBlockquote extracts text and raw markdown from a blockquote node.
func extractBlockquote(bq *ast.Blockquote, source []byte) (text, raw string) {
	var textBuf, rawBuf strings.Builder

	for child := bq.FirstChild(); child != nil; child = child.NextSibling() {
		childText := collectText(child, source)
		textBuf.WriteString(childText)
		textBuf.WriteString("\n")

		for _, line := range strings.Split(childText, "\n") {
			rawBuf.WriteString("> " + line + "\n")
		}
	}

	return strings.TrimRight(textBuf.String(), "\n"), rawBuf.String()
}
