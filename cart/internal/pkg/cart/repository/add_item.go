package repository

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
)

func (r *CartStorageRepository) AddItem(ctx context.Context, userId int64, cartItem *model.CartItem) error {
	if userId < 1 {
		return customerrors.InvalidUserId
	}
	if cartItem.SkuId < 1 {
		return customerrors.InvalidSkuId
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.cartStorage[userId]; !ok {
		r.cartStorage[userId] = make([]model.CartItem, 0)
	}

	for i, item := range r.cartStorage[userId] {
		if item.SkuId == cartItem.SkuId {
			r.cartStorage[userId][i].Count += cartItem.Count
			return nil
		}
	}

	r.cartStorage[userId] = append(r.cartStorage[userId], *cartItem)
	return nil
}
