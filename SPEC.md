# Build a golang CLI app

To build a golang CLI app that generates a bilingual PDF from a markdown file, using the specs and detalis below, in phases:

1. A plan for the app structure and the necessary libraries to use (e.g., for markdown parsing, translation, HTML generation, and PDF conversion). Consider using the Charm library for CLI argument parsing.

2. A plan to implement the app for local testing, automated and manual.

3. A plan to package the app for distribution, including handling dependencies and creating executable binaries for different platforms.

## Bilingual 2-Column PDF Generator

Converts a markdown document into a side-by-side bilingual PDF with the source language in the left column and its translation in the right column. Supports any language pair available through Google Translate. Defaults to French→Spanish. Corresponding paragraphs are vertically aligned. Provides user-friendly warnings if arguments are inaplicable or mismatched.

## Quick start

```bash
# French → Spanish (default)
bilingual_pdf my_doc.md

# Spanish → French
bilingual_pdf mi_doc.md --source es --target fr

# English → German
bilingual_pdf my_doc.md --source en --target de
```

## How it works

1. **Parse** the input markdown into structural blocks
2. **Translate** each block to the target language (automatically via Google Translate, or from a pre-translated file you supply)
3. **Render** a 2-column HTML table where each row pairs a source block with its translated counterpart
4. **Convert** the HTML to an A4 PDF

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

# Also save the generated HTML
bilingual_pdf document.md --html

# Get full help
bilingual_pdf --help

# Also save the translation markdown
bilingual_pdf document.md --save-translation

# List of principal supported language codes (for --source and --target)
bilingual_pdf --list-languages

```

**Default generated file names:** `<stem>.<source>.<target>.pdf` (or `.html`), or `<stem>.target>.md`. If the input name already ends with `.<source>.md`, the source suffix is not repeated (e.g. `doc.fr.md` → `doc.fr.es.pdf`, not `doc.fr.fr.es.pdf`).

## Input format

The input markdown file should contain standard markdown.

## Using a pre-translated file

If you prefer hand-edited translations over machine translation, provide a pre-translated markdown file with the **same structure** (same number and order of headings and paragraphs) as the source:

```bash
bilingual_pdf source_fr.md --translation source_es.md
```

The app warns if the block counts don't match and pads the shorter side with empty cells.
