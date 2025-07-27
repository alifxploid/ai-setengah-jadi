// Package ai provides AI service layer for chat and search functionality
package ai

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"gryt-backend/internal/config"
	"gryt-backend/internal/database"
	"gryt-backend/internal/models"
)

// Service provides AI-related business logic
type Service struct {
	client     *Client
	db         *database.DB
	config     *config.AIConfig
	chatRepo   *database.ChatRepository
	searchRepo *database.SearchRepository
	toolExec   *ToolExecutor
}

// NewService creates a new AI service
func NewService(client *Client, db *database.DB, config *config.AIConfig) *Service {
	chatRepo := database.NewChatRepository(db)
	searchRepo := database.NewSearchRepository(db)
	toolExec := NewToolExecutor()
	
	return &Service{
		client:     client,
		db:         db,
		config:     config,
		chatRepo:   chatRepo,
		searchRepo: searchRepo,
		toolExec:   toolExec,
	}
}

// GetDB returns the database connection
func (s *Service) GetDB() *database.DB {
	return s.db
}

// ServiceChatRequest represents a chat request from user
type ServiceChatRequest struct {
	SessionID string                `json:"session_id"`
	Message   string                `json:"message"`
	Files     []*multipart.FileHeader `json:"-"`
	Stream    bool                  `json:"stream,omitempty"`
}

// ServiceChatResponse represents a chat response
type ServiceChatResponse struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	Message   string    `json:"message"`
	Tokens    int       `json:"tokens"`
	CreatedAt time.Time `json:"created_at"`
}

// SearchRequest represents a search request
type SearchRequest struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

