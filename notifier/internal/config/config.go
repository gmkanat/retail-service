package config

import (
	"log"
	"os"
)

type Config struct {
	KafkaBrokers []string
	KafkaTopic   string
}

func Load() *Config {
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		log.Fatal("KAFKA_BROKERS is not set")
	}

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		log.Fatal("KAFKA_TOPIC is not set")
	}

	return &Config{
		KafkaBrokers: []string{kafkaBrokers},
		KafkaTopic:   kafkaTopic,
	}
}
