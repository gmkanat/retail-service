package stock

import (
	"context"
	"errors"
)

func (r *Repository) Reserve(ctx context.Context, sku uint32, count uint16) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	stock, ok := r.stocks[sku]
	if !ok {
		return errors.New("stock not found")
	}

	if stock.TotalCount-stock.Reserved < uint64(count) {
		return errors.New("insufficient stock")
	}

	stock.Reserved += uint64(count)
	return nil
}
