package order

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
)

func (s *Service) OrderInfo(ctx context.Context, orderID int64) (*model.Order, error) {
	return s.orderRepository.GetByID(ctx, orderID)
}
