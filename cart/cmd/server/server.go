package main

import (
	server2 "gitlab.ozon.dev/kanat_9999/homework/cart/internal/app/server"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/http/middleware"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/repository"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/service"
	service2 "gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/product/service"
	"log"
	"net/http"
)

func main() {
	log.Println("app starting")

	baseURL := "http://route256.pavl.uk:8080"
	token := "testtoken"

	productService := service2.NewProductService(baseURL, token)

	cartRepository := repository.NewCartStorageRepository()
	cartService := service.NewService(cartRepository, productService)
	server := server2.New(cartService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /user/{userId}/cart/{skuId}", server.AddItem)
	mux.HandleFunc("GET /user/{userId}/cart", server.GetCart)
	mux.HandleFunc("DELETE /user/{userId}/cart/{skuId}", server.RemoveItem)
	mux.HandleFunc("DELETE /user/{userId}/cart", server.ClearCart)

	logMux := middleware.LogMiddleware(mux)

	log.Println("server starting")
	if err := http.ListenAndServe(":8082", logMux); err != nil {
		panic(err)
	}
}
