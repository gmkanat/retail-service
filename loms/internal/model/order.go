package model

import "time"

type Order struct {
	OrderID   int64
	UserID    int64
	Status    OrderStatus
	Items     []Item
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Item struct {
	SKU   uint32
	Count uint16
}
