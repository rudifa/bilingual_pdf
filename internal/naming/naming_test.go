package naming

import "testing"

func TestOutputName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		source   string
		target   string
		explicit string
		want     string
	}{
		{
			name:   "basic",
			input:  "doc.md",
			source: "fr", target: "es",
			want: "doc.fr.es.pdf",
		},
		{
			name:   "source suffix dedup",
			input:  "doc.fr.md",
			source: "fr", target: "es",
			want: "doc.fr.es.pdf",
		},
		{
			name:     "explicit output",
			input:    "doc.md",
			source:   "fr", target: "es",
			explicit: "output.pdf",
			want:     "output.pdf",
		},
		{
			name:   "different language pair",
			input:  "doc.en.md",
			source: "en", target: "de",
			want: "doc.en.de.pdf",
		},
		{
			name:   "no dedup for different source",
			input:  "doc.fr.md",
			source: "en", target: "de",
			want: "doc.fr.en.de.pdf",
		},
		{
			name:   "with directory",
			input:  "/path/to/doc.md",
			source: "fr", target: "es",
			want: "/path/to/doc.fr.es.pdf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := OutputName(tt.input, tt.source, tt.target, tt.explicit)
			if got != tt.want {
				t.Errorf("OutputName(%q, %q, %q, %q) = %q, want %q",
					tt.input, tt.source, tt.target, tt.explicit, got, tt.want)
			}
		})
	}
}

func TestHTMLOutputName(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		source string
		target string
		want   string
	}{
		{
			name:   "basic",
			input:  "doc.md",
			source: "fr", target: "es",
			want: "doc.fr.es.html",
		},
		{
			name:   "source suffix dedup",
			input:  "doc.fr.md",
			source: "fr", target: "es",
			want: "doc.fr.es.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HTMLOutputName(tt.input, tt.source, tt.target, "")
			if got != tt.want {
				t.Errorf("HTMLOutputName(%q, %q, %q) = %q, want %q",
					tt.input, tt.source, tt.target, got, tt.want)
			}
		})
	}
}

func TestTranslationOutputName(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		source string
		target string
		want   string
	}{
		{
			name:   "basic",
			input:  "doc.md",
			source: "fr", target: "es",
			want: "doc.es.md",
		},
		{
			name:   "with source suffix",
			input:  "doc.fr.md",
			source: "fr", target: "es",
			want: "doc.es.md",
		},
		{
			name:   "different source - no dedup",
			input:  "doc.fr.md",
			source: "en", target: "de",
			want: "doc.fr.de.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TranslationOutputName(tt.input, tt.source, tt.target)
			if got != tt.want {
				t.Errorf("TranslationOutputName(%q, %q, %q) = %q, want %q",
					tt.input, tt.source, tt.target, got, tt.want)
			}
		})
	}
}
