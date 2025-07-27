// Package handlers provides HTTP handlers for the API
package handlers

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gryt-backend/internal/ai"
	"gryt-backend/internal/database"
	"gryt-backend/internal/middleware"
)

// AIHandler handles AI-related HTTP requests
type AIHandler struct {
	aiService *ai.Service
	db        *database.DB
}

// NewAIHandler creates a new AI handler
func NewAIHandler(aiService *ai.Service, db *database.DB) *AIHandler {
	return &AIHandler{
		aiService: aiService,
		db:        db,
	}
}

// RegisterRoutes registers AI routes
func (h *AIHandler) RegisterRoutes(r *gin.RouterGroup) {
	// Chat endpoints
	chat := r.Group("/chat")
	{
		chat.POST("/sessions", h.CreateChatSession)
		chat.GET("/sessions", h.GetChatSessions)
		chat.GET("/sessions/:session_id", h.GetChatSession)
		chat.DELETE("/sessions/:session_id", h.DeleteChatSession)
		chat.POST("/sessions/:session_id/messages", h.SendChatMessage)
		chat.GET("/sessions/:session_id/messages", h.GetChatMessages)
		chat.POST("/sessions/:session_id/stream", h.StreamChatMessage)
	}

	// Search endpoints
	search := r.Group("/search")
	{
		search.POST("/", h.SearchQuery)
		search.GET("/history", h.GetSearchHistory)
		search.GET("/:search_id", h.GetSearchResult)
	}

	// Models endpoint
	models := r.Group("/models")
	{
		models.GET("/", h.GetAvailableModels)
	}
}

// Chat Session Handlers

// CreateChatSessionRequest represents request to create a chat session
type CreateChatSessionRequest struct {
	Title string `json:"title,omitempty"`
}

// CreateChatSessionResponse represents response for creating a chat session
type CreateChatSessionResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateChatSession creates a new chat session
func (h *AIHandler) CreateChatSession(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreateChatSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Generate session ID
	sessionID := fmt.Sprintf("session-%d", time.Now().UnixNano())
	title := req.Title
	if title == "" {
		title = fmt.Sprintf("Chat Session %s", time.Now().Format("2006-01-02 15:04"))
	}

	// Save session to database
	if err := h.createChatSessionInDB(c.Request.Context(), userID, sessionID, title); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat session"})
		return
	}

	c.JSON(http.StatusCreated, CreateChatSessionResponse{
		ID:        sessionID,
		Title:     title,
		CreatedAt: time.Now(),
	})
}

// GetChatSessions gets all chat sessions for the user
func (h *AIHandler) GetChatSessions(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		} // Placeholder
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get sessions from database
	sessions, err := h.getChatSessionsFromDB(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

// GetChatSession gets a specific chat session
func (h *AIHandler) GetChatSession(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		} // Placeholder
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	// Get session from database
	session, err := h.getChatSessionFromDB(c.Request.Context(), userID, sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// DeleteChatSession deletes a chat session
func (h *AIHandler) DeleteChatSession(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		} // Placeholder
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	// Delete session from database
	if err := h.deleteChatSessionFromDB(c.Request.Context(), userID, sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chat session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat session deleted successfully"})
}

// Chat Message Handlers

// SendChatMessageRequest represents request to send a chat message
type SendChatMessageRequest struct {
	Message string `json:"message" binding:"required"`
	Stream  bool   `json:"stream,omitempty"`
}

// SendChatMessage sends a chat message and gets AI response
func (h *AIHandler) SendChatMessage(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		} // Placeholder
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	// Check if user has enough chat tokens
	if !h.checkChatTokens(c.Request.Context(), userID) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Insufficient chat tokens"})
		return
	}

	var req SendChatMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Parse multipart form for file uploads
	form, err := c.MultipartForm()
	var files []*multipart.FileHeader
	if err == nil && form != nil {
		files = form.File["files"]
	}

	// Create service request
	serviceReq := &ai.ServiceChatRequest{
		SessionID: sessionID,
		Message:   req.Message,
		Files:     files,
		Stream:    req.Stream,
	}

	// Process chat message
	response, err := h.aiService.ProcessChatMessage(c.Request.Context(), userID, serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to process message: %v", err)})
		return
	}

	// Deduct chat token
	h.deductChatToken(c.Request.Context(), userID)

	c.JSON(http.StatusOK, response)
}

