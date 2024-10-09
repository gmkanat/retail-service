package main

import (
	"context"
	"errors"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/product/roundtripper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	cartSvc, rateLimiter := setupServices(cfg)
	defer rateLimiter.Shutdown()

	srv := server.New(cartSvc)

	mux := setupRoutes(srv)
	logMux := middleware.LogMiddleware(mux)

	httpServer := &http.Server{
		Addr:    cfg.PortAddr,
		Handler: logMux,
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Server starting...")
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) && err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-stopChan
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully shutdown.")
}

func setupServices(cfg *config.Config) (*cartService.CartService, *roundtripper.CustomRateLimitedTransport) {
	httpClient, rateLimiter := createHTTPClient(cfg)
	productSvc := productService.NewProductService(cfg.BaseURL, cfg.Token, httpClient)

	cartRepository := repository.NewCartStorageRepository()

	lomsClient := createLomsClient(cfg.LomsAddr)
	lomsSvc := loms.NewClient(lomsClient)

	return cartService.NewService(cartRepository, productSvc, lomsSvc), rateLimiter
}

func createHTTPClient(cfg *config.Config) (*http.Client, *roundtripper.CustomRateLimitedTransport) {
	retryTransport := transport.NewRetryRoundTripper(http.DefaultTransport, cfg.MaxRetries, cfg.InitialBackoff)

	rateLimitedTransport, err := roundtripper.NewCustomRateLimitedTransport(retryTransport, cfg.RateLimit, cfg.BurstLimit)
	if err != nil {
		log.Fatalf("Failed to set rate limiter: %v", err)
	}

	return &http.Client{
		Transport: rateLimitedTransport,
	}, rateLimitedTransport
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
