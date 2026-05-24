package order_book

type PriceLevel struct {
	Price     uint64
	TotalQty  uint64 //aggregate quantity
	Head      *OrderNode
	Tail      *OrderNode
	NodeIndex map[uint64]*OrderNode
}

type OrderNode struct {
	OrderID   uint64
	Order     *Order
	PrevOrder *OrderNode
	NextOrder *OrderNode
}

func (p *PriceLevel) AddOrder(order *Order) {
	newNode := &OrderNode{
		OrderID:   order.ID,
		Order:     order,
		PrevOrder: nil,
		NextOrder: nil,
	}
	if p.Head == nil {
		p.Head = newNode
	} else {
		p.Tail.NextOrder = newNode
		newNode.PrevOrder = p.Tail
	}
	p.Tail = newNode
	p.NodeIndex[order.ID] = newNode
	p.TotalQty += order.Quantity
}

func (p *PriceLevel) CancelOrder(OrderId uint64) bool {
	OrderNode, ok := p.NodeIndex[OrderId]
	if !ok {
		return false
	}
	p.TotalQty -= OrderNode.Order.Quantity
	OrderNode.Order.Status = Cancelled
	p.RemoveNode(OrderNode)
	return true
}

func (p *PriceLevel) RemoveNode(orderNode *OrderNode) {

	if orderNode.PrevOrder != nil {
		orderNode.PrevOrder.NextOrder = orderNode.NextOrder
	} else {
		p.Head = orderNode.NextOrder
	}
	if orderNode.NextOrder != nil {
		orderNode.NextOrder.PrevOrder = orderNode.PrevOrder
	} else {
		p.Tail = orderNode.PrevOrder
	}
	orderNode.PrevOrder = nil
	orderNode.NextOrder = nil
	delete(p.NodeIndex, orderNode.OrderID)
	orderNode.Order = nil
}
