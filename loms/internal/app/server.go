package app

import (
	"context"
	"errors"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	servicepb "gitlab.ozon.dev/kanat_9999/homework/loms/pkg/api/proto/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ servicepb.LomsServer = (*Service)(nil)

type OrderService interface {
	OrderCreate(ctx context.Context, userID int64, items []model.Item) (int64, error)
	OrderInfo(ctx context.Context, orderID int64) (*model.Order, error)
	OrderPay(ctx context.Context, orderID int64) error
	OrderCancel(ctx context.Context, orderID int64) error
}

type StockService interface {
	StocksInfo(ctx context.Context, sku uint32) (uint64, error)
}

type Service struct {
	servicepb.UnimplementedLomsServer
	orderService OrderService
	stockService StockService
}

func NewService(orderService OrderService, stockService StockService) *Service {
	return &Service{
		orderService: orderService,
		stockService: stockService,
	}
}

func (s *Service) OrderCreate(ctx context.Context, req *servicepb.OrderCreateRequest) (*servicepb.OrderCreateResponse, error) {
	items := make([]model.Item, 0, len(req.Info.Items))

	for _, item := range req.Info.Items {
		items = append(items, model.Item{
			SKU:   item.Sku,
			Count: uint16(item.Count),
		})
	}

	orderID, err := s.orderService.OrderCreate(ctx, req.UserId, items)
	if err != nil {
		if errors.Is(err, customerrors.ErrOrderStatusFailed) {
			return nil, status.Errorf(codes.FailedPrecondition, "failed to create order")
		}
		if errors.Is(err, customerrors.ErrInvalidUserId) {
			return nil, status.Errorf(codes.FailedPrecondition, "invalid user ID")
		}
		return nil, status.Errorf(codes.Internal, "failed to create order")
	}

	return &servicepb.OrderCreateResponse{OrderId: orderID}, nil
}

func (s *Service) OrderInfo(ctx context.Context, req *servicepb.OrderInfoRequest) (*servicepb.OrderInfoResponse, error) {
	order, err := s.orderService.OrderInfo(ctx, req.OrderId)
	if err != nil {
		if errors.Is(err, customerrors.ErrOrderNotFound) {
			return nil, status.Errorf(codes.NotFound, "order not found")
		}
		if errors.Is(err, customerrors.ErrInvalidOrderId) {
			return nil, status.Errorf(codes.FailedPrecondition, "invalid order ID")
		}
		return nil, status.Errorf(codes.Internal, "failed to get order info")
	}

	items := make([]*servicepb.Item, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, &servicepb.Item{
			Sku:   item.SKU,
			Count: uint32(item.Count),
		})
	}

	return &servicepb.OrderInfoResponse{
		Status: order.Status.String(),
		Items:  items,
		User:   order.UserID,
	}, nil
}

func (s *Service) OrderPay(ctx context.Context, req *servicepb.OrderPayRequest) (*emptypb.Empty, error) {
	err := s.orderService.OrderPay(ctx, req.OrderId)
	if err != nil {
		if errors.Is(err, customerrors.ErrInvalidOrderId) {
			return nil, status.Errorf(codes.FailedPrecondition, "invalid order ID")
		}
		if errors.Is(err, customerrors.ErrOrderStatusAwaitingPayment) {
			return nil, status.Errorf(codes.FailedPrecondition, "order status is not awaiting payment")
		}
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Service) OrderCancel(ctx context.Context, req *servicepb.OrderCancelRequest) (*emptypb.Empty, error) {
	err := s.orderService.OrderCancel(ctx, req.OrderId)
	if err != nil {
		if errors.Is(err, customerrors.ErrInvalidOrderId) {
			return nil, status.Errorf(codes.FailedPrecondition, "invalid order ID")
		}
		if errors.Is(err, customerrors.ErrOrderStatusAwaitingPayment) {
			return nil, status.Errorf(codes.FailedPrecondition, "order status is not awaiting payment")
		}
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Service) StocksInfo(ctx context.Context, req *servicepb.StocksInfoRequest) (*servicepb.StocksInfoResponse, error) {
	stock, err := s.stockService.StocksInfo(ctx, req.Sku)
	if err != nil {
		return nil, err
	}

	return &servicepb.StocksInfoResponse{AvailableCount: int64(stock)}, nil
}
