package main

import (
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/app"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/repository"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/service"
	proto "gitlab.ozon.dev/kanat_9999/homework/loms/pkg/api/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	initialStocks := []model.Stock{
		{SKU: 773297411, TotalCount: 150, Reserved: 10},
		{SKU: 1002, TotalCount: 200, Reserved: 20},
		{SKU: 1003, TotalCount: 250, Reserved: 30},
		{SKU: 1004, TotalCount: 300, Reserved: 40},
		{SKU: 1005, TotalCount: 350, Reserved: 50},
	}

	orderRepo := repository.NewOrderRepository()
	stockRepo := repository.NewStockRepository(initialStocks)
	orderService := service.NewOrderService(orderRepo, stockRepo)
	stockService := service.NewStockService(stockRepo)

	server := app.NewService(orderService, stockService)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	proto.RegisterLomsServer(grpcServer, server)

	log.Println("LOMS gRPC server running on port 50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
