package stock

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
)

func (r *Repository) ReserveRemove(ctx context.Context, sku uint32, count uint16) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	stock, ok := r.stocks[sku]
	if !ok {
		return customerrors.ErrStockNotFound
	}

	if stock.TotalCount-stock.Reserved < uint64(count) {
		return customerrors.ErrInsufficientStock
	}

	stock.TotalCount -= uint64(count)
	stock.Reserved -= uint64(count)

	return nil
}
