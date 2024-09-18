package service_test

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/mocks"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/service"
)

func TestCartService_RemoveItem(t *testing.T) {
	mc := minimock.NewController(t)
	repoMock := mocks.NewCartRepositoryMock(mc)
	cartService := service.NewService(repoMock, nil)

	ctx := context.Background()
	userID := int64(123)
	skuID := int64(1000)

	t.Run("Invalid userId", func(t *testing.T) {
		err := cartService.RemoveItem(ctx, 0, skuID)
		require.Error(t, err)
		require.Equal(t, customerrors.InvalidUserId, err)
	})

	t.Run("Invalid skuId", func(t *testing.T) {
		err := cartService.RemoveItem(ctx, userID, 0)
		require.Error(t, err)
		require.Equal(t, customerrors.InvalidSkuId, err)
	})

	t.Run("Success", func(t *testing.T) {
		cartItem := model.CartItem{
			SkuId: 1000,
			Name:  "Test Product",
			Price: 100,
			Count: 2,
		}

		repoMock.GetCartMock.Expect(ctx, userID).Return([]model.CartItem{cartItem}, nil)
		repoMock.RemoveItemMock.Expect(ctx, userID, skuID).Return(nil)

		items, err := repoMock.GetCart(ctx, userID)
		require.Len(t, items, 1)
		require.NoError(t, err)
		require.NotEmpty(t, items)

		err = cartService.RemoveItem(ctx, userID, skuID)
		require.NoError(t, err)
	})
}
