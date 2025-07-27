# 🔥 AI Chat Website + Golang Backend Rules

## 🚨 DEPLOYMENT INFO - VPS ONLY!
⚠️ **DOMAIN WAJIB**: https://lipdev.id (NO LOCALHOST!)
⚠️ **FRONTEND**: Next.js 15.4.4 + Tailwind CSS v4
⚠️ **BACKEND**: Golang 1.24.5 + Gin Framework
⚠️ **NGINX**: Bebas akses, update pake `sed` command
⚠️ **ENV ONLY**: Zero hardcode, semua pake environment variables

## 🎯 Core Mission
Lo adalah backend engineer terhandal yang bikin AI chat website pake Next.js 15.4.4 + Golang 1.24.5. Kode harus:
- ⚡ **BLAZING FAST**: Optimasi Next.js & Golang sampe menggila (100x faster performance)
- 🔒 **ULTRA SECURE**: Lapisan keamanan di setiap konfigurasi
- 📦 **MODULAR AF**: Struktur terpisah, max 500 lines per file
- 🌐 **VPS READY**: Domain lipdev.id, no localhost bullshit

## 🛠️ Tech Stack Wajib
- **Frontend**: Next.js 15.4.4 + React 19.1.0 + TypeScript
- **Styling**: Tailwind CSS v4.1.11 (latest)
- **Backend**: Golang 1.24.5 + Gin Framework + performance optimizations
- **Database**: MySQL + Knex.js migrations
- **AI**: Vercel AI Gateway integration
- **Security**: Environment-based config only
- **Domain**: https://lipdev.id (VPS deployment)
- **Proxy**: Nginx (bebas akses, update via sed)
- **Build**: Turbopack untuk development, optimized production build

## 🔐 Security Rules (NON-NEGOTIABLE)
1. **ZERO HARDCODE**: Semua config via process.env
2. **ENV VALIDATION**: Validate semua env vars on startup
3. **INPUT SANITIZATION**: Filter semua user input
4. **RATE LIMITING**: Anti-spam & DDoS protection
5. **ERROR HANDLING**: Never expose internal errors
6. **HTTPS ONLY**: SSL/TLS everywhere
7. **API KEY ROTATION**: Support key rotation
8. **WEBHOOK SECURITY**: Validate Telegram webhook dengan secret token

### Security Practices

#### Golang Backend Security
```go
// internal/security/validation.go
package security

import (
    "errors"
    "regexp"
    "strings"
    "sync"
    "time"
    "unicode/utf8"

    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
)

// Input validation
type MessageValidator struct {
    MaxLength int
}

func NewMessageValidator() *MessageValidator {
    return &MessageValidator{
        MaxLength: 4096,
    }
}

func (v *MessageValidator) ValidateMessage(text string, userID int64) error {
    if text == "" {
        return errors.New("message cannot be empty")
    }
    
    if !utf8.ValidString(text) {
        return errors.New("invalid UTF-8 encoding")
    }
    
    if len(text) > v.MaxLength {
        return errors.New("message too long")
    }
    
    if userID <= 0 {
        return errors.New("invalid user ID")
    }
    
    return nil
}

// Rate limiting
type RateLimiter struct {
    limiters map[int64]*rate.Limiter
    mu       sync.RWMutex
}

func NewRateLimiter() *RateLimiter {
    return &RateLimiter{
        limiters: make(map[int64]*rate.Limiter),
    }
}

func (rl *RateLimiter) GetLimiter(userID int64) *rate.Limiter {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    limiter, exists := rl.limiters[userID]
    if !exists {
        // 30 requests per minute
        limiter = rate.NewLimiter(rate.Every(time.Minute/30), 30)
        rl.limiters[userID] = limiter
    }
    
    return limiter
}

func (rl *RateLimiter) Allow(userID int64) bool {
    return rl.GetLimiter(userID).Allow()
}

// Input sanitization
func SanitizeInput(input string) string {
    // Remove dangerous characters
    re := regexp.MustCompile(`[<>"'&\x00-\x1f\x7f-\x9f]`)
    sanitized := re.ReplaceAllString(input, "")
    
    // Trim whitespace
    sanitized = strings.TrimSpace(sanitized)
    
    // Limit length
    if len(sanitized) > 4096 {
        sanitized = sanitized[:4096]
    }
    
    return sanitized
}

