package grpc_suite

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/config"
	loms "gitlab.ozon.dev/kanat_9999/homework/loms/pkg/api/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

type PostgresSuit struct {
	suite.Suite
	client   loms.LomsClient
	conn     *grpc.ClientConn
	cancel   context.CancelFunc
	pool     *pgxpool.Pool
	orderIDs []int64
	sleepDur time.Duration
}

func (s *PostgresSuit) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	s.sleepDur = 300 * time.Millisecond
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

func (s *PostgresSuit) TearDownSuite() {
	s.conn.Close()
	s.cancel()
	s.pool.Close()
}

func (s *PostgresSuit) SetupMigration() {
	addStockDataQuery := `
		INSERT INTO stocks.stocks (id, available, reserved) VALUES
		($1, $2, $3)
		ON CONFLICT DO NOTHING
	`

	//	just generate from 100 to 110
	// 	100 - 10 reserved, 90 available
	// 	101 - 11 reserved, 89 available
	// 	102 - 12 reserved, 88 available
	// 	103 - 13 reserved, 87 available
	// 	104 - 14 reserved, 86 available
	// 	105 - 15 reserved, 85 available
	cnt := 10
	for i := 100; i < 106; i++ {
		_, err := s.pool.Exec(context.Background(), addStockDataQuery, i, 100-cnt, cnt)
		if err != nil {
			s.T().Fatal(err)
		}
		cnt += 1

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
	ctx := context.Background()
	err := s.pool.QueryRow(ctx, createOrderQuery, 1, "AwaitingPayment").Scan(&firstOrderID)
	if err != nil {
		s.T().Fatal(err)
	}

	err = s.pool.QueryRow(ctx, createOrderQuery, 1, "AwaitingPayment").Scan(&secondOrderID)
	if err != nil {
		s.T().Fatal(err)
	}

	_, err = s.pool.Exec(ctx, createItemQuery, firstOrderID, 100, 10)

	if err != nil {
		s.T().Fatal(err)
	}

	_, err = s.pool.Exec(ctx, createItemQuery, firstOrderID, 100, 20)

	if err != nil {
		s.T().Fatal(err)
	}

	_, err = s.pool.Exec(ctx, createItemQuery, secondOrderID, 101, 30)
	if err != nil {
		s.T().Fatal(err)
	}

	s.orderIDs = append(s.orderIDs, firstOrderID, secondOrderID)
}

func (s *PostgresSuit) TearDownMigration() {
	deleteOrderQuery := `
		DELETE FROM orders.orders
		WHERE id = $1
	`
	deleteItemQuery := `
		DELETE FROM orders.order_items
		WHERE order_id = $1
	`

	for _, id := range s.orderIDs {
		_, err := s.pool.Exec(context.Background(), deleteItemQuery, id)
		if err != nil {
			s.T().Fatal(err)
		}
		_, err = s.pool.Exec(context.Background(), deleteOrderQuery, id)
		if err != nil {
			s.T().Fatal(err)
		}
	}

	deleteStockQuery := `
		DELETE FROM stocks.stocks
		WHERE id >= 100 AND id < 106
	`
	_, err := s.pool.Exec(context.Background(), deleteStockQuery)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *PostgresSuit) TestOrderCreate() {
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
				UserId: 1,
				Info: &loms.OrderInfo{
					Items: []*loms.Item{
						{
							Sku:   100,
							Count: 10,
						},
					},
				},
			},
			expectResponse: &loms.OrderInfoResponse{
				Status: "AwaitingPayment",
				User:   1,
				Items: []*loms.Item{
					{
						Sku:   100,
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

			time.Sleep(s.sleepDur) // in order not to read uncommitted replica data

			require.NoError(t, err)
			infoResp, err := s.client.OrderInfo(context.Background(), &loms.OrderInfoRequest{OrderId: resp.OrderId})
			require.NoError(t, err)
			s.True(proto.Equal(tt.expectResponse, infoResp))
		})
	}
}

func (s *PostgresSuit) TestOrderPay() {
	s.SetupMigration()
	defer s.TearDownMigration()

	items := []*loms.Item{
		{
			Sku:   100,
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
					UserId: 1,
					Info: &loms.OrderInfo{
						Items: items,
					},
				})
				require.NoError(s.T(), err)
				fmt.Println(orderID)
				return &loms.OrderPayRequest{OrderId: orderID.OrderId}
			},
			expectedErr: nil,
			expectResp: &loms.OrderInfoResponse{
				Status: "Payed",
				User:   1,
				Items:  items,
			},
		},
		{
			name: "OrderPay with invalid order ID",
			req: func() *loms.OrderPayRequest {
				return &loms.OrderPayRequest{OrderId: 0}
			},
			expectedErr: fmt.Errorf("rpc error: code = FailedPrecondition desc = invalid order ID"),
			expectResp:  nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			payReq := tt.req()
			_, err := s.client.OrderPay(context.Background(), payReq)
			if tt.expectedErr != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), tt.expectedErr.Error(), err.Error())
				return
			}
			require.NoError(s.T(), err)

			time.Sleep(s.sleepDur) // Sleep to avoid reading uncommitted data from replica
			orderResp, err := s.client.OrderInfo(context.Background(), &loms.OrderInfoRequest{OrderId: payReq.OrderId})
			fmt.Println(orderResp)
			fmt.Println(tt.expectResp)
			require.NoError(s.T(), err)
			s.True(proto.Equal(tt.expectResp, orderResp))
		})
	}
}
