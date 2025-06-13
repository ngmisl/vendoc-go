package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func handleError(w http.ResponseWriter, r *http.Request, err error, userMessage string, statusCode int) {
	log.Printf("Error in %s %s: %v", r.Method, r.URL.Path, err)

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// Return HTML fragment for HTMX
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(statusCode)
		fmt.Fprintf(w, `
			<div class="message error-message">
				<div class="message-header">
					<span class="error-badge">❌ Error</span>
				</div>
				<div class="message-content">%s</div>
			</div>
		`, template.HTMLEscapeString(userMessage))
		return
	}

	// Regular HTTP error response
	http.Error(w, userMessage, statusCode)
}

func handleSuccess(w http.ResponseWriter, r *http.Request, message string) {
	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
			<div class="message success-message">
				<div class="message-header">
					<span class="success-badge">✅ Success</span>
				</div>
				<div class="message-content">%s</div>
			</div>
		`, template.HTMLEscapeString(message))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, message)
}