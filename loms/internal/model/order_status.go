package model

import "errors"

type OrderStatus int

const (
	OrderStatusNew OrderStatus = iota
	OrderStatusAwaitingPayment
	OrderStatusFailed
	OrderStatusPayed
	OrderStatusCancelled
)

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusNew:
		return "New"
	case OrderStatusAwaitingPayment:
		return "AwaitingPayment"
	case OrderStatusPayed:
		return "Paid"
	case OrderStatusCancelled:
		return "Cancelled"
	case OrderStatusFailed:
		return "Failed"
	default:
		return "unknown"
	}
}

func OrderStatusFromString(s string) (OrderStatus, error) {
	switch s {
	case "New":
		return OrderStatusNew, nil
	case "AwaitingPayment":
		return OrderStatusAwaitingPayment, nil
	case "Paid":
		return OrderStatusPayed, nil
	case "Cancelled":
		return OrderStatusCancelled, nil
	case "Failed":
		return OrderStatusFailed, nil
	default:
		return 0, errors.New("invalid order status")
	}
}
