package model

import "fmt"

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
		return "new"
	case OrderStatusAwaitingPayment:
		return "awaiting payment"
	case OrderStatusFailed:
		return "failed"
	case OrderStatusPayed:
		return "payed"
	case OrderStatusCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}

func ParseOrderStatus(status string) (OrderStatus, error) {
	switch status {
	case "new":
		return OrderStatusNew, nil
	case "awaiting payment":
		return OrderStatusAwaitingPayment, nil
	case "failed":
		return OrderStatusFailed, nil
	case "payed":
		return OrderStatusPayed, nil
	case "cancelled":
		return OrderStatusCancelled, nil
	default:
		return -1, fmt.Errorf("invalid order status: %s", status)
	}
}
