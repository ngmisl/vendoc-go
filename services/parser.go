package services

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/nguyenthenguyen/docx"
)

const MaxFileSize = 10 * 1024 * 1024 // 10MB

type DocumentParser struct{}

func NewDocumentParser() *DocumentParser {
	return &DocumentParser{}
}

func (p *DocumentParser) ParseDocument(filename string, content []byte) (string, error) {
	if len(content) > MaxFileSize {
		return "", fmt.Errorf("file size exceeds maximum limit of %d bytes", MaxFileSize)
	}

	ext := strings.ToLower(filepath.Ext(filename))
	
	switch ext {
	case ".pdf":
		return p.parsePDF(content)
	case ".docx":
		return p.parseDOCX(content)
	case ".txt":
		return p.parseTXT(content)
	default:
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}
}

func (p *DocumentParser) parsePDF(content []byte) (string, error) {
	reader, err := pdf.NewReader(strings.NewReader(string(content)), int64(len(content)))
	if err != nil {
		return "", fmt.Errorf("failed to create PDF reader: %w", err)
	}

	var text strings.Builder
	numPages := reader.NumPage()
	
	for i := 1; i <= numPages; i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}
		
		pageText, err := page.GetPlainText(nil)
		if err != nil {
			// Continue with other pages if one fails
			continue
		}
		
		text.WriteString(pageText)
		text.WriteString("\n")
	}

	result := text.String()
	if strings.TrimSpace(result) == "" {
		return "", fmt.Errorf("no text content found in PDF")
	}

	return result, nil
}

func (p *DocumentParser) parseDOCX(content []byte) (string, error) {
	reader := strings.NewReader(string(content))
	doc, err := docx.ReadDocxFromMemory(reader, int64(len(content)))
	if err != nil {
		return "", fmt.Errorf("failed to read DOCX: %w", err)
	}

	text := doc.Editable().GetContent()
	if strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("no text content found in DOCX")
	}

	return text, nil
}

func (p *DocumentParser) parseTXT(content []byte) (string, error) {
	text := string(content)
	if strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("text file is empty")
	}

	return text, nil
}

func (p *DocumentParser) GetSupportedTypes() []string {
	return []string{".pdf", ".docx", ".txt"}
}

func (p *DocumentParser) IsSupported(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	supported := p.GetSupportedTypes()
	
	for _, supportedExt := range supported {
		if ext == supportedExt {
			return true
		}
	}
	
	return false
}