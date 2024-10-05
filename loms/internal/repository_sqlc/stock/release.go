package stock

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
	pgstocksqry "gitlab.ozon.dev/kanat_9999/homework/loms/internal/repository_sqlc/stock/queries"
	"log"
)

func (r *Repository) Release(ctx context.Context, sku uint32, count uint16) error {
	writer, err := r.cluster.GetWriter(ctx)
	if err != nil {
		log.Printf("Failed to get writer: %v", err)
		return err
	}

	tx, err := writer.Begin(ctx)
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return err
	}

	defer tx.Rollback(ctx)

	q := pgstocksqry.New(writer).WithTx(tx)

	reserved, err := q.GetReservedStock(ctx, int64(sku))

	if err != nil {
		log.Printf("Failed to get reserved stock: %v", err)
		return err
	}

	if uint16(reserved) < count {
		log.Printf("Not enough reserved stock for sku %d", sku)
		return customerrors.ErrNotEnoughReservedStock
	}

	err = q.ReleaseStock(ctx, pgstocksqry.ReleaseStockParams{
		ID:        int64(sku),
		Available: int64(count),
	})
	if err != nil {
		log.Printf("Failed to release stock: %v", err)
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return err
	}

	return nil
}
