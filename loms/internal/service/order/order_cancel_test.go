package order_test

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/mocks/order"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/mocks/stock"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	service "gitlab.ozon.dev/kanat_9999/homework/loms/internal/service/order"
	"testing"
)

func TestOrderService_OrderCancel(t *testing.T) {
	mc := minimock.NewController(t)

	orderRepoMock := order.NewRepositoryMock(mc)
	stockRepoMock := stock.NewRepositoryMock(mc)

	orderService := service.NewOrderService(orderRepoMock, stockRepoMock)

	ctx := context.Background()
	orderID := int64(1)
	userID := int64(1)
	curOrder := &model.Order{
		OrderID: orderID,
		UserID:  userID,
		Status:  model.OrderStatusAwaitingPayment,
		Items: []model.Item{
			{SKU: 1001, Count: 2},
		},
	}

	t.Run("add item", func(t *testing.T) {
		orderRepoMock.CreateMock.Expect(ctx, userID, curOrder.Items).Return(orderID, nil)
		stockRepoMock.ReserveMock.Expect(ctx, curOrder.Items[0].SKU, curOrder.Items[0].Count).Return(nil)
		orderRepoMock.SetStatusMock.Expect(ctx, orderID, model.OrderStatusAwaitingPayment).Return(nil)
		createdOrderID, err := orderService.OrderCreate(ctx, curOrder.UserID, curOrder.Items)
		require.NoError(t, err)
		require.Equal(t, createdOrderID, curOrder.OrderID)
	})

	t.Run("cancel order", func(t *testing.T) {
		orderRepoMock.GetByIDMock.Expect(ctx, orderID).Return(curOrder, nil)
		stockRepoMock.ReleaseMock.Expect(ctx, curOrder.Items[0].SKU, curOrder.Items[0].Count).Return(nil)
		orderRepoMock.SetStatusMock.Expect(ctx, orderID, model.OrderStatusCancelled).Return(nil)
		err := orderService.OrderCancel(ctx, orderID)
		require.NoError(t, err)
	})
}

func TestOrderService_OrderCancelError(t *testing.T) {
	mc := minimock.NewController(t)

	orderRepoMock := order.NewRepositoryMock(mc)
	stockRepoMock := stock.NewRepositoryMock(mc)

	orderService := service.NewOrderService(orderRepoMock, stockRepoMock)

	ctx := context.Background()
	var orderID int64 = 2
	userID := int64(1)
	curOrder := &model.Order{
		OrderID: orderID,
		UserID:  userID,
		Status:  model.OrderStatusAwaitingPayment,
		Items: []model.Item{
			{SKU: 1625903, Count: 2},
		},
	}

	t.Run("add item", func(t *testing.T) {
		orderRepoMock.CreateMock.Expect(ctx, userID, curOrder.Items).Return(orderID, nil)
		stockRepoMock.ReserveMock.Expect(ctx, curOrder.Items[0].SKU, curOrder.Items[0].Count).Return(nil)
		orderRepoMock.SetStatusMock.Expect(ctx, orderID, model.OrderStatusAwaitingPayment).Return(nil)
		createdOrderID, err := orderService.OrderCreate(ctx, curOrder.UserID, curOrder.Items)
		require.NoError(t, err)
		require.Equal(t, createdOrderID, curOrder.OrderID)
	})

	t.Run("order not awaiting payment", func(t *testing.T) {
		curOrder.Status = model.OrderStatusCancelled
		orderRepoMock.GetByIDMock.Expect(ctx, orderID).Return(curOrder, nil)
		err := orderService.OrderCancel(ctx, orderID)
		require.Error(t, err, customerrors.ErrOrderStatusAwaitingPayment)
	})
}
