// Package ai provides Vercel AI Gateway integration
// Supports streaming chat, multimodal content, tool calls, and file attachments
package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gryt-backend/internal/config"
)

// Client represents Vercel AI Gateway client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new AI Gateway client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		baseURL: "https://ai-gateway.vercel.sh",
		apiKey:  cfg.AI.APIKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Message represents a chat message
type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

// TextContent represents simple text content
type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ImageContent represents image content
type ImageContent struct {
	Type     string   `json:"type"`
	ImageURL ImageURL `json:"image_url"`
}

// ImageURL represents image URL with detail level
type ImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

// FileContent represents file attachment
type FileContent struct {
	Type string `json:"type"`
	File File   `json:"file"`
}

// File represents file attachment data
type File struct {
	Data      string `json:"data"`
	MediaType string `json:"media_type"`
	Filename  string `json:"filename"`
}

// Tool represents a function tool definition
type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

// ToolFunction represents function definition
type ToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// ToolCall represents a tool call in response
type ToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Function ToolCallFunction `json:"function"`
}

// ToolCallFunction represents function call details
type ToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ChatRequest represents chat completion request
type ChatRequest struct {
	Model             string      `json:"model"`
	Messages          []Message   `json:"messages"`
	Stream            bool        `json:"stream,omitempty"`
	Temperature       *float64    `json:"temperature,omitempty"`
	MaxTokens         *int        `json:"max_tokens,omitempty"`
	TopP              *float64    `json:"top_p,omitempty"`
	FrequencyPenalty  *float64    `json:"frequency_penalty,omitempty"`
	PresencePenalty   *float64    `json:"presence_penalty,omitempty"`
	Stop              interface{} `json:"stop,omitempty"`
	Tools             []Tool      `json:"tools,omitempty"`
	ToolChoice        interface{} `json:"tool_choice,omitempty"`
}

// ChatResponse represents chat completion response
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   *Usage   `json:"usage,omitempty"`
}

// Choice represents a response choice
type Choice struct {
	Index        int          `json:"index"`
	Message      *Message     `json:"message,omitempty"`
	Delta        *Message     `json:"delta,omitempty"`
	FinishReason *string      `json:"finish_reason"`
	ToolCalls    []ToolCall   `json:"tool_calls,omitempty"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// StreamChunk represents streaming response chunk
type StreamChunk struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

// ErrorResponse represents API error response
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param,omitempty"`
		Code    string `json:"code,omitempty"`
	} `json:"error"`
}

// CreateChatCompletion creates a non-streaming chat completion
func (c *Client) CreateChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	req.Stream = false
	
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("API error (status %d): failed to decode error response", resp.StatusCode)
		}
		return nil, fmt.Errorf("API error: %s (type: %s, code: %s)", errResp.Error.Message, errResp.Error.Type, errResp.Error.Code)
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &chatResp, nil
}

// CreateChatCompletionStream creates a streaming chat completion
func (c *Client) CreateChatCompletionStream(ctx context.Context, req *ChatRequest) (<-chan StreamChunk, <-chan error) {
	req.Stream = true
	
	chunkChan := make(chan StreamChunk, 10)
	errorChan := make(chan error, 1)

	go func() {
		defer close(chunkChan)
		defer close(errorChan)

		reqBody, err := json.Marshal(req)
		if err != nil {
			errorChan <- fmt.Errorf("failed to marshal request: %w", err)
			return
		}

		httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/chat/completions", bytes.NewReader(reqBody))
		if err != nil {
			errorChan <- fmt.Errorf("failed to create request: %w", err)
			return
		}

		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Accept", "text/event-stream")

		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			errorChan <- fmt.Errorf("failed to send request: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			var errResp ErrorResponse
			if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
				errorChan <- fmt.Errorf("API error (status %d): failed to decode error response", resp.StatusCode)
				return
			}
			errorChan <- fmt.Errorf("API error: %s (type: %s, code: %s)", errResp.Error.Message, errResp.Error.Type, errResp.Error.Code)
			return
		}

		// Parse Server-Sent Events
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				errorChan <- ctx.Err()
				return
			default:
			}

			line := scanner.Text()
			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				
				if data == "[DONE]" {
					return
				}

				var chunk StreamChunk
				if err := json.Unmarshal([]byte(data), &chunk); err != nil {
					// Skip malformed chunks
					continue
				}

				chunkChan <- chunk
			}
		}

		if err := scanner.Err(); err != nil {
			errorChan <- fmt.Errorf("failed to read stream: %w", err)
		}
	}()

	return chunkChan, errorChan
}

// GetModels retrieves available models
func (c *Client) GetModels(ctx context.Context) ([]string, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/v1/models", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d)", resp.StatusCode)
	}

	var modelsResp struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	models := make([]string, len(modelsResp.Data))
	for i, model := range modelsResp.Data {
		models[i] = model.ID
	}

	return models, nil
}

// Helper functions for creating different types of messages

// NewTextMessage creates a simple text message
func NewTextMessage(role, content string) Message {
	return Message{
		Role:    role,
		Content: content,
	}
}

// NewMultiModalMessage creates a multimodal message with text and images
func NewMultiModalMessage(role string, contents []interface{}) Message {
	return Message{
		Role:    role,
		Content: contents,
	}
}

// NewTextContent creates text content for multimodal messages
func NewTextContent(text string) TextContent {
	return TextContent{
		Type: "text",
		Text: text,
	}
}

// NewImageContent creates image content for multimodal messages
func NewImageContent(imageData, mimeType string, detail string) ImageContent {
	if detail == "" {
		detail = "auto"
	}
	return ImageContent{
		Type: "image_url",
		ImageURL: ImageURL{
			URL:    fmt.Sprintf("data:%s;base64,%s", mimeType, imageData),
			Detail: detail,
		},
	}
}

// NewFileContent creates file content for multimodal messages
func NewFileContent(fileData, mediaType, filename string) FileContent {
	return FileContent{
		Type: "file",
		File: File{
			Data:      fileData,
			MediaType: mediaType,
			Filename:  filename,
		},
	}
}

// NewTool creates a tool definition
func NewTool(name, description string, parameters map[string]interface{}) Tool {
	return Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        name,
			Description: description,
			Parameters:  parameters,
		},
	}
}