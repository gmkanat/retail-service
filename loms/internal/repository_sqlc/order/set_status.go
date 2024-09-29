package order

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	pgordersqry "gitlab.ozon.dev/kanat_9999/homework/loms/internal/repository_sqlc/order/queries"
	"log"
)

func (r *Repository) SetStatus(ctx context.Context, orderID int64, status model.OrderStatus) error {
	writer, err := r.cluster.GetWriter(ctx)
	if err != nil {
		log.Printf("Failed to get writer: %v", err)
		return err
	}

	q := pgordersqry.New(writer)

	err = q.UpdateOrderStatus(ctx, pgordersqry.UpdateOrderStatusParams{
		ID:        orderID,
		Status:    pgordersqry.OrdersOrderStatus(status.String()),
		UpdatedAt: currentTimestamp(),
	})
	if err != nil {
		log.Printf("Failed to set order status: %v", err)
		return err
	}
	return nil
}
