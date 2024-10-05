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
	var items []model.Item

	rows, err := reader.Query(ctx, `
		SELECT
			o.id, o.user_id, o.status, o.created_at, o.updated_at,
			COALESCE(oi.sku_id, 0) AS sku_id, COALESCE(oi.count, 0) AS count
		FROM
			orders.orders o
		LEFT JOIN
			orders.order_items oi ON o.id = oi.order_id
		WHERE
			o.id = $1;
	`, orderID)

	if err != nil {
		log.Printf("Failed to get order and items: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Item
		if err = rows.Scan(&order.OrderID, &order.UserID, &status, &order.CreatedAt, &order.UpdatedAt, &item.SKU, &item.Count); err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Failed during rows iteration: %v", err)
		return nil, err
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
