package handlers

import (
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
	}{
		Title: "Private Doc Analyzer",
	}

	if err := templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}