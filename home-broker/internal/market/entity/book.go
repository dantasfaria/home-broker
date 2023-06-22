package entity

import (
	"container/heap"
	"hash"
	"sync"
)

type Book struct {
	Order         []*Order
	Transactions  []*Transaction
	OrdersChan    chan *Order
	OrdersChanOut chan *Order
	Wg            *sync.WaitGroup
}

func NewBook(orderChan chan *Order, orderChanOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Order:         []*Order{},
		Transactions:  []*Transaction{},
		OrdersChan:    orderChan,
		OrdersChanOut: orderChanOut,
		Wg:            wg,
	}
}

func (b *Book) Trade() {
	buyOrders := NewOrderQueue()
	sellOrders := NewOrderQueue()

	heap.Init(buyOrders)
	heap.Init(sellOrders)

	for order := range b.OrdersChan {
		if order.OrderType == "BUY" {
			buyOrders.Push(order)

			if sellOrders.Len() > 0 && sellOrders[0].Price <= order.Price
			sellOrder := sellOrders.Pop().(*Order)

			if sellOrder.PendingShares > 0 {
				transaction := NewTransaction(sellOrder, order, order.Shares, sellOrder.Price)
				b.AddTransaction(transaction, b.Wg)
				sellOrder.Transactions = append(sellOrder.Transactions, transaction)
				order.Transactions = append(order.Transactions, transaction)
				b.OrdersChanOut <- sellOrder
				b.OrdersChanOut <- order

				if sellOrder.PendingShares > 0 {
					sellOrders.Push(sellOrder)
				}
			}
		} 
}
