package service_test

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/mocks"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	service "gitlab.ozon.dev/kanat_9999/homework/loms/internal/service/order"
	"testing"
)

func TestService_OrderCreate(t *testing.T) {
	mc := minimock.NewController(t)

	orderRepoMock := mocks.NewOrderRepositoryMock(mc)
	stockRepoMock := mocks.NewStockRepositoryMock(mc)

	orderService := service.NewOrderService(orderRepoMock, stockRepoMock)

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
		orderRepoMock.CreateMock.Expect(ctx, userID, order.Items).Return(orderID, nil)
		stockRepoMock.ReserveMock.Expect(ctx, order.Items[0].SKU, order.Items[0].Count).Return(nil)
		orderRepoMock.SetStatusMock.Expect(ctx, orderID, model.OrderStatusAwaitingPayment).Return(nil)
		createdOrderID, err := orderService.OrderCreate(ctx, order.UserID, order.Items)
		require.NoError(t, err)
		require.Equal(t, createdOrderID, order.OrderID)
	})
}

func TestService_OrderCreateError(t *testing.T) {
	mc := minimock.NewController(t)

	orderRepoMock := mocks.NewOrderRepositoryMock(mc)
	stockRepoMock := mocks.NewStockRepositoryMock(mc)

	orderService := service.NewOrderService(orderRepoMock, stockRepoMock)

	ctx := context.Background()
	orderID := int64(1)
	userID := int64(1)
	order := &model.Order{
		OrderID: orderID,
		UserID:  userID,
		Status:  model.OrderStatusAwaitingPayment,
		Items: []model.Item{
			{SKU: 9999, Count: 2},
		},
	}

	t.Run("add item", func(t *testing.T) {
		orderRepoMock.CreateMock.Expect(ctx, userID, order.Items).Return(orderID, nil)
		stockRepoMock.ReserveMock.Expect(ctx, order.Items[0].SKU, order.Items[0].Count).Return(customerrors.ErrOrderStatusFailed)
		orderRepoMock.SetStatusMock.Expect(ctx, orderID, model.OrderStatusFailed).Return(nil)
		createdOrderID, err := orderService.OrderCreate(ctx, order.UserID, order.Items)
		require.Error(t, err, customerrors.ErrOrderStatusFailed)
		require.Equal(t, createdOrderID, order.OrderID)
	})
}
