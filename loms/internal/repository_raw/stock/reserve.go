package stock

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
	"log"
)

func (r *Repository) Reserve(ctx context.Context, sku uint32, count uint16) error {
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

	var available, reserved uint64
	err = tx.QueryRow(ctx, `
		SELECT available, reserved
		FROM stocks.stocks
		WHERE id = $1`, sku).Scan(&available, &reserved)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Stock %d not found", sku)
		return customerrors.ErrStockNotFound
	}
	if err != nil {
		log.Printf("Failed to get stock: %v", err)
		return err
	}

	if available < uint64(count) {
		log.Printf("Not enough available stock for sku %d", sku)
		return customerrors.ErrInsufficientStock
	}

	_, err = tx.Exec(ctx, `
		UPDATE stocks.stocks
		SET available = available - $1, reserved = reserved + $1
		WHERE id = $2`, count, sku)
	if err != nil {
		log.Printf("Failed to update stock: %v", err)
		return err
	}

	return tx.Commit(ctx)
}