// Security middleware
func SecurityMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Security headers
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Header("Content-Security-Policy", "default-src 'self'")
        
        c.Next()
    }
}
```

#### Next.js Frontend Security
```typescript
// lib/security.ts
export const sanitizeInput = (input: string): string => {
  return input
    .replace(/[<>"'&]/g, '')
    .trim()
    .substring(0, 4096);
};

export const validateMessage = (message: string): boolean => {
  if (!message || message.trim().length === 0) return false;
  if (message.length > 4096) return false;
  return true;
};

// Rate limiting untuk client-side
export class ClientRateLimit {
  private requests: number[] = [];
  private readonly limit = 30; // per minute
  
  canMakeRequest(): boolean {
    const now = Date.now();
    this.requests = this.requests.filter(time => now - time < 60000);
    
    if (this.requests.length >= this.limit) {
      return false;
    }
    
    this.requests.push(now);
    return true;
  }
}
```

## 📁 Struktur Project (Modular)
```
/gryt/
├── /app/                    # Next.js 15 App Router
│   ├── /api/               # Next.js API routes
│   ├── /chat/              # Chat pages
│   ├── /auth/              # Authentication pages
│   ├── /dashboard/         # Dashboard pages
│   ├── globals.css         # Global styles
│   ├── layout.tsx          # Root layout
│   └── page.tsx            # Home page
├── /components/             # React components
│   ├── /ui/                # Reusable UI components
│   ├── /chat/              # Chat-specific components
│   ├── /auth/              # Auth components
│   └── /layout/            # Layout components
├── /lib/                    # Utility libraries
│   ├── /utils/             # Helper functions
│   ├── /hooks/             # Custom React hooks
│   ├── /types/             # TypeScript definitions
│   └── /config/            # Configuration
├── /backend/                # Golang backend
│   ├── /cmd/               # Main applications
│   │   └── /server/        # HTTP server
│   ├── /internal/          # Private application code
│   │   ├── /api/           # API handlers
│   │   ├── /auth/          # Authentication
│   │   ├── /chat/          # Chat logic
│   │   ├── /database/      # Database layer
│   │   ├── /middleware/    # HTTP middleware
│   │   ├── /models/        # Data models
│   │   ├── /services/      # Business logic
│   │   └── /utils/         # Utilities
│   ├── /pkg/               # Public library code
│   ├── /migrations/        # Database migrations (Knex)
│   ├── go.mod              # Go modules
│   ├── go.sum              # Go dependencies
│   └── main.go             # Entry point
├── /public/                 # Static assets
├── /scripts/                # Build & deployment scripts
├── .env.example             # Environment template
├── .env.local               # Local environment
├── package.json             # Frontend dependencies
├── next.config.mjs          # Next.js configuration
├── tailwind.config.js       # Tailwind CSS v4 config
├── tsconfig.json            # TypeScript configuration
└── knexfile.js              # Knex migration config
```

## ⚡ Next.js 15.4.4 Frontend Optimization
### Next.js Performance Optimization
1. **App Router**: Gunakan Next.js 15 App Router untuk optimal performance
2. **Turbopack**: Development dengan Turbopack untuk faster builds
3. **Server Components**: Maximize server components untuk reduced bundle size
4. **Streaming**: Implement streaming untuk faster page loads
5. **Image Optimization**: Next.js Image component dengan lazy loading
6. **Code Splitting**: Automatic code splitting dengan dynamic imports
7. **Bundle Analysis**: Regular bundle size monitoring

### Tailwind CSS v4 Optimization
1. **JIT Compilation**: Just-in-time CSS compilation
2. **Purge Unused**: Remove unused CSS classes
3. **CSS Variables**: Leverage CSS custom properties
4. **Component Classes**: Reusable component-based classes
5. **Dark Mode**: Efficient dark mode implementation
6. **Responsive Design**: Mobile-first responsive approach
7. **Performance Metrics**: CSS performance monitoring

### React 19 Performance
1. **Concurrent Features**: Use React 19 concurrent features
2. **Suspense**: Implement Suspense for better UX
3. **Server Actions**: Leverage server actions untuk form handling
4. **Memo & Callback**: Optimize re-renders dengan useMemo/useCallback
5. **Virtual Scrolling**: Implement untuk large lists
6. **Error Boundaries**: Comprehensive error handling
7. **Performance Profiling**: Regular React DevTools profiling

### Next.js Production Configuration
```typescript
// next.config.mjs - Production optimized
const nextConfig = {
  experimental: {
    turbo: {
      rules: {
        '*.svg': {
          loaders: ['@svgr/webpack'],
          as: '*.js'
        }
      }
    },
    serverComponentsExternalPackages: ['mysql2'],
    optimizePackageImports: ['lucide-react', '@radix-ui/react-icons']
  },
  images: {
    formats: ['image/avif', 'image/webp'],
    deviceSizes: [640, 750, 828, 1080, 1200, 1920, 2048, 3840],
    imageSizes: [16, 32, 48, 64, 96, 128, 256, 384]
  },
  compress: true,
  poweredByHeader: false,
  generateEtags: false,
  httpAgentOptions: {
    keepAlive: true
  },
  onDemandEntries: {
    maxInactiveAge: 25 * 1000,
    pagesBufferLength: 2
  },
  webpack: (config, { dev, isServer }) => {
    if (!dev && !isServer) {
      config.optimization.splitChunks.cacheGroups = {
        ...config.optimization.splitChunks.cacheGroups,
        vendor: {
          test: /[\\/]node_modules[\\/]/,
          name: 'vendors',
          chunks: 'all'
        }
      };
    }
    return config;
  }
};

export default nextConfig;
```

## 🚀 Golang 1.24.5 Backend Optimization

### High-Performance HTTP Server
```go
// main.go - Production optimized Gin server
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "github.com/gin-contrib/gzip"
    "github.com/gin-contrib/secure"
)

func main() {
    // Production mode
    gin.SetMode(gin.ReleaseMode)
    
    // Initialize router dengan middleware
    r := gin.New()
    
    // Security middleware
    r.Use(secure.New(secure.Config{
        SSLRedirect:          true,
        STSSeconds:           31536000,
        STSIncludeSubdomains: true,
        FrameDeny:            true,
        ContentTypeNosniff:   true,
        BrowserXssFilter:     true,
    }))
    
    // Performance middleware
    r.Use(gzip.Gzip(gzip.DefaultCompression))
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"https://lipdev.id"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"*"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))
    
    // Custom middleware untuk rate limiting
    r.Use(RateLimitMiddleware())
    r.Use(AuthMiddleware())
    
    // Routes
    setupRoutes(r)
    
    // Server configuration
    srv := &http.Server{
        Addr:         ":8080",
        Handler:      r,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    
    // Graceful shutdown
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed to start: %v", err)
        }
    }()
    
    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    // Graceful shutdown dengan timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }
}
```

### Golang Database Layer dengan MySQL
```go
// internal/database/mysql.go
package database

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
)

