package service

import "context"

func (s *StockService) StocksInfo(ctx context.Context, sku uint32) (uint64, error) {
	return s.stockRepository.GetBySKU(ctx, sku)
}
