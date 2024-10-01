package stock

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
	pgstocksqry "gitlab.ozon.dev/kanat_9999/homework/loms/internal/repository_sqlc/stock/queries"
	"log"
)

func (r *Repository) ReserveRemove(ctx context.Context, sku uint32, count uint16) error {
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

	skuRow, err := q.GetStockBySKU(ctx, int64(sku))
	if err != nil {
		log.Printf("Failed to get stock: %v", err)
		return err
	}

	if uint16(skuRow.Reserved) < count {
		log.Printf("Not enough reserved stock for sku %d", sku)
		return customerrors.ErrNotEnoughReservedStock
	}

	err = q.ReserveRemoveStock(ctx, pgstocksqry.ReserveRemoveStockParams{
		ID:       int64(sku),
		Reserved: int64(count),
	})

	if err != nil {
		log.Printf("Failed to reserve and remove stock: %v", err)
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return err
	}
	return nil
}
