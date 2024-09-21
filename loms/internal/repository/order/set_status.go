package order

import (
	"context"
	"errors"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
)

func (r *Repository) SetStatus(ctx context.Context, orderID int64, status model.OrderStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.orders[orderID]
	if !ok {
		return errors.New("order not found")
	}

	order.Status = status
	return nil
}
