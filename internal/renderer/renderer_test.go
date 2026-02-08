package renderer

import (
	"html/template"
	"strings"
	"testing"
)

func TestRender_BasicOutput(t *testing.T) {
	data := TemplateData{
		Title:       "Test Document",
		SourceLabel: "French",
		TargetLabel: "Spanish",
		Pairs: []BlockPair{
			{
				Source: template.HTML("<h1>Bonjour</h1>"),
				Target: template.HTML("<h1>Hola</h1>"),
			},
			{
				Source: template.HTML("<p>Le monde est beau.</p>"),
				Target: template.HTML("<p>El mundo es hermoso.</p>"),
			},
		},
	}

	html, err := Render(data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Check basic structure
	checks := []string{
		"<!DOCTYPE html>",
		"<title>Test Document</title>",
		"French",
		"Spanish",
		"<h1>Bonjour</h1>",
		"<h1>Hola</h1>",
		"<p>Le monde est beau.</p>",
		"<p>El mundo es hermoso.</p>",
		"<table>",
		"</table>",
	}

	for _, check := range checks {
		if !strings.Contains(html, check) {
			t.Errorf("rendered HTML should contain %q", check)
		}
	}
}

func TestRender_EmptyPairs(t *testing.T) {
	data := TemplateData{
		Title:       "Empty",
		SourceLabel: "French",
		TargetLabel: "Spanish",
		Pairs:       []BlockPair{},
	}

	html, err := Render(data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if !strings.Contains(html, "<tbody>") {
		t.Error("should contain tbody even with no pairs")
	}
}

func TestRender_HTMLEscaping(t *testing.T) {
	// The template.HTML type should NOT escape the content
	data := TemplateData{
		Title:       "Test",
		SourceLabel: "EN",
		TargetLabel: "FR",
		Pairs: []BlockPair{
			{
				Source: template.HTML("<p>Hello <strong>world</strong></p>"),
				Target: template.HTML("<p>Bonjour <strong>monde</strong></p>"),
			},
		},
	}

	html, err := Render(data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Should contain unescaped HTML tags
	if !strings.Contains(html, "<strong>world</strong>") {
		t.Error("template.HTML content should not be escaped")
	}
}
