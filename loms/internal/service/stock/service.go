package stock

import "context"

type Repository interface {
	Reserve(ctx context.Context, sku uint32, count uint16) error
	Release(ctx context.Context, sku uint32, count uint16) error
	ReserveRemove(ctx context.Context, sku uint32, count uint16) error
	GetBySKU(ctx context.Context, sku uint32) (uint64, error)
}

type Service struct {
	stockRepository Repository
}

func NewStockService(stockRepository Repository) *Service {
	return &Service{
		stockRepository: stockRepository,
	}
}
