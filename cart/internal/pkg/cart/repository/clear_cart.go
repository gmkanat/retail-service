package repository

import (
	"context"
)

func (r *CartStorageRepository) ClearCart(ctx context.Context, userId int64) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.cartStorage, userId)

	return nil
}
