package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	orderRepository "gitlab.ozon.dev/kanat_9999/homework/loms/internal/repository_raw/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"

	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/app"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/config"
	stockRepository "gitlab.ozon.dev/kanat_9999/homework/loms/internal/repository_raw/stock"
	orderService "gitlab.ozon.dev/kanat_9999/homework/loms/internal/service/order"
	stockService "gitlab.ozon.dev/kanat_9999/homework/loms/internal/service/stock"
	"gitlab.ozon.dev/kanat_9999/homework/loms/middleware"
	proto "gitlab.ozon.dev/kanat_9999/homework/loms/pkg/api/proto/v1"
)

func main() {
	cfg := config.Load()

	server := setupServices(cfg)

	go startGRPCServer(cfg, server)
	startHTTPServer(cfg, server)
}

func setupServices(cfg *config.AppConfig) *app.Service {
	conn, err := pgx.Connect(context.Background(), cfg.DataBaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	orderRepo := orderRepository.NewRepository(conn)
	stockRepo := stockRepository.NewRepository(conn)

	orderSvc := orderService.NewOrderService(orderRepo, stockRepo)
	stockSvc := stockService.NewStockService(stockRepo)

	return app.NewService(orderSvc, stockSvc)
}

func startGRPCServer(cfg *config.AppConfig, server *app.Service) {
	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", cfg.GRPCPort, err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.LoggingInterceptor),
	)

	reflection.Register(grpcServer)
	proto.RegisterLomsServer(grpcServer, server)

	log.Printf("LOMS gRPC server running on port %s\n", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func startHTTPServer(cfg *config.AppConfig, server *app.Service) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := proto.RegisterLomsHandlerFromEndpoint(ctx, mux, cfg.GRPCPort, opts)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	httpMux := setupHTTPHandlers(cfg, mux)

	gwServer := &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: httpMux,
	}

	log.Printf("LOMS HTTP server running on port %s\n", cfg.HTTPPort)
	log.Fatal(gwServer.ListenAndServe())
}

func setupHTTPHandlers(cfg *config.AppConfig, gwmux *runtime.ServeMux) *http.ServeMux {
	httpMux := http.NewServeMux()

	httpMux.Handle("/swagger/", httpSwagger.WrapHandler)
	httpMux.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, cfg.SwaggerFile)
	})
	httpMux.Handle("/", middleware.LogMiddleware(gwmux))

	return httpMux
}
