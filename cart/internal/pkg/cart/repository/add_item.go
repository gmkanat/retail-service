package repository

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
)

func (r *CartStorageRepository) AddItem(ctx context.Context, userId int64, cartItem *model.CartItem) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.cartStorage[userId]; !ok {
		r.cartStorage[userId] = make(map[int64]*model.CartItem)
	}

	if item, exists := r.cartStorage[userId][cartItem.SkuId]; exists {
		item.Count += cartItem.Count
	} else {
		r.cartStorage[userId][cartItem.SkuId] = cartItem
	}

	return nil
}
