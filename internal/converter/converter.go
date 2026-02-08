package converter

import (
	"fmt"
	"io"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/ysmood/gson"
)

// Convert takes an HTML string and produces a PDF as bytes.
// It uses headless Chrome via the Rod library.
func Convert(htmlContent string) ([]byte, error) {
	// Try to find or download a browser
	u, err := launcher.New().Headless(true).Launch()
	if err != nil {
		return nil, fmt.Errorf("launching browser: %w", err)
	}

	browser := rod.New().ControlURL(u)
	if err := browser.Connect(); err != nil {
		return nil, fmt.Errorf("connecting to browser: %w", err)
	}
	defer browser.MustClose()

	page, err := browser.Page(proto.TargetCreateTarget{URL: "about:blank"})
	if err != nil {
		return nil, fmt.Errorf("creating page: %w", err)
	}

	// Set the HTML content directly
	if err := page.SetDocumentContent(htmlContent); err != nil {
		return nil, fmt.Errorf("setting document content: %w", err)
	}

	// Wait for the page to be stable
	if err := page.WaitStable(300e6); err != nil { // 300ms
		return nil, fmt.Errorf("waiting for page stability: %w", err)
	}

	// Generate PDF with A4 dimensions (210mm x 297mm = 8.27in x 11.69in)
	reader, err := page.PDF(&proto.PagePrintToPDF{
		PaperWidth:           gson.Num(8.27),
		PaperHeight:          gson.Num(11.69),
		MarginTop:            gson.Num(0), // margins handled by CSS @page
		MarginBottom:         gson.Num(0),
		MarginLeft:           gson.Num(0),
		MarginRight:          gson.Num(0),
		PrintBackground:      true,
		PreferCSSPageSize:    true,
		GenerateDocumentOutline: true,
	})
	if err != nil {
		return nil, fmt.Errorf("generating PDF: %w", err)
	}

	pdfBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("reading PDF data: %w", err)
	}

	return pdfBytes, nil
}
