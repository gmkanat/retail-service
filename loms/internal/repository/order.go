package repository

import (
	"context"
	"errors"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	"sync"
)

type OrderRepository struct {
	mu     sync.RWMutex
	orders map[int64]*model.Order
	nextID int64
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders: make(map[int64]*model.Order),
		nextID: 1,
	}
}

func (r *OrderRepository) OrderCreate(ctx context.Context, userID int64, items []model.Item) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order := &model.Order{
		OrderID: r.nextID,
		UserID:  userID,
		Status:  "new",
		Items:   items,
	}

	r.orders[order.OrderID] = order
	r.nextID++
	return order.OrderID, nil
}

func (r *OrderRepository) SetStatus(ctx context.Context, orderID int64, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.orders[orderID]
	if !ok {
		return errors.New("order not found")
	}

	order.Status = status
	return nil
}

func (r *OrderRepository) OrderInfo(ctx context.Context, orderID int64) (*model.Order, error) {
	order, ok := r.orders[orderID]
	if !ok {
		return nil, customerrors.ErrOrderNotFound
	}
	return order, nil
}
