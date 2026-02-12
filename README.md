# Bilingual 2-Column PDF Generator

Converts a markdown document into a side-by-side bilingual PDF with the source language in the left column and its translation in the right column. Supports any language pair available through Google Translate. Defaults to French→Spanish. Corresponding paragraphs are vertically aligned.

## Quick start

```bash
# French → Spanish (default)
bilingual_pdf my_doc.md

# Spanish → French
bilingual_pdf mi_doc.md \
    --source es --target fr

# English → German
bilingual_pdf my_doc.md \
    --source en --target de
```

## Installation

1. download the zip file appropriate for your platform from the [Releases](https://github.com/rudifa/bilingual-cd/releases/latest) page
2. unzip and move the `bilingual_pdf` executable into a directory which is on your system PATH
3. if needed, handle the security settings for the executable `bilingual_pdf`

On a Mac computer, you can allow running the app by removing the quarantine attribute:

```bash
xattr -d com.apple.quarantine /path/to/bilingual_pdf
```

## Usage

```bash
# French → Spanish (default)
bilingual_pdf document.md

# Any language pair
bilingual_pdf document.md \
    --source en --target de

# Use a pre-translated markdown file
bilingual_pdf document.md \
    --translation document_es.md

# Specify output filename
bilingual_pdf document.md \
    -o bilingual.pdf

# Choose font size:
# small, medium (default), or large
bilingual_pdf document.md \
    --font-size small

# Also save the intermediate HTML
# (useful for debugging)
bilingual_pdf document.md --html

# Get full help
bilingual_pdf --help

# Also save the translation markdown
bilingual_pdf document.md \
    --save-translation

# List of supported language codes
# (for --source and --target)
bilingual_pdf --list-languages

```

**Default output filename:** `<stem>.<source>.<target>.pdf` (or `.html` with `--html`). If the input already ends with `.<source>.md`, the source suffix is not repeated (e.g. `doc.fr.md` → `doc.fr.es.pdf`, not `doc.fr.fr.es.pdf`).

## Input format

The input markdown should contain simple text, optionally formatted with headings, paragraphs, lists, code blocks, blockquotes, horizontal rules and web links. For example:

```markdown
# Main Title

## Section

A paragraph of text. Multiple lines
in the source are joined into a single
paragraph.

Another paragraph, separated by a blank line.

A web link:
[OpenAI](https://www.openai.com)
```

The app does not support more complex markdown features, notably tables and images.

## Using a pre-translated file

If you prefer hand-edited translations over machine translation, provide a pre-translated markdown file with the **same structure** (same number and order of headings and paragraphs) as the source:

```bash
bilingual_pdf source_fr.md \
    --translation source_es.md
```

The app warns if the block counts don't match and pads the shorter side with empty cells.

## For developers only

### How it works

1. **Parse** the input markdown into structural blocks (headings and paragraphs)
2. **Translate** each block to the target language (automatically via Google Translate, or using a pre-translated file you supply)
3. **Render** a 2-column HTML table where each row pairs a source block with its translated counterpart
4. **Convert** the HTML to an A4 PDF

### Build, test and deploy

Use the go build, test and install commands or use the `Makefile` targets.

Use the resulting `bilingual_pdf` CLI tool to run the app on the sample markdown files in `testdata/`. The generated PDFs and intermediate HTML files are saved in the same directory. You can inspect these to understand how the app works and to debug any issues.

Run `.scripts/smoketest.sh` to verify that the app runs successfully with valid arguments and that it fails with invalid arguments.

```bash
# quick tests without network access
./.scripts/smoketest.sh

# include tests using the translation API
./.scripts/smoketest.sh --full

# keep generated files for inspection
./.scripts/smoketest.sh --full --keep

# remove generated files
./.scripts/smoketest.sh --clean
```
