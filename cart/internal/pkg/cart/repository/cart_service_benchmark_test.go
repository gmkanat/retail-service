package repository_test

import (
	"context"
	"testing"

	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/repository"
)

func BenchmarkAddItem(b *testing.B) {
	repo := repository.NewCartStorageRepository()
	ctx := context.Background()

	userId := int64(1)
	cartItem := &model.CartItem{
		SkuId: 1000,
		Name:  "Кроссовки Nike JORDAN",
		Count: 2,
		Price: 200,
	}

	for i := 0; i < b.N; i++ {
		err := repo.AddItem(ctx, userId, cartItem)
		if err != nil {
			b.Fatalf("failed to add item: %v", err)
		}
	}
}
