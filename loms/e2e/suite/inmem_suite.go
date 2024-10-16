package suite

import (
	"context"
	"github.com/stretchr/testify/suite"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	loms "gitlab.ozon.dev/kanat_9999/homework/loms/pkg/api/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const (
	UserID      = 123
	FirstSKUID  = 773297411
	SecondSKUID = 1002
	ThirdSKUID  = 1003
	FourthSKUID = 1004
	PaySKUID    = 1005
	CancelSKUID = 2956315
)

type InMemSuite struct {
	suite.Suite
	client loms.LomsClient
	conn   *grpc.ClientConn
	cancel context.CancelFunc
}

func (s *InMemSuite) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	conn, err := grpc.DialContext(ctx, ":50051", grpc.WithInsecure())
	if err != nil {
		s.T().Fatal(err)
	}

	s.conn = conn
	s.client = loms.NewLomsClient(conn)
}

func (s *InMemSuite) removeOrders(skuIDs ...int64) {
	for _, skuID := range skuIDs {
		_, err := s.client.OrderCancel(context.Background(), &loms.OrderCancelRequest{OrderId: skuID})
		s.Require().NoError(err)
	}
}
func (s *InMemSuite) TearDownSuite() {
	s.conn.Close()
	s.cancel()
}

var setupOrderReq = &loms.OrderCreateRequest{
	UserId: UserID,
	Info: &loms.OrderInfo{
		Items: []*loms.Item{
			{
				Sku:   FirstSKUID,
				Count: 10,
			},
			{
				Sku:   SecondSKUID,
				Count: 20,
			},
		},
	},
}

func (s *InMemSuite) TestOrderCreate() {
	tests := []struct {
		name           string
		req            *loms.OrderCreateRequest
		cleanup        func(...int64)
		expectError    error
		expectResponse *loms.OrderInfoResponse
	}{
		{
			name:    "OrderCreate",
			req:     setupOrderReq,
			cleanup: func(ids ...int64) {},
			expectResponse: &loms.OrderInfoResponse{
				Status: model.OrderStatusAwaitingPayment.String(),
				User:   UserID,
				Items: []*loms.Item{
					{
						Sku:   FirstSKUID,
						Count: 10,
					},
					{
						Sku:   SecondSKUID,
						Count: 20,
					},
				},
			},
		},
		{
			name:        "OrderCreate with invalid SKU",
			req:         &loms.OrderCreateRequest{UserId: UserID, Info: &loms.OrderInfo{Items: []*loms.Item{{Sku: 0, Count: 10}}}},
			cleanup:     func(ids ...int64) {},
			expectError: status.Errorf(codes.FailedPrecondition, "failed to create order"),
		},
		{
			name:        "OrderCreate with invalid user",
			req:         &loms.OrderCreateRequest{UserId: 0, Info: &loms.OrderInfo{Items: []*loms.Item{{Sku: FirstSKUID, Count: 10}}}},
			cleanup:     func(ids ...int64) {},
			expectError: status.Errorf(codes.FailedPrecondition, "invalid user ID"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			createResp, err := s.client.OrderCreate(context.Background(), tt.req)
			s.ErrorIs(err, tt.expectError)
			if tt.expectError == nil {
				orderResp, err := s.client.OrderInfo(context.Background(), &loms.OrderInfoRequest{
					OrderId: createResp.OrderId,
				})

				s.Require().NoError(err)
				s.True(proto.Equal(tt.expectResponse, orderResp))

				tt.cleanup(createResp.OrderId)
			}
		})
	}
}

// TestOrderPay used PaySKUID, we can't roll back it
func (s *InMemSuite) TestOrderPay() {
	var OrderId int64
	items := []*loms.Item{
		{
			Sku:   PaySKUID,
			Count: 10,
		},
	}

	tests := []struct {
		name        string
		req         func() *loms.OrderPayRequest
		expectedErr error
		expectResp  *loms.OrderInfoResponse
	}{
		{
			name: "OrderPay",
			req: func() *loms.OrderPayRequest {
				orderID, err := s.client.OrderCreate(context.Background(), &loms.OrderCreateRequest{
					UserId: UserID,
					Info: &loms.OrderInfo{
						Items: items,
					},
				})
				OrderId = orderID.OrderId
				s.Require().NoError(err)
				return &loms.OrderPayRequest{OrderId: orderID.OrderId}
			},
			expectedErr: nil,
			expectResp: &loms.OrderInfoResponse{
				Status: model.OrderStatusPayed.String(),
				User:   UserID,
				Items:  items,
			},
		},
		{
			name: "OrderPay with invalid order ID",
			req: func() *loms.OrderPayRequest {
				return &loms.OrderPayRequest{OrderId: 0}
			},
			expectedErr: status.Errorf(codes.FailedPrecondition, "invalid order ID"),
			expectResp:  nil,
		},
		{
			name: "OrderPay with not status awaiting payment",
			req: func() *loms.OrderPayRequest {
				return &loms.OrderPayRequest{OrderId: OrderId}
			},
			expectedErr: status.Errorf(codes.FailedPrecondition, customerrors.ErrOrderStatusAwaitingPayment.Error()),
			expectResp:  nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			payReq := tt.req()
			_, err := s.client.OrderPay(context.Background(), payReq)
			s.ErrorIs(err, tt.expectedErr)
			if tt.expectedErr == nil {
				orderResp, err := s.client.OrderInfo(context.Background(), &loms.OrderInfoRequest{
					OrderId: OrderId,
				})
				s.Require().NoError(err)
				s.True(proto.Equal(tt.expectResp, orderResp))
			}
		})
	}
}

