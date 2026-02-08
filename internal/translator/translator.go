package translator

import (
	"fmt"
	"io"
	"time"

	googletrans "github.com/Conight/go-googletrans"
)

// Translator translates a slice of text blocks from source to target language.
type Translator interface {
	Translate(blocks []string, source, target string) ([]string, error)
}

// GoogleTranslator implements Translator using the free Google Translate API.
type GoogleTranslator struct {
	Delay    time.Duration // delay between API calls for rate limiting
	Progress io.Writer     // if non-nil, progress is printed here
}

// NewGoogleTranslator creates a GoogleTranslator with sensible defaults.
func NewGoogleTranslator(progress io.Writer) *GoogleTranslator {
	return &GoogleTranslator{
		Delay:    100 * time.Millisecond,
		Progress: progress,
	}
}

func (g *GoogleTranslator) Translate(blocks []string, source, target string) ([]string, error) {
	t := googletrans.New()
	results := make([]string, len(blocks))

	for i, block := range blocks {
		if block == "" {
			results[i] = ""
			continue
		}

		if g.Progress != nil {
			fmt.Fprintf(g.Progress, "Translating block %d/%d...\r", i+1, len(blocks))
		}

		result, err := t.Translate(block, source, target)
		if err != nil {
			return nil, fmt.Errorf("translating block %d: %w", i, err)
		}
		results[i] = result.Text

		// Rate limiting delay between requests
		if i < len(blocks)-1 && g.Delay > 0 {
			time.Sleep(g.Delay)
		}
	}

	if g.Progress != nil {
		fmt.Fprintf(g.Progress, "Translated %d blocks.          \n", len(blocks))
	}

	return results, nil
}
