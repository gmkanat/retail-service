package service

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	"log"
)

func (s *Service) OrderCreate(ctx context.Context, userID int64, items []model.Item) (int64, error) {
	orderID, err := s.orderRepository.Create(ctx, userID, items)
	if err != nil {
		return 0, err
	}

	var reservedItems []model.Item

	for _, item := range items {
		reserveStockErr := s.stockRepository.Reserve(ctx, item.SKU, item.Count)
		if reserveStockErr != nil {
			log.Printf("failed to reserve stock for order %d: %v", orderID, reserveStockErr)

			// rollback stock reservation
			for _, reservedItem := range reservedItems {
				rollbackErr := s.stockRepository.Release(ctx, reservedItem.SKU, reservedItem.Count)
				if rollbackErr != nil {
					log.Printf("rollback failed for SKU: %d", reservedItem.SKU)
				}
			}

			updateOrderErr := s.orderRepository.SetStatus(ctx, orderID, model.OrderStatusFailed)
			if updateOrderErr != nil {
				log.Printf("failed to update order status to failed: %v", updateOrderErr)
			}
			return orderID, customerrors.ErrOrderStatusFailed
		}
		reservedItems = append(reservedItems, item)
	}

	err = s.orderRepository.SetStatus(ctx, orderID, model.OrderStatusAwaitingPayment)
	if err != nil {
		return orderID, err
	}

	return orderID, nil
}
