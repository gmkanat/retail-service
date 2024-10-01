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
	log.Printf("Getting order %d", orderID)
	orderItems, err := q.GetOrderWithItems(ctx, orderID)
	if err != nil {
		log.Printf("Failed to get order and items: %v", err)
		return nil, err
	}

	var order model.Order
	var items []model.Item
	var status string

	for _, orderItem := range orderItems {
		order = model.Order{
			UserID:    orderItem.UserID,
			CreatedAt: orderItem.CreatedAt.Time,
			UpdatedAt: orderItem.UpdatedAt.Time,
		}
		items = append(items, model.Item{
			SKU:   uint32(orderItem.SkuID),
			Count: uint16(orderItem.Count),
		})
		status = string(orderItem.Status)
	}

	order.Status, err = model.OrderStatusFromString(status)
	if err != nil {
		log.Printf("Invalid order status: %v", err)
		return nil, err
	}

	log.Printf("Order %d found", orderID)

	order.Items = items
	return &order, nil
}
