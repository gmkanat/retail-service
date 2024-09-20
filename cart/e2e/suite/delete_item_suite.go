package suite

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/app/server"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/repository"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/service"
	productService "gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/product/service"
	"net/http"
	"net/http/httptest"
	"strings"
)

type DeleteItemSuite struct {
	suite.Suite
	server *server.Server
	router *http.ServeMux
}

func (s *DeleteItemSuite) SetupSuite() {
	repo := repository.NewCartStorageRepository()
	productSvc := productService.NewProductService("http://route256.pavl.uk:8080", "testtoken", &http.Client{})
	cartService := service.NewService(repo, productSvc)
	srv := server.New(cartService)
	s.router = setupRouter(srv)
}

func (s *DeleteItemSuite) TestDeleteItem() {
	userId := int64(23)
	skuId := int64(773297411)
	count := uint16(2)

	s.addItemToCart(userId, skuId, count)
	s.removeItemFromCart(userId, skuId)
	s.verifyCartIsEmpty(userId)
}

func (s *DeleteItemSuite) addItemToCart(userId, skuId int64, count uint16) {
	url := fmt.Sprintf("/user/%d/cart/%d", userId, skuId)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(fmt.Sprintf(`{"count": %d}`, count)))
	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	s.Require().Equal(http.StatusOK, rr.Code)
}

func (s *DeleteItemSuite) removeItemFromCart(userId, skuId int64) {
	url := fmt.Sprintf("/user/%d/cart/%d", userId, skuId)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	s.Require().Equal(http.StatusNoContent, rr.Code)
}

func (s *DeleteItemSuite) verifyCartIsEmpty(userId int64) {
	url := fmt.Sprintf("/user/%d/cart", userId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	s.Require().Equal(http.StatusNotFound, rr.Code)
}
