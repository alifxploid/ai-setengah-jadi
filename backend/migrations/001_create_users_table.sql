-- Migration: Create users table
-- Created: 2024-01-01
-- Description: Create users table with access key authentication and token limits

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255),
    access_key VARCHAR(255) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    chat_tokens INT DEFAULT 10,
    search_tokens INT DEFAULT 100,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_users_email (email),
    INDEX idx_users_access_key (access_key),
    INDEX idx_users_active (is_active),
    INDEX idx_users_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert default admin user
INSERT INTO users (
    id, 
    email, 
    name, 
    access_key, 
    is_active, 
    chat_tokens, 
    search_tokens
) VALUES (
    'admin-001', 
    'admin@lipdev.id', 
    'Admin User', 
    'ADMIN_ACCESS_KEY_2024', 
    TRUE, 
    1000, 
    1000
) ON DUPLICATE KEY UPDATE 
    updated_at = CURRENT_TIMESTAMP;