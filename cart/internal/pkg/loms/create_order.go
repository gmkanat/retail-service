package loms

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
)

func (c *Client) OrderCreate(ctx context.Context, userId int64, items []model.CartItem) (orderId int64, err error) {
	resp, err := c.LomsClient.OrderCreate(ctx, ToOrderCreateRequest(userId, items))

	if err != nil {
		return 0, err
	}

	return resp.OrderId, nil
}
