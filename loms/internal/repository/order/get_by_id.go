package order

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
)

func (r *Repository) GetByID(ctx context.Context, orderID int64) (*model.Order, error) {
	order, ok := r.orders[orderID]
	if !ok {
		return nil, customerrors.ErrOrderNotFound
	}
	return order, nil
}
