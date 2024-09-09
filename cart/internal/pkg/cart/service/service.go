package service

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/product/service"
)

type CartRepository interface {
	AddItem(_ context.Context, userId int64, cartItem *model.CartItem) error
	RemoveItem(_ context.Context, userId int64, sku int64) error
	GetCart(_ context.Context, userId int64) ([]model.CartItem, error)
	ClearCart(_ context.Context, userId int64) error
}

type ProductService interface {
	GetProduct(ctx context.Context, skuID int64) (*service.Product, error)
}

type CartService struct {
	repository     CartRepository
	productService ProductService
}

func NewService(repository CartRepository, productService ProductService) *CartService {
	return &CartService{
		repository:     repository,
		productService: productService,
	}
}
