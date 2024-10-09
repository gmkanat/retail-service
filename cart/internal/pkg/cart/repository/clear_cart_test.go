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

func TestCartStorageRepository_ClearCart(t *testing.T) {
	t.Parallel()

	repo := repository.NewCartStorageRepository()
	ctx := context.Background()

	userId := int64(1)
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

	t.Run("verify item added to cart", func(t *testing.T) {
		items, err := repo.GetCart(ctx, userId)
		require.NoError(t, err)
		require.Len(t, items, 1)
	})

	t.Run("clear cart", func(t *testing.T) {
		err := repo.ClearCart(ctx, userId)
		require.NoError(t, err)
	})

	t.Run("verify cart is empty", func(t *testing.T) {
		items, err := repo.GetCart(ctx, userId)
		require.NoError(t, err)
		require.Len(t, items, 0)
	})
}

func TestCartStorageRepository_ClearCart_Concurrent(t *testing.T) {
	defer goleak.VerifyNone(t)

	repo := repository.NewCartStorageRepository()

	ctx := context.Background()

	item := &model.CartItem{
		SkuId: int64(1000),
		Name:  "Кроссовки Nike JORDAN",
		Count: 1,
		Price: 200,
	}

	err := repo.AddItem(ctx, 1, item)
	assert.NoError(t, err)

	wg := sync.WaitGroup{}
	goroutinesCount := 10
	for i := 0; i < goroutinesCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			err := repo.ClearCart(ctx, 1)
			assert.NoError(t, err)
		}()
	}

	wg.Wait()

	items, err := repo.GetCart(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, items, 0)
}
