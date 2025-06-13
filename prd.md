# Private Document Analyzer - Product Requirements Document

## Product Overview

### Vision
A mobile-first, privacy-focused document analysis tool that enables professionals to leverage AI for confidential documents without data retention risks. Built for speed and simplicity using Go and HTMX.

### Target Users
- Healthcare professionals analyzing patient records
- Lawyers reviewing confidential contracts
- Financial advisors processing sensitive documents
- Small business owners handling proprietary information

### Core Value Proposition
"Analyze confidential documents with AI in 30 seconds, with guaranteed zero data retention"

## Technical Architecture

### Stack Selection Rationale
- **Go 1.23+**: Single binary deployment, minimal memory footprint, built-in concurrency
- **HTMX**: Zero JavaScript build complexity, instant page updates, 14KB total size
- **Fly.io**: Global edge deployment, generous free tier (3 shared VMs, 3GB storage)
- **Venice AI**: Privacy-focused API with no data retention

### System Architecture
```
┌─────────────────┐
│   Mobile User   │
└────────┬────────┘
         │ HTTPS
┌────────▼────────┐
│   Fly.io Edge   │
│  (Global CDN)   │
└────────┬────────┘
         │
┌────────▼────────┐
│   Go Server     │
│  - HTMX Routes  │
│  - File Handler │
│  - Venice Client│
└────────┬────────┘
         │
┌────────▼────────┐
│   Venice API    │
│ (llama-3.3-70b) │
└─────────────────┘
```

## User Stories

### Critical Path (MVP)
1. **As a lawyer**, I want to upload a contract PDF from my phone and ask questions about specific clauses
2. **As a doctor**, I want to analyze patient records and get summaries without HIPAA violations
3. **As a user**, I want clear confirmation that my data is deleted after each session
4. **As a professional**, I want to quickly execute common analysis tasks without typing prompts

### Phase 2
5. **As a team lead**, I want to share analysis results via secure link (24-hour expiry)
6. **As a power user**, I want to analyze multiple documents in one session
7. **As an enterprise user**, I want SSO and audit logs

## UI/UX Requirements

### Mobile-First Design Principles
- **Touch targets**: Minimum 44x44px
- **Single column layout**: No horizontal scrolling
- **Progressive disclosure**: Show only essential elements
- **Speed**: <3 second initial load, <500ms interactions

### Page Structure

#### 1. Landing Page (`/`)
```html
<div class="container">
  <header>
    <h1>🔒 Private Doc Analyzer</h1>
    <p>AI analysis • Zero retention • 30 seconds</p>
  </header>
  
  <div class="upload-zone" 
       hx-post="/upload" 
       hx-encoding="multipart/form-data"
       hx-target="#main">
    <input type="file" name="document" accept=".pdf,.docx,.txt">
    <label>Tap to upload confidential document</label>
  </div>
  
  <div class="trust-signals">
    <div>✓ No data stored</div>
    <div>✓ End-to-end encrypted</div>
    <div>✓ Auto-delete in 30 min</div>
  </div>
</div>
```

#### 2. Analysis Page (`/analyze/{session-id}`)
```html
<div class="analysis-container">
  <div class="document-info">
    <span>📄 contract-2024.pdf</span>
    <button hx-delete="/session/abc123">End & Delete</button>
  </div>
  
  <!-- Quick Task Buttons -->
  <div class="task-grid">
    <button hx-post="/task/abc123" 
            hx-vals='{"task": "summarize"}'
            hx-target="#chat-messages">
      📝 Summarize
    </button>
    <button hx-post="/task/abc123" 
            hx-vals='{"task": "key-points"}'
            hx-target="#chat-messages">
      🎯 Key Points
    </button>
    <button hx-post="/task/abc123" 
            hx-vals='{"task": "risks"}'
            hx-target="#chat-messages">
      ⚠️ Find Risks
    </button>
    <button hx-post="/task/abc123" 
            hx-vals='{"task": "action-items"}'
            hx-target="#chat-messages">
      ✅ Action Items
    </button>
  </div>
  
  <div id="chat-messages">
    <!-- Messages inserted here -->
  </div>
  
  <form hx-post="/chat/abc123" 
        hx-target="#chat-messages" 
        hx-swap="beforeend">
    <input type="text" name="message" 
           placeholder="Ask about your document...">
    <button type="submit">→</button>
  </form>
</div>
```