// SearchResponse represents a search response
type SearchResponse struct {
	ID        string                 `json:"id"`
	Query     string                 `json:"query"`
	Results   []SearchResult         `json:"results"`
	Tokens    int                    `json:"tokens"`
	CreatedAt time.Time              `json:"created_at"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// SearchResult represents a single search result
type SearchResult struct {
	Title       string                 `json:"title"`
	Content     string                 `json:"content"`
	URL         string                 `json:"url,omitempty"`
	Score       float64                `json:"score,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ProcessChatMessage processes a chat message with AI
func (s *Service) ProcessChatMessage(ctx context.Context, userID string, req *ServiceChatRequest) (*ServiceChatResponse, error) {
	// Get conversation history
	history, err := s.getChatHistory(ctx, req.SessionID, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat history: %w", err)
	}

	// Build messages array
	messages := make([]Message, 0, len(history)+1)
	
	// Add system message
	systemMsg := s.buildSystemMessage()
	messages = append(messages, systemMsg)
	
	// Add conversation history
	for _, msg := range history {
		messages = append(messages, NewTextMessage(msg.Role, msg.Content.Text))
	}

	// Process user message with files if any
	userMessage, err := s.buildUserMessage(req.Message, req.Files)
	if err != nil {
		return nil, fmt.Errorf("failed to build user message: %w", err)
	}
	messages = append(messages, userMessage)

	// Create AI request
	aiReq := &ChatRequest{
		Model:            s.config.Model,
		Messages:         messages,
		Temperature:      &s.config.Temperature,
		MaxTokens:        &s.config.MaxTokens,
		TopP:             &s.config.TopP,
		FrequencyPenalty: &s.config.FrequencyPenalty,
		PresencePenalty:  &s.config.PresencePenalty,
		Stream:           req.Stream,
		Tools:            s.getAvailableTools(),
		ToolChoice:       "auto",
	}
	
	// Add stop sequences if configured
	if len(s.config.Stop) > 0 {
		aiReq.Stop = s.config.Stop
	}

	// Call AI API
	aiResp, err := s.client.CreateChatCompletion(ctx, aiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI response: %w", err)
	}

	if len(aiResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	choice := aiResp.Choices[0]
	responseContent := ""
	
	// Handle tool calls if present
	if len(choice.ToolCalls) > 0 {
		toolResults, err := s.handleToolCalls(ctx, choice.ToolCalls)
		if err != nil {
			return nil, fmt.Errorf("failed to handle tool calls: %w", err)
		}
		responseContent = toolResults
	} else if choice.Message != nil {
		if content, ok := choice.Message.Content.(string); ok {
			responseContent = content
		}
	}

	// Calculate tokens used
	tokens := 0
	if aiResp.Usage != nil {
		tokens = aiResp.Usage.TotalTokens
	}

	// Save messages to database
	if err := s.saveChatMessages(ctx, userID, req.SessionID, req.Message, responseContent, tokens); err != nil {
		return nil, fmt.Errorf("failed to save chat messages: %w", err)
	}

	return &ServiceChatResponse{
		ID:        aiResp.ID,
		SessionID: req.SessionID,
		Message:   responseContent,
		Tokens:    tokens,
		CreatedAt: time.Now(),
	}, nil
}

// ProcessChatMessageStream processes a chat message with streaming response
func (s *Service) ProcessChatMessageStream(ctx context.Context, userID string, req *ServiceChatRequest) (<-chan string, <-chan error) {
	responseChan := make(chan string, 10)
	errorChan := make(chan error, 1)

	go func() {
		defer close(responseChan)
		defer close(errorChan)

		// Get conversation history
		history, err := s.getChatHistory(ctx, req.SessionID, 10)
		if err != nil {
			errorChan <- fmt.Errorf("failed to get chat history: %w", err)
			return
		}

		// Build messages array
		messages := make([]Message, 0, len(history)+1)
		
		// Add system message
		systemMsg := s.buildSystemMessage()
		messages = append(messages, systemMsg)
		
		// Add conversation history
		for _, msg := range history {
			messages = append(messages, NewTextMessage(msg.Role, msg.Content.Text))
		}

		// Process user message
		userMessage, err := s.buildUserMessage(req.Message, req.Files)
		if err != nil {
			errorChan <- fmt.Errorf("failed to build user message: %w", err)
			return
		}
		messages = append(messages, userMessage)

		// Create AI request
		aiReq := &ChatRequest{
			Model:       s.config.Model,
			Messages:    messages,
			Temperature: &s.config.Temperature,
			MaxTokens:   &s.config.MaxTokens,
			Stream:      true,
			Tools:       s.getAvailableTools(),
			ToolChoice:  "auto",
		}

		// Call AI API with streaming
		chunkChan, errChan := s.client.CreateChatCompletionStream(ctx, aiReq)
		
		fullResponse := ""
		totalTokens := 0

		for {
			select {
			case chunk, ok := <-chunkChan:
				if !ok {
					// Stream finished, save to database
					if err := s.saveChatMessages(ctx, userID, req.SessionID, req.Message, fullResponse, totalTokens); err != nil {
						errorChan <- fmt.Errorf("failed to save chat messages: %w", err)
					}
					return
				}

				if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
					if content, ok := chunk.Choices[0].Delta.Content.(string); ok {
						fullResponse += content
						responseChan <- content
					}
				}

			case err := <-errChan:
				if err != nil {
					errorChan <- err
					return
				}

			case <-ctx.Done():
				errorChan <- ctx.Err()
				return
			}
		}
	}()

	return responseChan, errorChan
}

// ProcessSearchQuery processes a search query with AI
func (s *Service) ProcessSearchQuery(ctx context.Context, userID string, req *SearchRequest) (*SearchResponse, error) {
	// Build search prompt
	searchPrompt := s.buildSearchPrompt(req.Query)
	
	messages := []Message{
		NewTextMessage("system", "You are a helpful search assistant. Provide comprehensive and accurate search results."),
		NewTextMessage("user", searchPrompt),
	}

	// Create AI request with search tools
	aiReq := &ChatRequest{
		Model:       s.config.Model,
		Messages:    messages,
		Temperature: &s.config.Temperature,
		MaxTokens:   &s.config.MaxTokens,
		Tools:       s.getSearchTools(),
		ToolChoice:  "auto",
	}

	// Call AI API
	aiResp, err := s.client.CreateChatCompletion(ctx, aiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI response: %w", err)
	}

	if len(aiResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	// Process search results
	results := s.processSearchResults(aiResp.Choices[0])
	
	// Calculate tokens used
	tokens := 1 // Default 1 token for search
	if aiResp.Usage != nil {
		tokens = aiResp.Usage.TotalTokens
	}

	// Save search query to database
	searchID, err := s.saveSearchQuery(ctx, userID, req.Query, results, tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to save search query: %w", err)
	}

	return &SearchResponse{
		ID:        searchID,
		Query:     req.Query,
		Results:   results,
		Tokens:    tokens,
		CreatedAt: time.Now(),
	}, nil
}

// Helper methods

func (s *Service) buildSystemMessage() Message {
	systemPrompt := s.config.SystemPrompt
	if systemPrompt == "" {
		// Fallback system prompt if not configured
		systemPrompt = "You are GRYT, an advanced AI assistant powered by Vercel AI Gateway. You are designed to be helpful, accurate, and efficient. You have access to various tools and can handle multiple types of content including text, images, and documents. Always provide clear, concise, and helpful responses while maintaining a professional yet friendly tone."
	}

	return NewTextMessage("system", systemPrompt)
}

func (s *Service) buildUserMessage(message string, files []*multipart.FileHeader) (Message, error) {
	if len(files) == 0 {
		return NewTextMessage("user", message), nil
	}

	// Build multimodal message with files
	contents := []interface{}{
		NewTextContent(message),
	}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return Message{}, fmt.Errorf("failed to open file %s: %w", fileHeader.Filename, err)
		}
		defer file.Close()

		fileData, err := io.ReadAll(file)
		if err != nil {
			return Message{}, fmt.Errorf("failed to read file %s: %w", fileHeader.Filename, err)
		}

		fileBase64 := base64.StdEncoding.EncodeToString(fileData)
		contentType := fileHeader.Header.Get("Content-Type")

		// Handle different file types
		if strings.HasPrefix(contentType, "image/") {
			contents = append(contents, NewImageContent(fileBase64, contentType, "auto"))
		} else {
			contents = append(contents, NewFileContent(fileBase64, contentType, fileHeader.Filename))
		}
	}

	return NewMultiModalMessage("user", contents), nil
}

