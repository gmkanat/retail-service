package service

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/service/stock"
)

type OrderRepository interface {
	Create(ctx context.Context, userID int64, items []model.Item) (int64, error)
	SetStatus(ctx context.Context, orderID int64, status model.OrderStatus) error
	GetByID(ctx context.Context, orderID int64) (*model.Order, error)
}

type Service struct {
	orderRepository OrderRepository
	stockRepository service.StockRepository
}

func NewOrderService(orderRepository OrderRepository, stockRepository service.StockRepository) *Service {
	return &Service{
		orderRepository: orderRepository,
		stockRepository: stockRepository,
	}
}
