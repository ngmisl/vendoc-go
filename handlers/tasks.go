package handlers

import (
	"fmt"
	"net/http"
	"time"

	"vendoc/services"
	"vendoc/tasks"
)

func ExecuteTask(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("session")
	if sessionID == "" {
		handleError(w, r, fmt.Errorf("missing session ID"), "Session ID required", http.StatusBadRequest)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		handleError(w, r, err, "Invalid form data", http.StatusBadRequest)
		return
	}

	taskTypeStr := r.FormValue("task")
	if !tasks.IsValidTaskType(taskTypeStr) {
		handleError(w, r, fmt.Errorf("invalid task type"), "Invalid task type", http.StatusBadRequest)
		return
	}

	taskType := tasks.TaskType(taskTypeStr)

	// Get session and document
	session, err := services.GetSession(sessionID)
	if err != nil {
		handleError(w, r, err, "Session not found or expired. Please upload a new document.", http.StatusNotFound)
		return
	}

	// Get task prompt
	prompt := tasks.GetTaskPrompt(taskType)

	// Query Venice
	venice := services.NewVeniceClient()
	response, err := venice.Query(prompt, session.DocumentContent)
	if err != nil {
		handleError(w, r, err, "AI analysis failed. Please try again.", http.StatusInternalServerError)
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
		TaskResult   string
		TaskLabel    string
		Timestamp    string
		UserMessage  string
		ChatResponse string
	}{
		SessionID:    session.ID,
		Filename:     session.Filename,
		Title:        "Analyze Document: " + session.Filename,
		ExpiresIn:    expiresIn,
		TaskResult:   response,
		TaskLabel:    tasks.GetTaskLabel(taskType),
		Timestamp:    time.Now().Format("3:04 PM"),
		UserMessage:  "",
		ChatResponse: "",
	}

	// Return full page with task result
	if err := templates.ExecuteTemplate(w, "task_result.html", data); err != nil {
		handleError(w, r, err, "Failed to render page", http.StatusInternalServerError)
		return
	}
}