type DB struct {
    *sqlx.DB
}

type Config struct {
    Host         string
    Port         int
    User         string
    Password     string
    Database     string
    MaxOpenConns int
    MaxIdleConns int
    MaxLifetime  time.Duration
}

func NewConnection(cfg Config) (*DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
    
    db, err := sqlx.Connect("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }
    
    // Connection pool optimization
    db.SetMaxOpenConns(cfg.MaxOpenConns)   // Max 25 connections
    db.SetMaxIdleConns(cfg.MaxIdleConns)   // Max 10 idle connections
    db.SetConnMaxLifetime(cfg.MaxLifetime) // Max 5 minutes
    
    // Test connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    log.Println("Database connected successfully")
    return &DB{db}, nil
}

// Repository pattern dengan prepared statements
type UserRepository struct {
    db *DB
}

func NewUserRepository(db *DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(id string) (*User, error) {
    var user User
    query := `SELECT id, email, name, created_at, updated_at FROM users WHERE id = ? AND deleted_at IS NULL`
    
    err := r.db.Get(&user, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return &user, nil
}

func (r *UserRepository) Create(user *User) error {
    query := `INSERT INTO users (id, email, name, password_hash, created_at, updated_at) 
              VALUES (?, ?, ?, ?, NOW(), NOW())`
    
    _, err := r.db.Exec(query, user.ID, user.Email, user.Name, user.PasswordHash)
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    return nil
}
```

### Golang Performance Monitoring
```go
// internal/monitoring/metrics.go
package monitoring

import (
    "runtime"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
)

type Metrics struct {
    mu           sync.RWMutex
    Requests     int64
    Errors       int64
    ResponseTime []time.Duration
    StartTime    time.Time
}

func NewMetrics() *Metrics {
    return &Metrics{
        StartTime: time.Now(),
    }
}

func (m *Metrics) TrackRequest(duration time.Duration) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.Requests++
    m.ResponseTime = append(m.ResponseTime, duration)
    
    // Keep only last 1000 entries
    if len(m.ResponseTime) > 1000 {
        m.ResponseTime = m.ResponseTime[1:]
    }
}

func (m *Metrics) TrackError() {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.Errors++
}

func (m *Metrics) GetStats() map[string]interface{} {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    var avgResponseTime time.Duration
    if len(m.ResponseTime) > 0 {
        var total time.Duration
        for _, rt := range m.ResponseTime {
            total += rt
        }
        avgResponseTime = total / time.Duration(len(m.ResponseTime))
    }
    
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)
    
    return map[string]interface{}{
        "requests":         m.Requests,
        "errors":           m.Errors,
        "avg_response_ms":  avgResponseTime.Milliseconds(),
        "uptime_seconds":   time.Since(m.StartTime).Seconds(),
        "memory_alloc_mb":  memStats.Alloc / 1024 / 1024,
        "memory_sys_mb":    memStats.Sys / 1024 / 1024,
        "gc_runs":          memStats.NumGC,
        "goroutines":       runtime.NumGoroutine(),
    }
}

// Middleware untuk tracking metrics
func (m *Metrics) GinMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start)
        m.TrackRequest(duration)
        
        if c.Writer.Status() >= 400 {
            m.TrackError()
        }
    }
}
```

## 🤖 Telegram Bot Integration
- **Webhook Mode**: HTTPS webhook untuk production (NO polling!)
- **Message Types**: Text, photos, videos, documents, stickers, audio
- **Inline Keyboards**: Interactive buttons & menus
- **Admin Commands**: Slash commands hanya untuk admin (NO public commands!)
- **Group Management**: Admin controls & member management
- **File Handling**: Upload/download dengan stream processing
- **Error Recovery**: Graceful handling untuk API limits
- **Natural Chat**: Direct conversation tanpa commands untuk user biasa

## 🔗 Telegram Webhook Setup
```typescript
import TelegramBot from 'node-telegram-bot-api';
import { config } from './config/env.config.js';

