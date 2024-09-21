package config

import (
	"os"
)

type AppConfig struct {
	Port      string `json:"port"`
	StockFile string `json:"stock_file"`
}

func Load() *AppConfig {
	port := os.Getenv("LOMS_PORT")
	if port == "" {
		port = ":50051"
	}

	stockFile := os.Getenv("LOMS_STOCK_FILE")
	if stockFile == "" {
		stockFile = "internal/config/stock-data.json"
	}

	return &AppConfig{
		Port:      port,
		StockFile: stockFile,
	}
}
