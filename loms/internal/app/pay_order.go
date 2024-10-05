package app

import (
	"context"
	"errors"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
	servicepb "gitlab.ozon.dev/kanat_9999/homework/loms/pkg/api/proto/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

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
