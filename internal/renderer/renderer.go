package renderer

import (
	"bytes"
	"html/template"
)

// BlockPair holds a source block and its translated counterpart as HTML.
type BlockPair struct {
	Source template.HTML
	Target template.HTML
}

// TemplateData holds all data for the HTML template.
type TemplateData struct {
	Title       string
	SourceLabel string
	TargetLabel string
	Pairs       []BlockPair
}

// Render produces a complete HTML document with a 2-column table layout.
func Render(data TemplateData) (string, error) {
	tmpl, err := template.New("bilingual").Parse(htmlTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
