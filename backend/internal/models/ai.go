package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// ChatSession represents a chat session
type ChatSession struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ChatMessage represents a chat message
type ChatMessage struct {
	ID        string          `json:"id" db:"id"`
	SessionID string          `json:"session_id" db:"session_id"`
	UserID    string          `json:"user_id" db:"user_id"`
	Role      string          `json:"role" db:"role"` // user, assistant, system
	Content   MessageContent  `json:"content" db:"content"`
	Metadata  MessageMetadata `json:"metadata" db:"metadata"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}

// MessageContent represents the content of a message (can be text, multimodal, etc.)
type MessageContent struct {
	Type        string                 `json:"type"` // text, multimodal
	Text        string                 `json:"text,omitempty"`
	Parts       []MessageContentPart   `json:"parts,omitempty"`
	Attachments []MessageAttachment    `json:"attachments,omitempty"`
	ToolCalls   []ToolCall            `json:"tool_calls,omitempty"`
	ToolResults []ToolResult          `json:"tool_results,omitempty"`
}

// MessageContentPart represents a part of multimodal content
type MessageContentPart struct {
	Type     string `json:"type"` // text, image_url, file
	Text     string `json:"text,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
	FileURL  string `json:"file_url,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
}

// MessageAttachment represents file attachments
type MessageAttachment struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
}

// ToolCall represents a function call made by the AI
type ToolCall struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"` // function
	Function ToolCallFunction       `json:"function"`
}

// ToolCallFunction represents the function details
type ToolCallFunction struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// ToolResult represents the result of a tool execution
type ToolResult struct {
	ToolCallID string                 `json:"tool_call_id"`
	Result     map[string]interface{} `json:"result"`
	Error      string                 `json:"error,omitempty"`
}

// MessageMetadata represents additional metadata for messages
type MessageMetadata struct {
	Model       string  `json:"model,omitempty"`
	TokensUsed  int     `json:"tokens_used,omitempty"`
	Latency     int64   `json:"latency,omitempty"` // in milliseconds
	Temperature float64 `json:"temperature,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
}

// SearchQuery represents a search query
type SearchQuery struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Query     string    `json:"query" db:"query"`
	Results   SearchResults `json:"results" db:"results"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// SearchResults represents search results
type SearchResults struct {
	Results []SearchResult `json:"results"`
	Total   int           `json:"total"`
}

// SearchResult represents a single search result
type SearchResult struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Score       float64 `json:"score"`
}

// Implement driver.Valuer and sql.Scanner for JSON fields

// Value implements driver.Valuer for MessageContent
func (mc MessageContent) Value() (driver.Value, error) {
	return json.Marshal(mc)
}

// Scan implements sql.Scanner for MessageContent
func (mc *MessageContent) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, mc)
	case string:
		return json.Unmarshal([]byte(v), mc)
	default:
		return fmt.Errorf("cannot scan %T into MessageContent", value)
	}
}

// Value implements driver.Valuer for MessageMetadata
func (mm MessageMetadata) Value() (driver.Value, error) {
	return json.Marshal(mm)
}

// Scan implements sql.Scanner for MessageMetadata
func (mm *MessageMetadata) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, mm)
	case string:
		return json.Unmarshal([]byte(v), mm)
	default:
		return fmt.Errorf("cannot scan %T into MessageMetadata", value)
	}
}

// Value implements driver.Valuer for SearchResults
func (sr SearchResults) Value() (driver.Value, error) {
	return json.Marshal(sr)
}

// Scan implements sql.Scanner for SearchResults
func (sr *SearchResults) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, sr)
	case string:
		return json.Unmarshal([]byte(v), sr)
	default:
		return fmt.Errorf("cannot scan %T into SearchResults", value)
	}
}