package handlers

import (
	"net/http"

	"vendoc/services"
)

func DeleteSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("session")
	if sessionID == "" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}

	// Delete the session
	if err := services.DeleteSession(sessionID); err != nil {
		http.Error(w, "Failed to delete session", http.StatusInternalServerError)
		return
	}

	// Redirect to home page
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}