func (s *InMemSuite) TestOrderCancel() {
	var OrderId int64
	items := []*loms.Item{
		{
			Sku:   CancelSKUID,
			Count: 10,
		},
	}
	tests := []struct {
		name        string
		req         func() *loms.OrderCancelRequest
		expectedErr error
		expectResp  *loms.OrderInfoResponse
	}{
		{
			name: "Order Cancel",
			req: func() *loms.OrderCancelRequest {
				orderID, err := s.client.OrderCreate(context.Background(), &loms.OrderCreateRequest{
					UserId: UserID,
					Info: &loms.OrderInfo{
						Items: items,
					},
				})
				OrderId = orderID.OrderId
				s.Require().NoError(err)
				return &loms.OrderCancelRequest{OrderId: orderID.OrderId}
			},
			expectedErr: nil,
			expectResp: &loms.OrderInfoResponse{
				Status: model.OrderStatusCancelled.String(),
				User:   UserID,
				Items:  items,
			},
		},
		{
			name: "Order Cancel with invalid order ID",
			req: func() *loms.OrderCancelRequest {
				return &loms.OrderCancelRequest{OrderId: 0}
			},
			expectedErr: status.Errorf(codes.FailedPrecondition, "invalid order ID"),
		},
		{
			name: "Order Cancel with not status awaiting payment",
			req: func() *loms.OrderCancelRequest {
				return &loms.OrderCancelRequest{OrderId: OrderId}
			},
			expectedErr: status.Errorf(codes.FailedPrecondition, customerrors.ErrOrderStatusAwaitingPayment.Error()),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			cancelReq := tt.req()
			_, err := s.client.OrderCancel(context.Background(), cancelReq)
			s.ErrorIs(err, tt.expectedErr)
			if tt.expectedErr == nil {
				orderResp, err := s.client.OrderInfo(context.Background(), &loms.OrderInfoRequest{
					OrderId: OrderId,
				})
				s.Require().NoError(err)
				s.True(proto.Equal(tt.expectResp, orderResp))
			}
		})
	}
}

func (s *InMemSuite) TestOrderInfo() {
	tests := []struct {
		name        string
		req         *loms.OrderInfoRequest
		expectedErr error
	}{
		{
			name:        "OrderInfo",
			expectedErr: nil,
		},
		{
			name:        "OrderInfo with invalid order ID",
			expectedErr: status.Errorf(codes.FailedPrecondition, "invalid order ID"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			createResp, err := s.client.OrderCreate(context.Background(), &loms.OrderCreateRequest{
				UserId: UserID,
				Info: &loms.OrderInfo{
					Items: []*loms.Item{
						{
							Sku:   FirstSKUID,
							Count: 10,
						},
						{
							Sku:   SecondSKUID,
							Count: 20,
						},
					},
				},
			})
			s.Require().NoError(err)
			if tt.expectedErr == nil {
				tt.req = &loms.OrderInfoRequest{OrderId: createResp.OrderId}
			} else {
				tt.req = &loms.OrderInfoRequest{OrderId: 0}
			}

			_, err = s.client.OrderInfo(context.Background(), tt.req)
			s.ErrorIs(err, tt.expectedErr)
		})
	}
}

func (s *InMemSuite) StocksInfo() {
	tests := []struct {
		name         string
		req          *loms.StocksInfoRequest
		expectedErr  error
		expectedResp *loms.StocksInfoResponse
	}{
		{
			name:        "StocksInfo with ThirdSKUID",
			expectedErr: nil,
			expectedResp: &loms.StocksInfoResponse{
				AvailableCount: 220, // ThirdSKUID
			},
		},
		{
			name:        "StocksInfo with FourthSKUID",
			expectedErr: nil,
			expectedResp: &loms.StocksInfoResponse{
				AvailableCount: 260, // FourthSKUID
			},
		},
		{
			name:         "StocksInfo with invalid SKU",
			expectedErr:  status.Errorf(codes.FailedPrecondition, "invalid SKU"),
			expectedResp: nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp, err := s.client.StocksInfo(context.Background(), tt.req)
			s.ErrorIs(err, tt.expectedErr)
			if tt.expectedErr == nil {
				s.True(proto.Equal(tt.expectedResp, resp))
			}
		})
	}
}
