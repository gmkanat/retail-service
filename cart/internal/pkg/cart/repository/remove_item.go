package repository

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/customerrors"
)

func (r *CartStorageRepository) RemoveItem(ctx context.Context, userId int64, sku int64) error {

	if userId < 1 {
		return customerrors.InvalidUserId
	}

	if sku < 1 {
		return customerrors.InvalidSkuId
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.cartStorage[userId]; !ok {
		return nil
	}

	for i, item := range r.cartStorage[userId] {
		if item.SkuId == sku {
			r.cartStorage[userId] = append(r.cartStorage[userId][:i], r.cartStorage[userId][i+1:]...)
			return nil
		}
	}
	return nil
}
