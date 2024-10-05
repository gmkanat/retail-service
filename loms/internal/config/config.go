package config

import (
	"os"
)

type AppConfig struct {
	GRPCPort     string `json:"grpc_port"`
	HTTPPort     string `json:"http_port"`
	StockFile    string `json:"stock_file"`
	SwaggerFile  string `json:"swagger_file"`
	MasterDBURL  string `json:"master_db_url"`
	ReplicaDBURL string `json:"replica_db_url"`
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

	masterDBURL := os.Getenv("LOMS_MASTER_DB_URL")
	if masterDBURL == "" {
		masterDBURL = "postgres://loms_user:loms_password@pg_master:5432/loms_db"
	}

	replicaDBURL := os.Getenv("LOMS_REPLICA_DB_URL")
	if replicaDBURL == "" {
		replicaDBURL = "postgres://loms_user:loms_password@pg_slave:5432/loms_db"
	}

	return &AppConfig{
		GRPCPort:     grpcPort,
		HTTPPort:     httpPort,
		StockFile:    stockFile,
		SwaggerFile:  swaggerFile,
		MasterDBURL:  masterDBURL,
		ReplicaDBURL: replicaDBURL,
	}
}
