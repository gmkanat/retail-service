package service_test

import (
	"context"
	"errors"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/service"
	productService "gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/product/service"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/mocks"
)

func TestCartService_AddItem(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)

	repoMock := mocks.NewCartRepositoryMock(mc)
	productMock := mocks.NewProductServiceMock(mc)
	lomsClientMock := mocks.NewLomsClientMock(mc)
	cartService := service.NewService(repoMock, productMock, lomsClientMock)

	ctx := context.Background()

	t.Run("Invalid userId", func(t *testing.T) {
		t.Parallel()
		err := cartService.AddItem(ctx, 0, 1, 1)
		require.Error(t, err)
		require.Equal(t, customerrors.InvalidUserId, err)
	})

	t.Run("Invalid skuId", func(t *testing.T) {
		t.Parallel()
		err := cartService.AddItem(ctx, 1, 0, 1)
		require.Error(t, err)
		require.Equal(t, customerrors.InvalidSkuId, err)
	})

	t.Run("Invalid count", func(t *testing.T) {
		t.Parallel()
		err := cartService.AddItem(ctx, 1, 1, 0)
		require.Error(t, err)
		require.Equal(t, customerrors.InvalidCount, err)
	})

	t.Run("ProductService error", func(t *testing.T) {
		t.Parallel()
		productMock.GetProductMock.Expect(ctx, int64(1000)).Return(nil, errors.New("product not found"))

		err := cartService.AddItem(ctx, 1, 1000, 10)
		require.Error(t, err)
		require.Equal(t, "product not found", err.Error())
	})

	t.Run("Success", func(t *testing.T) {
		product := productService.Product{
			Name:  "Кроссовки Nike JORDAN",
			Price: 200,
		}
		productMock.GetProductMock.Expect(ctx, int64(1002)).Return(&product, nil)
		lomsClientMock.GetStockMock.Expect(ctx, int64(1002)).Return(int64(10), nil)

		cartItem := &model.CartItem{
			SkuId: 1002,
			Name:  "Кроссовки Nike JORDAN",
			Count: 2,
			Price: 200,
		}

		repoMock.AddItemMock.Expect(ctx, int64(123), cartItem).Return(nil)

		err := cartService.AddItem(ctx, 123, 1002, 2)
		require.NoError(t, err)
	})
}
