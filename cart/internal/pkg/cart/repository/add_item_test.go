package repository_test

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/mocks"
	"go.uber.org/goleak"
	"sync"
	"testing"

	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/repository"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestCartStorageRepository_AddItem(t *testing.T) {
	t.Parallel()

	repo := repository.NewCartStorageRepository()
	ctx := context.Background()

	userId := int64(123)
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
		require.Equal(t, 1, len(items))
		require.Equal(t, cartItem.SkuId, items[0].SkuId)
		require.Equal(t, cartItem.Name, items[0].Name)
		require.Equal(t, cartItem.Count, items[0].Count)
		require.Equal(t, cartItem.Price, items[0].Price)
	})

	t.Run("add item to cart with the same SkuId", func(t *testing.T) {
		err := repo.AddItem(ctx, userId, cartItem)
		require.NoError(t, err)
	})

	t.Run("get cart after adding same item", func(t *testing.T) {
		items, err := repo.GetCart(ctx, userId)
		require.NoError(t, err)
		require.Len(t, items, 1)
		require.Equal(t, 4, int(items[0].Count))
	})
}

func TestCartStorageRepository_AddItem_Concurrent(t *testing.T) {
	defer goleak.VerifyNone(t)

	mc := minimock.NewController(t)
	repoMock := mocks.NewCartRepositoryMock(mc)

	ctx := context.Background()
	itemCount := 200
	expectedItems := make([]model.CartItem, itemCount)

	for i := 0; i < itemCount; i++ {
		curItem := &model.CartItem{
			SkuId: int64(1000),
			Name:  "Кроссовки Nike JORDAN",
			Count: 1,
			Price: 200,
		}
		expectedItems[i] = *curItem
		repoMock.AddItemMock.Expect(ctx, int64(1), curItem).Return(nil)
	}

	wg := sync.WaitGroup{}
	for i := 0; i < itemCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := repoMock.AddItem(ctx, 1, &expectedItems[i])
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	repoMock.GetCartMock.Expect(ctx, int64(1)).Return(expectedItems, nil)
	items, err := repoMock.GetCart(ctx, 1)

	assert.NoError(t, err)
	assert.Len(t, items, itemCount)
}
