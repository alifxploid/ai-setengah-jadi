package services

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"gryt-backend/internal/ai"
	"gryt-backend/internal/config"
	"gryt-backend/internal/database"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

// Services contains all service dependencies
type Services struct {
	Auth      *AuthService
	Chat      *ChatService
	Search    *SearchService
	RateLimit *RateLimitService
	AI        *ai.Service
}

// NewServices creates new services instance
func NewServices(db *database.DB, cfg *config.Config) *Services {
	userRepo := database.NewUserRepository(db)
	chatRepo := database.NewChatRepository(db)
	searchRepo := database.NewSearchRepository(db)

	// Initialize AI client and service
	aiClient := ai.NewClient(cfg)
	aiService := ai.NewService(aiClient, db, &cfg.AI)

	return &Services{
		Auth:      NewAuthService(userRepo, cfg.Auth),
		Chat:      NewChatService(chatRepo, userRepo, cfg),
		Search:    NewSearchService(searchRepo, userRepo, cfg),
		RateLimit: NewRateLimitService(cfg.Limits.RateLimitPerMinute),
		AI:        aiService,
	}
}

// AuthService handles authentication
type AuthService struct {
	userRepo *database.UserRepository
	config   config.AuthConfig
}

func NewAuthService(userRepo *database.UserRepository, cfg config.AuthConfig) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
	}
}

func (s *AuthService) ValidateAccessKey(accessKey string) (*database.User, error) {
	if accessKey == "" {
		return nil, errors.New("access key is required")
	}

	log.Printf("[DEBUG] ValidateAccessKey called with: %s", accessKey)
	user, err := s.userRepo.GetByAccessKey(accessKey)
	if err != nil {
		log.Printf("[DEBUG] ValidateAccessKey error: %v", err)
		return nil, fmt.Errorf("failed to validate access key: %w", err)
	}

	if user == nil {
		log.Printf("[DEBUG] ValidateAccessKey: user not found for access key: %s", accessKey)
		return nil, errors.New("invalid access key")
	}

	log.Printf("[DEBUG] ValidateAccessKey: user found: %+v", user)
	return user, nil
}

func (s *AuthService) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.config.TokenExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", errors.New("invalid user_id in token")
		}
		return userID, nil
	}

	return "", errors.New("invalid token")
}

func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *AuthService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AuthService) AuthenticateUser(email, password string) (*database.User, error) {
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check if password hash exists
	if !user.PasswordHash.Valid || user.PasswordHash.String == "" {
		return nil, errors.New("password authentication not available for this user")
	}

	if !s.CheckPassword(password, user.PasswordHash.String) {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func (s *AuthService) GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.config.RefreshExpiry).Unix(),
		"iat":     time.Now().Unix(),
		"type":    "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return tokenString, nil
}

func (s *AuthService) ValidateRefreshToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check if it's a refresh token
		if tokenType, exists := claims["type"]; !exists || tokenType != "refresh" {
			return "", errors.New("not a refresh token")
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", errors.New("invalid user_id in refresh token")
		}
		return userID, nil
	}

	return "", errors.New("invalid refresh token")
}

func (s *AuthService) GetUserByID(userID string) (*database.User, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

// ChatService handles chat operations
type ChatService struct {
	chatRepo *database.ChatRepository
	userRepo *database.UserRepository
	config   *config.Config
}

func NewChatService(chatRepo *database.ChatRepository, userRepo *database.UserRepository, cfg *config.Config) *ChatService {
	return &ChatService{
		chatRepo: chatRepo,
		userRepo: userRepo,
		config:   cfg,
	}
}

func (s *ChatService) CreateSession(userID, title string) (*database.ChatSession, error) {
	session := &database.ChatSession{
		ID:       generateID(),
		UserID:   userID,
		Title:    title,
		IsActive: true,
	}

	err := s.chatRepo.CreateSession(session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

func (s *ChatService) SendMessage(userID, sessionID, content string) (*database.ChatMessage, error) {
	// Check if user has enough chat tokens
	err := s.userRepo.DecrementChatTokens(userID)
	if err != nil {
		return nil, fmt.Errorf("insufficient chat tokens: %w", err)
	}

	// Create user message
	userMessage := &database.ChatMessage{
		ID:        generateID(),
		SessionID: sessionID,
		UserID:    userID,
		Role:      "user",
		Content:   content,
		Tokens:    1,
	}

	err = s.chatRepo.CreateMessage(userMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}

	// TODO: Integrate with AI service untuk generate response
	// For now, return a simple response
	aiResponse := &database.ChatMessage{
		ID:        generateID(),
		SessionID: sessionID,
		UserID:    userID,
		Role:      "assistant",
		Content:   "This is a placeholder AI response. AI integration will be implemented next.",
		Tokens:    0,
	}

	err = s.chatRepo.CreateMessage(aiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to save AI response: %w", err)
	}

	return aiResponse, nil
}

func (s *ChatService) GetSessionMessages(sessionID string) ([]database.ChatMessage, error) {
	return s.chatRepo.GetSessionMessages(sessionID)
}

func (s *ChatService) GetUserSessions(userID string) ([]database.ChatSession, error) {
	return s.chatRepo.GetUserSessions(userID)
}

// SearchService handles search operations
type SearchService struct {
	searchRepo *database.SearchRepository
	userRepo   *database.UserRepository
	config     *config.Config
}

func NewSearchService(searchRepo *database.SearchRepository, userRepo *database.UserRepository, cfg *config.Config) *SearchService {
	return &SearchService{
		searchRepo: searchRepo,
		userRepo:   userRepo,
		config:     cfg,
	}
}

func (s *SearchService) Search(userID, query string) (*database.SearchQuery, error) {
	// Check if user has enough search tokens
	err := s.userRepo.DecrementSearchTokens(userID)
	if err != nil {
		return nil, fmt.Errorf("insufficient search tokens: %w", err)
	}

	// TODO: Implement actual search logic
	// For now, return placeholder results
	results := fmt.Sprintf("Search results for: %s\n\nThis is a placeholder search result. Search integration will be implemented next.", query)

	searchQuery := &database.SearchQuery{
		ID:      generateID(),
		UserID:  userID,
		Query:   query,
		Results: results,
		Tokens:  1,
	}

	err = s.searchRepo.CreateQuery(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to save search query: %w", err)
	}

	return searchQuery, nil
}

func (s *SearchService) GetUserQueries(userID string, limit int) ([]database.SearchQuery, error) {
	return s.searchRepo.GetUserQueries(userID, limit)
}

// RateLimitService handles rate limiting
type RateLimitService struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

func NewRateLimitService(requestsPerMinute int) *RateLimitService {
	return &RateLimitService{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Every(time.Minute / time.Duration(requestsPerMinute)),
		burst:    requestsPerMinute,
	}
}

func (s *RateLimitService) GetLimiter(key string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	limiter, exists := s.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(s.rate, s.burst)
		s.limiters[key] = limiter
	}

	return limiter
}

func (s *RateLimitService) Allow(key string) bool {
	return s.GetLimiter(key).Allow()
}

// Helper function untuk generate ID
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}