package api

import (
	"log"
	"net/http"

	"gryt-backend/internal/database"
	"gryt-backend/internal/handlers"
	"gryt-backend/internal/middleware"
	"gryt-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine, services *services.Services, db *database.DB) {
	// Health check
	r.GET("/health", healthCheck)

	// API group
	api := r.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/validate-key", validateAccessKey(services.Auth))
			auth.POST("/login", login(services.Auth))
			auth.POST("/refresh", refreshToken(services.Auth))
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.Auth(services.Auth))
		{
			// Initialize AI handler
			aiHandler := handlers.NewAIHandler(services.AI, db)
			
			// AI routes with rate limiting
			ai := protected.Group("/ai")
			ai.Use(middleware.AIRateLimit())
			{
				aiHandler.RegisterRoutes(ai)
			}

			// Legacy chat routes (keeping for backward compatibility)
			chat := protected.Group("/chat")
			{
				chat.POST("/sessions", createChatSession(services.Chat))
				chat.GET("/sessions", getUserSessions(services.Chat))
				chat.POST("/sessions/:sessionId/messages", sendMessage(services.Chat))
				chat.GET("/sessions/:sessionId/messages", getSessionMessages(services.Chat))
			}

			// Search routes
			search := protected.Group("/search")
			{
				search.POST("/", performSearch(services.Search))
				search.GET("/history", getSearchHistory(services.Search))
			}

			// User routes
			user := protected.Group("/user")
			{
				user.GET("/profile", getUserProfile(services.Auth))
				user.GET("/tokens", getUserTokens(services.Auth))
			}
		}
	}
}

// Health check endpoint
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Server is running",
		"version": "1.0.0",
	})
}

// Auth handlers
func validateAccessKey(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[DEBUG] validateAccessKey handler called")
		var req struct {
			AccessKey string `json:"access_key" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("[DEBUG] validateAccessKey: JSON bind error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format",
				"message": "access_key is required",
			})
			return
		}

		log.Printf("[DEBUG] validateAccessKey: calling ValidateAccessKey with: %s", req.AccessKey)
		user, err := authService.ValidateAccessKey(req.AccessKey)
		if err != nil || user == nil {
			log.Printf("[DEBUG] validateAccessKey: validation failed, err: %v, user: %v", err, user)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid access key",
				"message": "Please check your access key and try again",
			})
			return
		}

		// Generate JWT token
		token, err := authService.GenerateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate token",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Access key validated successfully",
			"token":   token,
			"user": gin.H{
				"id":            user.ID,
				"name":          user.Name,
				"email":         user.Email,
				"chat_tokens":   user.ChatTokens,
				"search_tokens": user.SearchTokens,
			},
		})
	}
}

func login(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format",
			})
			return
		}

		// Authenticate user with email/password
		user, err := authService.AuthenticateUser(req.Email, req.Password)
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
				"message": "Please check your email and password",
			})
			return
		}

		// Generate JWT token
		token, err := authService.GenerateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate token",
			})
			return
		}

		// Generate refresh token
		refreshToken, err := authService.GenerateRefreshToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate refresh token",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   token,
			"refresh_token": refreshToken,
			"user": gin.H{
				"id":            user.ID,
				"name":          user.Name,
				"email":         user.Email,
				"chat_tokens":   user.ChatTokens,
				"search_tokens": user.SearchTokens,
			},
		})
	}
}

func refreshToken(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format",
				"message": "refresh_token is required",
			})
			return
		}

		// Validate refresh token and get user ID
		userID, err := authService.ValidateRefreshToken(req.RefreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired refresh token",
			})
			return
		}

		// Generate new access token
		newToken, err := authService.GenerateToken(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate new token",
			})
			return
		}

		// Generate new refresh token
		newRefreshToken, err := authService.GenerateRefreshToken(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate new refresh token",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Token refreshed successfully",
			"token":   newToken,
			"refresh_token": newRefreshToken,
		})
	}
}

// Chat handlers
func createChatSession(chatService *services.ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		var req struct {
			Title string `json:"title" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format",
			})
			return
		}

		session, err := chatService.CreateSession(userID, req.Title)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create chat session",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Chat session created successfully",
			"session": session,
		})
	}
}

func getUserSessions(chatService *services.ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		sessions, err := chatService.GetUserSessions(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get user sessions",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"sessions": sessions,
		})
	}
}

func sendMessage(chatService *services.ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		sessionID := c.Param("sessionId")
		if sessionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Session ID is required",
			})
			return
		}

		var req struct {
			Content string `json:"content" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format",
			})
			return
		}

		response, err := chatService.SendMessage(userID, sessionID, req.Content)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Message sent successfully",
			"response": response,
		})
	}
}

func getSessionMessages(chatService *services.ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("sessionId")
		if sessionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Session ID is required",
			})
			return
		}

		messages, err := chatService.GetSessionMessages(sessionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get session messages",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"messages": messages,
		})
	}
}

// Search handlers
func performSearch(searchService *services.SearchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		var req struct {
			Query string `json:"query" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format",
			})
			return
		}

		result, err := searchService.Search(userID, req.Query)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Search completed successfully",
			"result":  result,
		})
	}
}

func getSearchHistory(searchService *services.SearchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		queries, err := searchService.GetUserQueries(userID, 50)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get search history",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"queries": queries,
		})
	}
}

// User handlers
func getUserProfile(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		// Get user profile from database
		user, err := authService.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get user profile",
			})
			return
		}

		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user": gin.H{
				"id":            user.ID,
				"name":          user.Name,
				"email":         user.Email,
				"is_active":     user.IsActive,
				"chat_tokens":   user.ChatTokens,
				"search_tokens": user.SearchTokens,
				"created_at":    user.CreatedAt,
				"updated_at":    user.UpdatedAt,
			},
		})
	}
}

func getUserTokens(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		// Get user token counts from database
		user, err := authService.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get user tokens",
			})
			return
		}

		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
			"tokens": gin.H{
				"chat_tokens":   user.ChatTokens,
				"search_tokens": user.SearchTokens,
			},
			"message": "Token counts retrieved successfully",
		})
	}
}