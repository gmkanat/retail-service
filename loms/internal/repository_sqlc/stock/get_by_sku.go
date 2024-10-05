package stock

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
	pgstocksqry "gitlab.ozon.dev/kanat_9999/homework/loms/internal/repository_sqlc/stock/queries"
	"log"
)

func (r *Repository) GetBySKU(ctx context.Context, sku uint32) (uint64, error) {
	reader, err := r.cluster.GetReader(ctx)
	if err != nil {
		log.Printf("Failed to get reader: %v", err)
		return 0, err
	}

	q := pgstocksqry.New(reader)

	stock, err := q.GetStockBySKU(ctx, int64(sku))

	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Stock %d not found", sku)
		return 0, customerrors.ErrStockNotFound
	}

	if err != nil {
		log.Printf("Failed to get stock by SKU: %v", err)
		return 0, err
	}

	return uint64(stock.Available), nil
}
