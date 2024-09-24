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

type LomsClient interface {
	OrderCreate(ctx context.Context, userID int64, items []model.CartItem) (orderID int64, err error)
	GetStock(ctx context.Context, skuID int64) (int64, error)
}

type CartService struct {
	repository     CartRepository
	productService ProductService
	lomsClient     LomsClient
}

func NewService(repository CartRepository, productService ProductService, lomsClient LomsClient) *CartService {
	return &CartService{
		repository:     repository,
		productService: productService,
		lomsClient:     lomsClient,
	}
}
