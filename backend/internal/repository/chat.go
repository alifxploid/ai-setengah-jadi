package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ChatMessage represents a chat message
type ChatMessage struct {
	ID        int       `json:"id"`
	SessionID string    `json:"session_id"`
	UserID    int       `json:"user_id"`
	Message   string    `json:"message"`
	Response  string    `json:"response"`
	Tokens    int       `json:"tokens"`
	CreatedAt time.Time `json:"created_at"`
	IsDeleted bool      `json:"-"`
}

// ChatSession represents a chat session
type ChatSession struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsActive  bool      `json:"is_active"`
}

// ChatRepository handles chat database operations
type ChatRepository struct {
	db *sql.DB
}

// NewChatRepository creates a new chat repository
func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

// CreateSession creates a new chat session
func (r *ChatRepository) CreateSession(ctx context.Context, userID int, title string) (*ChatSession, error) {
	sessionID := uuid.New().String()
	
	query := `
		INSERT INTO chat_sessions (id, user_id, title, created_at, updated_at, is_active) 
		VALUES (?, ?, ?, NOW(), NOW(), 1)
	`
	
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, sessionID, userID, title)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	session := &ChatSession{
		ID:        sessionID,
		UserID:    userID,
		Title:     title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsActive:  true,
	}

	return session, nil
}

// GetSession retrieves a chat session by ID and user ID (security check)
func (r *ChatRepository) GetSession(ctx context.Context, sessionID string, userID int) (*ChatSession, error) {
	query := `
		SELECT id, user_id, title, created_at, updated_at, is_active 
		FROM chat_sessions 
		WHERE id = ? AND user_id = ? AND is_active = 1
	`
	
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var session ChatSession
	err = stmt.QueryRowContext(ctx, sessionID, userID).Scan(
		&session.ID,
		&session.UserID,
		&session.Title,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Session not found
		}
		return nil, fmt.Errorf("failed to scan session: %w", err)
	}

	return &session, nil
}

// SaveMessage saves a chat message with response
func (r *ChatRepository) SaveMessage(ctx context.Context, sessionID string, userID int, message, response string, tokens int) (*ChatMessage, error) {
	// First verify session belongs to user
	session, err := r.GetSession(ctx, sessionID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify session: %w", err)
	}
	if session == nil {
		return nil, fmt.Errorf("session not found or access denied")
	}

	query := `
		INSERT INTO chat_messages (session_id, user_id, message, response, tokens, created_at, is_deleted) 
		VALUES (?, ?, ?, ?, ?, NOW(), 0)
	`
	
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, sessionID, userID, message, response, tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to save message: %w", err)
	}

	messageID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get message ID: %w", err)
	}

	// Update session timestamp
	err = r.UpdateSessionTimestamp(ctx, sessionID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	chatMessage := &ChatMessage{
		ID:        int(messageID),
		SessionID: sessionID,
		UserID:    userID,
		Message:   message,
		Response:  response,
		Tokens:    tokens,
		CreatedAt: time.Now(),
		IsDeleted: false,
	}

	return chatMessage, nil
}

// GetMessages retrieves chat messages for a session with pagination
func (r *ChatRepository) GetMessages(ctx context.Context, sessionID string, userID int, limit, offset int) ([]*ChatMessage, error) {
	// First verify session belongs to user
	session, err := r.GetSession(ctx, sessionID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify session: %w", err)
	}
	if session == nil {
		return nil, fmt.Errorf("session not found or access denied")
	}

	query := `
		SELECT id, session_id, user_id, message, response, tokens, created_at 
		FROM chat_messages 
		WHERE session_id = ? AND user_id = ? AND is_deleted = 0 
		ORDER BY created_at ASC 
		LIMIT ? OFFSET ?
	`
	
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, sessionID, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []*ChatMessage
	for rows.Next() {
		var msg ChatMessage
		err = rows.Scan(
			&msg.ID,
			&msg.SessionID,
			&msg.UserID,
			&msg.Message,
			&msg.Response,
			&msg.Tokens,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, &msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return messages, nil
}

// UpdateSessionTimestamp updates session's updated_at timestamp
func (r *ChatRepository) UpdateSessionTimestamp(ctx context.Context, sessionID string, userID int) error {
	query := `UPDATE chat_sessions SET updated_at = NOW() WHERE id = ? AND user_id = ? AND is_active = 1`
	
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, sessionID, userID)
	if err != nil {
		return fmt.Errorf("failed to update session timestamp: %w", err)
	}

	return nil
}

// GetUserSessions retrieves all sessions for a user
func (r *ChatRepository) GetUserSessions(ctx context.Context, userID int, limit, offset int) ([]*ChatSession, error) {
	query := `
		SELECT id, user_id, title, created_at, updated_at, is_active 
		FROM chat_sessions 
		WHERE user_id = ? AND is_active = 1 
		ORDER BY updated_at DESC 
		LIMIT ? OFFSET ?
	`
	
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*ChatSession
	for rows.Next() {
		var session ChatSession
		err = rows.Scan(
			&session.ID,
			&session.UserID,
			&session.Title,
			&session.CreatedAt,
			&session.UpdatedAt,
			&session.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, &session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return sessions, nil
}