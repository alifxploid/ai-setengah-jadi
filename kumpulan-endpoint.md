# 🔥 GRYT Backend API Endpoints Documentation

## 🚀 Base URL
```
https://lipdev.id
```

## 🔐 Authentication

### JWT Token Information
- **JWT Secret**: `sEJmU2nN5D1oQrEm07i13UURnSXBIVIz42bHFz7pMAP_pxoSPEiJgXJgyAM0xXCcxNJx3xSSFeXyeGXAGUACLQ==`
- **Access Key**: `NBzbgYqdG6oErJPOKzi4JkaFK3eka8C5TPcz4uLikuY=`
- **Token Expiry**: 24 hours
- **Refresh Token Expiry**: 168 hours (7 days)

### Sample JWT Token (Valid for Testing)
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJleHAiOjE3Mzc5NzQ0MzQsImlhdCI6MTczNzg4ODAzNH0.8vQxK5rZYmJ9X8wL3nP2mR4sT6uI0oE7cA1bF9dG3hK
```

---

## 📋 Complete Endpoint List

### 1. Health Check

#### GET /health
**Description**: Check server status

```bash
curl -X GET "https://lipdev.id/health" \
  -H "Content-Type: application/json"
```

**Response**:
```json
{
  "status": "ok",
  "message": "Server is running",
  "version": "1.0.0"
}
```

---

## 🔐 Authentication Endpoints

### 2. Validate Access Key

#### POST /api/auth/validate-key
**Description**: Validate access key and get JWT token

```bash
curl -X POST "https://lipdev.id/api/auth/validate-key" \
  -H "Content-Type: application/json" \
  -d '{
    "access_key": "NBzbgYqdG6oErJPOKzi4JkaFK3eka8C5TPcz4uLikuY="
  }'
```

**Response**:
```json
{
  "message": "Access key validated successfully",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "user-123",
    "name": "John Doe",
    "email": "john@example.com",
    "chat_tokens": 10,
    "search_tokens": 100
  }
}
```

### 3. Login (Email/Password)

#### POST /api/auth/login
**Description**: Login with email and password (Not implemented yet)

```bash
curl -X POST "https://lipdev.id/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Response**:
```json
{
  "error": "Email/password login not implemented yet",
  "message": "Please use access key authentication"
}
```

### 4. Refresh Token

#### POST /api/auth/refresh
**Description**: Refresh JWT token (Not implemented yet)

```bash
curl -X POST "https://lipdev.id/api/auth/refresh" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response**:
```json
{
  "error": "Token refresh not implemented yet"
}
```

---

## 🤖 AI Endpoints (Protected)

### 5. Create Chat Session

#### POST /api/ai/chat/sessions
**Description**: Create a new AI chat session

```bash
curl -X POST "https://lipdev.id/api/ai/chat/sessions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJleHAiOjE3Mzc5NzQ0MzQsImlhdCI6MTczNzg4ODAzNH0.8vQxK5rZYmJ9X8wL3nP2mR4sT6uI0oE7cA1bF9dG3hK" \
  -d '{
    "title": "My AI Chat Session"
  }'
```

**Response**:
```json
{
  "id": "session-1737888034123456789",
  "title": "My AI Chat Session",
  "created_at": "2025-01-26T10:30:00Z"
}
```

### 6. Get Chat Sessions

#### GET /api/ai/chat/sessions
**Description**: Get all chat sessions for the user

```bash
curl -X GET "https://lipdev.id/api/ai/chat/sessions" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJleHAiOjE3Mzc5NzQ0MzQsImlhdCI6MTczNzg4ODAzNH0.8vQxK5rZYmJ9X8wL3nP2mR4sT6uI0oE7cA1bF9dG3hK"
```

**Response**:
```json
{
  "sessions": [
    {
      "id": "session-123",
      "title": "My AI Chat Session",
      "created_at": "2025-01-26T10:30:00Z",
      "updated_at": "2025-01-26T10:30:00Z"
    }
  ]
}
```

### 7. Get Specific Chat Session

#### GET /api/ai/chat/sessions/:session_id
**Description**: Get details of a specific chat session

```bash
curl -X GET "https://lipdev.id/api/ai/chat/sessions/session-123" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJleHAiOjE3Mzc5NzQ0MzQsImlhdCI6MTczNzg4ODAzNH0.8vQxK5rZYmJ9X8wL3nP2mR4sT6uI0oE7cA1bF9dG3hK"
```

### 8. Delete Chat Session

#### DELETE /api/ai/chat/sessions/:session_id
**Description**: Delete a chat session

```bash
curl -X DELETE "https://lipdev.id/api/ai/chat/sessions/session-123" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJleHAiOjE3Mzc5NzQ0MzQsImlhdCI6MTczNzg4ODAzNH0.8vQxK5rZYmJ9X8wL3nP2mR4sT6uI0oE7cA1bF9dG3hK"
```

**Response**:
```json
{
  "message": "Chat session deleted successfully"
}
```

### 9. Send Chat Message

#### POST /api/ai/chat/sessions/:session_id/messages
**Description**: Send a message to AI and get response

```bash
curl -X POST "https://lipdev.id/api/ai/chat/sessions/session-123/messages" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJleHAiOjE3Mzc5NzQ0MzQsImlhdCI6MTczNzg4ODAzNH0.8vQxK5rZYmJ9X8wL3nP2mR4sT6uI0oE7cA1bF9dG3hK" \
  -d '{
    "message": "Hello, how are you?",
    "stream": false
  }'
