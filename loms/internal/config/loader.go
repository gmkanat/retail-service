package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	"os"
)

func LoadStocks(fileName string) ([]model.Stock, error) {
	filePath, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read stock file: %w", err)
	}

	var stocks []model.Stock
	if err := json.Unmarshal(filePath, &stocks); err != nil {
		return nil, fmt.Errorf("failed to parse stock JSON data: %w", err)
	}

	if len(stocks) == 0 {
		return nil, errors.New("no stock data found in the file")
	}

	return stocks, nil
}
