-- Create chat_sessions table
CREATE TABLE IF NOT EXISTS chat_sessions (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    title VARCHAR(500) NOT NULL DEFAULT 'New Chat',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_chat_sessions_user_id (user_id),
    INDEX idx_chat_sessions_updated_at (updated_at)
);

-- Create chat_messages table
CREATE TABLE IF NOT EXISTS chat_messages (
    id VARCHAR(255) PRIMARY KEY,
    session_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    role ENUM('user', 'assistant', 'system') NOT NULL,
    content JSON NOT NULL,
    metadata JSON DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_chat_messages_session_id (session_id),
    INDEX idx_chat_messages_user_id (user_id),
    INDEX idx_chat_messages_created_at (created_at),
    FOREIGN KEY (session_id) REFERENCES chat_sessions(id) ON DELETE CASCADE
);

-- Create search_queries table
CREATE TABLE IF NOT EXISTS search_queries (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    query TEXT NOT NULL,
    results JSON DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_search_queries_user_id (user_id),
    INDEX idx_search_queries_created_at (created_at),
    FULLTEXT INDEX idx_search_queries_query (query)
);

-- Create ai_usage_stats table for tracking usage and rate limiting
CREATE TABLE IF NOT EXISTS ai_usage_stats (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    endpoint VARCHAR(100) NOT NULL, -- chat, search, etc.
    tokens_used INT DEFAULT 0,
    requests_count INT DEFAULT 1,
    date DATE NOT NULL,
    hour TINYINT NOT NULL, -- 0-23
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_user_endpoint_date_hour (user_id, endpoint, date, hour),
    INDEX idx_ai_usage_user_date (user_id, date),
    INDEX idx_ai_usage_endpoint_date (endpoint, date)
);

-- Create ai_model_configs table for dynamic model configuration
CREATE TABLE IF NOT EXISTS ai_model_configs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    model_name VARCHAR(100) NOT NULL UNIQUE,
    display_name VARCHAR(200) NOT NULL,
    provider VARCHAR(50) NOT NULL, -- openai, anthropic, google, etc.
    max_tokens INT NOT NULL DEFAULT 4096,
    temperature DECIMAL(3,2) NOT NULL DEFAULT 0.70,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    cost_per_1k_input_tokens DECIMAL(10,6) DEFAULT 0,
    cost_per_1k_output_tokens DECIMAL(10,6) DEFAULT 0,
    supports_streaming BOOLEAN NOT NULL DEFAULT TRUE,
    supports_tools BOOLEAN NOT NULL DEFAULT TRUE,
    supports_vision BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_ai_models_enabled (enabled),
    INDEX idx_ai_models_provider (provider)
);

-- Insert default model configurations
INSERT INTO ai_model_configs (model_name, display_name, provider, max_tokens, temperature, supports_vision, supports_tools) VALUES
('gpt-4o', 'GPT-4o', 'openai', 4096, 0.7, TRUE, TRUE),
('gpt-4o-mini', 'GPT-4o Mini', 'openai', 2048, 0.7, TRUE, TRUE),
('claude-3-5-sonnet-20241022', 'Claude 3.5 Sonnet', 'anthropic', 4096, 0.7, TRUE, TRUE),
('gemini-1.5-pro', 'Gemini 1.5 Pro', 'google', 4096, 0.7, TRUE, TRUE)
ON DUPLICATE KEY UPDATE
    display_name = VALUES(display_name),
    max_tokens = VALUES(max_tokens),
    temperature = VALUES(temperature),
    supports_vision = VALUES(supports_vision),
    supports_tools = VALUES(supports_tools),
    updated_at = CURRENT_TIMESTAMP;

-- Create ai_tool_configs table for dynamic tool configuration
CREATE TABLE IF NOT EXISTS ai_tool_configs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tool_name VARCHAR(100) NOT NULL UNIQUE,
    display_name VARCHAR(200) NOT NULL,
    description TEXT,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    config JSON DEFAULT NULL, -- Tool-specific configuration
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_ai_tools_enabled (enabled)
);

-- Insert default tool configurations
INSERT INTO ai_tool_configs (tool_name, display_name, description, enabled, config) VALUES
('web_search', 'Web Search', 'Search the internet for current information', TRUE, JSON_OBJECT('max_results', 10, 'timeout', 30)),
('calculator', 'Calculator', 'Perform mathematical calculations', TRUE, JSON_OBJECT('precision', 10)),
('weather', 'Weather', 'Get current weather information', TRUE, JSON_OBJECT('units', 'metric', 'timeout', 15)),
('translation', 'Translation', 'Translate text between languages', TRUE, JSON_OBJECT('max_length', 5000)),
('image_analysis', 'Image Analysis', 'Analyze and describe images', TRUE, JSON_OBJECT('max_size_mb', 10, 'supported_formats', JSON_ARRAY('jpg', 'jpeg', 'png', 'webp')))
ON DUPLICATE KEY UPDATE
    display_name = VALUES(display_name),
    description = VALUES(description),
    config = VALUES(config),
    updated_at = CURRENT_TIMESTAMP;

-- Create indexes for better performance
CREATE INDEX idx_chat_sessions_user_updated ON chat_sessions(user_id, updated_at DESC);
CREATE INDEX idx_chat_messages_session_created ON chat_messages(session_id, created_at ASC);
CREATE INDEX idx_search_queries_user_created ON search_queries(user_id, created_at DESC);