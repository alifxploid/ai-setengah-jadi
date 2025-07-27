package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gryt-backend/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func NewConnection(cfg config.DatabaseConfig) (*DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Connection pool optimization
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")
	return &DB{db}, nil
}

// User represents a user in the system
type User struct {
	ID           string         `db:"id" json:"id"`
	Email        string         `db:"email" json:"email"`
	Name         string         `db:"name" json:"name"`
	PasswordHash sql.NullString `db:"password_hash" json:"-"`
	AccessKey    string         `db:"access_key" json:"access_key"`
	IsActive     bool           `db:"is_active" json:"is_active"`
	ChatTokens   int            `db:"chat_tokens" json:"chat_tokens"`
	SearchTokens int            `db:"search_tokens" json:"search_tokens"`
	CreatedAt    time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at" json:"updated_at"`
}

// ChatSession represents a chat session
type ChatSession struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Title     string    `db:"title" json:"title"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// ChatMessage represents a chat message
type ChatMessage struct {
	ID        string    `db:"id" json:"id"`
	SessionID string    `db:"session_id" json:"session_id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Role      string    `db:"role" json:"role"` // user, assistant, system
	Content   string    `db:"content" json:"content"`
	Tokens    int       `db:"tokens" json:"tokens"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// SearchQuery represents a search query
type SearchQuery struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Query     string    `db:"query" json:"query"`
	Results   string    `db:"results" json:"results"`
	Tokens    int       `db:"tokens" json:"tokens"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// UserRepository handles user database operations
type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByAccessKey(accessKey string) (*User, error) {
	var user User
	query := `SELECT id, email, name, password_hash, access_key, is_active, chat_tokens, search_tokens, created_at, updated_at 
			  FROM users WHERE access_key = ? AND is_active = 1`

	log.Printf("[DEBUG] GetByAccessKey called with accessKey: %s", accessKey)
	err := r.db.Get(&user, query, accessKey)
	if err != nil {
		log.Printf("[DEBUG] GetByAccessKey error: %v", err)
		if err == sql.ErrNoRows {
			log.Printf("[DEBUG] No rows found for access key: %s", accessKey)
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by access key: %w", err)
	}

	log.Printf("[DEBUG] User found: %+v", user)
	return &user, nil
}

func (r *UserRepository) GetByID(id string) (*User, error) {
	var user User
	query := `SELECT id, email, name, password_hash, access_key, is_active, chat_tokens, search_tokens, created_at, updated_at 
			  FROM users WHERE id = ? AND is_active = 1`

	err := r.db.Get(&user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	var user User
	query := `SELECT id, email, name, password_hash, access_key, is_active, chat_tokens, search_tokens, created_at, updated_at 
			  FROM users WHERE email = ? AND is_active = 1`

	log.Printf("[DEBUG] GetByEmail called with email: %s", email)
	err := r.db.Get(&user, query, email)
	if err != nil {
		log.Printf("[DEBUG] GetByEmail error: %v", err)
		if err == sql.ErrNoRows {
			log.Printf("[DEBUG] No rows found for email: %s", email)
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	log.Printf("[DEBUG] User found by email: %+v", user)
	return &user, nil
}

func (r *UserRepository) UpdateTokens(userID string, chatTokens, searchTokens int) error {
	query := `UPDATE users SET chat_tokens = ?, search_tokens = ?, updated_at = NOW() WHERE id = ?`

	_, err := r.db.Exec(query, chatTokens, searchTokens, userID)
	if err != nil {
		return fmt.Errorf("failed to update user tokens: %w", err)
	}

	return nil
}

func (r *UserRepository) DecrementChatTokens(userID string) error {
	query := `UPDATE users SET chat_tokens = chat_tokens - 1, updated_at = NOW() 
			  WHERE id = ? AND chat_tokens > 0`

	result, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to decrement chat tokens: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("insufficient chat tokens")
	}

	return nil
}

func (r *UserRepository) DecrementSearchTokens(userID string) error {
	query := `UPDATE users SET search_tokens = search_tokens - 1, updated_at = NOW() 
			  WHERE id = ? AND search_tokens > 0`

	result, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to decrement search tokens: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("insufficient search tokens")
	}

	return nil
}

// ChatRepository handles chat database operations
type ChatRepository struct {
	db *DB
}

func NewChatRepository(db *DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) CreateSession(session *ChatSession) error {
	query := `INSERT INTO chat_sessions (id, user_id, title, is_active, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, NOW(), NOW())`

	_, err := r.db.Exec(query, session.ID, session.UserID, session.Title, session.IsActive)
	if err != nil {
		return fmt.Errorf("failed to create chat session: %w", err)
	}

	return nil
}

func (r *ChatRepository) GetUserSessions(userID string) ([]ChatSession, error) {
	var sessions []ChatSession
	query := `SELECT id, user_id, title, is_active, created_at, updated_at 
			  FROM chat_sessions WHERE user_id = ? ORDER BY updated_at DESC`

	err := r.db.Select(&sessions, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}

	return sessions, nil
}

func (r *ChatRepository) CreateMessage(message *ChatMessage) error {
	query := `INSERT INTO chat_messages (id, session_id, user_id, role, content, tokens, created_at) 
			  VALUES (?, ?, ?, ?, ?, ?, NOW())`

	_, err := r.db.Exec(query, message.ID, message.SessionID, message.UserID, 
		message.Role, message.Content, message.Tokens)
	if err != nil {
		return fmt.Errorf("failed to create chat message: %w", err)
	}

	return nil
}

func (r *ChatRepository) GetSessionMessages(sessionID string) ([]ChatMessage, error) {
	var messages []ChatMessage
	query := `SELECT id, session_id, user_id, role, content, tokens, created_at 
			  FROM chat_messages WHERE session_id = ? ORDER BY created_at ASC`

	err := r.db.Select(&messages, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session messages: %w", err)
	}

	return messages, nil
}

// SearchRepository handles search database operations
type SearchRepository struct {
	db *DB
}

func NewSearchRepository(db *DB) *SearchRepository {
	return &SearchRepository{db: db}
}

func (r *SearchRepository) CreateQuery(query *SearchQuery) error {
	sql := `INSERT INTO search_queries (id, user_id, query, results, tokens, created_at) 
			VALUES (?, ?, ?, ?, ?, NOW())`

	_, err := r.db.Exec(sql, query.ID, query.UserID, query.Query, query.Results, query.Tokens)
	if err != nil {
		return fmt.Errorf("failed to create search query: %w", err)
	}

	return nil
}

func (r *SearchRepository) GetUserQueries(userID string, limit int) ([]SearchQuery, error) {
	var queries []SearchQuery
	sql := `SELECT id, user_id, query, results, tokens, created_at 
			FROM search_queries WHERE user_id = ? ORDER BY created_at DESC LIMIT ?`

	err := r.db.Select(&queries, sql, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user queries: %w", err)
	}

	return queries, nil
}