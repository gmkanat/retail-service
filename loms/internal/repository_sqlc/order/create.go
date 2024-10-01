package order

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	pgordersqry "gitlab.ozon.dev/kanat_9999/homework/loms/internal/repository_sqlc/order/queries"
	"log"
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

	q := pgordersqry.New(writer).WithTx(tx)

	orderID, err := q.CreateOrder(ctx, pgordersqry.CreateOrderParams{
		UserID:    userID,
		Status:    "New",
		CreatedAt: currentTimestamp(),
		UpdatedAt: currentTimestamp(),
	})

	if err != nil {
		log.Printf("failed to insert order: %v", err)
		return 0, err
	}

	for _, item := range items {
		err = q.InsertOrderItem(ctx, pgordersqry.InsertOrderItemParams{
			SkuID:     int64(item.SKU),
			OrderID:   orderID,
			Count:     int64(item.Count),
			CreatedAt: currentTimestamp(),
			UpdatedAt: currentTimestamp(),
		})
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
