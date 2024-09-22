package main

import (
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/app"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/config"
	orderRepository "gitlab.ozon.dev/kanat_9999/homework/loms/internal/repository/order"
	stockRepository "gitlab.ozon.dev/kanat_9999/homework/loms/internal/repository/stock"
	orderService "gitlab.ozon.dev/kanat_9999/homework/loms/internal/service/order"
	stockService "gitlab.ozon.dev/kanat_9999/homework/loms/internal/service/stock"
	"gitlab.ozon.dev/kanat_9999/homework/loms/middleware"
	proto "gitlab.ozon.dev/kanat_9999/homework/loms/pkg/api/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	cfg := config.Load()
	initialStocks, err := config.LoadStocks(cfg.StockFile)
	if err != nil {
		log.Fatalf("failed to load stock data: %v", err)
	}

	orderRepo := orderRepository.NewOrderRepository()
	stockRepo := stockRepository.NewStockRepository(initialStocks)
	orderSvc := orderService.NewOrderService(orderRepo, stockRepo)
	stockSvc := stockService.NewStockService(stockRepo)

	server := app.NewService(orderSvc, stockSvc)

	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", cfg.Port, err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.LoggingInterceptor),
	)

	reflection.Register(grpcServer)

	proto.RegisterLomsServer(grpcServer, server)

	log.Printf("LOMS gRPC server running on port %s\n", cfg.Port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
