package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	BaseURL        string
	Token          string
	Port           string
	MaxRetries     int
	InitialBackoff time.Duration
}

func Load() *Config {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		log.Fatalf("BASE_URL not set")
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatalf("TOKEN not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT not set")
	}

	maxRetriesStr := os.Getenv("MAX_RETRIES")
	maxRetries, err := strconv.Atoi(maxRetriesStr)
	if err != nil {
		log.Fatalf("Invalid MAX_RETRIES value: %v", err)
	}

	initialBackoffStr := os.Getenv("INITIAL_BACKOFF")
	initialBackoff, err := time.ParseDuration(initialBackoffStr)
	if err != nil {
		log.Fatalf("Invalid INITIAL_BACKOFF value: %v", err)
	}

	return &Config{
		BaseURL:        baseURL,
		Token:          token,
		Port:           port,
		MaxRetries:     maxRetries,
		InitialBackoff: initialBackoff,
	}
}
