package app

import (
	"context"
	"errors"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
	servicepb "gitlab.ozon.dev/kanat_9999/homework/loms/pkg/api/proto/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
		Status:    order.Status.String(),
		Items:     items,
		User:      order.UserID,
		CreatedAt: timestamppb.New(order.CreatedAt),
		UpdatedAt: timestamppb.New(order.UpdatedAt),
	}, nil
}
