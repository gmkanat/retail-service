package repository_test

import (
	"context"
	"go.uber.org/goleak"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/repository"
)

func TestCartStorageRepository_GetCart(t *testing.T) {
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

	t.Run("get empty cart", func(t *testing.T) {
		items, err := repo.GetCart(ctx, userId)
		require.NoError(t, err)
		require.Len(t, items, 0)
	})

	t.Run("add item to cart", func(t *testing.T) {
		err := repo.AddItem(ctx, userId, cartItem)
		require.NoError(t, err)
	})

	t.Run("get cart after adding item", func(t *testing.T) {
		items, err := repo.GetCart(ctx, userId)
		require.NoError(t, err)
		require.NotEmpty(t, items)
		require.Len(t, items, 1)
		require.Equal(t, cartItem.SkuId, items[0].SkuId)
		require.Equal(t, cartItem.Name, items[0].Name)
		require.Equal(t, cartItem.Count, items[0].Count)
		require.Equal(t, cartItem.Price, items[0].Price)
	})
}

func TestCartStorageRepository_GetCart_Concurrent(t *testing.T) {
	defer goleak.VerifyNone(t)

	repo := repository.NewCartStorageRepository()
	ctx := context.Background()

	userId := int64(1)
	cartItem := &model.CartItem{
		SkuId: 1000,
		Name:  "Кроссовки Nike JORDAN",
		Count: 2,
		Price: 200,
	}

	err := repo.AddItem(ctx, userId, cartItem)
	require.NoError(t, err)

	goroutineCount := 10
	wg := sync.WaitGroup{}

	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			items, err := repo.GetCart(ctx, userId)
			require.NoError(t, err)
			require.NotEmpty(t, items)
			require.Len(t, items, 1)
			require.Equal(t, cartItem.SkuId, items[0].SkuId)
			require.Equal(t, cartItem.Name, items[0].Name)
			require.Equal(t, cartItem.Count, items[0].Count)
			require.Equal(t, cartItem.Price, items[0].Price)
		}()
	}

	wg.Wait()
}