```

**With File Upload**:
```bash
curl -X POST "https://lipdev.id/api/ai/chat/sessions/session-123/messages" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJleHAiOjE3Mzc5NzQ0MzQsImlhdCI6MTczNzg4ODAzNH0.8vQxK5rZYmJ9X8wL3nP2mR4sT6uI0oE7cA1bF9dG3hK" \
  -F "message=Analyze this image" \
  -F "files=@/path/to/image.jpg"
```

### 10. Get Chat Messages

#### GET /api/ai/chat/sessions/:session_id/messages
**Description**: Get messages from a chat session with pagination

```bash
curl -X GET "https://lipdev.id/api/ai/chat/sessions/session-123/messages?limit=20&offset=0" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJleHAiOjE3Mzc5NzQ0MzQsImlhdCI6MTczNzg4ODAzNH0.8vQxK5rZYmJ9X8wL3nP2mR4sT6uI0oE7cA1bF9dG3hK"
```

**Response**:
```json
{
  "messages": [
    {
      "id": "msg-123",
      "role": "user",
      "content": "Hello, how are you?",
      "created_at": "2025-01-26T10:30:00Z"
    },
    {
      "id": "msg-124",
      "role": "assistant",
      "content": "Hello! I'm doing well, thank you for asking. How can I help you today?",
      "created_at": "2025-01-26T10:30:01Z"
    }
  ]
}
```

### 11. Stream Chat Message

#### POST /api/ai/chat/sessions/:session_id/stream
**Description**: Send a message and get streaming response (Server-Sent Events)

```bash
curl -X POST "https://lipdev.id/api/ai/chat/sessions/session-123/stream" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJleHAiOjE3Mzc5NzQ0MzQsImlhdCI6MTczNzg4ODAzNH0.8vQxK5rZYmJ9X8wL3nP2mR4sT6uI0oE7cA1bF9dG3hK" \
  -d '{
    "message": "Tell me a story",
    "stream": true
  }'
```

**Response** (Server-Sent Events):
```
data: {"type": "message", "content": "Once upon"}

data: {"type": "message", "content": " a time"}

data: {"type": "done"}
```

### 12. Search Query

#### POST /api/ai/search/
**Description**: Perform AI-powered search

```bash
curl -X POST "https://lipdev.id/api/ai/search/" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJleHAiOjE3Mzc5NzQ0MzQsImlhdCI6MTczNzg4ODAzNH0.8vQxK5rZYmJ9X8wL3nP2mR4sT6uI0oE7cA1bF9dG3hK" \
  -d '{
    "query": "latest AI developments",
    "limit": 10
  }'
