package stock

import "context"

func (s *Service) StocksInfo(ctx context.Context, sku uint32) (uint64, error) {
	return s.stockRepository.GetBySKU(ctx, sku)
}