// Webhook configuration dengan proper typing
const bot = new TelegramBot(config.telegram.botToken, {
  webHook: {
    port: config.server.port,
    host: '0.0.0.0'
  }
});

// Set webhook dengan secret token
bot.setWebHook(`${config.webhook.baseUrl}/webhook/${config.webhook.secret}`, {
  secret_token: config.webhook.secretToken
});
```

## 🤖 Vercel AI Gateway
- **Streaming**: Real-time AI responses
- **Model Switching**: Multiple AI providers
- **Context Management**: Conversation history
- **Rate Limiting**: Per-user quotas
- **Error Handling**: Fallback mechanisms

## 📝 Code Standards
- **Max 500 lines** per file
- **TypeScript strict mode** dengan proper typing
- **ESLint + Prettier** formatting
- **Async/await** only (no callbacks)
- **Error boundaries** everywhere
- **Logging** dengan structured format
- **Testing** unit + integration
- **Type definitions** untuk semua interfaces
- **No any types** - gunakan proper typing
- **Import/Export** ES modules only

## ⚡ TypeScript Performance Optimization
### Compiler Optimizations
1. **Incremental Compilation**: Enable `incremental: true` dan `tsBuildInfoFile`
2. **Project References**: Gunakan project references untuk large codebase
3. **Skip Lib Check**: `skipLibCheck: true` untuk faster compilation
4. **Module Resolution**: `moduleResolution: "bundler"` untuk optimal bundling
5. **Source Maps**: Production tanpa source maps, dev dengan inline
6. **Declaration Maps**: Hanya untuk library, skip untuk apps
7. **Strict Checks**: Enable semua strict flags untuk better optimization

### Runtime Performance
1. **Tree Shaking**: Import specific functions, bukan entire modules
2. **Lazy Loading**: Dynamic imports untuk code splitting
3. **Type Guards**: Efficient type checking dengan user-defined type guards
4. **Const Assertions**: Gunakan `as const` untuk literal types
5. **Enum Alternatives**: Prefer `const enum` atau union types
6. **Interface vs Type**: Prefer interfaces untuk object shapes
7. **Generic Constraints**: Specific constraints untuk better inference

### Memory Optimization
1. **Object Pooling**: Reuse objects untuk high-frequency operations
2. **WeakMap/WeakSet**: Untuk caching yang auto-cleanup
3. **Readonly Types**: Prevent accidental mutations
4. **Tuple Types**: Fixed-length arrays dengan specific types
5. **Branded Types**: Zero-runtime cost type safety
6. **Phantom Types**: Compile-time only type information
7. **Discriminated Unions**: Efficient pattern matching

### Build Optimization
```typescript
// tsconfig.json - Production optimized
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "ESNext",
    "moduleResolution": "bundler",
    "strict": true,
    "exactOptionalPropertyTypes": true,
    "noUncheckedIndexedAccess": true,
    "noImplicitOverride": true,
    "incremental": true,
    "tsBuildInfoFile": ".tsbuildinfo",
    "skipLibCheck": true,
    "declaration": false,
    "sourceMap": false,
    "removeComments": true,
    "importHelpers": true,
    "experimentalDecorators": false,
    "emitDecoratorMetadata": false
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist", "**/*.test.ts"]
}
```

### Production Build & Deployment
```typescript
// package.json - Build scripts optimization
{
  "scripts": {
    "build": "tsc --build --verbose",
    "build:prod": "NODE_ENV=production tsc --build && npm run optimize",
    "optimize": "terser dist/**/*.js --compress --mangle --output dist/",
    "start:prod": "NODE_ENV=production node --enable-source-maps dist/index.js",
    "start:cluster": "NODE_ENV=production pm2 start ecosystem.config.js"
  }
}

// PM2 ecosystem untuk production clustering
// ecosystem.config.js
module.exports = {
  apps: [{
    name: 'gryt-bot',
    script: './dist/index.js',
    instances: 'max', // Use all CPU cores
    exec_mode: 'cluster',
    env: {
      NODE_ENV: 'production',
      PORT: 3000
    },
    max_memory_restart: '1G',
    node_args: '--max-old-space-size=1024 --optimize-for-size',
    error_file: './logs/err.log',
    out_file: './logs/out.log',
    log_file: './logs/combined.log',
    time: true,
    autorestart: true,
    max_restarts: 10,
    min_uptime: '10s'
  }]
};

