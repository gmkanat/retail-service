package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/repository"
)

func TestCartStorageRepository_AddItem(t *testing.T) {
	repo := repository.NewCartStorageRepository()
	ctx := context.Background()

	userId := int64(123)
	cartItem := &model.CartItem{
		SkuId: 1000,
		Name:  "Кроссовки Nike JORDAN",
		Count: 2,
		Price: 200,
	}

	err := repo.AddItem(ctx, userId, cartItem)
	require.NoError(t, err)

	items, err := repo.GetCart(ctx, userId)
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, cartItem.SkuId, items[0].SkuId)
	require.Equal(t, cartItem.Name, items[0].Name)
	require.Equal(t, cartItem.Price, items[0].Price)
	require.Equal(t, cartItem.Count, items[0].Count)

	err = repo.AddItem(ctx, userId, cartItem)
	require.NoError(t, err)

	items, err = repo.GetCart(ctx, userId)
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, 4, int(items[0].Count))
}
