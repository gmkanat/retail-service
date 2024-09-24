package loms

import (
	loms "gitlab.ozon.dev/kanat_9999/homework/cart/pkg/api/proto/v1"
)

type Client struct {
	loms.LomsClient
}

func NewClient(client loms.LomsClient) *Client {
	return &Client{client}
}
