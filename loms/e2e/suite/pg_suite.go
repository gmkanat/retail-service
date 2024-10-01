package grpc_suite

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/config"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	loms "gitlab.ozon.dev/kanat_9999/homework/loms/pkg/api/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

type PostgresSuite struct {
	suite.Suite
	client   loms.LomsClient
	conn     *grpc.ClientConn
	cancel   context.CancelFunc
	pool     *pgxpool.Pool
	orderIDs []int64
	sleepDur time.Duration
}

func (s *PostgresSuite) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	s.sleepDur = 100 * time.Millisecond
	cfg := config.Load()

	conn, err := grpc.DialContext(ctx, cfg.GRPCPort, grpc.WithInsecure())
	if err != nil {
		s.T().Fatal(err)
	}

	s.conn = conn

	s.client = loms.NewLomsClient(conn)

	pool, err := pgxpool.New(ctx, cfg.MasterDBURL)
	if err != nil {
		s.T().Fatal(err)
	}

	s.pool = pool
}

func (s *PostgresSuite) TearDownSuite() {
	s.conn.Close()
	s.cancel()
	s.pool.Close()
}

// SetupMigration
// 12345678 -> total 200, 50 reserved, 150 available,
// 23456789 -> total 100, 70 reserved, 30 available
// 34567890 -> total 50, 20 reserved, 30 available
// we have 2 orders,
func (s *PostgresSuite) SetupMigration() {
	ctx := context.Background()

	userID := 999
	createStockQuery := `
		INSERT INTO stocks.stocks (id, available, reserved) VALUES
		(12345678, 150, 50),
		(23456789, 30, 70),
		(34567890, 30, 20)
		ON CONFLICT (id) DO UPDATE SET available = EXCLUDED.available, reserved = EXCLUDED.reserved
	`
	_, err := s.pool.Exec(ctx, createStockQuery)
	if err != nil {
		s.T().Fatal(err)
	}

	createOrderQuery := `
		INSERT INTO orders.orders (user_id, status) VALUES
		($1, $2)
		RETURNING id
	`
	createItemQuery := `
		INSERT INTO orders.order_items (order_id, sku_id, count) VALUES
		($1, $2, $3)
	`

	var firstOrderID, secondOrderID int64

	// first order
	err = s.pool.QueryRow(ctx, createOrderQuery, userID, model.OrderStatusAwaitingPayment.String()).Scan(&firstOrderID)
	if err != nil {
		s.T().Fatal(err)
	}

	_, err = s.pool.Exec(ctx, createItemQuery, firstOrderID, 12345678, 10)

	if err != nil {
		s.T().Fatal(err)
	}

	_, err = s.pool.Exec(ctx, createItemQuery, firstOrderID, 23456789, 20)

	if err != nil {
		s.T().Fatal(err)
	}

	// second order
	err = s.pool.QueryRow(ctx, createOrderQuery, userID, model.OrderStatusAwaitingPayment.String()).Scan(&secondOrderID)
	if err != nil {
		s.T().Fatal(err)
	}

	_, err = s.pool.Exec(ctx, createItemQuery, secondOrderID, 34567890, 5)
	if err != nil {
		s.T().Fatal(err)
	}

	s.orderIDs = append(s.orderIDs, firstOrderID, secondOrderID)
}

func (s *PostgresSuite) TearDownMigration() {
	ctx := context.Background()

	deleteOrderItemsQuery := `
		DELETE FROM orders.order_items WHERE order_id = $1
	`
	deleteOrdersQuery := `
		DELETE FROM orders.orders WHERE id = $1
	`
	deleteStocksQuery := `
		DELETE FROM stocks.stocks WHERE id IN (12345678, 23456789, 34567890)
	`

	for _, orderID := range s.orderIDs {
		_, err := s.pool.Exec(ctx, deleteOrderItemsQuery, orderID)
		if err != nil {
			s.T().Fatal(err)
		}

		_, err = s.pool.Exec(ctx, deleteOrdersQuery, orderID)
		if err != nil {
			s.T().Fatal(err)
		}
	}

	_, err := s.pool.Exec(ctx, deleteStocksQuery)
	if err != nil {
		s.T().Fatal(err)
	}

	s.orderIDs = nil
}

