package main

import (
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/app/server"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/config"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/http/middleware"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/http/transport"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/repository"
	cartService "gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/service"
	productService "gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/product/service"
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
	cartSvc := cartService.NewService(cartRepository, productSvc)
	srv := server.New(cartSvc)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /user/{userId}/cart/{skuId}", srv.AddItem)
	mux.HandleFunc("GET /user/{userId}/cart", srv.GetCart)
	mux.HandleFunc("DELETE /user/{userId}/cart/{skuId}", srv.RemoveItem)
	mux.HandleFunc("DELETE /user/{userId}/cart", srv.ClearCart)

	logMux := middleware.LogMiddleware(mux)

	log.Println("server starting")
	if err := http.ListenAndServe(":8082", logMux); err != nil {
		panic(err)
	}
}
