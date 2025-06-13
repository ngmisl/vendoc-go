package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"vendoc/services"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	log.Println("Upload handler started")

	// Parse multipart form with 10MB limit
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Printf("ERROR parsing multipart form: %v", err)
		handleError(w, r, err, "File too large (max 10MB) or invalid format", http.StatusBadRequest)
		return
	}
	log.Println("Multipart form parsed successfully")

	// Get the file from form
	file, header, err := r.FormFile("document")
	if err != nil {
		log.Printf("ERROR getting form file: %v", err)
		handleError(w, r, err, "Please select a file to upload", http.StatusBadRequest)
		return
	}
	defer file.Close()
	log.Printf("File received: %s, Size: %d bytes", header.Filename, header.Size)

	// Check file type
	parser := services.NewDocumentParser()
	if !parser.IsSupported(header.Filename) {
		supported := parser.GetSupportedTypes()
		log.Printf("ERROR unsupported file type: %s", header.Filename)
		handleError(w, r, fmt.Errorf("unsupported file type: %s", header.Filename),
			fmt.Sprintf("Unsupported file type. Please upload: %v", supported), http.StatusBadRequest)
		return
	}
	log.Println("File type is supported")

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		log.Printf("ERROR reading file content: %v", err)
		handleError(w, r, err, "Failed to read uploaded file", http.StatusInternalServerError)
		return
	}
	log.Printf("File content read successfully, %d bytes", len(content))

	// Validate file size
	if len(content) == 0 {
		log.Println("ERROR file is empty")
		handleError(w, r, fmt.Errorf("empty file"), "File appears to be empty", http.StatusBadRequest)
		return
	}

	// Parse document
	log.Println("Parsing document...")
	documentText, err := parser.ParseDocument(header.Filename, content)
	if err != nil {
		log.Printf("ERROR parsing document: %v", err)
		handleError(w, r, err, fmt.Sprintf("Failed to parse document: %s", err.Error()), http.StatusBadRequest)
		return
	}
	log.Printf("Document parsed successfully, extracted %d characters", len(documentText))

	// Validate parsed content
	if len(documentText) < 10 {
		log.Printf("ERROR insufficient content after parsing: %d chars", len(documentText))
		handleError(w, r, fmt.Errorf("insufficient content"), "Document contains too little text to analyze", http.StatusBadRequest)
		return
	}

	// Create session
	log.Println("Creating session...")
	session, err := services.CreateSession(header.Filename, documentText)
	if err != nil {
		log.Printf("ERROR creating session: %v", err)
		handleError(w, r, err, "Failed to create analysis session", http.StatusInternalServerError)
		return
	}
	log.Printf("Session created successfully: %s", session.ID)

	// Redirect to analysis page using standard HTTP redirect
	redirectURL := fmt.Sprintf("/analyze/%s", session.ID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	log.Printf("Redirecting to: %s", redirectURL)
}