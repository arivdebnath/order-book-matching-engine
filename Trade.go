package order_book

type Trade struct {
	MakerOrderId uint64
	TakerOrderId uint64
	Quantity     uint64
	Price        int64
	Timestamp    int64
}
