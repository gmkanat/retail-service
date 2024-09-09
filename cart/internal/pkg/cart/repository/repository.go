package repository

import (
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
	"sync"
)

type CartStorageRepository struct {
	cartStorage map[int64]map[int64]*model.CartItem
	mutex       sync.RWMutex
}

func NewCartStorageRepository() *CartStorageRepository {
	return &CartStorageRepository{
		cartStorage: make(map[int64]map[int64]*model.CartItem),
	}
}