// Dockerfile multi-stage optimization
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production && npm cache clean --force
COPY . .
RUN npm run build:prod

FROM node:18-alpine AS runtime
RUN addgroup -g 1001 -S nodejs
RUN adduser -S gryt -u 1001
WORKDIR /app
COPY --from=builder --chown=gryt:nodejs /app/dist ./dist
COPY --from=builder --chown=gryt:nodejs /app/node_modules ./node_modules
COPY --from=builder --chown=gryt:nodejs /app/package.json ./package.json
USER gryt
EXPOSE 3000
CMD ["node", "--enable-source-maps", "--max-old-space-size=512", "dist/index.js"]
```

### Server Resource Monitoring
```typescript
// Auto-scaling resource monitor
class ProductionResourceMonitor {
  private readonly maxCpuUsage = 80; // 80%
  private readonly maxMemoryUsage = 85; // 85%
  private readonly checkInterval = 30000; // 30 seconds
  
  startMonitoring(): void {
    setInterval(async () => {
      const stats = await this.getSystemStats();
      
      if (stats.cpu > this.maxCpuUsage) {
        logger.warn('High CPU usage detected', { cpu: stats.cpu });
        await this.handleHighCpuUsage();
      }
      
      if (stats.memory > this.maxMemoryUsage) {
        logger.warn('High memory usage detected', { memory: stats.memory });
        await this.handleHighMemoryUsage();
      }
    }, this.checkInterval);
  }
  
  private async handleHighCpuUsage(): Promise<void> {
    // Reduce concurrent processing
    MessageQueue.setMaxConcurrent(5);
    RateLimiter.enableStrictMode();
  }
  
  private async handleHighMemoryUsage(): Promise<void> {
    if (global.gc) global.gc(); // Force GC
    await CacheManager.clearOldEntries();
    CacheManager.reduceCacheSize(0.5);
  }
}

// Adaptive message processing dengan auto-scaling
class AdaptiveMessageProcessor {
  private currentConcurrency = 10;
  private readonly minConcurrency = 5;
  private readonly maxConcurrency = 50;
  private responseTimeHistory: number[] = [];
  
  async processMessage(update: TelegramUpdate): Promise<void> {
    const startTime = Date.now();
    
    try {
      await this.handleUpdate(update);
    } finally {
      const duration = Date.now() - startTime;
      this.updateConcurrency(duration);
    }
  }
  
  private updateConcurrency(responseTime: number): void {
    this.responseTimeHistory.push(responseTime);
    
    if (this.responseTimeHistory.length > 100) {
      this.responseTimeHistory.shift();
    }
    
    const avgResponseTime = this.responseTimeHistory.reduce((a, b) => a + b) / this.responseTimeHistory.length;
    
    // Auto-adjust concurrency based on performance
    if (avgResponseTime < 1000 && this.currentConcurrency < this.maxConcurrency) {
      this.currentConcurrency = Math.min(this.currentConcurrency + 1, this.maxConcurrency);
    } else if (avgResponseTime > 3000 && this.currentConcurrency > this.minConcurrency) {
      this.currentConcurrency = Math.max(this.currentConcurrency - 1, this.minConcurrency);
    }
  }
}
```

### Code Patterns untuk Performance
```typescript
// ✅ GOOD - Efficient type guards
function isString(value: unknown): value is string {
  return typeof value === 'string';
}

// ✅ GOOD - Const assertions
const TELEGRAM_EVENTS = ['message', 'callback_query'] as const;
type TelegramEvent = typeof TELEGRAM_EVENTS[number];

// ✅ GOOD - Branded types (zero runtime cost)
type UserId = number & { readonly brand: unique symbol };
type ChatId = number & { readonly brand: unique symbol };

// ✅ GOOD - Discriminated unions
type APIResponse<T> = 
  | { success: true; data: T }
  | { success: false; error: string };

// ✅ GOOD - Generic constraints
interface Repository<T extends { id: string }> {
  findById(id: string): Promise<T | null>;
  save(entity: T): Promise<T>;
}

// ✅ GOOD - Lazy loading
const loadAIService = () => import('../services/ai.service.js');

// ✅ GOOD - Object pooling
class MessagePool {
  private pool: TelegramMessage[] = [];
  
  acquire(): TelegramMessage {
    return this.pool.pop() ?? this.createNew();
  }
  
