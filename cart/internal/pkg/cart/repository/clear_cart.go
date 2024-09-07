package repository

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/customerrors"
)

func (r *CartStorageRepository) ClearCart(ctx context.Context, userId int64) error {
	if userId < 1 {
		return customerrors.InvalidUserId
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.cartStorage, userId)

	return nil
}
