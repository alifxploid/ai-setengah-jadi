package repository

import (
	"database/sql"
	"fmt"
	"time"

	"gryt-backend/internal/models"
)

// AIRepository handles database operations for AI-related entities
type AIRepository struct {
	db *sql.DB
}

// NewAIRepository creates a new AI repository
func NewAIRepository(db *sql.DB) *AIRepository {
	return &AIRepository{db: db}
}

// Chat Session Operations

// CreateChatSession creates a new chat session
func (r *AIRepository) CreateChatSession(session *models.ChatSession) error {
	query := `
		INSERT INTO chat_sessions (id, user_id, title, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	
	_, err := r.db.Exec(query, session.ID, session.UserID, session.Title, session.CreatedAt, session.UpdatedAt)
	return err
}

// GetChatSession retrieves a chat session by ID
func (r *AIRepository) GetChatSession(sessionID, userID string) (*models.ChatSession, error) {
	query := `
		SELECT id, user_id, title, created_at, updated_at
		FROM chat_sessions
		WHERE id = ? AND user_id = ?
	`
	
	session := &models.ChatSession{}
	err := r.db.QueryRow(query, sessionID, userID).Scan(
		&session.ID,
		&session.UserID,
		&session.Title,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return session, nil
}

// GetUserChatSessions retrieves all chat sessions for a user
func (r *AIRepository) GetUserChatSessions(userID string, limit, offset int) ([]*models.ChatSession, error) {
	query := `
		SELECT id, user_id, title, created_at, updated_at
		FROM chat_sessions
		WHERE user_id = ?
		ORDER BY updated_at DESC
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var sessions []*models.ChatSession
	for rows.Next() {
		session := &models.ChatSession{}
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.Title,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	
	return sessions, nil
}

// UpdateChatSession updates a chat session
func (r *AIRepository) UpdateChatSession(session *models.ChatSession) error {
	query := `
		UPDATE chat_sessions
		SET title = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`
	
	_, err := r.db.Exec(query, session.Title, session.UpdatedAt, session.ID, session.UserID)
	return err
}

// DeleteChatSession deletes a chat session and all its messages
func (r *AIRepository) DeleteChatSession(sessionID, userID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Delete messages first
	_, err = tx.Exec("DELETE FROM chat_messages WHERE session_id = ?", sessionID)
	if err != nil {
		return err
	}
	
	// Delete session
	_, err = tx.Exec("DELETE FROM chat_sessions WHERE id = ? AND user_id = ?", sessionID, userID)
	if err != nil {
		return err
	}
	
	return tx.Commit()
}

// Chat Message Operations

// CreateChatMessage creates a new chat message
func (r *AIRepository) CreateChatMessage(message *models.ChatMessage) error {
	query := `
		INSERT INTO chat_messages (id, session_id, user_id, role, content, metadata, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.Exec(query,
		message.ID,
		message.SessionID,
		message.UserID,
		message.Role,
		message.Content,
		message.Metadata,
		message.CreatedAt,
	)
	return err
}

// GetChatMessages retrieves messages for a chat session
func (r *AIRepository) GetChatMessages(sessionID, userID string, limit, offset int) ([]*models.ChatMessage, error) {
	query := `
		SELECT id, session_id, user_id, role, content, metadata, created_at
		FROM chat_messages
		WHERE session_id = ? AND user_id = ?
		ORDER BY created_at ASC
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.db.Query(query, sessionID, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var messages []*models.ChatMessage
	for rows.Next() {
		message := &models.ChatMessage{}
		err := rows.Scan(
			&message.ID,
			&message.SessionID,
			&message.UserID,
			&message.Role,
			&message.Content,
			&message.Metadata,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	
	return messages, nil
}

// GetChatMessage retrieves a specific chat message
func (r *AIRepository) GetChatMessage(messageID, userID string) (*models.ChatMessage, error) {
	query := `
		SELECT id, session_id, user_id, role, content, metadata, created_at
		FROM chat_messages
		WHERE id = ? AND user_id = ?
	`
	
	message := &models.ChatMessage{}
	err := r.db.QueryRow(query, messageID, userID).Scan(
		&message.ID,
		&message.SessionID,
		&message.UserID,
		&message.Role,
		&message.Content,
		&message.Metadata,
		&message.CreatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return message, nil
}

// Search Query Operations

// CreateSearchQuery creates a new search query record
func (r *AIRepository) CreateSearchQuery(query *models.SearchQuery) error {
	sqlQuery := `
		INSERT INTO search_queries (id, user_id, query, results, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	
	_, err := r.db.Exec(sqlQuery, query.ID, query.UserID, query.Query, query.Results, query.CreatedAt)
	return err
}

// GetSearchHistory retrieves search history for a user
func (r *AIRepository) GetSearchHistory(userID string, limit, offset int) ([]*models.SearchQuery, error) {
	query := `
		SELECT id, user_id, query, results, created_at
		FROM search_queries
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var queries []*models.SearchQuery
	for rows.Next() {
		searchQuery := &models.SearchQuery{}
		err := rows.Scan(
			&searchQuery.ID,
			&searchQuery.UserID,
			&searchQuery.Query,
			&searchQuery.Results,
			&searchQuery.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		queries = append(queries, searchQuery)
	}
	
	return queries, nil
}

// GetSearchQuery retrieves a specific search query
func (r *AIRepository) GetSearchQuery(queryID, userID string) (*models.SearchQuery, error) {
	query := `
		SELECT id, user_id, query, results, created_at
		FROM search_queries
		WHERE id = ? AND user_id = ?
	`
	
	searchQuery := &models.SearchQuery{}
	err := r.db.QueryRow(query, queryID, userID).Scan(
		&searchQuery.ID,
		&searchQuery.UserID,
		&searchQuery.Query,
		&searchQuery.Results,
		&searchQuery.CreatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return searchQuery, nil
}

// Statistics and Analytics

// GetUserChatStats retrieves chat statistics for a user
func (r *AIRepository) GetUserChatStats(userID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Total sessions
	var totalSessions int
	err := r.db.QueryRow("SELECT COUNT(*) FROM chat_sessions WHERE user_id = ?", userID).Scan(&totalSessions)
	if err != nil {
		return nil, err
	}
	stats["total_sessions"] = totalSessions
	
	// Total messages
	var totalMessages int
	err = r.db.QueryRow("SELECT COUNT(*) FROM chat_messages WHERE user_id = ?", userID).Scan(&totalMessages)
	if err != nil {
		return nil, err
	}
	stats["total_messages"] = totalMessages
	
	// Total searches
	var totalSearches int
	err = r.db.QueryRow("SELECT COUNT(*) FROM search_queries WHERE user_id = ?", userID).Scan(&totalSearches)
	if err != nil {
		return nil, err
	}
	stats["total_searches"] = totalSearches
	
	// Recent activity (last 7 days)
	var recentSessions int
	weekAgo := time.Now().AddDate(0, 0, -7)
	err = r.db.QueryRow("SELECT COUNT(*) FROM chat_sessions WHERE user_id = ? AND created_at >= ?", userID, weekAgo).Scan(&recentSessions)
	if err != nil {
		return nil, err
	}
	stats["recent_sessions"] = recentSessions
	
	return stats, nil
}

// CleanupOldData removes old chat sessions and messages based on retention policy
func (r *AIRepository) CleanupOldData(retentionDays int) error {
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)
	
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Delete old messages
	_, err = tx.Exec("DELETE FROM chat_messages WHERE created_at < ?", cutoffDate)
	if err != nil {
		return fmt.Errorf("failed to delete old messages: %w", err)
	}
	
	// Delete old sessions
	_, err = tx.Exec("DELETE FROM chat_sessions WHERE created_at < ?", cutoffDate)
	if err != nil {
		return fmt.Errorf("failed to delete old sessions: %w", err)
	}
	
	// Delete old search queries
	_, err = tx.Exec("DELETE FROM search_queries WHERE created_at < ?", cutoffDate)
	if err != nil {
		return fmt.Errorf("failed to delete old search queries: %w", err)
	}
	
	return tx.Commit()
}