// StreamChatMessage handles streaming chat messages
func (h *AIHandler) StreamChatMessage(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		} // Placeholder
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	// Check if user has enough chat tokens
	if !h.checkChatTokens(c.Request.Context(), userID) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Insufficient chat tokens"})
		return
	}

	var req SendChatMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Parse multipart form for file uploads
	form, err := c.MultipartForm()
	var files []*multipart.FileHeader
	if err == nil && form != nil {
		files = form.File["files"]
	}

	// Create service request
	serviceReq := &ai.ServiceChatRequest{
		SessionID: sessionID,
		Message:   req.Message,
		Files:     files,
		Stream:    true,
	}

	// Set headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// Process streaming chat message
	responseChan, errorChan := h.aiService.ProcessChatMessageStream(c.Request.Context(), userID, serviceReq)

	// Deduct chat token
	h.deductChatToken(c.Request.Context(), userID)

	// Stream response
	for {
		select {
		case chunk, ok := <-responseChan:
			if !ok {
				// Stream finished
				c.SSEvent("done", "")
				return
			}
			c.SSEvent("message", chunk)
			c.Writer.Flush()

		case err := <-errorChan:
			if err != nil {
				c.SSEvent("error", err.Error())
				return
			}

		case <-c.Request.Context().Done():
			return
		}
	}
}

// GetChatMessages gets messages for a chat session
func (h *AIHandler) GetChatMessages(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		} // Placeholder
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	// Parse pagination parameters
	limit := 50 // default
	offset := 0 // default

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get messages from database
	messages, err := h.getChatMessagesFromDB(c.Request.Context(), userID, sessionID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

// Search Handlers

// SearchQueryRequest represents a search request
type SearchQueryRequest struct {
	Query string `json:"query" binding:"required"`
	Limit int    `json:"limit,omitempty"`
}

// SearchQuery processes a search query
func (h *AIHandler) SearchQuery(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		} // Placeholder
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user has enough search tokens
	if !h.checkSearchTokens(c.Request.Context(), userID) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Insufficient search tokens"})
		return
	}

	var req SearchQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate query
	if len(strings.TrimSpace(req.Query)) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query must be at least 3 characters long"})
		return
	}

	// Create service request
	serviceReq := &ai.SearchRequest{
		Query: req.Query,
		Limit: req.Limit,
	}

	// Process search query
	response, err := h.aiService.ProcessSearchQuery(c.Request.Context(), userID, serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to process search: %v", err)})
		return
	}

	// Deduct search token
	h.deductSearchToken(c.Request.Context(), userID)

	c.JSON(http.StatusOK, response)
}

// GetSearchHistory gets search history for the user
func (h *AIHandler) GetSearchHistory(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		} // Placeholder
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse pagination parameters
	limit := 20 // default
	offset := 0 // default

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get search history from database
	history, err := h.getSearchHistoryFromDB(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get search history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}

