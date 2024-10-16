package order

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
)

func (s *Service) OrderCancel(ctx context.Context, orderID int64) error {
	if orderID <= 0 {
		return customerrors.ErrInvalidOrderId
	}

	order, err := s.orderRepository.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	return s.txManager.WithRepeatableReadTx(ctx, func(c context.Context) error {
		if order.Status != model.OrderStatusAwaitingPayment {
			return customerrors.ErrOrderStatusAwaitingPayment
		}

		for _, item := range order.Items {
			releaseErr := s.stockRepository.Release(c, item.SKU, item.Count)
			if releaseErr != nil {
				return releaseErr
			}
		}

		err = s.orderRepository.SetStatus(c, orderID, model.OrderStatusCancelled)
		if err != nil {
			return err
		}

		return nil
	})
}
