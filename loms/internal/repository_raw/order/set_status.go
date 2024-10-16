package order

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/transaction"
	"time"
)

func (r *Repository) SetStatus(ctx context.Context, orderID int64, status model.OrderStatus) error {
	tx, ok := transaction.GetTx(ctx)
	if !ok {
		return fmt.Errorf("transaction not found in context")
	}

	_, err := tx.Exec(ctx,
		`UPDATE orders.orders SET status = $1, updated_at = $2 WHERE id = $3`,
		status, time.Now(), orderID,
	)
	if err != nil {
		return err
	}

	if err = r.notifier.CreateEvent(ctx, model.Event{
		OrderID: orderID,
		Status:  status.String(),
	}); err != nil {
		return err
	}

	return nil
}
