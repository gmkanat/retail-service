package repository

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
)

func (r *CartStorageRepository) GetCart(ctx context.Context, userId int64) ([]model.CartItem, error) {

	if userId < 1 {
		return nil, customerrors.InvalidUserId
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if _, ok := r.cartStorage[userId]; !ok {
		return []model.CartItem{}, nil
	}

	return r.cartStorage[userId], nil
}
