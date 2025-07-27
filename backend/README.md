# 🚀 Gryt Chat Backend - Golang 1.24.5

Backend API untuk AI Chat Website menggunakan Golang 1.24.5 + Gin Framework dengan arsitektur modular dan performa tinggi.

## 🛠️ Tech Stack

- **Language**: Golang 1.24.5
- **Framework**: Gin Web Framework
- **Database**: MySQL 8.0+ dengan connection pooling
- **Authentication**: JWT + Access Key
- **Security**: Rate limiting, input sanitization, CORS
- **Deployment**: Docker + VPS (lipdev.id)

## 📁 Struktur Project

```
/backend/
├── /cmd/                    # Main applications
├── /internal/               # Private application code
│   ├── /api/               # API handlers & routes
│   ├── /config/            # Configuration management
│   ├── /database/          # Database layer & models
│   ├── /middleware/        # HTTP middleware
│   └── /services/          # Business logic services
├── /migrations/             # Database migrations (SQL)
├── .env.example             # Environment template
├── Dockerfile               # Multi-stage Docker build
├── go.mod                   # Go modules
├── go.sum                   # Go dependencies
├── knexfile.js              # Database migration config
├── main.go                  # Application entry point
└── README.md                # Documentation
```

## 🔧 Setup & Installation

### Prerequisites

- Go 1.24.5+
- MySQL 8.0+
- Node.js (untuk Knex migrations)

### 1. Clone & Setup

```bash
cd backend
cp .env.example .env
# Edit .env dengan konfigurasi yang sesuai
```

### 2. Install Dependencies

```bash
go mod tidy
npm install -g knex
```

### 3. Database Setup

```bash
# Buat database MySQL
mysql -u root -p -e "CREATE DATABASE gryt_chat CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# Jalankan migrations
knex migrate:latest
```

### 4. Run Development Server

```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## 🔐 Environment Variables

### Server Configuration
```env
PORT=8080
ENVIRONMENT=development
FRONTEND_URL=https://lipdev.id
DOMAIN=lipdev.id
```

### Database Configuration
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=gryt_chat
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=10
DB_MAX_LIFETIME=5m
```

### Authentication
```env
JWT_SECRET=your_super_secret_jwt_key_here_minimum_32_characters
ACCESS_KEY=your_admin_access_key_here
JWT_TOKEN_EXPIRY=24h
JWT_REFRESH_EXPIRY=168h
```

### Token Limits (Configurable)
```env
CHAT_TOKENS_PER_USER=10
SEARCH_TOKENS_PER_USER=100
RATE_LIMIT_PER_MINUTE=30
```

### AI Configuration (Vercel AI Gateway)
```env
AI_API_KEY=your_vercel_ai_api_key_here
AI_BASE_URL=https://ai-gateway.vercel.sh/v1
AI_MODEL=anthropic/claude-sonnet-4
AI_MAX_TOKENS=4096
AI_TEMPERATURE=0.7
AI_TIMEOUT=60s
AI_SYSTEM_PROMPT="You are GRYT, an advanced AI assistant..."

# Advanced AI Parameters
AI_TOP_P=1.0
AI_FREQUENCY_PENALTY=0.0
AI_PRESENCE_PENALTY=0.0
AI_STOP_SEQUENCES=
```

## 🔗 API Endpoints

### Health Check
```
GET /health
```

### Authentication
```
POST /api/auth/validate-key    # Validate access key & get JWT
POST /api/auth/login           # Email/password login (future)
POST /api/auth/refresh         # Refresh JWT token (future)
```

### Chat (Protected)
```
POST /api/chat/sessions                    # Create new chat session
GET  /api/chat/sessions                    # Get user's chat sessions
POST /api/chat/sessions/:id/messages       # Send message to session
GET  /api/chat/sessions/:id/messages       # Get session messages
```

### Search (Protected)
```
POST /api/search/              # Perform search query
GET  /api/search/history       # Get search history
```

### User (Protected)
```
GET /api/user/profile          # Get user profile
GET /api/user/tokens           # Get token counts
```

## 🔐 Authentication Flow

### 1. Access Key Validation
```bash
curl -X POST http://localhost:8080/api/auth/validate-key \
  -H "Content-Type: application/json" \
  -d '{"access_key": "ADMIN_ACCESS_KEY_2024"}'
```

### 2. Use JWT Token
```bash
curl -X GET http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🚀 Production Deployment

### Docker Build
```bash
docker build -t gryt-backend .
docker run -p 8080:8080 --env-file .env gryt-backend
```

### VPS Deployment
```bash
# Build untuk production
CGO_ENABLED=0 GOOS=linux go build -o gryt-backend

# Upload ke VPS
scp gryt-backend user@lipdev.id:/opt/gryt/
scp .env user@lipdev.id:/opt/gryt/

# Setup systemd service
sudo systemctl enable gryt-backend
sudo systemctl start gryt-backend
```

## 📊 Performance Features

- **Connection Pooling**: MySQL connection pool dengan max 25 connections
- **Rate Limiting**: 30 requests per minute per IP
- **Middleware Optimization**: Gzip compression, security headers
- **Graceful Shutdown**: 30 detik timeout untuk graceful shutdown
- **Health Checks**: Built-in health check endpoint
- **Structured Logging**: JSON formatted logs untuk production

## 🔒 Security Features

- **JWT Authentication**: Secure token-based authentication
- **Access Key System**: Admin-controlled access keys
- **Rate Limiting**: Anti-spam & DDoS protection
- **Input Sanitization**: XSS & injection prevention
- **CORS Protection**: Restricted cross-origin requests
- **Security Headers**: Comprehensive security headers
- **Token Limits**: Per-user chat & search token limits

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Benchmark tests
go test -bench=. ./...
```

## 📝 Database Schema

### Users Table
- `id`: Primary key (VARCHAR 36)
- `email`: Unique email address
- `name`: User display name
- `access_key`: Unique access key untuk authentication
- `chat_tokens`: Available chat tokens (default: 10)
- `search_tokens`: Available search tokens (default: 100)
- `is_active`: Account status

### Chat Sessions Table
- `id`: Primary key
- `user_id`: Foreign key ke users
- `title`: Session title
- `is_active`: Session status

### Chat Messages Table
- `id`: Primary key
- `session_id`: Foreign key ke chat_sessions
- `user_id`: Foreign key ke users
- `role`: Message role (user/assistant/system)
- `content`: Message content
- `tokens`: Token cost

### Search Queries Table
- `id`: Primary key
- `user_id`: Foreign key ke users
- `query`: Search query text
- `results`: Search results (JSON)
- `tokens`: Token cost

## 🔧 Development

### Code Standards
- Max 500 lines per file
- Proper error handling
- Structured logging
- Input validation
- Security best practices

### Git Workflow
```bash
git checkout -b feature/new-feature
# Make changes
git commit -m "feat: add new feature"
git push origin feature/new-feature
```

## 📞 Support

Untuk pertanyaan atau masalah, silakan buat issue di repository ini atau hubungi tim development.

---

**Built with ❤️ untuk lipdev.id**