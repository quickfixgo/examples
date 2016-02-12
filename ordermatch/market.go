package main

import (
	"fmt"
	"github.com/quickfixgo/quickfix/enum"
	"sort"
	"time"
)

type orderList struct {
	orders []*Order
	sortBy func(o1, o2 *Order) bool
}

func (l orderList) Len() int { return len(l.orders) }
func (l orderList) Swap(i, j int) {
	l.orders[i], l.orders[j] = l.orders[j], l.orders[i]
}
func (l orderList) Less(i, j int) bool { return l.sortBy(l.orders[i], l.orders[j]) }

func (l *orderList) Insert(order *Order) {
	l.orders = append(l.orders, order)
	sort.Sort(l)
}

func (l *orderList) Remove(clordID string) (order *Order) {
	for i := 0; i < len(l.orders); i++ {
		if l.orders[i].ClOrdID == clordID {
			order = l.orders[i]
			l.orders = append(l.orders[:i], l.orders[i+1:]...)
			return
		}
	}

	return
}

func bids() (b orderList) {
	b.sortBy = func(i, j *Order) bool {
		if i.Price > j.Price {
			return true
		} else if i.Price < j.Price {
			return false
		}

		return i.insertTime.Before(j.insertTime)
	}

	return
}

func offers() (o orderList) {
	o.sortBy = func(i, j *Order) bool {
		if i.Price < j.Price {
			return true
		} else if i.Price < j.Price {
			return false
		}

		return i.insertTime.Before(j.insertTime)
	}

	return
}

//Market is a simple CLOB
type Market struct {
	Bids   orderList
	Offers orderList
}

//NewMarket returns an initialized Market instance
func NewMarket() *Market {
	return &Market{bids(), offers()}
}

func (m Market) Display() {
	fmt.Println("BIDS:")
	fmt.Println("-----")
	fmt.Println()

	for _, bid := range m.Bids.orders {
		fmt.Printf("%+v\n", bid)
	}

	fmt.Println()
	fmt.Println("OFFERS:")

	for _, offer := range m.Offers.orders {
		fmt.Printf("%+v\n", offer)
	}
}

func (m *Market) Insert(order Order) {
	order.insertTime = time.Now()
	if order.Side == enum.Side_BUY {
		m.Bids.Insert(&order)
	} else {
		m.Offers.Insert(&order)
	}
}
func (m *Market) Cancel(clordID, side string) (order *Order) {
	if side == enum.Side_BUY {
		order = m.Bids.Remove(clordID)
	} else {
		order = m.Offers.Remove(clordID)
	}

	if order != nil {
		order.Cancel()
	}

	return
}

func (m *Market) Match() (matched []Order) {
	for m.Bids.Len() > 0 && m.Offers.Len() > 0 {
		bestBid := m.Bids.orders[0]
		bestOffer := m.Offers.orders[0]

		price := bestOffer.Price
		quantity := bestBid.OpenQuantity()
		if offerQuant := bestOffer.OpenQuantity(); offerQuant < quantity {
			quantity = offerQuant
		}

		bestBid.Execute(price, quantity)
		bestOffer.Execute(price, quantity)

		matched = append(matched, *bestBid, *bestOffer)

		if bestBid.IsClosed() {
			m.Bids.orders = m.Bids.orders[1:]
		}

		if bestOffer.IsClosed() {
			m.Offers.orders = m.Offers.orders[1:]
		}
	}

	return
}
