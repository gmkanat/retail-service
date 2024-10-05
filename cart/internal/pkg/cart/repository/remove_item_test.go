package repository_test

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/mocks"
	"go.uber.org/goleak"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/repository"
)

func TestCartStorageRepository_RemoveItem(t *testing.T) {
	t.Parallel()

	repo := repository.NewCartStorageRepository()
	ctx := context.Background()

	userId := int64(1)
	skuId := int64(1000)
	cartItem := &model.CartItem{
		SkuId: 1000,
		Name:  "Кроссовки Nike JORDAN",
		Count: 2,
		Price: 200,
	}

	t.Run("add item to cart", func(t *testing.T) {
		err := repo.AddItem(ctx, userId, cartItem)
		require.NoError(t, err)
	})

	t.Run("verify item in cart before removal", func(t *testing.T) {
		items, err := repo.GetCart(ctx, userId)
		require.NoError(t, err)
		require.Len(t, items, 1)
	})

	t.Run("remove item from cart", func(t *testing.T) {
		err := repo.RemoveItem(ctx, userId, skuId)
		require.NoError(t, err)
	})

	t.Run("verify cart is empty", func(t *testing.T) {
		items, err := repo.GetCart(ctx, userId)
		require.NoError(t, err)
		require.Len(t, items, 0)
	})
}

func TestCartStorageRepository_RemoveItem_Concurrent(t *testing.T) {
	defer goleak.VerifyNone(t)

	mc := minimock.NewController(t)
	repoMock := mocks.NewCartRepositoryMock(mc)

	ctx := context.Background()

	item := &model.CartItem{
		SkuId: int64(1000),
		Name:  "Кроссовки Nike JORDAN",
		Count: 1,
		Price: 200,
	}

	repoMock.AddItemMock.Expect(ctx, int64(1), item).Return(nil)
	repoMock.RemoveItemMock.Expect(ctx, int64(1), int64(1000)).Return(nil)

	wg := sync.WaitGroup{}
	goroutinesCount := 100
	for i := 0; i < goroutinesCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if i%2 == 0 {
				err := repoMock.AddItem(ctx, 1, item)
				assert.NoError(t, err)
			} else {
				_ = repoMock.RemoveItem(ctx, 1, 1000)
			}
		}()
	}

	wg.Wait()

	repoMock.GetCartMock.Expect(ctx, int64(1)).Return([]model.CartItem{}, nil)
	items, err := repoMock.GetCart(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, items, 0)
}
