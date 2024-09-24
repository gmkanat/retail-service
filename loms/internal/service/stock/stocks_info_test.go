package stock_test

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/mocks/stock"
	service "gitlab.ozon.dev/kanat_9999/homework/loms/internal/service/stock"
	"testing"
)

func TestService_StocksInfo(t *testing.T) {
	mc := minimock.NewController(t)

	stockRepoMock := stock.NewRepositoryMock(mc)

	stockService := service.NewStockService(stockRepoMock)

	ctx := context.Background()

	sku := uint32(1001)

	t.Run("get stock info", func(t *testing.T) {
		stockRepoMock.GetBySKUMock.Expect(ctx, sku).Return(uint64(10), nil)
		availableCount, err := stockService.StocksInfo(ctx, sku)
		require.NoError(t, err)
		require.Equal(t, availableCount, uint64(10))
	})
}
