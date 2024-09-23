package main

import (
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/app/server"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/config"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/http/middleware"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/http/transport"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/repository"
	cartService "gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/service"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/loms"
	productService "gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/product/service"
	proto "gitlab.ozon.dev/kanat_9999/homework/cart/pkg/api/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()

	log.Println("app starting")

	retryTransport := transport.NewRetryRoundTripper(http.DefaultTransport, cfg.MaxRetries, cfg.InitialBackoff)
	httpClient := &http.Client{Transport: retryTransport}
	productSvc := productService.NewProductService(cfg.BaseURL, cfg.Token, httpClient)

	cartRepository := repository.NewCartStorageRepository()

	conn, err := grpc.Dial(cfg.LomsAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to LOMS service: %v", err)
	}
	defer conn.Close()

	lomsClient := proto.NewLomsClient(conn)
	lomsSvc := loms.NewClient(lomsClient)

	cartSvc := cartService.NewService(cartRepository, productSvc, lomsSvc)
	srv := server.New(cartSvc)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /user/{userId}/cart/{skuId}", srv.AddItem)
	mux.HandleFunc("GET /user/{userId}/cart", srv.GetCart)
	mux.HandleFunc("DELETE /user/{userId}/cart/{skuId}", srv.RemoveItem)
	mux.HandleFunc("DELETE /user/{userId}/cart", srv.ClearCart)
	mux.HandleFunc("POST /cart/checkout", srv.Checkout)
	logMux := middleware.LogMiddleware(mux)

	log.Println("server starting")

	if err := http.ListenAndServe(cfg.PortAddr, logMux); err != nil {
		log.Fatal(err)
	}
}
