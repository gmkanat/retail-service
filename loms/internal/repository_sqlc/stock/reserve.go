package stock

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
	pgstocksqry "gitlab.ozon.dev/kanat_9999/homework/loms/internal/repository_sqlc/stock/queries"
	"log"
)

func (r *Repository) Reserve(ctx context.Context, sku uint32, count uint16) error {
	writer, err := r.cluster.GetWriter(ctx)
	if err != nil {
		log.Printf("Failed to get writer: %v", err)
		return err
	}

	q := pgstocksqry.New(writer)

	tx, err := writer.Begin(ctx)
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return err
	}

	defer tx.Rollback(ctx)

	skuRow, err := q.GetStockBySKU(ctx, int64(sku))
	if err != nil {
		log.Printf("Failed to get stock: %v", err)
		return err
	}

	if uint16(skuRow.Available) < count {
		log.Printf("Not enough available stock for sku %d", sku)
		return customerrors.ErrInsufficientStock
	}

	err = q.ReserveStock(ctx, pgstocksqry.ReserveStockParams{
		ID:        int64(sku),
		Available: int64(count),
	})
	if err != nil {
		log.Printf("Failed to reserve stock: %v", err)
		return err
	}

	return nil
}