// GetSearchResult gets a specific search result
func (h *AIHandler) GetSearchResult(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	searchID := c.Param("search_id")
	if searchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search ID is required"})
		return
	}

	// Get search result from database
	result, err := h.getSearchResultFromDB(c.Request.Context(), userID, searchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Search result not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetAvailableModels gets available AI models
func (h *AIHandler) GetAvailableModels(c *gin.Context) {
	// This would typically call the AI service to get available models
	models := []map[string]interface{}{
		{
			"id":          "gpt-4o",
			"name":        "GPT-4o",
			"description": "Most capable GPT-4 model",
			"context":     128000,
		},
		{
			"id":          "gpt-4o-mini",
			"name":        "GPT-4o Mini",
			"description": "Faster and more affordable GPT-4 model",
			"context":     128000,
		},
		{
			"id":          "claude-3-5-sonnet-20241022",
			"name":        "Claude 3.5 Sonnet",
			"description": "Anthropic's most capable model",
			"context":     200000,
		},
	}

	c.JSON(http.StatusOK, gin.H{"models": models})
}

// Helper methods for database operations

func (h *AIHandler) createChatSessionInDB(ctx context.Context, userID, sessionID, title string) error {
	// Use database repository to create chat session
	chatRepo := database.NewChatRepository(h.db)
	session := &database.ChatSession{
		ID:        sessionID,
		UserID:    userID,
		Title:     title,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return chatRepo.CreateSession(session)
}

func (h *AIHandler) getChatSessionsFromDB(ctx context.Context, userID string) ([]interface{}, error) {
	// Use database repository to get chat sessions
	chatRepo := database.NewChatRepository(h.db)
	sessions, err := chatRepo.GetUserSessions(userID)
	if err != nil {
		return nil, err
	}
	
	// Convert to interface slice
	result := make([]interface{}, len(sessions))
	for i, session := range sessions {
		result[i] = session
	}
	return result, nil
}

func (h *AIHandler) getChatSessionFromDB(ctx context.Context, userID, sessionID string) (interface{}, error) {
	// Implementation would query specific chat session
	return map[string]interface{}{}, nil
}

func (h *AIHandler) deleteChatSessionFromDB(ctx context.Context, userID, sessionID string) error {
	// Implementation would soft delete chat session
	return nil
}

func (h *AIHandler) getChatMessagesFromDB(ctx context.Context, userID, sessionID string, limit, offset int) ([]interface{}, error) {
	// Use database repository to get chat messages
	chatRepo := database.NewChatRepository(h.db)
	
	// First verify the session belongs to the user
	sessions, err := chatRepo.GetUserSessions(userID)
	if err != nil {
		return nil, err
	}
	
	// Check if session exists and belongs to user
	sessionExists := false
	for _, session := range sessions {
		if session.ID == sessionID {
			sessionExists = true
			break
		}
	}
	
	if !sessionExists {
		return nil, fmt.Errorf("session not found or access denied")
	}
	
	// For now, return empty messages since we don't have GetMessages method
	// This would need to be implemented in ChatRepository
	return []interface{}{}, nil
}

func (h *AIHandler) getSearchHistoryFromDB(ctx context.Context, userID string, limit, offset int) ([]interface{}, error) {
	// Use database repository to get search history
	searchRepo := database.NewSearchRepository(h.db)
	queries, err := searchRepo.GetUserQueries(userID, limit)
	if err != nil {
		return nil, err
	}
	
	// Convert to interface slice
	result := make([]interface{}, len(queries))
	for i, query := range queries {
		result[i] = query
	}
	return result, nil
}

func (h *AIHandler) getSearchResultFromDB(ctx context.Context, userID, searchID string) (interface{}, error) {
	// Use database repository to get specific search result
	searchRepo := database.NewSearchRepository(h.db)
	queries, err := searchRepo.GetUserQueries(userID, 100) // Get more to find specific one
	if err != nil {
		return nil, err
	}
	
	// Find the specific search by ID
	for _, query := range queries {
		if query.ID == searchID {
			return query, nil
		}
	}
	
	return nil, fmt.Errorf("search result not found")
}

func (h *AIHandler) checkChatTokens(ctx context.Context, userID string) bool {
	// Use database repository to check user's remaining chat tokens
	userRepo := database.NewUserRepository(h.db)
	user, err := userRepo.GetByID(userID)
	if err != nil || user == nil {
		return false
	}
	return user.ChatTokens > 0
}

func (h *AIHandler) checkSearchTokens(ctx context.Context, userID string) bool {
	// Use database repository to check user's remaining search tokens
	userRepo := database.NewUserRepository(h.db)
	user, err := userRepo.GetByID(userID)
	if err != nil || user == nil {
		return false
	}
	return user.SearchTokens > 0
}

func (h *AIHandler) deductChatToken(ctx context.Context, userID string) error {
	// Use database repository to deduct one chat token from user
	userRepo := database.NewUserRepository(h.db)
	return userRepo.DecrementChatTokens(userID)
}

func (h *AIHandler) deductSearchToken(ctx context.Context, userID string) error {
	// Use database repository to deduct one search token from user
	userRepo := database.NewUserRepository(h.db)
	return userRepo.DecrementSearchTokens(userID)
}