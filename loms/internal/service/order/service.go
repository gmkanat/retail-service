package order

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/service/stock"
)

type Repository interface {
	Create(ctx context.Context, userID int64, items []model.Item) (int64, error)
	SetStatus(ctx context.Context, orderID int64, status model.OrderStatus) error
	GetByID(ctx context.Context, orderID int64) (*model.Order, error)
}

type Service struct {
	orderRepository Repository
	stockRepository stock.Repository
}

func NewOrderService(orderRepository Repository, stockRepository stock.Repository) *Service {
	return &Service{
		orderRepository: orderRepository,
		stockRepository: stockRepository,
	}
}
