package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	BaseURL        string
	LomsAddr       string
	Token          string
	PortAddr       string
	MaxRetries     int
	InitialBackoff time.Duration
	RateLimit      int
	BurstLimit     int
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

	lomsAddr := os.Getenv("LOMS_ADDR")
	if lomsAddr == "" {
		log.Fatalf("LOMS_ADDR not set")
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

	rateLimitStr := os.Getenv("RATE_LIMIT")
	rateLimit, err := strconv.Atoi(rateLimitStr)
	if err != nil {
		log.Fatalf("Invalid RATE_LIMIT value: %v", err)
	}

	burstLimitStr := os.Getenv("BURST_LIMIT")
	burstLimit, err := strconv.Atoi(burstLimitStr)
	if err != nil {
		log.Fatalf("Invalid BURST_LIMIT value: %v", err)
	}

	return &Config{
		BaseURL:        baseURL,
		Token:          token,
		PortAddr:       port,
		MaxRetries:     maxRetries,
		InitialBackoff: initialBackoff,
		LomsAddr:       lomsAddr,
		RateLimit:      rateLimit,
		BurstLimit:     burstLimit,
	}
}
