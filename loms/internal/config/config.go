package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type AppConfig struct {
	GRPCPort           string
	HTTPPort           string
	StockFile          string
	SwaggerFile        string
	MasterDBURL        string
	ReplicaDBURL       string
	Brokers            []string
	Topic              string
	Tick               time.Duration
	BatchSize          int
	OutboxPollInterval time.Duration
}

func Load() *AppConfig {
	grpcPort := getEnv("LOMS_GRPC_PORT")
	httpPort := getEnv("LOMS_HTTP_PORT")
	stockFile := getEnv("LOMS_STOCK_FILE")
	swaggerFile := getEnv("LOMS_SWAGGER_FILE")
	masterDBURL := getEnv("LOMS_MASTER_DB_URL")
	replicaDBURL := getEnv("LOMS_SLAVE_DB_URL")
	brokers := []string{getEnv("KAFKA_BROKERS")}
	topic := getEnv("KAFKA_TOPIC")
	tick := parseDuration(getEnv("DISPATCHER_TICK"))
	batchSize := parseInt(getEnv("DISPATCHER_BATCH_SIZE"))
	outboxPollInterval := parseDuration(getEnv("LOMS_OUTBOX_POLL_INTERVAL"))

	return &AppConfig{
		GRPCPort:           grpcPort,
		HTTPPort:           httpPort,
		StockFile:          stockFile,
		SwaggerFile:        swaggerFile,
		MasterDBURL:        masterDBURL,
		ReplicaDBURL:       replicaDBURL,
		Brokers:            brokers,
		Topic:              topic,
		Tick:               tick,
		BatchSize:          batchSize,
		OutboxPollInterval: outboxPollInterval,
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}

func parseDuration(value string) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("Invalid duration format for %s: %v", value, err)
	}
	return duration
}

func parseInt(value string) int {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Invalid integer format for %s: %v", value, err)
	}
	return parsed
}
