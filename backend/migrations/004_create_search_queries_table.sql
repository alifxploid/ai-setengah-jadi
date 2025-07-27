-- Migration: Create search_queries table
-- Created: 2024-01-01
-- Description: Create search queries table for storing user search history

CREATE TABLE IF NOT EXISTS search_queries (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    query TEXT NOT NULL,
    results LONGTEXT,
    tokens INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_search_queries_user_id (user_id),
    INDEX idx_search_queries_created_at (created_at),
    FULLTEXT idx_search_queries_query (query),
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;