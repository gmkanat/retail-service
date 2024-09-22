package config

import (
	"os"
)

type AppConfig struct {
	GRPCPort  string `json:"grpc_port"`
	HTTPPort  string `json:"http_port"`
	StockFile string `json:"stock_file"`
}

func Load() *AppConfig {
	grpcPort := os.Getenv("LOMS_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = ":50051"
	}

	httpPort := os.Getenv("LOMS_HTTP_PORT")
	if httpPort == "" {
		httpPort = ":8082"
	}

	stockFile := os.Getenv("LOMS_STOCK_FILE")
	if stockFile == "" {
		stockFile = "internal/config/stock-data.json"
	}

	return &AppConfig{
		GRPCPort:  grpcPort,
		HTTPPort:  httpPort,
		StockFile: stockFile,
	}
}
