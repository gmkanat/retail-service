package service

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
)

func (s *CartService) GetCart(ctx context.Context, userId int64) (*model.GetCartResponse, error) {
	if userId < 1 {
		return nil, customerrors.InvalidUserId
	}

	items, err := s.repository.GetCart(ctx, userId)
	if err != nil {
		return nil, err
	}

	totalPrice := uint32(0)
	for _, item := range items {
		product, err := s.productService.GetProduct(ctx, item.SkuId)
		if err != nil {
			return nil, err
		}

		totalPrice += product.Price * uint32(item.Count)
	}

	return &model.GetCartResponse{
		Items:      items,
		TotalPrice: totalPrice,
	}, nil

}
