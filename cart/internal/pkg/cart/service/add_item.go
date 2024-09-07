package service

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
)

func (s *CartService) AddItem(ctx context.Context, userId int64, skuId int64, count uint16) error {
	if userId < 1 {
		return customerrors.InvalidUserId
	}
	if skuId < 1 {
		return customerrors.InvalidSkuId
	}

	if count < 1 {
		return customerrors.InvalidCount
	}

	product, err := s.productService.GetProduct(ctx, skuId)
	if err != nil {
		return err
	}

	cartItem := &model.CartItem{
		SkuId: skuId,
		Name:  product.Name,
		Count: count,
		Price: product.Price,
	}

	return s.repository.AddItem(ctx, userId, cartItem)
}
