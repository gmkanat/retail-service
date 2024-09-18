package suite

import (
	"context"
	"github.com/stretchr/testify/suite"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/repository"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/service"
	productService "gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/product/service"
	"net/http"
)

type DeleteItemSuite struct {
	suite.Suite
	cartService *service.CartService
}

func (s *DeleteItemSuite) SetupSuite() {
	repo := repository.NewCartStorageRepository()
	productSvc := productService.NewProductService("http://route256.pavl.uk:8080", "testtoken", &http.Client{})
	s.cartService = service.NewService(repo, productSvc)
}

func (s *DeleteItemSuite) TestDeleteItem() {
	ctx := context.Background()

	userId := int64(123)
	item1 := model.CartItem{
		SkuId: 773297411,
		Count: 2,
		Name:  "Кроссовки Nike JORDAN",
		Price: 2202,
	}

	item2 := model.CartItem{
		SkuId: 3596599,
		Count: 1,
		Name:  "Невербальная коммуникация. Психология и право",
		Price: 3386,
	}

	err := s.cartService.AddItem(ctx, userId, item1.SkuId, item1.Count)
	s.Require().NoError(err)

	err = s.cartService.AddItem(ctx, userId, item2.SkuId, item2.Count)
	s.Require().NoError(err)

	items, err := s.cartService.GetCart(ctx, userId)
	s.Require().NoError(err)
	s.Require().Len(items.Items, 2)

	err = s.cartService.RemoveItem(ctx, userId, item1.SkuId)
	s.Require().NoError(err)

	items, err = s.cartService.GetCart(ctx, userId)
	s.Require().NoError(err)
	s.Require().Len(items.Items, 1)
	s.Require().Equal(item2, items.Items[0])

	s.Require().Equal(item2.Price*uint32(item2.Count), items.TotalPrice)

	err = s.cartService.RemoveItem(ctx, userId, item2.SkuId)
	s.Require().NoError(err)

	items, err = s.cartService.GetCart(ctx, userId)
	s.Require().NoError(err)
	s.Require().Len(items.Items, 0)
}
