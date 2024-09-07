package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/customerrors"
	"net/http"
)

func (s *ProductService) GetProduct(ctx context.Context, skuID int64) (*Product, error) {
	if skuID <= 0 {
		return nil, fmt.Errorf("skuID must be greater than 0")
	}

	url := fmt.Sprintf("%s/get_product", s.baseURL)
	reqBody := struct {
		Token string `json:"token"`
		SKU   int64  `json:"sku"`
	}{
		Token: s.token,
		SKU:   skuID,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized")
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, customerrors.SkuNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var product Product
	if err = json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &product, nil
}