```

**Response**:
```json
{
  "search_id": "search-123",
  "query": "latest AI developments",
  "results": [
    {
      "title": "AI Development News",
      "url": "https://example.com/ai-news",
      "snippet": "Latest developments in AI technology..."
    }
  ],
  "created_at": "2025-01-26T10:30:00Z"
}
```

### 13. Get Search History

#### GET /api/ai/search/history
**Description**: Get user's search history with pagination

```bash
curl -X GET "https://lipdev.id/api/ai/search/history?limit=20&offset=0" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJleHAiOjE3Mzc5NzQ0MzQsImlhdCI6MTczNzg4ODAzNH0.8vQxK5rZYmJ9X8wL3nP2mR4sT6uI0oE7cA1bF9dG3hK"
```

**Response**:
```json
{
  "history": [
    {
      "id": "search-123",
      "query": "latest AI developments",
      "created_at": "2025-01-26T10:30:00Z"
    }
  ]
}
```

### 14. Get Search Result

#### GET /api/ai/search/:search_id
**Description**: Get specific search result by ID

```bash
curl -X GET "https://lipdev.id/api/ai/search/search-123" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJleHAiOjE3Mzc5NzQ0MzQsImlhdCI6MTczNzg4ODAzNH0.8vQxK5rZYmJ9X8wL3nP2mR4sT6uI0oE7cA1bF9dG3hK"
```

### 15. Get Available Models

#### GET /api/ai/models/
**Description**: Get list of available AI models

```bash
curl -X GET "https://lipdev.id/api/ai/models/" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJleHAiOjE3Mzc5NzQ0MzQsImlhdCI6MTczNzg4ODAzNH0.8vQxK5rZYmJ9X8wL3nP2mR4sT6uI0oE7cA1bF9dG3hK"
```

**Response**:
```json
{
  "models": [
    {
      "id": "gpt-4o",
      "name": "GPT-4o",
      "description": "Most capable GPT-4 model",
      "context": 128000
    },
    {
      "id": "gpt-4o-mini",
      "name": "GPT-4o Mini",
      "description": "Faster and more affordable GPT-4 model",
      "context": 128000
    },
    {
      "id": "claude-3-5-sonnet-20241022",
      "name": "Claude 3.5 Sonnet",
      "description": "Anthropic's most capable model",
      "context": 200000
    }
  ]
}
```

---

## 💬 Legacy Chat Endpoints (Protected)

### 16. Create Legacy Chat Session

#### POST /api/chat/sessions
**Description**: Create a new chat session (legacy endpoint)

```bash
curl -X POST "https://lipdev.id/api/chat/sessions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJleHAiOjE3Mzc5NzQ0MzQsImlhdCI6MTczNzg4ODAzNH0.8vQxK5rZYmJ9X8wL3nP2mR4sT6uI0oE7cA1bF9dG3hK" \
  -d '{
    "title": "Legacy Chat Session"
  }'
```

### 17. Get Legacy User Sessions

#### GET /api/chat/sessions
**Description**: Get all chat sessions for user (legacy endpoint)

```bash
curl -X GET "https://lipdev.id/api/chat/sessions" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJleHAiOjE3Mzc5NzQ0MzQsImlhdCI6MTczNzg4ODAzNH0.8vQxK5rZYmJ9X8wL3nP2mR4sT6uI0oE7cA1bF9dG3hK"
```

### 18. Send Legacy Message

#### POST /api/chat/sessions/:sessionId/messages
**Description**: Send message to legacy chat session

```bash
curl -X POST "https://lipdev.id/api/chat/sessions/session-123/messages" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJleHAiOjE3Mzc5NzQ0MzQsImlhdCI6MTczNzg4ODAzNH0.8vQxK5rZYmJ9X8wL3nP2mR4sT6uI0oE7cA1bF9dG3hK" \
  -d '{
    "content": "Hello from legacy endpoint"
  }'
```

### 19. Get Legacy Session Messages

#### GET /api/chat/sessions/:sessionId/messages
**Description**: Get messages from legacy chat session

```bash
curl -X GET "https://lipdev.id/api/chat/sessions/session-123/messages" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJleHAiOjE3Mzc5NzQ0MzQsImlhdCI6MTczNzg4ODAzNH0.8vQxK5rZYmJ9X8wL3nP2mR4sT6uI0oE7cA1bF9dG3hK"
```

---

## 🔍 Legacy Search Endpoints (Protected)

### 20. Perform Legacy Search

#### POST /api/search/
**Description**: Perform search (legacy endpoint)

```bash
curl -X POST "https://lipdev.id/api/search/" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJleHAiOjE3Mzc5NzQ0MzQsImlhdCI6MTczNzg4ODAzNH0.8vQxK5rZYmJ9X8wL3nP2mR4sT6uI0oE7cA1bF9dG3hK" \
  -d '{
    "query": "search term"
  }'
