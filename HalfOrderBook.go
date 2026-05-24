package order_book

import "fmt"

type HalfOrderBook struct {
	Side         Side
	Levels       []*PriceLevel
	bestPriceIdx int // -1 during initialization or for empty Levels
	MinPrice     uint64
	MaxPrice     uint64
}

func NewHalfOrderBook(side Side, minPrice uint64, maxPrice uint64) (*HalfOrderBook, error) {
	size := maxPrice - minPrice + 1
	if size <= 0 {
		return nil, fmt.Errorf("MaxPrice is %d and MinPrice is %d", maxPrice, minPrice)
	}
	hb := &HalfOrderBook{
		Side:     side,
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		Levels:   make([]*PriceLevel, size),
	}
	return hb, nil
}

func (hb *HalfOrderBook) priceToIdx(price uint64) (int, error) { // for price lookup
	if price < hb.MinPrice || price > hb.MaxPrice {
		return -1, fmt.Errorf("price %d is out of range", price)
	}
	return int(price - hb.MinPrice), nil
}

func (hb *HalfOrderBook) Add(order *Order) error {
	idx, err := hb.priceToIdx(order.Price)
	if err != nil {
		return err
	}
	if hb.Levels[idx] == nil {
		hb.Levels[idx] = &PriceLevel{
			Price:     order.Price,
			NodeIndex: make(map[uint64]*OrderNode),
		}
	}
	hb.Levels[idx].AddOrder(order)
	if hb.bestPriceIdx == -1 {
		hb.bestPriceIdx = idx
	} else if hb.Side == Buy && hb.bestPriceIdx < idx {
		hb.bestPriceIdx = idx
	} else if hb.Side == Sell && hb.bestPriceIdx > idx {
		hb.bestPriceIdx = idx
	}
	return nil
}

func (hb *HalfOrderBook) CancelOrder(orderId uint64, price uint64) error {
	if price < hb.MinPrice || price > hb.MaxPrice {
		return fmt.Errorf("price %d is out of range", price)
	}
	idx := int(price - hb.MinPrice)
	hb.Levels[idx].CancelOrder(orderId)
	if hb.Levels[idx].TotalQty == 0 && hb.bestPriceIdx == idx {
		if hb.Side == Sell {
			for i := idx + 1; i < len(hb.Levels); i++ {
				if hb.Levels[i] != nil {
					hb.bestPriceIdx = i
					break
				}
			}
		} else {
			for i := idx - 1; i >= 0; i-- {
				if hb.Levels[i] != nil {
					hb.bestPriceIdx = i
					break
				}
			}
		}
		if hb.bestPriceIdx == idx {
			hb.bestPriceIdx = -1
		}
		hb.Levels[idx] = nil
	}
	return nil
}
