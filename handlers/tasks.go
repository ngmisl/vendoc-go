package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"vendoc/services"
	"vendoc/tasks"
)

func ExecuteTask(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("session")
	if sessionID == "" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	taskTypeStr := r.FormValue("task")
	if !tasks.IsValidTaskType(taskTypeStr) {
		http.Error(w, "Invalid task type", http.StatusBadRequest)
		return
	}

	taskType := tasks.TaskType(taskTypeStr)

	// Get session and document
	session, err := services.GetSession(sessionID)
	if err != nil {
		http.Error(w, "Session not found or expired", http.StatusNotFound)
		return
	}

	// Get task prompt
	prompt := tasks.GetTaskPrompt(taskType)

	// Query Venice
	venice := services.NewVeniceClient()
	response, err := venice.Query(prompt, session.DocumentContent)
	if err != nil {
		http.Error(w, "Analysis failed", http.StatusInternalServerError)
		return
	}

	// Return formatted response as HTML fragment
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<div class="message ai-message">
			<div class="message-header">
				<span class="task-badge">%s</span>
				<span class="timestamp">%s</span>
			</div>
			<div class="message-content">%s</div>
		</div>
	`, tasks.GetTaskLabel(taskType), 
	   time.Now().Format("3:04 PM"), 
	   template.HTMLEscapeString(response))
}