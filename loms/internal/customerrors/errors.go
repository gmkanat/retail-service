package customerrors

import "errors"

var (
	ErrOrderNotFound              = errors.New("order not found")
	ErrOrderStatusFailed          = errors.New("order status failed")
	ErrStockNotFound              = errors.New("stock not found")
	ErrNotEnoughReservedStock     = errors.New("not enough reserved stock")
	ErrInsufficientStock          = errors.New("insufficient stock")
	ErrInvalidOrderId             = errors.New("invalid order id")
	ErrInvalidUserId              = errors.New("invalid user id")
	ErrOrderStatusAwaitingPayment = errors.New("order status is not awaiting payment")
)
