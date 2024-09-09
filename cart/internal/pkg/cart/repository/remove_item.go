package repository

import (
	"context"
)

func (r *CartStorageRepository) RemoveItem(ctx context.Context, userId int64, sku int64) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.cartStorage[userId]; !ok {
		return nil
	}

	if _, exists := r.cartStorage[userId][sku]; exists {
		delete(r.cartStorage[userId], sku)
	}

	return nil
}