func (s *Service) buildSearchPrompt(query string) string {
	return fmt.Sprintf(`Please search for information about: "%s"

Provide comprehensive search results including:
1. Relevant information and facts
2. Multiple perspectives if applicable
3. Recent developments or updates
4. Reliable sources when possible

Format the response as structured search results.`, query)
}

func (s *Service) getAvailableTools() []Tool {
	return []Tool{
		NewTool(
			"web_search",
			"Search the web for current information",
			map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Search query",
					},
				},
				"required": []string{"query"},
			},
		),
		NewTool(
			"calculate",
			"Perform mathematical calculations",
			map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"expression": map[string]interface{}{
						"type":        "string",
						"description": "Mathematical expression to calculate",
					},
				},
				"required": []string{"expression"},
			},
		),
		NewTool(
			"get_current_time",
			"Get the current date and time",
			map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		),
	}
}

func (s *Service) getSearchTools() []Tool {
	return []Tool{
		NewTool(
			"web_search",
			"Search the web for information",
			map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Search query",
					},
					"num_results": map[string]interface{}{
						"type":        "integer",
						"description": "Number of results to return",
						"default":     10,
					},
				},
				"required": []string{"query"},
			},
		),
	}
}

func (s *Service) handleToolCalls(ctx context.Context, toolCalls []ToolCall) (string, error) {
	results := make([]string, 0, len(toolCalls))

	for _, toolCall := range toolCalls {
		result, err := s.executeToolCall(ctx, toolCall)
		if err != nil {
			return "", fmt.Errorf("failed to execute tool %s: %w", toolCall.Function.Name, err)
		}
		results = append(results, result)
	}

	return strings.Join(results, "\n\n"), nil
}

func (s *Service) executeToolCall(ctx context.Context, toolCall ToolCall) (string, error) {
	switch toolCall.Function.Name {
	case "web_search":
		return "Search functionality would be implemented here", nil
	case "calculate":
		return "Calculation functionality would be implemented here", nil
	case "get_current_time":
		return fmt.Sprintf("Current time: %s", time.Now().Format(time.RFC3339)), nil
	default:
		return "", fmt.Errorf("unknown tool: %s", toolCall.Function.Name)
	}
}

func (s *Service) processSearchResults(choice Choice) []SearchResult {
	// This would process the AI response and extract structured search results
	// For now, return a simple result
	results := []SearchResult{
		{
			Title:   "AI Generated Result",
			Content: "Search results would be processed here",
			Score:   0.95,
		},
	}

	return results
}

// Database helper methods (these would interact with the database layer)

func (s *Service) getChatHistory(ctx context.Context, sessionID string, limit int) ([]*models.ChatMessage, error) {
	// This would fetch chat history from database
	return []*models.ChatMessage{}, nil
}

func (s *Service) saveChatMessages(ctx context.Context, userID, sessionID, userMessage, aiResponse string, tokens int) error {
	// Save user message
	userMsg := &database.ChatMessage{
		ID:        generateID(),
		SessionID: sessionID,
		Role:      "user",
		Content:   userMessage,
		Tokens:    0,
		CreatedAt: time.Now(),
	}
	if err := s.chatRepo.CreateMessage(userMsg); err != nil {
		return fmt.Errorf("failed to save user message: %w", err)
	}

	// Save AI response
	aiMsg := &database.ChatMessage{
		ID:        generateID(),
		SessionID: sessionID,
		Role:      "assistant",
		Content:   aiResponse,
		Tokens:    tokens,
		CreatedAt: time.Now(),
	}
	if err := s.chatRepo.CreateMessage(aiMsg); err != nil {
		return fmt.Errorf("failed to save AI message: %w", err)
	}

	return nil
}

func (s *Service) saveSearchQuery(ctx context.Context, userID, query string, results []SearchResult, tokens int) (string, error) {
	// Convert results to JSON
	resultsJSON, err := json.Marshal(results)
	if err != nil {
		return "", fmt.Errorf("failed to marshal results: %w", err)
	}

	searchQuery := &database.SearchQuery{
		ID:        generateID(),
		UserID:    userID,
		Query:     query,
		Results:   string(resultsJSON),
		Tokens:    tokens,
		CreatedAt: time.Now(),
	}

	if err := s.searchRepo.CreateQuery(searchQuery); err != nil {
		return "", fmt.Errorf("failed to save search query: %w", err)
	}

	return searchQuery.ID, nil
}

// generateID creates a random ID
func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}