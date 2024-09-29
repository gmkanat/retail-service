package stock

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
	"log"
)

func (r *Repository) GetBySKU(ctx context.Context, sku uint32) (uint64, error) {
	reader, err := r.cluster.GetReader(ctx)
	if err != nil {
		log.Printf("Failed to get reader: %v", err)
		return 0, err
	}

	var available, reserved uint64

	err = reader.QueryRow(ctx, `
		SELECT available, reserved
		FROM stocks.stocks
		WHERE id = $1`, sku).Scan(&available, &reserved)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Stock %d not found", sku)
		return 0, customerrors.ErrStockNotFound
	}
	if err != nil {
		log.Printf("Failed to get stock: %v", err)
		return 0, err
	}

	return available, nil
}
