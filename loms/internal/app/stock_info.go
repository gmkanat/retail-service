package app

import (
	"context"
	servicepb "gitlab.ozon.dev/kanat_9999/homework/loms/pkg/api/proto/v1"
)

func (s *Service) StocksInfo(ctx context.Context, req *servicepb.StocksInfoRequest) (*servicepb.StocksInfoResponse, error) {
	stock, err := s.stockService.StocksInfo(ctx, req.Sku)
	if err != nil {
		return nil, err
	}

	return &servicepb.StocksInfoResponse{AvailableCount: int64(stock)}, nil
}
