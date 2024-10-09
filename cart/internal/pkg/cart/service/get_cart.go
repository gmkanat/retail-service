package service

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/errgroup"
	"log"
)

func (s *CartService) GetCart(ctx context.Context, userId int64) (*model.GetCartResponse, error) {
	if userId < 1 {
		return nil, customerrors.InvalidUserId
	}

	items, err := s.repository.GetCart(ctx, userId)
	if err != nil {
		return nil, err
	}

	totalPriceCh := make(chan uint32, len(items))

	g, gCtx := errgroup.WithContext(ctx)

	for i := range items {
		item := items[i]
		g.Go(func() error {
			log.Printf("getting product for skuId: %d", item.SkuId)

			product, err := s.productService.GetProduct(gCtx, item.SkuId)
			if err != nil {
				return err
			}
			totalPriceCh <- product.Price * uint32(item.Count)
			return nil
		})
	}

	if err = g.Wait(); err != nil {
		return nil, err
	}

	close(totalPriceCh)

	var totalPrice uint32
	for price := range totalPriceCh {
		totalPrice += price
	}

	return &model.GetCartResponse{
		Items:      items,
		TotalPrice: totalPrice,
	}, nil
}