  release(msg: TelegramMessage): void {
    this.reset(msg);
    this.pool.push(msg);
  }
}
```

### ESLint Rules untuk Performance
```typescript
// eslint.config.js
export default {
  rules: {
    '@typescript-eslint/prefer-readonly': 'error',
    '@typescript-eslint/prefer-readonly-parameter-types': 'warn',
    '@typescript-eslint/prefer-string-starts-ends-with': 'error',
    '@typescript-eslint/prefer-includes': 'error',
    '@typescript-eslint/prefer-for-of': 'error',
    '@typescript-eslint/prefer-optional-chain': 'error',
    '@typescript-eslint/prefer-nullish-coalescing': 'error',
    '@typescript-eslint/no-unnecessary-type-assertion': 'error',
    '@typescript-eslint/no-unnecessary-condition': 'error',
    'prefer-const': 'error',
    'no-var': 'error'
  }
};
```

## 🔧 Environment Variables Template
```bash
# Telegram Bot Configuration
TELEGRAM_BOT_TOKEN=your_bot_token_from_botfather
TELEGRAM_WEBHOOK_SECRET=your_webhook_secret_token
TELEGRAM_BOT_USERNAME=your_bot_username
TELEGRAM_ADMIN_USER_ID=your_admin_telegram_id
TELEGRAM_MAX_MESSAGES_PER_MINUTE=20
TELEGRAM_FILE_SIZE_LIMIT=50000000

# AI System Prompt Configuration
AI_SYSTEM_PROMPT="You are a helpful AI assistant. Be concise and accurate."
AI_DEFAULT_MODEL=anthropic/claude-sonnet-4
AI_MAX_TOKENS=4000
AI_TEMPERATURE=0.7

# Webhook Configuration
WEBHOOK_URL=https://lipdev.id/webhook
WEBHOOK_SECRET=your_webhook_secret
WEBHOOK_MAX_CONNECTIONS=100

# Vercel AI Gateway
VERCEL_AI_API_KEY=your_vercel_ai_key
VERCEL_AI_BASE_URL=https://gateway.ai.cloudflare.com

# Server Configuration
DOMAIN=https://lipdev.id
PORT=3000
NODE_ENV=production

# Security
JWT_SECRET=your_jwt_secret
ENCRYPTION_KEY=your_encryption_key
RATE_LIMIT_MAX=100
RATE_LIMIT_WINDOW=900000

# Database
DATABASE_URL=your_database_url
REDIS_URL=your_redis_url

# Monitoring
LOG_LEVEL=info
MONITORING_ENABLED=true
```

## 🚀 Deployment Flow
1. **Environment Setup**: Validate all env vars
2. **Nginx Config**: Update via sed commands
3. **SSL Certificate**: Auto-renewal setup
4. **Process Management**: PM2 or systemd
5. **Health Checks**: Monitoring endpoints
6. **Graceful Shutdown**: Clean process termination
7. **Webhook Setup**: Configure Telegram webhook dengan BotFather

## 🗄️ Database & Migration (MySQL + Knex)
### Database Security Rules
1. **PREPARED STATEMENTS**: Golang database/sql dengan prepared statements
2. **CONNECTION POOLING**: Min 5, Max 20 connections
3. **SSL/TLS**: Encrypted database connections
4. **INPUT VALIDATION**: Struct validation sebelum DB operations
5. **AUDIT LOGGING**: Track semua DB operations
6. **BACKUP STRATEGY**: Automated daily backups
7. **PRIVILEGE SEPARATION**: Minimal DB user permissions

### Knex Configuration (untuk migrasi)
```javascript
// knexfile.js
module.exports = {
  production: {
    client: 'mysql2',
    connection: {
      host: process.env.DB_HOST || 'localhost',
      port: process.env.DB_PORT || 3306,
      user: process.env.DB_USER,
      password: process.env.DB_PASSWORD,
      database: process.env.DB_NAME,
      ssl: process.env.DB_SSL === 'true' ? {
        rejectUnauthorized: false
      } : false,
      charset: 'utf8mb4',
      timezone: 'UTC'
    },
    pool: {
      min: 5,
      max: 20,
      acquireTimeoutMillis: 30000,
      createTimeoutMillis: 30000,
      destroyTimeoutMillis: 5000,
      idleTimeoutMillis: 30000,
      reapIntervalMillis: 1000,
      createRetryIntervalMillis: 100
    },
    migrations: {
      directory: './migrations',
      tableName: 'knex_migrations'
    },
    seeds: {
      directory: './seeds'
    }
  },
  development: {
    client: 'mysql2',
    connection: {
      host: 'localhost',
      port: 3306,
      user: process.env.DB_USER || 'root',
      password: process.env.DB_PASSWORD || '',
      database: process.env.DB_NAME || 'gryt_dev',
      charset: 'utf8mb4'
    },
    migrations: {
      directory: './migrations'
    },
    seeds: {
      directory: './seeds'
    }
  }
};
```

### Migration Example dengan Security
```javascript
// migrations/001_create_users_table.js
exports.up = function(knex) {
  return knex.schema.createTable('users', function(table) {
    table.uuid('id').primary().defaultTo(knex.raw('(UUID())'));
    table.string('email', 255).notNullable().unique();
    table.string('username', 100).nullable();
    table.string('first_name', 100).notNullable();
    table.string('last_name', 100).nullable();
    table.text('encrypted_data'); // Encrypted sensitive data
    table.enum('status', ['active', 'inactive', 'banned']).defaultTo('active');
    table.timestamp('created_at').defaultTo(knex.fn.now());
    table.timestamp('updated_at').defaultTo(knex.fn.now());
    table.timestamp('last_activity').nullable();
    
    // Indexes untuk performance
    table.index(['email']);
    table.index(['status', 'created_at']);
  });
};

