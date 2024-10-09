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

func TestCartService_Checkout(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)
	cartRepoMock := mocks.NewCartRepositoryMock(mc)
	lomsClientMock := mocks.NewLomsClientMock(mc)

	cartSvc := service.NewService(cartRepoMock, nil, lomsClientMock)

	ctx := context.Background()
	userId := int64(123)
	orderId := int64(456)
	cartItems := []model.CartItem{
		{SkuId: 1001, Name: "Item 1", Count: 2, Price: 500},
		{SkuId: 1002, Name: "Item 2", Count: 1, Price: 300},
	}

	t.Run("invalid userId", func(t *testing.T) {
		t.Parallel()
		invalidUserID := int64(0)
		_, err := cartSvc.Checkout(ctx, invalidUserID)
		require.ErrorIs(t, err, customerrors.InvalidUserId)
	})

	t.Run("successful checkout", func(t *testing.T) {
		cartRepoMock.GetCartMock.Expect(ctx, userId).Return(cartItems, nil)
		lomsClientMock.OrderCreateMock.Expect(ctx, userId, cartItems).Return(orderId, nil)
		cartRepoMock.ClearCartMock.Expect(ctx, userId).Return(nil)

		resultOrderId, err := cartSvc.Checkout(ctx, userId)
		require.NoError(t, err)
		require.Equal(t, orderId, resultOrderId)
	})
}
