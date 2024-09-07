package service

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/customerrors"
)

func (s *CartService) ClearCart(ctx context.Context, userId int64) error {
	if userId < 1 {
		return customerrors.InvalidUserId
	}

	return s.repository.ClearCart(ctx, userId)
}
