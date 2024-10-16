package stock

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/transaction"
	"log"
)

func (r *Repository) ReserveRemove(ctx context.Context, sku uint32, count uint16) error {
	tx, ok := transaction.GetTx(ctx)
	if !ok {
		return fmt.Errorf("transaction not found in context")
	}

	var reserved uint64
	err := tx.QueryRow(ctx, `
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
		return customerrors.ErrInsufficientStock
	}

	_, err = tx.Exec(ctx, `
		UPDATE stocks.stocks
		SET reserved = reserved - $1
		WHERE id = $2`, count, sku)
	if err != nil {
		log.Printf("Failed to update stock: %v", err)
		return err
	}

	return nil
}
