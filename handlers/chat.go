package handlers

import (
	"fmt"
	"net/http"
	"time"

	"vendoc/services"
)

func Chat(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("session")
	if sessionID == "" {
		handleError(w, r, fmt.Errorf("missing session ID"), "Session ID required", http.StatusBadRequest)
		return
	}

	// Get session
	session, err := services.GetSession(sessionID)
	if err != nil {
		handleError(w, r, err, "Session not found or expired. Please upload a new document.", http.StatusNotFound)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		handleError(w, r, err, "Invalid form data", http.StatusBadRequest)
		return
	}

	message := r.FormValue("message")
	if message == "" {
		handleError(w, r, fmt.Errorf("empty message"), "Please enter a question", http.StatusBadRequest)
		return
	}

	// Validate message length
	if len(message) > 1000 {
		handleError(w, r, fmt.Errorf("message too long"), "Question is too long (max 1000 characters)", http.StatusBadRequest)
		return
	}

	// Query Venice AI
	venice := services.NewVeniceClient()
	response, err := venice.Query(message, session.DocumentContent)
	if err != nil {
		handleError(w, r, err, "AI analysis failed. Please try again or rephrase your question.", http.StatusInternalServerError)
		return
	}

	// Calculate remaining time for display
	remaining := time.Until(session.ExpiresAt).Minutes()
	expiresIn := fmt.Sprintf("%.0f minutes", remaining)

	// Prepare data for template
	data := struct {
		SessionID    string
		Filename     string
		Title        string
		ExpiresIn    string
		UserMessage  string
		ChatResponse string
		Timestamp    string
	}{
		SessionID:    session.ID,
		Filename:     session.Filename,
		Title:        "Analyze Document: " + session.Filename,
		ExpiresIn:    expiresIn,
		UserMessage:  message,
		ChatResponse: response,
		Timestamp:    time.Now().Format("3:04 PM"),
	}

	// Return full page with chat result
	if err := templates.ExecuteTemplate(w, "task_result.html", data); err != nil {
		handleError(w, r, err, "Failed to render page", http.StatusInternalServerError)
		return
	}
}