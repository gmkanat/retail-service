package service

import "context"

type StockRepository interface {
	Reserve(ctx context.Context, sku uint32, count uint16) error
	Release(ctx context.Context, sku uint32, count uint16) error
	ReserveRemove(ctx context.Context, sku uint32, count uint16) error
	GetBySKU(ctx context.Context, sku uint32) (uint64, error)
}

type StockService struct {
	stockRepository StockRepository
}

func NewStockService(stockRepository StockRepository) *StockService {
	return &StockService{
		stockRepository: stockRepository,
	}
}
