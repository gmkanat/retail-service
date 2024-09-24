package stock

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
)

func (r *Repository) Release(ctx context.Context, sku uint32, count uint16) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	stock, ok := r.stocks[sku]
	if !ok {
		return customerrors.ErrStockNotFound
	}

	if stock.Reserved < uint64(count) {
		return customerrors.ErrNotEnoughReservedStock
	}

	stock.Reserved -= uint64(count)
	return nil
}
