package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"vendoc/services"
)

func Analyze(w http.ResponseWriter, r *http.Request) {
	log.Println("Analyze handler started")
	sessionID := r.PathValue("session")
	if sessionID == "" {
		log.Println("ERROR: Missing session ID in analyze request")
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}
	log.Printf("Analyzing session: %s", sessionID)

	session, err := services.GetSession(sessionID)
	if err != nil {
		log.Printf("ERROR: Could not get session '%s': %v", sessionID, err)
		http.Error(w, "Session not found or expired", http.StatusNotFound)
		return
	}
	log.Println("Session retrieved successfully")

	// Calculate remaining time for display
	remaining := time.Until(session.ExpiresAt).Minutes()
	expiresIn := fmt.Sprintf("%.0f minutes", remaining)

	data := struct {
		SessionID string
		Filename  string
		Title     string
		ExpiresIn string
	}{
		SessionID: session.ID,
		Filename:  session.Filename,
		Title:     "Analyze Document: " + session.Filename,
		ExpiresIn: expiresIn,
	}

	log.Println("Executing analyze template...")
	if err := templates.ExecuteTemplate(w, "analyze.html", data); err != nil {
		log.Printf("ERROR executing analyze template: %v", err)
		handleError(w, r, err, "Template execution failed", http.StatusInternalServerError)
		return
	}
	log.Println("Analyze template executed successfully")
}