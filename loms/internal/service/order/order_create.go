package order

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	"log"
)

func (s *Service) OrderCreate(ctx context.Context, userID int64, items []model.Item) (int64, error) {
	if userID <= 0 {
		return 0, customerrors.ErrInvalidUserId
	}

	var orderID int64

	err := s.txManager.WithRepeatableReadTx(ctx, func(c context.Context) error {
		id, err := s.orderRepository.Create(c, userID, items)
		if err != nil {
			return err
		}
		orderID = id

		for _, item := range items {
			if err := s.stockRepository.Reserve(c, item.SKU, item.Count); err != nil {
				return err
			}
		}

		if err := s.orderRepository.SetStatus(c, orderID, model.OrderStatusAwaitingPayment); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if err := s.txManager.WithRepeatableReadTx(ctx, func(c context.Context) error {
			return s.orderRepository.SetStatus(c, orderID, model.OrderStatusFailed)
		}); err != nil {
			log.Printf("failed to update order status to failed: %v", err)
		}
		return orderID, customerrors.ErrOrderStatusFailed
	}

	return orderID, nil
}
