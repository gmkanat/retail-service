package order

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	"log"
	"time"
)

func (r *Repository) SetStatus(ctx context.Context, orderID int64, status model.OrderStatus) error {
	log.Printf("Set order %d status to %s", orderID, status)
	_, err := r.db.Exec(ctx,
		`UPDATE orders.orders SET status = $1, updated_at = $2 WHERE id = $3`,
		status, time.Now(), orderID,
	)
	if err != nil {
		return err
	}
	log.Printf("Order %d status set to %s", orderID, status)
	return nil
}
