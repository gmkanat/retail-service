package order

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
)

func (s *Service) OrderPay(ctx context.Context, orderID int64) error {
	if orderID <= 0 {
		return customerrors.ErrInvalidOrderId
	}

	order, err := s.orderRepository.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != model.OrderStatusAwaitingPayment {
		return customerrors.ErrOrderStatusAwaitingPayment
	}

	for _, item := range order.Items {
		reserveRemoveStockErr := s.stockRepository.ReserveRemove(ctx, item.SKU, item.Count)
		if reserveRemoveStockErr != nil {
			return reserveRemoveStockErr
		}
	}

	err = s.orderRepository.SetStatus(ctx, orderID, model.OrderStatusPayed)
	if err != nil {
		return err
	}

	return nil
}
