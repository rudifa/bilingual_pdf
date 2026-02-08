package naming

import (
	"path/filepath"
	"strings"
)

// stem returns the filename without the .md extension.
func stem(inputPath string) string {
	base := filepath.Base(inputPath)
	return strings.TrimSuffix(base, ".md")
}

// OutputName computes the PDF output filename.
// If explicitOutput is set, it is returned as-is.
// Otherwise: <stem>.<source>.<target>.pdf, with smart dedup
// (if stem already ends with ".<source>", it is not repeated).
func OutputName(inputPath, sourceLang, targetLang, explicitOutput string) string {
	if explicitOutput != "" {
		return explicitOutput
	}
	s := stem(inputPath)
	dir := filepath.Dir(inputPath)

	sourceSuffix := "." + sourceLang
	if strings.HasSuffix(s, sourceSuffix) {
		return filepath.Join(dir, s+"."+targetLang+".pdf")
	}
	return filepath.Join(dir, s+"."+sourceLang+"."+targetLang+".pdf")
}

// HTMLOutputName computes the HTML output filename (same logic, .html extension).
func HTMLOutputName(inputPath, sourceLang, targetLang, explicitOutput string) string {
	pdf := OutputName(inputPath, sourceLang, targetLang, explicitOutput)
	return strings.TrimSuffix(pdf, ".pdf") + ".html"
}

// TranslationOutputName computes the translated markdown output filename.
// Format: <stem>.<target>.md, with dedup if stem already ends with .<source>.
func TranslationOutputName(inputPath, sourceLang, targetLang string) string {
	s := stem(inputPath)
	dir := filepath.Dir(inputPath)

	sourceSuffix := "." + sourceLang
	if strings.HasSuffix(s, sourceSuffix) {
		// stem is e.g. "doc.fr", so translation is "doc.es.md"
		base := strings.TrimSuffix(s, sourceSuffix)
		return filepath.Join(dir, base+"."+targetLang+".md")
	}
	return filepath.Join(dir, s+"."+targetLang+".md")
}