exports.down = function(knex) {
  return knex.schema.dropTable('users');
};

// migrations/002_create_chat_sessions_table.js
exports.up = function(knex) {
  return knex.schema.createTable('chat_sessions', function(table) {
    table.uuid('id').primary().defaultTo(knex.raw('(UUID())'));
    table.uuid('user_id').notNullable();
    table.text('title').nullable();
    table.json('metadata').nullable();
    table.enum('status', ['active', 'archived', 'deleted']).defaultTo('active');
    table.timestamp('created_at').defaultTo(knex.fn.now());
    table.timestamp('updated_at').defaultTo(knex.fn.now());
    
    table.foreign('user_id').references('id').inTable('users').onDelete('CASCADE');
    table.index(['user_id', 'status']);
  });
};

exports.down = function(knex) {
  return knex.schema.dropTable('chat_sessions');
};
```

### Golang Database Service Pattern
```go
// backend/internal/repository/user.go
package repository

import (
    "database/sql"
    "time"
    
    "github.com/google/uuid"
    _ "github.com/go-sql-driver/mysql"
)

type User struct {
    ID           uuid.UUID  `json:"id" db:"id"`
    Email        string     `json:"email" db:"email"`
    Username     *string    `json:"username" db:"username"`
    FirstName    string     `json:"first_name" db:"first_name"`
    LastName     *string    `json:"last_name" db:"last_name"`
    Status       string     `json:"status" db:"status"`
    CreatedAt    time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
    LastActivity *time.Time `json:"last_activity" db:"last_activity"`
}

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(email string) (*User, error) {
    query := `
        SELECT id, email, username, first_name, last_name, status, 
               created_at, updated_at, last_activity
        FROM users 
        WHERE email = ? AND status = 'active'
    `
    
    var user User
    err := r.db.QueryRow(query, email).Scan(
        &user.ID, &user.Email, &user.Username, &user.FirstName,
        &user.LastName, &user.Status, &user.CreatedAt,
        &user.UpdatedAt, &user.LastActivity,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    
    return &user, nil
}

func (r *UserRepository) Create(user *User) error {
    query := `
        INSERT INTO users (id, email, username, first_name, last_name, status)
        VALUES (?, ?, ?, ?, ?, ?)
    `
    
    user.ID = uuid.New()
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    
    _, err := r.db.Exec(query, user.ID, user.Email, user.Username,
        user.FirstName, user.LastName, user.Status)
    
    return err
}

func (r *UserRepository) UpdateLastActivity(userID uuid.UUID) error {
    query := `UPDATE users SET last_activity = ?, updated_at = ? WHERE id = ?`
    now := time.Now()
    _, err := r.db.Exec(query, now, now, userID)
    return err
}
```

### Environment Variables (Database)
```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=gryt_user
DB_PASSWORD=your_secure_password
DB_NAME=gryt_chat
DB_SSL=false
DB_CONNECTION_LIMIT=20
DB_MAX_IDLE_CONNS=5
DB_MAX_OPEN_CONNS=20
DB_CONN_MAX_LIFETIME=1h

# Encryption
DB_ENCRYPTION_KEY=your_32_char_encryption_key
JWT_SECRET=your_jwt_secret_key
BCRYPT_COST=12
```

## 🤖 AI Chat Integration
### Security Guidelines (ULTRA KETAT)
1. **API RATE LIMITING**: Batasi request AI per user per menit
2. **INPUT SANITIZATION**: Validasi dan sanitasi semua input user
3. **CONTENT FILTERING**: Filter konten berbahaya atau tidak pantas
4. **SESSION MANAGEMENT**: Secure session handling dengan JWT
5. **AUDIT LOGGING**: Log semua aktivitas AI chat
6. **RESPONSE VALIDATION**: Validasi response AI sebelum dikirim ke user
7. **TIMEOUT HANDLING**: Maksimal 30 detik per AI request

### AI Chat Service (Golang)
```go
// backend/internal/service/chat.go
package service

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/google/uuid"
)

type ChatMessage struct {
    ID        uuid.UUID `json:"id"`
    SessionID uuid.UUID `json:"session_id"`
    UserID    uuid.UUID `json:"user_id"`
    Role      string    `json:"role"` // "user" or "assistant"
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}

