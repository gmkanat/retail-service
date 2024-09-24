package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/mocks"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/service"
	productService "gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/product/service"
)

func TestCartService_GetCart(t *testing.T) {
	mc := minimock.NewController(t)
	repoMock := mocks.NewCartRepositoryMock(mc)
	productMock := mocks.NewProductServiceMock(mc)

	cartService := service.NewService(repoMock, productMock, nil)

	ctx := context.Background()

	t.Run("Invalid userId", func(t *testing.T) {
		_, err := cartService.GetCart(ctx, 0)
		require.Error(t, err)
		require.Equal(t, customerrors.InvalidUserId, err)
	})

	t.Run("Empty cart", func(t *testing.T) {
		repoMock.GetCartMock.Expect(ctx, int64(123)).Return([]model.CartItem{}, nil)

		cartResponse, err := cartService.GetCart(ctx, int64(123))
		require.NoError(t, err)
		require.Empty(t, cartResponse.Items)
		require.Equal(t, uint32(0), cartResponse.TotalPrice)
	})

	t.Run("Success", func(t *testing.T) {
		cartItem := model.CartItem{
			SkuId: 1000,
			Name:  "Кроссовки Nike JORDAN",
			Price: 100,
			Count: 2,
		}

		repoMock.GetCartMock.Expect(ctx, int64(123)).Return([]model.CartItem{cartItem}, nil)

		product := productService.Product{
			Name:  "Кроссовки Nike JORDAN",
			Price: 100,
		}
		productMock.GetProductMock.Expect(ctx, int64(1000)).Return(&product, nil)

		cartResponse, err := cartService.GetCart(ctx, int64(123))
		require.NoError(t, err)
		require.NotNil(t, cartResponse)
		require.Len(t, cartResponse.Items, 1)
		require.Equal(t, uint32(200), cartResponse.TotalPrice)
	})

	t.Run("ProductService error", func(t *testing.T) {
		cartItem := model.CartItem{
			SkuId: 1000,
			Name:  "Test Product",
			Price: 100,
			Count: 2,
		}

		repoMock.GetCartMock.Expect(ctx, int64(123)).Return([]model.CartItem{cartItem}, nil)

		productMock.GetProductMock.Expect(ctx, int64(1000)).Return(nil, errors.New("product not found"))

		_, err := cartService.GetCart(ctx, int64(123))
		require.Error(t, err)
		require.Equal(t, "product not found", err.Error())
	})
}
