package stock

import (
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	"sync"
)

type Repository struct {
	mu     sync.RWMutex
	stocks map[uint32]*model.Stock
}

func NewStockRepository(initialData []model.Stock) *Repository {
	stocks := make(map[uint32]*model.Stock, 8)
	for _, stock := range initialData {
		stocks[stock.SKU] = &stock
	}
	return &Repository{
		stocks: stocks,
	}
}
