package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type LLMService struct {
	apiKey string
	client *http.Client
}

type ChatGPTRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatGPTResponse struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int64     `json:"created"`
	Model   string    `json:"model"`
	Choices []Choice  `json:"choices"`
	Usage   Usage     `json:"usage"`
	Error   *APIError `json:"error,omitempty"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

func NewLLMService() *LLMService {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		// Try to read from a config file or use a default
		apiKey = "" // Will need to be set by user
	}

	return &LLMService{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// SetAPIKey allows setting the API key at runtime
func (s *LLMService) SetAPIKey(apiKey string) {
	s.apiKey = apiKey
}

// GenerateSummary generates a summary of the given text using ChatGPT 3.5 Turbo
func (s *LLMService) GenerateSummary(documentText string, maxLength int) (string, error) {
	if s.apiKey == "" {
		return "", errors.New("OpenAI API key not configured, please set OPENAI_API_KEY environment variable")
	}

	if documentText == "" {
		return "", errors.New("document text is empty")
	}

	// Truncate if too long (ChatGPT has token limits)
	if len(documentText) > 12000 {
		documentText = documentText[:12000] + "..."
	}

	prompt := fmt.Sprintf(`Please provide a concise summary of the following document in %d words or less. 
Focus on the main ideas, key findings, and important conclusions.

Document:
%s

Summary:`, maxLength, documentText)

	request := ChatGPTRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant that creates concise and accurate summaries of documents.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   500,
		Temperature: 0.5,
	}

	return s.makeRequest(request)
}

// GenerateTags generates relevant tags for the given document
func (s *LLMService) GenerateTags(documentText string, maxTags int) ([]string, error) {
	if s.apiKey == "" {
		return nil, errors.New("OpenAI API key not configured")
	}

	if documentText == "" {
		return nil, errors.New("document text is empty")
	}

	// Truncate if too long
	if len(documentText) > 8000 {
		documentText = documentText[:8000] + "..."
	}

	prompt := fmt.Sprintf(`Generate %d relevant keywords/tags for the following document. 
Return only the tags separated by commas, nothing else.

Document:
%s

Tags:`, maxTags, documentText)

	request := ChatGPTRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant that generates relevant keywords and tags for documents.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   100,
		Temperature: 0.3,
	}

	response, err := s.makeRequest(request)
	if err != nil {
		return nil, err
	}

	// Parse comma-separated tags
	tags := []string{}
	for _, tag := range bytes.Split([]byte(response), []byte(",")) {
		trimmed := bytes.TrimSpace(tag)
		if len(trimmed) > 0 {
			tags = append(tags, string(trimmed))
		}
	}

	return tags, nil
}

// AnswerQuestion uses ChatGPT to answer questions about a document
func (s *LLMService) AnswerQuestion(documentText string, question string) (string, error) {
	if s.apiKey == "" {
		return "", errors.New("OpenAI API key not configured")
	}

	if documentText == "" {
		return "", errors.New("document text is empty")
	}

	if question == "" {
		return "", errors.New("question is empty")
	}

	// Truncate if too long
	if len(documentText) > 10000 {
		documentText = documentText[:10000] + "..."
	}

	prompt := fmt.Sprintf(`Based on the following document, please answer this question:

Question: %s

Document:
%s

Answer:`, question, documentText)

	request := ChatGPTRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant that answers questions about documents accurately and concisely.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   300,
		Temperature: 0.7,
	}

	return s.makeRequest(request)
}

// makeRequest sends a request to the OpenAI API
func (s *LLMService) makeRequest(request ChatGPTRequest) (string, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var chatResponse ChatGPTResponse
	if err := json.Unmarshal(body, &chatResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for API errors
	if chatResponse.Error != nil {
		return "", fmt.Errorf("OpenAI API error: %s (type: %s, code: %s)",
			chatResponse.Error.Message,
			chatResponse.Error.Type,
			chatResponse.Error.Code)
	}

	// Check if we got a valid response
	if len(chatResponse.Choices) == 0 {
		return "", errors.New("no response choices returned from API")
	}

	return chatResponse.Choices[0].Message.Content, nil
}

// GenerateDescription generates a description based on document name and type
func (s *LLMService) GenerateDescription(documentName string, documentType string, fileName string) (string, error) {
	if s.apiKey == "" {
		return "", errors.New("OpenAI API key not configured")
	}

	if documentName == "" {
		return "", errors.New("document name is empty")
	}

	prompt := fmt.Sprintf(`Based on the following document information, generate a brief, professional description (2-3 sentences) that would be suitable for a document management system.

Document Name: %s
Document Type: %s
File Name: %s

Generate a description that explains what this document likely contains and its purpose. Be concise and professional.

Description:`, documentName, documentType, fileName)

	request := ChatGPTRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant that creates professional document descriptions for a research institute's document management system.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   150,
		Temperature: 0.7,
	}

	return s.makeRequest(request)
}
