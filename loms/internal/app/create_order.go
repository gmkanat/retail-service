package app

import (
	"context"
	"errors"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	servicepb "gitlab.ozon.dev/kanat_9999/homework/loms/pkg/api/proto/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
