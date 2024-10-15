package order

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	"log"
	"time"
)

func (r *Repository) SetStatus(ctx context.Context, orderID int64, status model.OrderStatus) error {
	writer, err := r.cluster.GetWriter(ctx)
	if err != nil {
		log.Printf("Failed to get writer: %v", err)
		return err
	}

	tx, err := writer.Begin(ctx)
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()

	log.Printf("Set order %d status to %s", orderID, status)
	_, err = tx.Exec(ctx,
		`UPDATE orders.orders SET status = $1, updated_at = $2 WHERE id = $3`,
		status, time.Now(), orderID,
	)
	if err != nil {
		return err
	}

	log.Printf("Order %d status set to %s", orderID, status)

	if err = r.notifier.CreateEvent(ctx, tx, model.Event{
		OrderID: orderID,
		Status:  status.String(),
	}); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("failed to commit transaction: %v", err)
		return err
	}

	return nil
}
