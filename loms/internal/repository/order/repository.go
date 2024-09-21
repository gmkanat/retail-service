package order

import (
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	"sync"
)

type Repository struct {
	mu     sync.RWMutex
	orders map[int64]*model.Order
	nextID int64
}

func NewOrderRepository() *Repository {
	return &Repository{
		orders: make(map[int64]*model.Order),
		nextID: 1,
	}
}
