package service

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/customerrors"
)

func (s *CartService) Checkout(ctx context.Context, userId int64) (orderId int64, error error) {
	if userId < 1 {
		return 0, customerrors.InvalidUserId
	}

	cartItems, err := s.repository.GetCart(ctx, userId)
	if err != nil {
		return 0, err
	}

	orderId, err = s.lomsClient.OrderCreate(ctx, userId, cartItems)
	if err != nil {
		return 0, err
	}

	err = s.repository.ClearCart(ctx, userId)

	return orderId, err
}
