package repository_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/repository"
	"testing"
)

func TestCartStorageRepository_ClearCart(t *testing.T) {
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

	items, err := repo.GetCart(ctx, userId)
	require.Len(t, items, 1)
	require.NoError(t, err)
	require.NotEmpty(t, items)

	err = repo.ClearCart(ctx, userId)
	require.NoError(t, err)

	items, err = repo.GetCart(ctx, userId)
	require.Len(t, items, 0)
	require.NoError(t, err)
	require.Empty(t, items)
}
