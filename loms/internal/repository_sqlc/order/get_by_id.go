package order

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	pgordersqry "gitlab.ozon.dev/kanat_9999/homework/loms/internal/repository_sqlc/order/queries"
	"log"
)

func (r *Repository) GetByID(ctx context.Context, orderID int64) (*model.Order, error) {
	reader, err := r.cluster.GetReader(ctx)
	if err != nil {
		log.Printf("Failed to get reader: %v", err)
		return nil, err
	}
	q := pgordersqry.New(reader)

	order, err := q.GetOrderById(ctx, orderID)
	if err != nil {
		log.Printf("Failed to get order: %v", err)
		return nil, err
	}

	items, err := q.GetOrderItems(ctx, orderID)
	if err != nil {
		log.Printf("Failed to get order items: %v", err)
		return nil, err
	}

	var orderItems []model.Item
	for _, item := range items {
		orderItems = append(orderItems, model.Item{
			SKU:   uint32(item.SkuID),
			Count: uint16(item.Count),
		})
	}

	orderStatus, err := model.OrderStatusFromString(string(order.Status))
	if err != nil {
		log.Printf("Invalid order status: %v", err)
		return nil, err
	}

	return &model.Order{
		OrderID:   order.ID,
		UserID:    order.UserID,
		Status:    orderStatus,
		CreatedAt: order.CreatedAt.Time,
		UpdatedAt: order.UpdatedAt.Time,
		Items:     orderItems,
	}, nil
}
