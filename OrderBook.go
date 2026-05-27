package order_book

import "sync"

type OrderBook struct {
	symbol   string
	rw       sync.RWMutex
	bids     *HalfOrderBook
	asks     *HalfOrderBook
	orderMap map[uint64]*OrderRecord
	seqNum   uint64
}

type OrderRecord struct {
	price int64
	side  Side
}

func NewOrderBook(symbol string, minPrice int64, maxPrice int64) *OrderBook {
	bids, err := NewHalfOrderBook(Buy, minPrice, maxPrice)
	if err != nil {
		return nil
	}

	asks, err := NewHalfOrderBook(Sell, minPrice, maxPrice)
	if err != nil {
		return nil
	}
	orderBook := &OrderBook{
		symbol:   symbol,
		bids:     bids,
		asks:     asks,
		orderMap: make(map[uint64]*OrderRecord),
	}
	return orderBook
}