```

### 21. Get Legacy Search History

#### GET /api/search/history
**Description**: Get search history (legacy endpoint)

```bash
curl -X GET "https://lipdev.id/api/search/history" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJleHAiOjE3Mzc5NzQ0MzQsImlhdCI6MTczNzg4ODAzNH0.8vQxK5rZYmJ9X8wL3nP2mR4sT6uI0oE7cA1bF9dG3hK"
```

---

## 👤 User Endpoints (Protected)

### 22. Get User Profile

#### GET /api/user/profile
**Description**: Get user profile information (To be implemented)

```bash
curl -X GET "https://lipdev.id/api/user/profile" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJleHAiOjE3Mzc5NzQ0MzQsImlhdCI6MTczNzg4ODAzNH0.8vQxK5rZYmJ9X8wL3nP2mR4sT6uI0oE7cA1bF9dG3hK"
```

**Response**:
```json
{
  "user_id": "user-123",
  "message": "User profile endpoint - to be implemented"
}
```

### 23. Get User Tokens

#### GET /api/user/tokens
**Description**: Get user token counts (To be implemented)

```bash
curl -X GET "https://lipdev.id/api/user/tokens" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci0xMjMiLCJleHAiOjE3Mzc5NzQ0MzQsImlhdCI6MTczNzg4ODAzNH0.8vQxK5rZYmJ9X8wL3nP2mR4sT6uI0oE7cA1bF9dG3hK"
```

**Response**:
```json
{
  "user_id": "user-123",
  "message": "User tokens endpoint - to be implemented"
}
```

---

## 📊 Summary

### Total Endpoints: 23

#### Public Endpoints (1):
- ✅ `GET /health` - Health check

#### Authentication Endpoints (3):
- ✅ `POST /api/auth/validate-key` - Validate access key
- ⚠️ `POST /api/auth/login` - Login (Not implemented)
- ⚠️ `POST /api/auth/refresh` - Refresh token (Not implemented)

#### AI Endpoints (11):
- ✅ `POST /api/ai/chat/sessions` - Create chat session
- ✅ `GET /api/ai/chat/sessions` - Get chat sessions
- ✅ `GET /api/ai/chat/sessions/:session_id` - Get specific session
- ✅ `DELETE /api/ai/chat/sessions/:session_id` - Delete session
- ✅ `POST /api/ai/chat/sessions/:session_id/messages` - Send message
- ✅ `GET /api/ai/chat/sessions/:session_id/messages` - Get messages
- ✅ `POST /api/ai/chat/sessions/:session_id/stream` - Stream message
- ✅ `POST /api/ai/search/` - Search query
- ✅ `GET /api/ai/search/history` - Search history
- ✅ `GET /api/ai/search/:search_id` - Get search result
- ✅ `GET /api/ai/models/` - Get available models

#### Legacy Chat Endpoints (4):
- ✅ `POST /api/chat/sessions` - Create legacy session
- ✅ `GET /api/chat/sessions` - Get legacy sessions
- ✅ `POST /api/chat/sessions/:sessionId/messages` - Send legacy message
- ✅ `GET /api/chat/sessions/:sessionId/messages` - Get legacy messages

#### Legacy Search Endpoints (2):
- ✅ `POST /api/search/` - Legacy search
- ✅ `GET /api/search/history` - Legacy search history

#### User Endpoints (2):
- ⚠️ `GET /api/user/profile` - User profile (To be implemented)
- ⚠️ `GET /api/user/tokens` - User tokens (To be implemented)

---

## 🔧 Testing Notes

### Rate Limiting:
- AI endpoints: 10 requests per minute
- General endpoints: 30 requests per minute
- Chat tokens: 10 per user
- Search tokens: 100 per user

### Authentication:
- All protected endpoints require `Authorization: Bearer <JWT_TOKEN>` header
- JWT tokens expire after 24 hours
- Use the provided access key to get a valid JWT token

### File Upload Support:
- AI chat endpoints support file uploads via multipart/form-data
- Supported file types: images, documents
- Max file size: 10MB

### Streaming Support:
- AI chat streaming uses Server-Sent Events (SSE)
- Content-Type: `text/event-stream`
- Real-time message chunks

### Error Handling:
- 400: Bad Request (invalid input)
- 401: Unauthorized (missing/invalid token)
- 403: Forbidden (insufficient permissions)
- 429: Too Many Requests (rate limit exceeded)
- 500: Internal Server Error

---

## 🚀 Postman Testing Flow

1. **Health Check**: Test server connectivity
2. **Validate Access Key**: Get JWT token
3. **Create Chat Session**: Start new AI conversation
4. **Send Messages**: Test AI chat functionality
5. **Test Streaming**: Verify real-time responses
6. **Search Functionality**: Test AI search
7. **File Upload**: Test multimodal capabilities
8. **Legacy Endpoints**: Verify backward compatibility

---

**Generated on**: 2025-01-26  
**Backend Version**: 1.0.0  
**Domain**: https://lipdev.id  
**Environment**: Development