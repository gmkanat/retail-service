package service

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	"log"
)

type OrderRepository interface {
	OrderCreate(ctx context.Context, userID int64, items []model.Item) (int64, error)
	SetStatus(ctx context.Context, orderID int64, status string) error
	OrderInfo(ctx context.Context, orderID int64) (*model.Order, error)
}

type OrderService struct {
	orderRepository OrderRepository
	stockRepository StockRepository
}

func NewOrderService(orderRepository OrderRepository, stockRepository StockRepository) *OrderService {
	return &OrderService{
		orderRepository: orderRepository,
		stockRepository: stockRepository,
	}
}

func (s *OrderService) OrderCreate(ctx context.Context, userID int64, items []model.Item) (int64, error) {
	orderID, err := s.orderRepository.OrderCreate(ctx, userID, items)
	if err != nil {
		return 0, err
	}

	var reservedItems []model.Item

	for _, item := range items {
		reserveStockErr := s.stockRepository.ReserveStock(ctx, item.SKU, item.Count)
		if reserveStockErr != nil {
			log.Printf("failed to reserve stock for order %d: %v", orderID, err)

			// rollback stock reservation
			for _, reservedItem := range reservedItems {
				rollbackErr := s.stockRepository.ReleaseStock(ctx, reservedItem.SKU, reservedItem.Count)
				if rollbackErr != nil {
					log.Printf("rollback failed for SKU: %d", reservedItem.SKU)
				}
			}

			updateOrderErr := s.orderRepository.SetStatus(ctx, orderID, "failed")
			if updateOrderErr != nil {
				log.Printf("failed to update order status to failed: %v", updateOrderErr)
			}
			return orderID, err
		}
		reservedItems = append(reservedItems, item)
	}

	err = s.orderRepository.SetStatus(ctx, orderID, "awaiting_payment")
	if err != nil {
		return orderID, err
	}

	return orderID, nil
}

func (s *OrderService) OrderInfo(ctx context.Context, orderID int64) (*model.Order, error) {
	return s.orderRepository.OrderInfo(ctx, orderID)
}

func (s *OrderService) OrderPay(ctx context.Context, orderID int64) error {
	order, err := s.orderRepository.OrderInfo(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != "awaiting_payment" {
		return nil
	}

	for _, item := range order.Items {
		reserveRemoveStockErr := s.stockRepository.ReserveRemoveStock(ctx, item.SKU, item.Count)
		if reserveRemoveStockErr != nil {
			return reserveRemoveStockErr
		}
	}

	err = s.orderRepository.SetStatus(ctx, orderID, "paid")
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) OrderCancel(ctx context.Context, orderID int64) error {
	order, err := s.orderRepository.OrderInfo(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != "awaiting_payment" {
		return nil
	}

	for _, item := range order.Items {
		releaseErr := s.stockRepository.ReleaseStock(ctx, item.SKU, item.Count)
		if releaseErr != nil {
			return releaseErr
		}
	}

	err = s.orderRepository.SetStatus(ctx, orderID, "cancelled")
	if err != nil {
		return err
	}

	return nil
}
