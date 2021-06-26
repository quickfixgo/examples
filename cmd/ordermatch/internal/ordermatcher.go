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

	"github.com/quickfixgo/enum"
)

type OrderMatcher struct {
	markets map[string]*Market
}

func NewOrderMatcher() *OrderMatcher {
	return &OrderMatcher{markets: make(map[string]*Market)}
}

func (m OrderMatcher) DisplayMarket(symbol string) {
	if market, ok := m.markets[symbol]; ok {
		market.Display()
		return
	}
	fmt.Println("================")
	fmt.Println("SYMBOL NOT FOUND")
	fmt.Println("================")
}

func (m OrderMatcher) Display() {
	hasMarkets := len(m.markets) > 0
	if hasMarkets {
		fmt.Println("===============")
		fmt.Println("ACTIVE SYMBOLS:")
		fmt.Println("===============")
		for symbol := range m.markets {
			fmt.Println(symbol)
		}
		return
	}
	fmt.Println("===========================")
	fmt.Println("THERE ARE NO ACTIVE SYMBOLS")
	fmt.Println("===========================")
}

func (m *OrderMatcher) Insert(order Order) {
	market, ok := m.markets[order.Symbol]
	if !ok {
		market = NewMarket()
		m.markets[order.Symbol] = market
	}

	market.Insert(order)
}

func (m *OrderMatcher) Cancel(clordID, symbol string, side enum.Side) *Order {
	market, ok := m.markets[symbol]
	if !ok {
		return nil
	}

	return market.Cancel(clordID, side)
}

func (m *OrderMatcher) Match(symbol string) []Order {
	market, ok := m.markets[symbol]
	if !ok {
		return []Order{}
	}

	return market.Match()
}
