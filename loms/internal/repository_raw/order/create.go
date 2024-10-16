package order

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/transaction"
	"log"
	"time"
)

func (r *Repository) Create(ctx context.Context, userID int64, items []model.Item) (int64, error) {
	tx, ok := transaction.GetTx(ctx)
	if !ok {
		return 0, fmt.Errorf("transaction not found in context")
	}

	var orderID int64
	err := tx.QueryRow(ctx,
		`INSERT INTO orders.orders (user_id, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4) RETURNING id`,
		userID, model.OrderStatusNew.String(), time.Now(), time.Now(),
	).Scan(&orderID)
	if err != nil {
		log.Printf("failed to insert order: %v", err)
		return 0, err
	}

	for _, item := range items {
		_, err = tx.Exec(ctx,
			`INSERT INTO orders.order_items (sku_id, order_id, count, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5)`,
			item.SKU, orderID, item.Count, time.Now(), time.Now(),
		)
		if err != nil {
			log.Printf("failed to insert order item: %v", err)
			return 0, err
		}
	}

	if err = r.notifier.CreateEvent(ctx, model.Event{
		OrderID: orderID,
		Status:  model.OrderStatusNew.String(),
	}); err != nil {
		log.Printf("failed to create event: %v", err)
		return 0, err
	}

	log.Printf("Order %d created", orderID)
	return orderID, nil
}