### Visual Design
- **Color palette**: 
  - Primary: `#1a1a2e` (Dark blue)
  - Accent: `#16213e` (Navy)
  - Success: `#0f3460` (Deep blue)
  - Background: `#f5f5f5`
- **Typography**: System fonts for fastest loading
- **Icons**: Unicode symbols only (no icon fonts)

## Implementation Plan

### Project Structure
```
privatedoc/
├── main.go
├── handlers/
│   ├── upload.go
│   ├── chat.go
│   ├── tasks.go
│   └── session.go
├── services/
│   ├── venice.go
│   ├── parser.go
│   └── storage.go
├── tasks/
│   └── prompts.go
├── templates/
│   ├── layout.html
│   ├── index.html
│   ├── analyze.html
│   └── components/
├── static/
│   └── style.css
├── Dockerfile
├── fly.toml
├── .env.example
└── go.mod
```

### Core Implementation Files

#### `main.go`
```go
package main

import (
    "embed"
    "html/template"
    "log"
    "net/http"
    "os"
    
    "github.com/gorilla/mux"
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
    
    r := mux.NewRouter()
    
    // Static files
    r.PathPrefix("/static/").Handler(
        http.StripPrefix("/static/", 
        http.FileServer(http.FS(staticFS))))
    
    // Routes
    r.HandleFunc("/", handlers.Home).Methods("GET")
    r.HandleFunc("/upload", handlers.Upload).Methods("POST")
    r.HandleFunc("/analyze/{session}", handlers.Analyze).Methods("GET")
    r.HandleFunc("/chat/{session}", handlers.Chat).Methods("POST")
    r.HandleFunc("/task/{session}", handlers.ExecuteTask).Methods("POST")
    r.HandleFunc("/session/{session}", handlers.DeleteSession).Methods("DELETE")
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}
```

#### `.env.example`
```env
# Venice AI Configuration
VENICE_API_KEY=your_venice_api_key_here

# Server Configuration
PORT=8080

# Session Configuration
SESSION_TIMEOUT_MINUTES=30
MAX_FILE_SIZE_MB=10
```

#### `services/venice.go`
```go
package services

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
)

const veniceAPI = "https://api.venice.ai/v1/chat/completions"

type VeniceClient struct {
    client *http.Client
    apiKey string
}

func NewVeniceClient() *VeniceClient {
    return &VeniceClient{
        client: &http.Client{},
        apiKey: os.Getenv("VENICE_API_KEY"),
    }
}

func (v *VeniceClient) Query(prompt, context string) (string, error) {
    systemPrompt := `You are a private document analyzer. Never store or remember document contents beyond this conversation. Provide accurate, professional analysis.`
    
    payload := map[string]interface{}{
        "model": "llama-3.3-70b",
        "messages": []map[string]string{
            {"role": "system", "content": systemPrompt},
            {"role": "user", "content": fmt.Sprintf("Document:\n%s\n\nRequest: %s", context, prompt)},
        },
        "temperature": 0.3,
        "max_tokens": 1500,
    }
    
    jsonData, _ := json.Marshal(payload)
    
    req, err := http.NewRequest("POST", veniceAPI, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", err
    }
    
    req.Header.Set("Authorization", "Bearer " + v.apiKey)
    req.Header.Set("Content-Type", "application/json")
    
    // Implementation continues...
}
```

