package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName  string
	AppEnv   string
	Port     string
	LogLevel string

	DatabaseURL      string
	DBMaxConnections int
	DBMaxIdleTime    time.Duration

	RedisURL      string
	RedisPassword string
	RedisDB       int

	ApolloAPIKey     string
	ClearbitAPIKey   string
	GoogleMapsAPIKey string

	JWTSecret     string
	JWTExpiration time.Duration

	RateLimitRPS   int // Requests per second
	RateLimitBurst int

	CacheTTL time.Duration

	CORSAllowedOrigins []string

	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
	HTTPIdleTimeout  time.Duration
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		AppName:  getEnv("APP_NAME", "sniper-api"),
		AppEnv:   getEnv("APP_ENV", "development"),
		Port:     getEnv("PORT", "8000"),
		LogLevel: getEnv("LOG_LEVEL", "info"),

		DatabaseURL:      getEnv("DATABASE_URL", ""),
		DBMaxConnections: getEnvAsInt("DB_MAX_CONNECTIONS", 25),
		DBMaxIdleTime:    getEnvAsDuration("DB_MAX_IDLE_TIME", 5*time.Minute),

		RedisURL:      getEnv("REDIS_URL", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		ApolloAPIKey:     getEnv("APOLLO_API_KEY", ""),
		ClearbitAPIKey:   getEnv("CLEARBIT_API_KEY", ""),
		GoogleMapsAPIKey: getEnv("GOOGLE_MAPS_API_KEY", ""),

		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-me"),
		JWTExpiration: getEnvAsDuration("JWT_EXPIRATION", 24*time.Hour),

		RateLimitRPS:   getEnvAsInt("RATE_LIMIT_RPS", 10),
		RateLimitBurst: getEnvAsInt("RATE_LIMIT_BURST", 20),

		CacheTTL: getEnvAsDuration("CACHE_TTL", 10*time.Minute),

		CORSAllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:5174"}),

		HTTPReadTimeout:  getEnvAsDuration("HTTP_READ_TIMEOUT", 10*time.Second),
		HTTPWriteTimeout: getEnvAsDuration("HTTP_WRITE_TIMEOUT", 10*time.Second),
		HTTPIdleTimeout:  getEnvAsDuration("HTTP_IDLE_TIMEOUT", 120*time.Second),
	}

	cfg.validate()

	return cfg
}

func (c *Config) validate() {
	if c.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	if c.AppEnv == "production" {
		if c.JWTSecret == "your-secret-key-change-me" {
			log.Fatal("JWT_SECRET must be set in production")
		}
		if c.ApolloAPIKey == "" {
			log.Fatal("APOLLO_API_KEY is required in production")
		}
	}
}

func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development" || c.AppEnv == "dev"
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production" || c.AppEnv == "prod"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: invalid value for %s, using default: %d", key, defaultValue)
		return defaultValue
	}

	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		log.Printf("Warning: invalid duration for %s, using default: %s", key, defaultValue)
		return defaultValue
	}

	return value
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	// Exemplo: "http://localhost:3000,https://app.com"
	return splitAndTrim(valueStr, ",")
}

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
