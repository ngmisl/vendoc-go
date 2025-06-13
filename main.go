package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"vendoc/handlers"
	"vendoc/middleware"
	"vendoc/services"

	"github.com/joho/godotenv"
)

//go:embed templates/*
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

func main() {
	godotenv.Load()

	// Validate required env vars
	if os.Getenv("VENICE_API_KEY") == "" {
		log.Fatal("VENICE_API_KEY environment variable is required")
	}

	// Initialize services
	services.InitStorage()

	// Parse templates
	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		log.Fatal("Failed to parse templates:", err)
	}
	handlers.SetTemplates(tmpl)

	// Routes using Go 1.22+ pattern matching
	mux := http.NewServeMux()

	// Application routes
	mux.HandleFunc("GET /", handlers.Home)
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	mux.HandleFunc("POST /upload", handlers.Upload)
	mux.HandleFunc("GET /analyze/{session}", handlers.Analyze)
	mux.HandleFunc("POST /chat/{session}", handlers.Chat)
	mux.HandleFunc("POST /task/{session}", handlers.ExecuteTask)
	mux.HandleFunc("DELETE /session/{session}", handlers.DeleteSession)
	mux.HandleFunc("POST /session/{session}/delete", handlers.DeleteSession)

	// Static files (must be after other routes)
	mux.HandleFunc("GET /static/{file...}", func(w http.ResponseWriter, r *http.Request) {
		// r.PathValue("file...") is not reliable. Using TrimPrefix instead.
		fileName := strings.TrimPrefix(r.URL.Path, "/static/")
		if fileName == "" {
			log.Printf("Attempted to access static directory root")
			http.NotFound(w, r)
			return
		}
		filePath := "static/" + fileName
		log.Printf("Serving static file: %s", filePath)

		fileBytes, err := staticFS.ReadFile(filePath)
		if err != nil {
			log.Printf("ERROR: Static file not found: %s", filePath)
			http.NotFound(w, r)
			return
		}

		var contentType string
		if strings.HasSuffix(filePath, ".css") {
			contentType = "text/css; charset=utf-8"
		} else if strings.HasSuffix(filePath, ".js") {
			contentType = "application/javascript; charset=utf-8"
		} else {
			contentType = http.DetectContentType(fileBytes)
		}
		w.Header().Set("Content-Type", contentType)
		log.Printf("Serving %s with Content-Type: %s", filePath, contentType)

		w.Write(fileBytes)
	})

	// Add security and rate limiting middleware
	handler := middleware.RateLimit(securityMiddleware(mux))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func securityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Security headers  
		csp := "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; font-src 'self';"
		w.Header().Set("Content-Security-Policy", csp)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		log.Printf("Request: %s %s, Setting CSP: %s", r.Method, r.URL.Path, csp)

		next.ServeHTTP(w, r)
	})
}