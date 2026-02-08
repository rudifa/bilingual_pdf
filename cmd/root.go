package cmd

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"bilingual_pdf/internal/converter"
	"bilingual_pdf/internal/languages"
	"bilingual_pdf/internal/naming"
	"bilingual_pdf/internal/parser"
	"bilingual_pdf/internal/renderer"
	"bilingual_pdf/internal/translator"

	"github.com/spf13/cobra"
)

var (
	sourceLang      string
	targetLang      string
	translationFile string
	outputFile      string
	saveHTML        bool
	saveTranslation bool
	listLanguages   bool
)

var rootCmd = &cobra.Command{
	Use:   "bilingual_pdf [input.md]",
	Short: "Generate a bilingual 2-column PDF from a markdown file",
	Long: `Converts a markdown document into a side-by-side bilingual PDF
with the source language in the left column and its translation
in the right column. Supports any language pair available through
Google Translate. Defaults to French → Spanish.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runPipeline,
}

func init() {
	rootCmd.Flags().StringVarP(&sourceLang, "source", "s", "fr", "source language code")
	rootCmd.Flags().StringVarP(&targetLang, "target", "t", "es", "target language code")
	rootCmd.Flags().StringVar(&translationFile, "translation", "", "path to pre-translated markdown file")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output PDF filename")
	rootCmd.Flags().BoolVar(&saveHTML, "html", false, "also save the generated HTML")
	rootCmd.Flags().BoolVar(&saveTranslation, "save-translation", false, "also save the translation markdown")
	rootCmd.Flags().BoolVar(&listLanguages, "list-languages", false, "list supported language codes")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runPipeline(cmd *cobra.Command, args []string) error {
	if listLanguages {
		languages.PrintSupported(os.Stdout)
		return nil
	}

	inputFile, err := validateArgs(args)
	if err != nil {
		return err
	}
	printWarnings()

	// 1. Read and parse input
	blocks, err := readAndParse(inputFile)
	if err != nil {
		return err
	}

	// 2. Translate
	translatedBlocks, err := translateAll(blocks)
	if err != nil {
		return err
	}

	// 3. Save translation markdown if requested
	if err := maybeSaveTranslation(inputFile, blocks, translatedBlocks); err != nil {
		return err
	}

	// 4. Render HTML
	pairs := buildPairs(blocks, translatedBlocks)
	htmlContent, err := renderer.Render(renderer.TemplateData{
		Title:       fmt.Sprintf("Bilingual: %s → %s", languages.Name(sourceLang), languages.Name(targetLang)),
		SourceLabel: languages.Name(sourceLang),
		TargetLabel: languages.Name(targetLang),
		Pairs:       pairs,
	})
	if err != nil {
		return fmt.Errorf("rendering HTML: %w", err)
	}

	// 5. Save HTML if requested
	if err := maybeSaveHTML(inputFile, htmlContent); err != nil {
		return err
	}

	// 6. Convert to PDF and write
	return generatePDF(inputFile, htmlContent)
}

func validateArgs(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("input markdown file is required (use --help for usage)")
	}
	inputFile := args[0]

	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return "", fmt.Errorf("input file not found: %s", inputFile)
	}
	if err := languages.Validate(sourceLang); err != nil {
		return "", fmt.Errorf("invalid source language: %w", err)
	}
	if err := languages.Validate(targetLang); err != nil {
		return "", fmt.Errorf("invalid target language: %w", err)
	}
	return inputFile, nil
}

func printWarnings() {
	if sourceLang == targetLang {
		fmt.Fprintf(os.Stderr, "Warning: source and target languages are the same (%s)\n", sourceLang)
	}
	if translationFile != "" && saveTranslation {
		fmt.Fprintln(os.Stderr, "Warning: --save-translation is ignored when --translation is provided")
	}
}

func readAndParse(inputFile string) ([]parser.Block, error) {
	source, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("reading input file: %w", err)
	}

	blocks, err := parser.Parse(source)
	if err != nil {
		return nil, fmt.Errorf("parsing markdown: %w", err)
	}

	if len(blocks) == 0 {
		fmt.Fprintln(os.Stderr, "Warning: input file contains no parseable blocks")
	}
	return blocks, nil
}

func translateAll(blocks []parser.Block) ([]parser.Block, error) {
	if translationFile != "" {
		return translateFromFile(blocks)
	}
	return translateWithGoogle(blocks)
}

func translateFromFile(blocks []parser.Block) ([]parser.Block, error) {
	ft := translator.NewFileTranslator(translationFile, os.Stderr)
	result, err := ft.TranslateBlocks(blocks)
	if err != nil {
		return nil, fmt.Errorf("reading translation file: %w", err)
	}
	return result, nil
}

func translateWithGoogle(blocks []parser.Block) ([]parser.Block, error) {
	gt := translator.NewGoogleTranslator(os.Stderr)

	texts := make([]string, len(blocks))
	for i, b := range blocks {
		if b.Kind != parser.BlockCodeBlock {
			texts[i] = b.Text
		}
	}

	translatedTexts, err := gt.Translate(texts, sourceLang, targetLang)
	if err != nil {
		return nil, fmt.Errorf("translating: %w", err)
	}

	return buildTranslatedBlocks(blocks, translatedTexts), nil
}

func buildTranslatedBlocks(blocks []parser.Block, translatedTexts []string) []parser.Block {
	result := make([]parser.Block, len(blocks))
	for i, b := range blocks {
		if b.Kind == parser.BlockCodeBlock {
			result[i] = b
			continue
		}
		md := reconstructMarkdown(b, translatedTexts[i])
		tBlocks, err := parser.Parse([]byte(md))
		if err != nil || len(tBlocks) == 0 {
			result[i] = parser.Block{
				Kind: b.Kind,
				Text: translatedTexts[i],
				HTML: "<p>" + template.HTMLEscapeString(translatedTexts[i]) + "</p>\n",
			}
		} else {
			result[i] = tBlocks[0]
		}
	}
	return result
}

func maybeSaveTranslation(inputFile string, blocks, translatedBlocks []parser.Block) error {
	if !saveTranslation || translationFile != "" {
		return nil
	}
	transPath := naming.TranslationOutputName(inputFile, sourceLang, targetLang)
	var mdBuf strings.Builder
	for i, b := range translatedBlocks {
		if i < len(blocks) {
			mdBuf.WriteString(reconstructMarkdown(blocks[i], b.Text))
		} else {
			mdBuf.WriteString(b.Text)
		}
		mdBuf.WriteString("\n")
	}
	if err := os.WriteFile(transPath, []byte(mdBuf.String()), 0644); err != nil {
		return fmt.Errorf("saving translation: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Saved translation: %s\n", transPath)
	return nil
}

func buildPairs(blocks, translatedBlocks []parser.Block) []renderer.BlockPair {
	maxLen := len(blocks)
	if len(translatedBlocks) > maxLen {
		maxLen = len(translatedBlocks)
	}

	pairs := make([]renderer.BlockPair, maxLen)
	for i := 0; i < maxLen; i++ {
		if i < len(blocks) {
			pairs[i].Source = template.HTML(blocks[i].HTML)
		}
		if i < len(translatedBlocks) {
			pairs[i].Target = template.HTML(translatedBlocks[i].HTML)
		}
	}
	return pairs
}

func maybeSaveHTML(inputFile, htmlContent string) error {
	if !saveHTML {
		return nil
	}
	htmlPath := naming.HTMLOutputName(inputFile, sourceLang, targetLang, outputFile)
	if err := os.WriteFile(htmlPath, []byte(htmlContent), 0644); err != nil {
		return fmt.Errorf("saving HTML: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Saved HTML: %s\n", htmlPath)
	return nil
}

func generatePDF(inputFile, htmlContent string) error {
	pdfBytes, err := converter.Convert(htmlContent)
	if err != nil {
		return fmt.Errorf("converting to PDF: %w", err)
	}

	pdfPath := naming.OutputName(inputFile, sourceLang, targetLang, outputFile)
	if err := os.WriteFile(pdfPath, pdfBytes, 0644); err != nil {
		return fmt.Errorf("writing PDF: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Saved PDF: %s\n", pdfPath)
	return nil
}

// reconstructMarkdown rebuilds markdown from a source block structure and translated text.
func reconstructMarkdown(sourceBlock parser.Block, translatedText string) string {
	switch sourceBlock.Kind {
	case parser.BlockHeading:
		prefix := strings.Repeat("#", sourceBlock.Level) + " "
		return prefix + strings.TrimSpace(translatedText)
	case parser.BlockBlockquote:
		lines := strings.Split(translatedText, "\n")
		var buf strings.Builder
		for _, line := range lines {
			buf.WriteString("> ")
			buf.WriteString(line)
			buf.WriteString("\n")
		}
		return buf.String()
	case parser.BlockCodeBlock:
		return sourceBlock.Raw
	case parser.BlockThematicBreak:
		return "---"
	case parser.BlockList:
		return translatedText
	default:
		return translatedText
	}
}
