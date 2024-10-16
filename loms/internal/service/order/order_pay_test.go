package order_test

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/mocks/order"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/mocks/stock"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/mocks/transaction"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	service "gitlab.ozon.dev/kanat_9999/homework/loms/internal/service/order"
	"testing"
)

func TestService_OrderPay(t *testing.T) {
	mc := minimock.NewController(t)

	orderRepoMock := order.NewRepositoryMock(mc)
	stockRepoMock := stock.NewRepositoryMock(mc)
	txManagerMock := transaction.NewTransactionManagerMock(mc)

	orderService := service.NewOrderService(orderRepoMock, stockRepoMock, txManagerMock)

	ctx := context.Background()
	orderID := int64(1)
	userID := int64(1)
	order := &model.Order{
		OrderID: orderID,
		UserID:  userID,
		Status:  model.OrderStatusAwaitingPayment,
		Items: []model.Item{
			{SKU: 1001, Count: 2},
		},
	}

	t.Run("add item", func(t *testing.T) {
		txManagerMock.WithRepeatableReadTxMock.Set(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
		orderRepoMock.CreateMock.Expect(ctx, userID, order.Items).Return(orderID, nil)
		stockRepoMock.ReserveMock.Expect(ctx, order.Items[0].SKU, order.Items[0].Count).Return(nil)
		orderRepoMock.SetStatusMock.Expect(ctx, orderID, model.OrderStatusAwaitingPayment).Return(nil)
		createdOrderID, err := orderService.OrderCreate(ctx, order.UserID, order.Items)
		require.NoError(t, err)
		require.Equal(t, createdOrderID, order.OrderID)
	})

	t.Run("pay for order", func(t *testing.T) {
		txManagerMock.WithRepeatableReadTxMock.Set(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
		orderRepoMock.GetByIDMock.Expect(ctx, orderID).Return(order, nil)
		stockRepoMock.ReserveRemoveMock.Expect(ctx, order.Items[0].SKU, order.Items[0].Count).Return(nil)
		orderRepoMock.SetStatusMock.Expect(ctx, orderID, model.OrderStatusPayed).Return(nil)
		err := orderService.OrderPay(ctx, orderID)
		require.NoError(t, err)
	})
}
