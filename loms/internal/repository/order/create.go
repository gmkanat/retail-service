package order

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
)

func (r *Repository) Create(ctx context.Context, userID int64, items []model.Item) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order := &model.Order{
		OrderID: r.nextID,
		UserID:  userID,
		Status:  model.OrderStatusNew,
		Items:   items,
	}

	r.orders[order.OrderID] = order
	r.nextID++
	return order.OrderID, nil
}
