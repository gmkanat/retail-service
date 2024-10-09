package repository_test

import (
	"context"
	"github.com/stretchr/testify/assert"
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

	repo := repository.NewCartStorageRepository()

	ctx := context.Background()

	goroutinesCount := 100

	wg := sync.WaitGroup{}
	for i := 0; i < goroutinesCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			item := &model.CartItem{
				SkuId: int64(1000 + i), // unique sku_id
				Name:  "Кроссовки Nike JORDAN",
				Count: 1,
				Price: 200,
			}

			err := repo.AddItem(ctx, 1, item)
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	for i := 0; i < goroutinesCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := repo.RemoveItem(ctx, 1, int64(1000+i))
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	items, err := repo.GetCart(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, items, 0)
}
