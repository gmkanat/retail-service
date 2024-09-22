package model

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
