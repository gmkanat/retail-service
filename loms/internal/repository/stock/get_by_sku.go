package stock

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
)

func (r *Repository) GetBySKU(ctx context.Context, sku uint32) (uint64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stock, ok := r.stocks[sku]
	if !ok {
		return 0, customerrors.ErrStockNotFound
	}

	return stock.TotalCount - stock.Reserved, nil
}
