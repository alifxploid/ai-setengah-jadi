package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"gryt-backend/internal/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Logger middleware untuk logging requests
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// Recovery middleware untuk panic recovery
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("Panic recovered: %v", recovered)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	})
}

// Security middleware untuk security headers
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}

// RateLimiter untuk rate limiting per IP
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}
}

func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = limiter
	}

	return limiter
}

func (rl *RateLimiter) Allow(ip string) bool {
	return rl.GetLimiter(ip).Allow()
}

// RateLimit middleware
func RateLimit(rateLimitService *services.RateLimitService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		log.Printf("[DEBUG] RateLimit check for IP: %s, Path: %s", ip, c.Request.URL.Path)
		if !rateLimitService.Allow(ip) {
			log.Printf("[DEBUG] Rate limit exceeded for IP: %s", ip)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"message": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}

		log.Printf("[DEBUG] Rate limit OK for IP: %s", ip)
		c.Next()
	}
}

// Auth middleware untuk authentication
func Auth(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth untuk public endpoints
		if isPublicEndpoint(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Try X-Access-Key header first
		accessKey := c.GetHeader("X-Access-Key")
		if accessKey != "" {
			log.Printf("[DEBUG] Auth middleware: using X-Access-Key: %s", accessKey)
			// Validate access key
			user, err := authService.ValidateAccessKey(accessKey)
			if err != nil || user == nil {
				log.Printf("[DEBUG] Auth middleware: X-Access-Key validation failed: %v", err)
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid access key",
					"message": "Please check your access key and try again",
				})
				c.Abort()
				return
			}

			// Set user info dalam context
			c.Set("user_id", user.ID)
			c.Set("user", user)
			c.Next()
			return
		}

		// Fallback to Bearer token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header or X-Access-Key required",
			})
			c.Abort()
			return
		}

		// Parse Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Validate token
		userID, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user ID dalam context
		c.Set("user_id", userID)
		c.Next()
	}
}

// AccessKey middleware untuk validasi access key
func AccessKey(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			AccessKey string `json:"access_key" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format",
				"message": "access_key is required",
			})
			c.Abort()
			return
		}

		// Validate access key
		user, err := authService.ValidateAccessKey(req.AccessKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid access key",
				"message": "Please check your access key and try again",
			})
			c.Abort()
			return
		}

		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Access key not found",
				"message": "Invalid or inactive access key",
			})
			c.Abort()
			return
		}

		// Set user dalam context
		c.Set("user", user)
		c.Next()
	}
}

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Allow specific origins
		allowedOrigins := []string{
			"https://lipdev.id",
			"http://localhost:3000", // untuk development
		}

		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				c.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Helper function untuk check public endpoints
func isPublicEndpoint(path string) bool {
	publicPaths := []string{
		"/health",
		"/api/auth/login",
		"/api/auth/validate-key",
	}

	for _, publicPath := range publicPaths {
		if strings.HasPrefix(path, publicPath) {
			return true
		}
	}

	return false
}

// Input sanitization middleware
func SanitizeInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sanitize query parameters
		for key, values := range c.Request.URL.Query() {
			for i, value := range values {
				c.Request.URL.Query()[key][i] = sanitizeString(value)
			}
		}

		c.Next()
	}
}

// Helper function untuk sanitize string
func sanitizeString(input string) string {
	// Remove dangerous characters
	input = strings.ReplaceAll(input, "<", "")
	input = strings.ReplaceAll(input, ">", "")
	input = strings.ReplaceAll(input, "'", "")
	input = strings.ReplaceAll(input, "\"", "")
	input = strings.ReplaceAll(input, "&", "")
	
	// Trim whitespace
	input = strings.TrimSpace(input)
	
	// Limit length
	if len(input) > 4096 {
		input = input[:4096]
	}
	
	return input
}

// AI-specific rate limiting
var (
	aiRateLimiter *RateLimiter
	aiRateLimiterOnce sync.Once
)

// InitAIRateLimit initializes AI rate limiter
func InitAIRateLimit(requestsPerMinute int) {
	aiRateLimiterOnce.Do(func() {
		aiRateLimiter = NewRateLimiter(rate.Limit(requestsPerMinute), requestsPerMinute*2)
	})
}

// AIRateLimit middleware for AI endpoints
func AIRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if aiRateLimiter == nil {
			// Rate limiter not initialized, allow request
			c.Next()
			return
		}

		// Get client identifier (prefer user ID over IP)
		clientID := c.ClientIP()
		if userID, exists := c.Get("user_id"); exists {
			if uid, ok := userID.(string); ok {
				clientID = "user:" + uid
			}
		}

		if !aiRateLimiter.Allow(clientID) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "AI rate limit exceeded",
				"message": "Too many AI requests, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}