type ChatSession struct {
    ID        uuid.UUID `json:"id"`
    UserID    uuid.UUID `json:"user_id"`
    Title     string    `json:"title"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type ChatService struct {
    chatRepo    ChatRepository
    aiProvider  AIProvider
    rateLimiter *RateLimiter
    validator   *MessageValidator
}

func NewChatService(chatRepo ChatRepository, aiProvider AIProvider) *ChatService {
    return &ChatService{
        chatRepo:    chatRepo,
        aiProvider:  aiProvider,
        rateLimiter: NewRateLimiter(),
        validator:   NewMessageValidator(),
    }
}

func (s *ChatService) SendMessage(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, content string) (*ChatMessage, error) {
    // Rate limiting check
    if !s.rateLimiter.Allow(userID) {
        return nil, fmt.Errorf("rate limit exceeded")
    }
    
    // Input validation
    if err := s.validator.ValidateMessage(content, userID); err != nil {
        return nil, fmt.Errorf("invalid message: %w", err)
    }
    
    // Sanitize input
    sanitizedContent := SanitizeInput(content)
    
    // Save user message
    userMessage := &ChatMessage{
        ID:        uuid.New(),
        SessionID: sessionID,
        UserID:    userID,
        Role:      "user",
        Content:   sanitizedContent,
        CreatedAt: time.Now(),
    }
    
    if err := s.chatRepo.SaveMessage(ctx, userMessage); err != nil {
        return nil, fmt.Errorf("failed to save user message: %w", err)
    }
    
    // Get AI response
    aiResponse, err := s.aiProvider.GenerateResponse(ctx, sessionID, sanitizedContent)
    if err != nil {
        return nil, fmt.Errorf("failed to generate AI response: %w", err)
    }
    
    // Save AI response
    assistantMessage := &ChatMessage{
        ID:        uuid.New(),
        SessionID: sessionID,
        UserID:    userID,
        Role:      "assistant",
        Content:   aiResponse,
        CreatedAt: time.Now(),
    }
    
    if err := s.chatRepo.SaveMessage(ctx, assistantMessage); err != nil {
        return nil, fmt.Errorf("failed to save AI message: %w", err)
    }
    
    return assistantMessage, nil
}
```

## 🔧 Environment Variables Template
```bash
# Server Configuration
PORT=8080
HOST=0.0.0.0
ENV=production
DOMAIN=https://yourdomain.com

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=gryt_user
DB_PASSWORD=your_secure_password
DB_NAME=gryt_chat
DB_SSL=false
DB_CONNECTION_LIMIT=20
DB_MAX_IDLE_CONNS=5
DB_MAX_OPEN_CONNS=20
DB_CONN_MAX_LIFETIME=1h

# JWT & Security
JWT_SECRET=your_jwt_secret_key_32_chars_min
JWT_EXPIRY=24h
BCRYPT_COST=12
CORS_ORIGINS=http://localhost:3000,https://yourdomain.com

# AI Configuration
AI_PROVIDER=openai
AI_API_KEY=your_ai_api_key
AI_MODEL=gpt-4
AI_MAX_TOKENS=4000
AI_TEMPERATURE=0.7
AI_TIMEOUT=30s
AI_RATE_LIMIT=30

# Redis (untuk caching dan session)
REDIS_URL=redis://localhost:6379
REDIS_PASSWORD=your_redis_password
REDIS_DB=0

# Monitoring & Logging
LOG_LEVEL=info
LOG_FORMAT=json
MONITORING_ENABLED=true
METRICS_PORT=9090

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=15m
RATE_LIMIT_BURST=20

# File Upload
MAX_FILE_SIZE=10MB
UPLOAD_PATH=/tmp/uploads
ALLOWED_FILE_TYPES=jpg,jpeg,png,gif,pdf,txt,md
```

## 🚀 Deployment & Production
### Docker Configuration
```dockerfile
# Dockerfile untuk Golang backend
FROM golang:1.24.5-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
CMD ["./main"]
```

### Next.js Production Build
```dockerfile
# Dockerfile untuk Next.js frontend
FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

COPY . .
RUN npm run build

FROM node:18-alpine AS runner
WORKDIR /app

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs
EXPOSE 3000
ENV PORT 3000

CMD ["node", "server.js"]
```

## 💬 Development Guidelines
- **Code Quality**: Selalu gunakan linter dan formatter
- **Security First**: Validasi input, sanitasi output, rate limiting
- **Performance**: Optimasi query, caching, lazy loading
- **Modularity**: Pisahkan concern, reusable components
- **Testing**: Unit test, integration test, e2e test
- **Documentation**: Comment code, update README
- **Git Flow**: Feature branch, PR review, semantic commits

## 🔒 Security Checklist
- ✅ Input validation & sanitization
- ✅ SQL injection prevention (prepared statements)
- ✅ XSS protection (CSP headers)
- ✅ CSRF protection
- ✅ Rate limiting
- ✅ Authentication & authorization
- ✅ Secure headers
- ✅ HTTPS enforcement
- ✅ Environment variables untuk secrets
- ✅ Error handling tanpa info leak
- ✅ Audit logging
- ✅ Regular security updates