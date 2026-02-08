package translator

import (
	"fmt"
	"io"
	"os"

	"bilingual_pdf/internal/parser"
)

// FileTranslator reads translations from a pre-translated markdown file.
type FileTranslator struct {
	Path    string
	Warn    io.Writer // where to print warnings (typically os.Stderr)
}

// NewFileTranslator creates a FileTranslator for the given file path.
func NewFileTranslator(path string, warn io.Writer) *FileTranslator {
	return &FileTranslator{
		Path: path,
		Warn: warn,
	}
}

func (f *FileTranslator) Translate(blocks []string, source, target string) ([]string, error) {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		return nil, fmt.Errorf("reading translation file: %w", err)
	}

	transBlocks, err := parser.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parsing translation file: %w", err)
	}

	srcCount := len(blocks)
	tgtCount := len(transBlocks)

	if srcCount != tgtCount {
		fmt.Fprintf(f.Warn, "Warning: block count mismatch — source has %d blocks, translation has %d blocks\n", srcCount, tgtCount)
	}

	// Build result, padding shorter side with empty strings
	maxLen := srcCount
	if tgtCount > maxLen {
		maxLen = tgtCount
	}

	results := make([]string, maxLen)
	for i := 0; i < maxLen; i++ {
		if i < tgtCount {
			results[i] = transBlocks[i].Text
		}
	}

	return results, nil
}

// TranslateBlocks translates parser.Block slices and returns translated Blocks
// with HTML already rendered. This is used when we need the full Block info.
func (f *FileTranslator) TranslateBlocks(sourceBlocks []parser.Block) ([]parser.Block, error) {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		return nil, fmt.Errorf("reading translation file: %w", err)
	}

	transBlocks, err := parser.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parsing translation file: %w", err)
	}

	srcCount := len(sourceBlocks)
	tgtCount := len(transBlocks)

	if srcCount != tgtCount {
		fmt.Fprintf(f.Warn, "Warning: block count mismatch — source has %d blocks, translation has %d blocks\n", srcCount, tgtCount)
	}

	// Pad to the longer length
	maxLen := srcCount
	if tgtCount > maxLen {
		maxLen = tgtCount
	}

	result := make([]parser.Block, maxLen)
	for i := 0; i < maxLen; i++ {
		if i < tgtCount {
			result[i] = transBlocks[i]
		}
	}

	return result, nil
}
