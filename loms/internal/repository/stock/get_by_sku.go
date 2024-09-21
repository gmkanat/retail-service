package stock

import (
	"context"
	"errors"
)

func (r *Repository) GetBySKU(ctx context.Context, sku uint32) (uint64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stock, ok := r.stocks[sku]
	if !ok {
		return 0, errors.New("stock not found")
	}

	return stock.TotalCount - stock.Reserved, nil
}
