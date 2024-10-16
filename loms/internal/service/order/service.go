package order

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
)

type Repository interface {
	Create(ctx context.Context, userID int64, items []model.Item) (int64, error)
	SetStatus(ctx context.Context, orderID int64, status model.OrderStatus) error
	GetByID(ctx context.Context, orderID int64) (*model.Order, error)
}

type StockRepository interface {
	Reserve(ctx context.Context, sku uint32, count uint16) error
	Release(ctx context.Context, sku uint32, count uint16) error
	ReserveRemove(ctx context.Context, sku uint32, count uint16) error
}

type TransactionManager interface {
	WithRepeatableReadTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type Service struct {
	orderRepository Repository
	stockRepository StockRepository
	txManager       TransactionManager
}

func NewOrderService(
	orderRepository Repository,
	stockRepository StockRepository,
	txManager TransactionManager,
) *Service {
	return &Service{
		orderRepository: orderRepository,
		stockRepository: stockRepository,
		txManager:       txManager,
	}
}
