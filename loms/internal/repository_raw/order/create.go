package order

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	"log"
	"time"
)

func (r *Repository) Create(ctx context.Context, userID int64, items []model.Item) (int64, error) {
	writer, err := r.cluster.GetWriter(ctx)
	if err != nil {
		log.Printf("Failed to get writer: %v", err)
		return 0, err
	}

	tx, err := writer.Begin(ctx)
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		return 0, err
	}
	defer tx.Rollback(ctx)

	var orderID int64
	err = tx.QueryRow(ctx,
		`INSERT INTO orders.orders (user_id, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4) RETURNING id`,
		userID, "New", time.Now(), time.Now(),
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

	if err = tx.Commit(ctx); err != nil {
		log.Printf("failed to commit transaction: %v", err)
		return 0, err
	}

	log.Printf("Order %d created", orderID)
	return orderID, nil
}