func (s *PostgresSuite) TestOrderCreate() {
	s.SetupMigration()
	defer s.TearDownMigration()

	tests := []struct {
		name           string
		req            *loms.OrderCreateRequest
		expectError    error
		expectResponse *loms.OrderInfoResponse
	}{
		{
			name: "OrderCreate",
			req: &loms.OrderCreateRequest{
				UserId: 999,
				Info: &loms.OrderInfo{
					Items: []*loms.Item{
						{
							Sku:   12345678,
							Count: 10,
						},
					},
				},
			},
			expectResponse: &loms.OrderInfoResponse{
				Status: "AwaitingPayment",
				User:   999,
				Items: []*loms.Item{
					{
						Sku:   12345678,
						Count: 10,
					},
				},
			},
		},
		{
			name: "OrderCreateWithInvalidUserID",
			req: &loms.OrderCreateRequest{
				UserId: 0,
				Info: &loms.OrderInfo{
					Items: []*loms.Item{
						{
							Sku:   100,
							Count: 10,
						},
					},
				},
			},
			expectError: fmt.Errorf("rpc error: code = FailedPrecondition desc = invalid user ID"),
		},
		{
			name: "OrderCreateWithInvalidSku",
			req: &loms.OrderCreateRequest{
				UserId: 1,
				Info: &loms.OrderInfo{
					Items: []*loms.Item{
						{
							Sku:   0,
							Count: 10,
						},
					},
				},
			},
			expectError: fmt.Errorf("rpc error: code = FailedPrecondition desc = failed to create order"),
		},
	}

	for _, tt := range tests {
		tt := tt
		s.T().Run(tt.name, func(t *testing.T) {
			resp, err := s.client.OrderCreate(context.Background(), tt.req)
			if tt.expectError != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectError.Error(), err.Error())
				return
			}
			s.orderIDs = append(s.orderIDs, resp.OrderId)
			time.Sleep(s.sleepDur) // in order not to read uncommitted replica data

			require.NoError(t, err)
			infoResp, err := s.client.OrderInfo(context.Background(), &loms.OrderInfoRequest{OrderId: resp.OrderId})
			require.NoError(t, err)
			require.Equal(t, tt.expectResponse.Status, infoResp.Status)
			require.Equal(t, tt.expectResponse.User, infoResp.User)
			require.Equal(t, tt.expectResponse.Items, infoResp.Items)
		})
	}
}

func (s *PostgresSuite) TestOrderPay() {
	s.SetupMigration()
	defer s.TearDownMigration()

	tests := []struct {
		name        string
		req         func() *loms.OrderPayRequest
		expectedErr error
		expectResp  *loms.OrderInfoResponse
	}{
		{
			name: "OrderPay",
			req: func() *loms.OrderPayRequest {
				return &loms.OrderPayRequest{
					OrderId: s.orderIDs[1],
				}
			},
			expectedErr: nil,
			expectResp: &loms.OrderInfoResponse{
				Status: "Paid",
				User:   999,
				Items: []*loms.Item{
					{
						Sku:   34567890,
						Count: 5,
					},
				},
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
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			time.Sleep(s.sleepDur) // in order not to read uncommitted replica data
			payReq := tt.req()
			_, err := s.client.OrderPay(context.Background(), payReq)
			s.ErrorIs(err, tt.expectedErr)

			time.Sleep(s.sleepDur) // in order not to read uncommitted replica data

			if err == nil {
				infoResp, err := s.client.OrderInfo(context.Background(), &loms.OrderInfoRequest{OrderId: payReq.OrderId})
				require.NoError(s.T(), err)
				require.Equal(s.T(), tt.expectResp.Status, infoResp.Status)
				require.Equal(s.T(), tt.expectResp.Items, infoResp.Items)
			}
		})
	}
}

