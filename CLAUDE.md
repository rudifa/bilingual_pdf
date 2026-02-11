# CLAUDE.md

## Project Overview

Bilingual-cd is a Go CLI tool that converts markdown documents into side-by-side bilingual PDFs. It uses Google Translate for automatic translation and supports any language pair.

### Build & Test Commands

```bash
make build              # Build binary (bilingual_pdf)
make test               # Run unit tests
make test-integration   # Run integration tests
make lint               # Run golangci-lint
make clean              # Remove build artifacts and generated files

# Smoke tests
./.scripts/smoketest.sh           # Quick tests (no network)
./.scripts/smoketest.sh --full    # Full tests including translation API
./.scripts/smoketest.sh --keep    # Keep generated files after tests
```

### Project Structure

```
cmd/root.go              # CLI setup (Cobra)
main.go                  # Entry point
internal/
  parser/                # Markdown parsing into blocks
  translator/            # Google Translate integration
  renderer/              # 2-column HTML table generation
  converter/             # HTML to PDF conversion (go-rod)
  naming/                # Output filename logic
  languages/             # Language code validation
testdata/                # Sample markdown files for testing
```

### Key Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/yuin/goldmark` - Markdown parsing
- `github.com/Conight/go-googletrans` - Google Translate API
- `github.com/go-rod/rod` - Headless Chrome for PDF generation

### Conventions

- Default language pair: French (fr) → Spanish (es)
- Output naming: `<stem>.<source>.<target>.pdf`
- Generated files go in `testdata/` for testing
- Version is injected via `-ldflags` from git tags

## From Report

## Workflow
When editing Python or Go files, always run existing tests after making changes. If tests don't exist yet, suggest creating them.

## Quality Checks
When making edits to formatting/rendering code (markdown, HTML, PDF), verify output by running the tool end-to-end on sample input before reporting completion.

## Code Style
When removing CLI arguments or refactoring function signatures, carefully check surrounding whitespace and indentation after the edit. Run a linter or formatter if available.

## Languages & Conventions
This project primarily uses Python and Go. For Python: use argparse for CLI, follow PEP 8. For Go: use standard library conventions and run `go vet`/`go test` after changes.

## Project-Specific Notes
When working with markdown-to-PDF/HTML conversion, be aware that goldmark's default security strips raw HTML blocks. Always test raw HTML passthrough when modifying the rendering pipeline.
