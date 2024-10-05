package app

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	servicepb "gitlab.ozon.dev/kanat_9999/homework/loms/pkg/api/proto/v1"
)

var _ servicepb.LomsServer = (*Service)(nil)

type OrderService interface {
	OrderCreate(ctx context.Context, userID int64, items []model.Item) (int64, error)
	OrderInfo(ctx context.Context, orderID int64) (*model.Order, error)
	OrderPay(ctx context.Context, orderID int64) error
	OrderCancel(ctx context.Context, orderID int64) error
}

type StockService interface {
	StocksInfo(ctx context.Context, sku uint32) (uint64, error)
}

type Service struct {
	servicepb.UnimplementedLomsServer
	orderService OrderService
	stockService StockService
}

func NewService(orderService OrderService, stockService StockService) *Service {
	return &Service{
		orderService: orderService,
		stockService: stockService,
	}
}
