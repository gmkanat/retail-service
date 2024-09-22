package suite

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/app/server"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/repository"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/service"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/loms"
	productService "gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/product/service"
	proto "gitlab.ozon.dev/kanat_9999/homework/cart/pkg/api/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

type ListItemSuite struct {
	suite.Suite
	server *server.Server
	router http.Handler
}

func (s *ListItemSuite) SetupSuite() {
	repo := repository.NewCartStorageRepository()
	productSvc := productService.NewProductService("http://route256.pavl.uk:8080", "testtoken", &http.Client{})

	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to LOMS service: %v", err)
	}

	lomsClient := proto.NewLomsClient(conn)
	lomsSvc := loms.NewClient(lomsClient)

	cartSvc := service.NewService(repo, productSvc, lomsSvc)
	srv := server.New(cartSvc)
	s.router = setupRouter(srv)
}

func (s *ListItemSuite) TestListItem() {
	userId := int64(123)

	item1 := model.CartItem{
		SkuId: 773297411,
		Count: 2,
		Name:  "Кроссовки Nike JORDAN",
		Price: 2202,
	}
	s.addItemToCart(userId, item1)
	s.verifyCartContent(userId, []model.CartItem{item1})
}

func (s *ListItemSuite) addItemToCart(userID int64, item model.CartItem) {
	url := fmt.Sprintf("/user/%d/cart/%d", userID, item.SkuId)
	body := strings.NewReader(fmt.Sprintf(`{"count": %d}`, item.Count))
	req, err := http.NewRequest(http.MethodPost, url, body)
	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	s.Require().Equal(http.StatusOK, rr.Code, "Failed to add item to cart")
}

func (s *ListItemSuite) verifyCartContent(userID int64, expectedItems []model.CartItem) {
	url := fmt.Sprintf("/user/%d/cart", userID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	s.Require().Equal(http.StatusOK, rr.Code, "failed to list cart")

	var response model.GetCartResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	s.Require().NoError(err)
	s.Require().Len(response.Items, len(expectedItems), "unexpected number of items in the cart")

	totalPrice := uint32(0)
	for i, item := range expectedItems {
		s.Require().Equal(item, response.Items[i], "item mismatch")
		totalPrice += item.Price * uint32(item.Count)
	}

	s.Require().Equal(totalPrice, response.TotalPrice, "price mismatch")
}
