// Copyright (c) quickfixengine.org  All rights reserved.
//
// This file may be distributed under the terms of the quickfixengine.org
// license as defined by quickfixengine.org and appearing in the file
// LICENSE included in the packaging of this file.
//
// This file is provided AS IS with NO WARRANTY OF ANY KIND, INCLUDING
// THE WARRANTY OF DESIGN, MERCHANTABILITY AND FITNESS FOR A
// PARTICULAR PURPOSE.
//
// See http://www.quickfixengine.org/LICENSE for licensing information.
//
// Contact ask@quickfixengine.org if any conditions of this licensing
// are not clear to you.

package internal

import (
	"fmt"
	"sort"
	"time"

	"github.com/quickfixgo/enum"
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
		switch i.Price.Cmp(j.Price) {
		case 1:
			return true
		case -1:
			return false
		}

		return i.insertTime.Before(j.insertTime)
	}

	return
}

func offers() (o orderList) {
	o.sortBy = func(i, j *Order) bool {
		switch i.Price.Cmp(j.Price) {
		case 1:
			return false
		case -1:
			return true
		}

		return i.insertTime.Before(j.insertTime)
	}

	return
}

// Market is a simple CLOB
type Market struct {
	Bids   orderList
	Offers orderList
}

// NewMarket returns an initialized Market instance
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

	fmt.Println("OFFERS:")
	fmt.Println("-----")
	fmt.Println()

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
func (m *Market) Cancel(clordID string, side enum.Side) (order *Order) {
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
		if offerQuant := bestOffer.OpenQuantity(); offerQuant.Cmp(quantity) == -1 {
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
