package order

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	"log"
)

func (r *Repository) GetByID(ctx context.Context, orderID int64) (*model.Order, error) {
	reader, err := r.cluster.GetReader(ctx)
	if err != nil {
		log.Printf("Failed to get reader: %v", err)
		return nil, err
	}

	var order model.Order
	var status string

	err = reader.QueryRow(ctx,
		`SELECT id, user_id, status, created_at, updated_at FROM orders.orders WHERE id = $1`,
		orderID,
	).Scan(&order.OrderID, &order.UserID, &status, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		log.Printf("Failed to get order: %v", err)
		return nil, err
	}

	rows, err := reader.Query(ctx,
		`SELECT sku_id, count FROM orders.order_items WHERE order_id = $1`,
		orderID,
	)
	if err != nil {
		log.Printf("Failed to get order items: %v", err)
		return nil, err
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		if err := rows.Scan(&item.SKU, &item.Count); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	orderStatus, err := model.OrderStatusFromString(status)
	if err != nil {
		log.Printf("Invalid order status received: %s", status)
		return nil, err
	}
	order.Status = orderStatus

	order.Items = items

	return &order, nil
}
