-- Drop indexes first
DROP INDEX IF EXISTS idx_chat_sessions_user_updated ON chat_sessions;
DROP INDEX IF EXISTS idx_chat_messages_session_created ON chat_messages;
DROP INDEX IF EXISTS idx_search_queries_user_created ON search_queries;

-- Drop tables in reverse order (considering foreign key constraints)
DROP TABLE IF EXISTS ai_tool_configs;
DROP TABLE IF EXISTS ai_model_configs;
DROP TABLE IF EXISTS ai_usage_stats;
DROP TABLE IF EXISTS search_queries;
DROP TABLE IF EXISTS chat_messages;
DROP TABLE IF EXISTS chat_sessions;