package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const veniceAPI = "https://api.venice.ai/api/v1/chat/completions"

type VeniceClient struct {
	client *http.Client
	apiKey string
}

// VeniceParameters defines model-specific options.
type VeniceParameters struct {
	IncludeVeniceSystemPrompt bool `json:"include_venice_system_prompt"`
}

// VeniceRequest is the full payload for the Venice API.
// Based on the user-provided example.
type VeniceRequest struct {
	Model            string            `json:"model"`
	Messages         []Message         `json:"messages"`
	Temperature      float64           `json:"temperature"`
	MaxTokens        int               `json:"max_tokens"`
	Stream           bool              `json:"stream"`
	VeniceParameters *VeniceParameters `json:"venice_parameters,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type VeniceResponse struct {
	Choices []Choice  `json:"choices"`
	Error   *APIError `json:"error,omitempty"`
}

type Choice struct {
	Message Message `json:"message"`
}

type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func NewVeniceClient() *VeniceClient {
	return &VeniceClient{
		client: &http.Client{},
		apiKey: os.Getenv("VENICE_API_KEY"),
	}
}

func (v *VeniceClient) Query(prompt, context string) (string, error) {
	systemPrompt := `You are a private document analyzer. You are an expert at extracting information from documents. Never store or remember document contents beyond this conversation. Provide accurate, professional analysis in a clear, structured format.`

	request := VeniceRequest{
		Model: "llama-3.3-70b", // Use default, robust model from user's list
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: fmt.Sprintf("Document Context:\n---\n%s\n---\n\nUser Request: %s", context, prompt)},
		},
		Temperature: 0.3,
		MaxTokens:   2048,
		Stream:      false, // Ensure non-streaming response to match decoder
		VeniceParameters: &VeniceParameters{
			IncludeVeniceSystemPrompt: false, // We provide our own system prompt
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", veniceAPI, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+v.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := v.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("venice API returned non-200 status: %s", resp.Status)
	}

	var veniceResp VeniceResponse
	if err := json.NewDecoder(resp.Body).Decode(&veniceResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if veniceResp.Error != nil {
		return "", fmt.Errorf("venice API error: %s", veniceResp.Error.Message)
	}

	if len(veniceResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices from Venice AI")
	}

	return veniceResp.Choices[0].Message.Content, nil
}