func (s *PostgresSuite) TestOrderCancel() {
	s.SetupMigration()
	defer s.TearDownMigration()

	tests := []struct {
		name        string
		req         func() *loms.OrderCancelRequest
		expectedErr error
	}{
		{
			name: "OrderCancel",
			req: func() *loms.OrderCancelRequest {
				return &loms.OrderCancelRequest{
					OrderId: s.orderIDs[0],
				}
			},
			expectedErr: nil,
		},
		{
			name: "OrderCancel with invalid order ID",
			req: func() *loms.OrderCancelRequest {
				return &loms.OrderCancelRequest{OrderId: 0}
			},
			expectedErr: status.Errorf(codes.FailedPrecondition, "invalid order ID"),
		},
		{
			name: "OrderCancel with not status awaiting payment",
			req: func() *loms.OrderCancelRequest {
				return &loms.OrderCancelRequest{OrderId: s.orderIDs[0]}
			},
			expectedErr: status.Errorf(codes.FailedPrecondition, "order status is not awaiting payment"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			time.Sleep(s.sleepDur) // in order not to read uncommitted replica data
			cancelReq := tt.req()
			_, err := s.client.OrderCancel(context.Background(), cancelReq)
			s.ErrorIs(err, tt.expectedErr)

			time.Sleep(s.sleepDur) // in order not to read uncommitted replica data

			if err == nil {
				infoResp, err := s.client.OrderInfo(context.Background(), &loms.OrderInfoRequest{OrderId: cancelReq.OrderId})
				require.NoError(s.T(), err)
				require.Equal(s.T(), "Cancelled", infoResp.Status)
			}
		})
	}
}

func (s *PostgresSuite) TestOrderInfo() {
	s.SetupMigration()
	defer s.TearDownMigration()

	tests := []struct {
		name        string
		req         *loms.OrderInfoRequest
		expectedErr error
		expectResp  *loms.OrderInfoResponse
	}{
		{
			name: "OrderInfo",
			req: &loms.OrderInfoRequest{
				OrderId: s.orderIDs[0],
			},
			expectResp: &loms.OrderInfoResponse{
				Status: "AwaitingPayment",
				User:   999,
				Items: []*loms.Item{
					{
						Sku:   12345678,
						Count: 10,
					},
					{
						Sku:   23456789,
						Count: 20,
					},
				},
			},
		},
		{
			name: "OrderInfo with invalid order ID",
			req: &loms.OrderInfoRequest{
				OrderId: 0,
			},
			expectedErr: status.Errorf(codes.FailedPrecondition, "invalid order ID"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			time.Sleep(s.sleepDur) // in order not to read uncommitted replica data
			infoResp, err := s.client.OrderInfo(context.Background(), tt.req)

			if tt.expectedErr != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), tt.expectedErr.Error(), err.Error())
				return
			}

			require.NoError(s.T(), err)
			require.Equal(s.T(), tt.expectResp.Status, infoResp.Status)
			require.Equal(s.T(), tt.expectResp.User, infoResp.User)
			require.Equal(s.T(), tt.expectResp.Items, infoResp.Items)
		})
	}
}

func (s *PostgresSuite) TestStocksInfo() {
	s.SetupMigration()
	defer s.TearDownMigration()

	tests := []struct {
		name        string
		req         *loms.StocksInfoRequest
		expectedErr error
		expectResp  *loms.StocksInfoResponse
	}{
		{
			name: "StocksInfo",
			req: &loms.StocksInfoRequest{
				Sku: 12345678,
			},
			expectResp: &loms.StocksInfoResponse{
				AvailableCount: 150,
			},
		},
		{
			name: "StocksInfo with invalid SKU",
			req: &loms.StocksInfoRequest{
				Sku: 0,
			},
			expectedErr: status.Errorf(codes.Unknown, "stock not found"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			time.Sleep(s.sleepDur) // in order not to read uncommitted replica data
			resp, err := s.client.StocksInfo(context.Background(), tt.req)
			s.ErrorIs(err, tt.expectedErr)
			if tt.expectedErr == nil {
				s.True(proto.Equal(tt.expectResp, resp))
			}
		})
	}
}
