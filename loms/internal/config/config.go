package config

import (
	"os"
)

type AppConfig struct {
	GRPCPort    string `json:"grpc_port"`
	HTTPPort    string `json:"http_port"`
	StockFile   string `json:"stock_file"`
	SwaggerFile string `json:"swagger_file"`
	DataBaseURL string `json:"database_url"`
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
	swaggerFile := os.Getenv("LOMS_SWAGGER_FILE")
	if swaggerFile == "" {
		swaggerFile = "../api/openapiv2/loms.swagger.json"
	}

	dataBaseURL := os.Getenv("LOMS_DATABASE_URL")
	if dataBaseURL == "" {
		dataBaseURL = "postgres://user:password@localhost:5432/route256"
	}

	return &AppConfig{
		GRPCPort:    grpcPort,
		HTTPPort:    httpPort,
		StockFile:   stockFile,
		SwaggerFile: swaggerFile,
		DataBaseURL: dataBaseURL,
	}
}
