package main

import (
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/app/server"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/config"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/http/middleware"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/http/transport"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/repository"
	cartService "gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/service"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/loms"
	productService "gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/product/service"
	proto "gitlab.ozon.dev/kanat_9999/homework/cart/pkg/api/proto/v1"
)

func main() {
	cfg := config.Load()
	log.Println("App starting")

	cartSvc := setupServices(cfg)
	srv := server.New(cartSvc)

	mux := setupRoutes(srv)
	logMux := middleware.LogMiddleware(mux)

	log.Println("Server starting")
	if err := http.ListenAndServe(cfg.PortAddr, logMux); err != nil {
		log.Fatal(err)
	}
}

func setupServices(cfg *config.Config) *cartService.CartService {
	httpClient := createHTTPClient(cfg)
	productSvc := productService.NewProductService(cfg.BaseURL, cfg.Token, httpClient)

	cartRepository := repository.NewCartStorageRepository()

	lomsClient := createLomsClient(cfg.LomsAddr)
	lomsSvc := loms.NewClient(lomsClient)

	return cartService.NewService(cartRepository, productSvc, lomsSvc)
}

func createHTTPClient(cfg *config.Config) *http.Client {
	retryTransport := transport.NewRetryRoundTripper(http.DefaultTransport, cfg.MaxRetries, cfg.InitialBackoff)
	return &http.Client{Transport: retryTransport}
}

func createLomsClient(lomsAddr string) proto.LomsClient {
	conn, err := grpc.Dial(lomsAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to LOMS service: %v", err)
	}
	return proto.NewLomsClient(conn)
}

func setupRoutes(srv *server.Server) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{userId}/cart/{skuId}", srv.AddItem)
	mux.HandleFunc("GET /user/{userId}/cart", srv.GetCart)
	mux.HandleFunc("DELETE /user/{userId}/cart/{skuId}", srv.RemoveItem)
	mux.HandleFunc("DELETE /user/{userId}/cart", srv.ClearCart)
	mux.HandleFunc("POST /cart/checkout", srv.Checkout)
	return mux
}
