package service

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/customerrors"
)

func (s *CartService) RemoveItem(ctx context.Context, userId int64, skuId int64) error {
	if userId < 1 {
		return customerrors.InvalidUserId
	}

	if skuId < 1 {
		return customerrors.InvalidSkuId
	}

	return s.repository.RemoveItem(ctx, userId, skuId)
}
