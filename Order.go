package order_book

import "sync/atomic"

type Order struct {
	ID             uint64
	Side           Side
	Type           OrderType
	Price          uint64 //will save it in cents
	Quantity       uint64
	FilledQuantity uint64
	Status         OrderStatus
	ReceiveTime    uint64 //epoch millis
}

type Side int64

const (
	Buy  Side = 0
	Sell Side = 1
)

type OrderStatus uint8

const (
	Open      OrderStatus = 0
	Cancelled OrderStatus = 1
	Filled    OrderStatus = 2
	Partial   OrderStatus = 3
)

type OrderType uint8

const (
	Market OrderType = 0
	Limit  OrderType = 1
)

var globalOrderID uint64

func NextOrderID() uint64 {
	return atomic.AddUint64(&globalOrderID, 1)
}
