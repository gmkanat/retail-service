package repository

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
)

func (r *CartStorageRepository) GetCart(ctx context.Context, userId int64) ([]model.CartItem, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	itemsMap, ok := r.cartStorage[userId]
	if !ok {
		return []model.CartItem{}, nil
	}

	var items []model.CartItem
	for _, item := range itemsMap {
		items = append(items, *item)
	}

	return items, nil
}
