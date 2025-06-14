[![CodeQL](https://github.com/ngmisl/vendoc-go/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/ngmisl/vendoc-go/actions/workflows/github-code-scanning/codeql)

# Vendoc Go: Private Document Analyzer

## Overview

Vendoc Go is a secure, private document analysis tool. Users can upload documents (PDF, DOCX, TXT) which are then processed in-memory and analyzed using the powerful Venice AI API. The application provides an interactive interface to get summaries, extract key points, identify risks, list action items, and chat with an AI about the document's content. All uploaded documents are automatically deleted after 30 minutes to ensure privacy.

## Features

*   **Secure Document Upload**: Supports PDF, DOCX, and TXT file formats.
*   **In-Memory Processing**: Documents are processed entirely in memory and are not stored persistently on the server.
*   **Automatic Deletion**: Sessions and document content are automatically deleted after 30 minutes of inactivity.
*   **Venice AI Integration**: Leverages the Venice AI API for advanced document understanding.
    *   **Summarization**: Get concise summaries of your documents.
    *   **Key Point Extraction**: Identify the most important points.
    *   **Risk Analysis**: Uncover potential risks and concerns.
    *   **Action Item Identification**: Extract actionable tasks and next steps.
*   **Interactive Chat**: Engage in a conversation with an AI assistant about the uploaded document.
*   **Modern Frontend**: Clean, responsive, and mobile-friendly user interface built with HTMX for dynamic interactions without full page reloads.
*   **Robust Security**: Implements Content Security Policy (CSP) and other security headers to protect against common web vulnerabilities.

## Technologies Used

*   **Backend**: Go (Golang)
*   **Frontend**: HTML, CSS, HTMX
*   **AI Services**: Venice AI API
*   **Document Parsing**: `go-textract` for PDF, DOCX; standard library for TXT.

## Prerequisites

*   Go 1.22 or later
*   A valid Venice AI API Key

## Getting Started

### 1. Environment Variables

Before running the application, you need to set up your environment variables. Copy the `.env.example` file to a new file named `.env`:

```bash
cp .env.example .env
```

Then, edit the `.env` file and add your Venice AI API Key:

```
VENICE_API_KEY=your_actual_venice_api_key_here
PORT=8080 # Optional, defaults to 8080
```

Replace `your_actual_venice_api_key_here` with your real API key.

### 2. Build and Run

To run the application, navigate to the project's root directory and execute:

```bash
go run main.go
```

The server will start, typically on `http://localhost:8080` (or the port specified in your `.env` file).

## Usage

1.  Open your web browser and navigate to `http://localhost:8080`.
2.  You will see the landing page. Click the "Upload Document" button or drag and drop a file onto the designated area.
3.  Supported file types are PDF, DOCX, and TXT.
4.  Once the document is uploaded and processed, you will be redirected to the analysis page.
5.  On the analysis page:
    *   Use the "Quick Tasks" buttons (Summarize, Key Points, Risks, Action Items) to get specific insights.
    *   Use the chat input at the bottom to ask specific questions about the document.
6.  The AI's responses will appear in the chat area.
7.  Your session and the document content will be automatically deleted after 30 minutes of inactivity.

## Security Considerations

*   **Content Security Policy (CSP)**: The application enforces a strict CSP to mitigate risks like XSS. Scripts and styles are generally restricted to `'self'`.
*   **Other Security Headers**: Includes `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, and `X-XSS-Protection: 1; mode=block` for enhanced security.
*   **No Persistent Storage**: Uploaded documents are processed in memory and are not saved to disk by the server.
*   **Temporary Sessions**: Document content and analysis sessions are temporary and expire after 30 minutes.