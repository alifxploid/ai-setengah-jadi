package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Limits   LimitsConfig
	AI       AIConfig
}

type ServerConfig struct {
	Port        string
	Environment string
	FrontendURL string
	Domain      string
}

type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	Database     string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

type AuthConfig struct {
	JWTSecret     string
	AccessKey     string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
}

type LimitsConfig struct {
	ChatTokensPerUser   int
	SearchTokensPerUser int
	RateLimitPerMinute  int
}

type AIConfig struct {
	// Vercel AI Gateway Configuration
	APIKey      string
	BaseURL     string
	Model       string
	MaxTokens   int
	Temperature float64
	Timeout     time.Duration
	
	// System Prompt Configuration
	SystemPrompt string
	
	// Advanced Parameters
	TopP             float64
	FrequencyPenalty float64
	PresencePenalty  float64
	Stop             []string
	
	// Model Configurations
	Models      map[string]ModelConfig
	
	// Tool Configurations
	Tools       ToolsConfig
	
	// Rate Limiting
	RateLimit   RateLimitConfig
	
	// Streaming
	Streaming   StreamingConfig
}

type ModelConfig struct {
	Name        string
	MaxTokens   int
	Temperature float64
	Enabled     bool
}

type ToolsConfig struct {
	WebSearch   ToolConfig
	Calculator  ToolConfig
	Weather     ToolConfig
	Translation ToolConfig
	ImageAnalysis ToolConfig
}

type ToolConfig struct {
	Enabled bool
	APIKey  string
	BaseURL string
}

type RateLimitConfig struct {
	RequestsPerMinute int
	RequestsPerHour   int
	RequestsPerDay    int
}

type StreamingConfig struct {
	Enabled     bool
	BufferSize  int
	Timeout     time.Duration
}

