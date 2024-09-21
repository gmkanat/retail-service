package repository

import (
	"context"
	"errors"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	"sync"
)

type StockRepository struct {
	mu     sync.RWMutex
	stocks map[uint32]*model.Stock
}

func NewStockRepository(initialData []model.Stock) *StockRepository {
	stocks := make(map[uint32]*model.Stock)
	for _, stock := range initialData {
		stocks[stock.SKU] = &stock
	}
	return &StockRepository{
		stocks: stocks,
	}
}

func (r *StockRepository) ReserveStock(ctx context.Context, sku uint32, count uint16) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	stock, ok := r.stocks[sku]
	if !ok {
		return errors.New("stock not found")
	}

	if stock.TotalCount-stock.Reserved < uint64(count) {
		return errors.New("insufficient stock")
	}

	stock.Reserved += uint64(count)
	return nil
}

func (r *StockRepository) ReleaseStock(ctx context.Context, sku uint32, count uint16) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	stock, ok := r.stocks[sku]
	if !ok {
		return errors.New("stock not found")
	}

	if stock.Reserved < uint64(count) {
		return errors.New("not enough reserved stock")
	}

	stock.Reserved -= uint64(count)
	return nil
}

func (r *StockRepository) ReserveRemoveStock(ctx context.Context, sku uint32, count uint16) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	stock, ok := r.stocks[sku]
	if !ok {
		return errors.New("stock not found")
	}

	if stock.TotalCount-stock.Reserved < uint64(count) {
		return errors.New("insufficient stock")
	}

	stock.TotalCount -= uint64(count)
	stock.Reserved -= uint64(count)

	return nil
}

func (r *StockRepository) StocksInfo(ctx context.Context, sku uint32) (uint64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stock, ok := r.stocks[sku]
	if !ok {
		return 0, errors.New("stock not found")
	}

	return stock.TotalCount - stock.Reserved, nil
}
