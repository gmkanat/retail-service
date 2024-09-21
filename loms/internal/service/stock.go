package service

import "context"

type StockRepository interface {
	ReserveStock(ctx context.Context, sku uint32, count uint16) error
	ReleaseStock(ctx context.Context, sku uint32, count uint16) error
	ReserveRemoveStock(ctx context.Context, sku uint32, count uint16) error
	StocksInfo(ctx context.Context, sku uint32) (uint64, error)
}

type StockService struct {
	stockRepository StockRepository
}

func NewStockService(stockRepository StockRepository) *StockService {
	return &StockService{
		stockRepository: stockRepository,
	}
}

func (s *StockService) StocksInfo(ctx context.Context, sku uint32) (uint64, error) {
	return s.stockRepository.StocksInfo(ctx, sku)
}
