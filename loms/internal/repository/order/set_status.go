package order

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
)

func (r *Repository) SetStatus(ctx context.Context, orderID int64, status model.OrderStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.orders[orderID]
	if !ok {
		return customerrors.ErrOrderNotFound
	}

	order.Status = status
	return nil
}
