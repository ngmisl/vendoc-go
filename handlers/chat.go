package handlers

import (
	"fmt"
	"html/template"
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

	// Return HTML fragment
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<div class="message user-message">
			<div class="message-header">
				<span class="user-badge">You</span>
				<span class="timestamp">%s</span>
			</div>
			<div class="message-content">%s</div>
		</div>
		<div class="message ai-message">
			<div class="message-header">
				<span class="ai-badge">AI Assistant</span>
				<span class="timestamp">%s</span>
			</div>
			<div class="message-content">%s</div>
		</div>
	`, time.Now().Format("3:04 PM"), 
	   template.HTMLEscapeString(message),
	   time.Now().Format("3:04 PM"),
	   template.HTMLEscapeString(response))
}