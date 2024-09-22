package service_test

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/mocks"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/service"
)

func TestCartService_ClearCart(t *testing.T) {
	mc := minimock.NewController(t)
	repoMock := mocks.NewCartRepositoryMock(mc)
	cartService := service.NewService(repoMock, nil, nil)

	ctx := context.Background()
	userID := int64(123)

	t.Run("Invalid userId", func(t *testing.T) {
		err := cartService.ClearCart(ctx, 0)
		require.Error(t, err)
		require.Equal(t, customerrors.InvalidUserId, err)
	})

	t.Run("Success", func(t *testing.T) {
		// add item, then clear cart
		cartItem := model.CartItem{
			SkuId: 1000,
			Name:  "Кроссовки Nike JORDAN",
			Count: 2,
			Price: 200,
		}
		repoMock.GetCartMock.Expect(ctx, userID).Return([]model.CartItem{cartItem}, nil)

		repoMock.ClearCartMock.Expect(ctx, userID).Return(nil)

		items, err := repoMock.GetCart(ctx, userID)
		require.Len(t, items, 1)
		require.NoError(t, err)
		require.NotEmpty(t, items)

		err = cartService.ClearCart(ctx, userID)
		require.NoError(t, err)

		repoMock.GetCartMock.Expect(ctx, userID).Return([]model.CartItem{}, nil)
		emptyItems, err := repoMock.GetCart(ctx, userID)
		require.NoError(t, err)
		require.Empty(t, emptyItems)
	})
}