#### `tasks/prompts.go`
```go
package tasks

type TaskType string

const (
    TaskSummarize   TaskType = "summarize"
    TaskKeyPoints   TaskType = "key-points"
    TaskRisks       TaskType = "risks"
    TaskActionItems TaskType = "action-items"
)

var TaskPrompts = map[TaskType]string{
    TaskSummarize: `Provide a comprehensive summary of this document in the following format:
    
    OVERVIEW:
    [2-3 sentence high-level summary]
    
    MAIN SECTIONS:
    [Bullet points of major sections/topics]
    
    KEY FINDINGS:
    [Most important discoveries or statements]
    
    CONCLUSION:
    [Brief wrap-up of document's purpose and outcome]`,
    
    TaskKeyPoints: `Extract and list the most important points from this document:
    
    1. [First key point with brief explanation]
    2. [Second key point with brief explanation]
    ... (continue for all major points)
    
    Focus on actionable information, critical deadlines, important numbers, and binding commitments.`,
    
    TaskRisks: `Analyze this document for potential risks, concerns, or red flags:
    
    LEGAL RISKS:
    - [Any legal vulnerabilities or unclear terms]
    
    FINANCIAL RISKS:
    - [Financial obligations or exposures]
    
    OPERATIONAL RISKS:
    - [Process or execution challenges]
    
    COMPLIANCE RISKS:
    - [Regulatory or policy concerns]
    
    RECOMMENDATIONS:
    - [Suggested mitigations for identified risks]`,
    
    TaskActionItems: `Extract all action items and next steps from this document:
    
    IMMEDIATE ACTIONS (Within 7 days):
    □ [Action item with responsible party if mentioned]
    □ [Action item with deadline if specified]
    
    SHORT-TERM ACTIONS (Within 30 days):
    □ [Action items]
    
    LONG-TERM ACTIONS (30+ days):
    □ [Action items]
    
    DEPENDENCIES:
    - [Items that require other actions to complete first]`,
}

func GetTaskPrompt(taskType TaskType) string {
    if prompt, exists := TaskPrompts[taskType]; exists {
        return prompt
    }
    return "Analyze this document and provide insights."
}
```

#### `handlers/tasks.go`
```go
package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
)

func ExecuteTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    sessionID := vars["session"]
    
    // Parse task type
    r.ParseForm()
    taskType := tasks.TaskType(r.FormValue("task"))
    
    // Get session and document
    session, err := storage.GetSession(sessionID)
    if err != nil {
        http.Error(w, "Session not found", http.StatusNotFound)
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
    `, GetTaskLabel(taskType), time.Now().Format("3:04 PM"), 
       template.HTMLEscapeString(response))
}
```

#### `static/style.css`
```css
:root {
    --primary: #1a1a2e;
    --accent: #16213e;
    --success: #0f3460;
    --bg: #f5f5f5;
    --text: #333;
}

* { box-sizing: border-box; }

body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
    margin: 0;
    background: var(--bg);
    color: var(--text);
    line-height: 1.6;
}

.container {
    max-width: 600px;
    margin: 0 auto;
    padding: 1rem;
}

/* Mobile-first styles */
.upload-zone {
    background: white;
    border: 3px dashed var(--accent);
    border-radius: 12px;
    padding: 3rem 1rem;
    text-align: center;
    cursor: pointer;
    transition: all 0.3s;
}

.upload-zone:active {
    transform: scale(0.98);
}

/* Task buttons */
.task-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 0.75rem;
    margin: 1rem 0;
}

.task-grid button {
    background: white;
    border: 2px solid var(--accent);
    border-radius: 8px;
    padding: 0.75rem;
    font-size: 0.9rem;
    cursor: pointer;
    transition: all 0.2s;
}

.task-grid button:active {
    transform: scale(0.95);
    background: var(--accent);
    color: white;
}

/* Chat messages */
.message {
    margin: 1rem 0;
    padding: 1rem;
    background: white;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.message-header {
    display: flex;
    justify-content: space-between;
    margin-bottom: 0.5rem;
    font-size: 0.85rem;
    color: #666;
}

.task-badge {
    background: var(--accent);
    color: white;
    padding: 0.2rem 0.5rem;
    border-radius: 4px;
    font-size: 0.75rem;
}

/* HTMX loading states */
.htmx-request .htmx-indicator {
    display: inline-block;
}

.htmx-settling {
    opacity: 0.5;
}

/* Chat interface */
.chat-messages {
    height: 50vh;
    overflow-y: auto;
    -webkit-overflow-scrolling: touch;
    padding: 1rem;
}

/* Form styling */
form {
    display: flex;
    gap: 0.5rem;
    margin-top: 1rem;
}

form input {
    flex: 1;
    padding: 0.75rem;
    border: 2px solid #ddd;
    border-radius: 8px;
    font-size: 1rem;
}

form button {
    background: var(--accent);
    color: white;
    border: none;
    padding: 0.75rem 1.5rem;
    border-radius: 8px;
    font-size: 1.2rem;
    cursor: pointer;
}

@media (max-width: 640px) {
    .task-grid {
        grid-template-columns: 1fr;
    }
    
    .chat-messages { 
        height: 40vh; 
    }
}
```

### Deployment Configuration

#### `fly.toml`
```toml
app = "private-doc-analyzer"
primary_region = "iad"

[build]
  builder = "paketobuildpacks/builder:base"
  buildpacks = ["gcr.io/paketo-buildpacks/go"]

[env]
  PORT = "8080"
  SESSION_TIMEOUT_MINUTES = "30"
  MAX_FILE_SIZE_MB = "10"

[[services]]
  http_checks = []
  internal_port = 8080
  protocol = "tcp"
  
  [[services.ports]]
    port = 80
    handlers = ["http"]
    
  [[services.ports]]
    port = 443
    handlers = ["tls", "http"]

[services.concurrency]
  hard_limit = 25
  soft_limit = 20
```

## Development Timeline

### Day 1: Core Functionality (8 hours)
- [ ] Project setup with Go modules
- [ ] Basic HTTP server with HTMX routes
- [ ] File upload handler with multipart parsing
- [ ] Venice API integration with env vars
- [ ] Document parsing (PDF, DOCX, TXT)
- [ ] In-memory session management
- [ ] Task system implementation

### Day 2: UI & Deployment (8 hours)
- [ ] Mobile-responsive HTML templates
- [ ] HTMX interactions for chat and tasks
- [ ] CSS styling with mobile focus
- [ ] Error handling and loading states
- [ ] Deploy to Fly.io
- [ ] Environment variable configuration

### Day 3: Polish & Launch (4 hours)
- [ ] Performance optimizations
- [ ] Security headers
- [ ] Usage analytics (privacy-preserving)
- [ ] Documentation
- [ ] Testing all task types

## Success Metrics

### Technical KPIs
- **Page Load Speed**: <3s on 3G mobile
- **Time to First Analysis**: <30s from upload
- **Task Execution Time**: <5s per task
- **Uptime**: 99.9% (Fly.io SLA)
- **Memory Usage**: <256MB per instance

### Business KPIs
- **Week 1**: 100 document analyses
- **Month 1**: 1,000 active users
- **Task Usage**: 70% of sessions use quick tasks
- **Retention**: 40% weekly active users

## Security & Compliance

### Data Handling
- No persistent storage of documents
- 30-minute automatic session expiry
- Memory-only processing
- No logging of document contents
- API keys stored in environment variables

### Infrastructure Security
- HTTPS-only with HSTS
- Content Security Policy headers
- Rate limiting (10 requests/minute)
- File size limit (10MB)
- Environment-based configuration

## Task Types and Use Cases

### Available Tasks
1. **Summarize**: Executive summary with key sections
2. **Key Points**: Bullet-point list of critical information
3. **Find Risks**: Identify potential issues and concerns
4. **Action Items**: Extract todos and next steps

### Industry-Specific Applications
- **Legal**: Contract review, clause identification, obligation tracking
- **Healthcare**: Patient summary, treatment plans, medication lists
- **Finance**: Risk assessment, compliance checks, transaction summaries
- **Business**: Meeting notes, proposal analysis, project documentation

## Go-to-Market Strategy

### Launch Channels
1. **ProductHunt**: "Privacy-first ChatPDF alternative"
2. **HackerNews**: Technical deep-dive post
3. **LinkedIn**: Target legal/healthcare professionals
4. **Reddit**: r/privacy, r/lawyers, r/medicine

### Positioning
"The only document analyzer that guarantees your confidential files are never stored or logged. Built by privacy advocates for professionals who can't compromise on security."

## Technical Advantages

### Why Go + HTMX?
1. **Single Binary**: Deploy anywhere, no dependencies
2. **Resource Efficient**: 10x less memory than Node.js
3. **No Build Step**: Change HTML, reload, done
4. **Mobile Performance**: Minimal JavaScript = fast phones
5. **Fly.io Compatible**: Perfect for edge deployment

This architecture enables launching a functional POC in 2-3 days while maintaining professional quality and preparing for scale.