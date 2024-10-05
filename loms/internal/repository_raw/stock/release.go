package stock

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
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

	var reserved uint64
	err = tx.QueryRow(ctx, `
		SELECT reserved
		FROM stocks.stocks
		WHERE id = $1`, sku).Scan(&reserved)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Stock %d not found", sku)
		return customerrors.ErrStockNotFound
	}

	if err != nil {
		return err
	}

	if reserved < uint64(count) {
		log.Printf("Not enough reserved stock for sku %d", sku)
		return customerrors.ErrNotEnoughReservedStock
	}

	_, err = tx.Exec(ctx, `
		UPDATE stocks.stocks
		SET reserved = reserved - $1
		WHERE id = $2`, count, sku)
	if err != nil {
		log.Printf("Failed to update stock: %v", err)
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return err
	}

	return nil
}
