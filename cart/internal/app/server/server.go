package server

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
)

type CartService interface {
	AddItem(ctx context.Context, userId int64, skuId int64, count uint16) error
	RemoveItem(ctx context.Context, userId int64, skuId int64) error
	GetCart(ctx context.Context, userId int64) (*model.GetCartResponse, error)
	ClearCart(ctx context.Context, userId int64) error
	Checkout(ctx context.Context, userId int64) (orderId int64, error error)
}

type Server struct {
	cartService CartService
}

func New(cartService CartService) *Server {
	return &Server{cartService: cartService}
}