func Load() (*Config, error) {
	// Validate required environment variables
	requiredEnvs := []string{
		"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME",
		"JWT_SECRET", "ACCESS_KEY", "AI_API_KEY",
	}

	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", env)
		}
	}

	// Parse database port
	dbPort, err := strconv.Atoi(getEnvOrDefault("DB_PORT", "3306"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %v", err)
	}

	// Parse connection limits
	maxOpenConns, _ := strconv.Atoi(getEnvOrDefault("DB_MAX_OPEN_CONNS", "25"))
	maxIdleConns, _ := strconv.Atoi(getEnvOrDefault("DB_MAX_IDLE_CONNS", "10"))
	maxLifetime, _ := time.ParseDuration(getEnvOrDefault("DB_MAX_LIFETIME", "5m"))

	// Parse token limits
	chatTokens, _ := strconv.Atoi(getEnvOrDefault("CHAT_TOKENS_PER_USER", "10"))
	searchTokens, _ := strconv.Atoi(getEnvOrDefault("SEARCH_TOKENS_PER_USER", "100"))
	rateLimitPerMin, _ := strconv.Atoi(getEnvOrDefault("RATE_LIMIT_PER_MINUTE", "30"))

	// Parse AI config
	maxTokens, _ := strconv.Atoi(getEnvOrDefault("AI_MAX_TOKENS", "4096"))
	temperature, _ := strconv.ParseFloat(getEnvOrDefault("AI_TEMPERATURE", "0.7"), 64)
	aiTimeout, _ := time.ParseDuration(getEnvOrDefault("AI_TIMEOUT", "30s"))
	
	// Parse advanced AI parameters
	topP, _ := strconv.ParseFloat(getEnvOrDefault("AI_TOP_P", "1.0"), 64)
	frequencyPenalty, _ := strconv.ParseFloat(getEnvOrDefault("AI_FREQUENCY_PENALTY", "0.0"), 64)
	presencePenalty, _ := strconv.ParseFloat(getEnvOrDefault("AI_PRESENCE_PENALTY", "0.0"), 64)
	stopSequences := strings.Split(getEnvOrDefault("AI_STOP_SEQUENCES", ""), ",")
	if len(stopSequences) == 1 && stopSequences[0] == "" {
		stopSequences = []string{}
	}
	
	// Parse rate limiting
	rpmLimit, _ := strconv.Atoi(getEnvOrDefault("AI_RATE_LIMIT_RPM", "60"))
	rphLimit, _ := strconv.Atoi(getEnvOrDefault("AI_RATE_LIMIT_RPH", "1000"))
	rpdLimit, _ := strconv.Atoi(getEnvOrDefault("AI_RATE_LIMIT_RPD", "10000"))
	
	// Parse streaming config
	streamingEnabled, _ := strconv.ParseBool(getEnvOrDefault("AI_STREAMING_ENABLED", "true"))
	streamingBufferSize, _ := strconv.Atoi(getEnvOrDefault("AI_STREAMING_BUFFER_SIZE", "1024"))
	streamingTimeout, _ := time.ParseDuration(getEnvOrDefault("AI_STREAMING_TIMEOUT", "60s"))

	// Parse JWT expiry
	tokenExpiry, _ := time.ParseDuration(getEnvOrDefault("JWT_TOKEN_EXPIRY", "24h"))
	refreshExpiry, _ := time.ParseDuration(getEnvOrDefault("JWT_REFRESH_EXPIRY", "168h"))

	return &Config{
		Server: ServerConfig{
			Port:        getEnvOrDefault("PORT", "8080"),
			Environment: getEnvOrDefault("ENVIRONMENT", "development"),
			FrontendURL: getEnvOrDefault("FRONTEND_URL", "https://lipdev.id"),
			Domain:      getEnvOrDefault("DOMAIN", "lipdev.id"),
		},
		Database: DatabaseConfig{
			Host:         os.Getenv("DB_HOST"),
			Port:         dbPort,
			User:         os.Getenv("DB_USER"),
			Password:     os.Getenv("DB_PASSWORD"),
			Database:     os.Getenv("DB_NAME"),
			MaxOpenConns: maxOpenConns,
			MaxIdleConns: maxIdleConns,
			MaxLifetime:  maxLifetime,
		},
		Auth: AuthConfig{
			JWTSecret:     os.Getenv("JWT_SECRET"),
			AccessKey:     os.Getenv("ACCESS_KEY"),
			TokenExpiry:   tokenExpiry,
			RefreshExpiry: refreshExpiry,
		},
		Limits: LimitsConfig{
			ChatTokensPerUser:   chatTokens,
			SearchTokensPerUser: searchTokens,
			RateLimitPerMinute:  rateLimitPerMin,
		},
		AI: AIConfig{
			APIKey:      os.Getenv("AI_API_KEY"),
			BaseURL:     getEnvOrDefault("AI_BASE_URL", "https://ai-gateway.vercel.sh/v1"),
			Model:       getEnvOrDefault("AI_MODEL", "anthropic/claude-sonnet-4"),
			MaxTokens:   maxTokens,
			Temperature: temperature,
			Timeout:     aiTimeout,
			SystemPrompt: getEnvOrDefault("AI_SYSTEM_PROMPT", "You are GRYT, an advanced AI assistant powered by Vercel AI Gateway. You are designed to be helpful, accurate, and efficient. You have access to various tools and can handle multiple types of content including text, images, and documents. Always provide clear, concise, and helpful responses while maintaining a professional yet friendly tone."),
			TopP:             topP,
			FrequencyPenalty: frequencyPenalty,
			PresencePenalty:  presencePenalty,
			Stop:             stopSequences,
			Models: map[string]ModelConfig{
				// OpenAI Models
				"openai/gpt-4o": {
					Name:        "openai/gpt-4o",
					MaxTokens:   4096,
					Temperature: 0.7,
					Enabled:     true,
				},
				"openai/gpt-4o-mini": {
					Name:        "openai/gpt-4o-mini",
					MaxTokens:   2048,
					Temperature: 0.7,
					Enabled:     true,
				},
				"openai/gpt-4-turbo": {
					Name:        "openai/gpt-4-turbo",
					MaxTokens:   4096,
					Temperature: 0.7,
					Enabled:     true,
				},
				// Anthropic Models
				"anthropic/claude-sonnet-4": {
					Name:        "anthropic/claude-sonnet-4",
					MaxTokens:   4096,
					Temperature: 0.7,
					Enabled:     true,
				},
				"anthropic/claude-3-5-sonnet-20241022": {
					Name:        "anthropic/claude-3-5-sonnet-20241022",
					MaxTokens:   4096,
					Temperature: 0.7,
					Enabled:     true,
				},
				"anthropic/claude-3-haiku-20240307": {
					Name:        "anthropic/claude-3-haiku-20240307",
					MaxTokens:   2048,
					Temperature: 0.7,
					Enabled:     true,
				},
				// Google Models
				"google/gemini-1.5-pro": {
					Name:        "google/gemini-1.5-pro",
					MaxTokens:   4096,
					Temperature: 0.7,
					Enabled:     true,
				},
				"google/gemini-1.5-flash": {
					Name:        "google/gemini-1.5-flash",
					MaxTokens:   2048,
					Temperature: 0.7,
					Enabled:     true,
				},
			},
			Tools: ToolsConfig{
				WebSearch: ToolConfig{
					Enabled: getEnvOrDefault("AI_TOOLS_WEB_SEARCH_ENABLED", "true") == "true",
					APIKey:  getEnvOrDefault("SEARCH_API_KEY", ""),
					BaseURL: getEnvOrDefault("SEARCH_API_URL", ""),
				},
				Calculator: ToolConfig{
					Enabled: getEnvOrDefault("AI_TOOLS_CALCULATOR_ENABLED", "true") == "true",
				},
				Weather: ToolConfig{
					Enabled: getEnvOrDefault("AI_TOOLS_WEATHER_ENABLED", "true") == "true",
					APIKey:  getEnvOrDefault("WEATHER_API_KEY", ""),
					BaseURL: getEnvOrDefault("WEATHER_API_URL", "https://api.openweathermap.org/data/2.5"),
				},
				Translation: ToolConfig{
					Enabled: getEnvOrDefault("AI_TOOLS_TRANSLATION_ENABLED", "true") == "true",
					APIKey:  getEnvOrDefault("TRANSLATION_API_KEY", ""),
					BaseURL: getEnvOrDefault("TRANSLATION_API_URL", ""),
				},
				ImageAnalysis: ToolConfig{
					Enabled: getEnvOrDefault("AI_TOOLS_IMAGE_ANALYSIS_ENABLED", "true") == "true",
				},
			},
			RateLimit: RateLimitConfig{
				RequestsPerMinute: rpmLimit,
				RequestsPerHour:   rphLimit,
				RequestsPerDay:    rpdLimit,
			},
			Streaming: StreamingConfig{
				Enabled:     streamingEnabled,
				BufferSize:  streamingBufferSize,
				Timeout:     streamingTimeout,
			},
		},
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}