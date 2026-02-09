# Bilingual 2-Column PDF Generator

Converts a markdown document into a side-by-side bilingual PDF with the source language in the left column and its translation in the right column. Supports any language pair available through Google Translate. Defaults to French→Spanish. Corresponding paragraphs are vertically aligned.

## Quick start

```bash
# French → Spanish (default)
bilingual_pdf my_doc.md

# Spanish → French
bilingual_pdf mi_doc.md --source es --target fr

# English → German
bilingual_pdf my_doc.md --source en --target de
```

## Installation

...

## Usage

```bash
# French → Spanish (default)
bilingual_pdf document.md

# Any language pair
bilingual_pdf document.md --source en --target de

# Use a pre-translated markdown file
bilingual_pdf document.md --translation document_es.md

# Specify output filename
bilingual_pdf document.md -o bilingual.pdf

# Also save the intermediate HTML (useful for debugging)
bilingual_pdf document.md --html

# Get full help
bilingual_pdf --help

# Also save the translation markdown
bilingual_pdf document.md --save-translation

# List of principal supported language codes (for --source and --target)
bilingual_pdf --list-languages

```

**Default output filename:** `<stem>.<source>.<target>.pdf` (or `.html` with `--html-only`). If the input already ends with `.<source>.md`, the source suffix is not repeated (e.g. `doc.fr.md` → `doc.fr.es.pdf`, not `doc.fr.fr.es.pdf`).

## Input format

The input markdown should contain simple text, optionally formatted with headings, paragraphs, lists, code blocks, blockquotes, horizontal rules and web links. For example:

```markdown
# Main Title

## Section

A paragraph of text. Multiple lines in the source
are joined into a single paragraph.

Another paragraph, separated by a blank line.

A web link example: [OpenAI](https://www.openai.com)
```

The app does not support more complex markdown features, notably tables and images.

## Using a pre-translated file

If you prefer hand-edited translations over machine translation, provide a pre-translated markdown file with the **same structure** (same number and order of headings and paragraphs) as the source:

```bash
bilingual_pdf source_fr.md --translation source_es.md
```

The app warns if the block counts don't match and pads the shorter side with empty cells.

## How it works

1. **Parse** the input markdown into structural blocks (headings and paragraphs)
2. **Translate** each block to the target language (automatically via Google Translate, or using a pre-translated file you supply)
3. **Render** a 2-column HTML table where each row pairs a source block with its translated counterpart
4. **Convert** the HTML to an A4 PDF
