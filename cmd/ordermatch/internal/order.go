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
	"time"

	"github.com/quickfixgo/enum"
	"github.com/shopspring/decimal"
)

type Order struct {
	ClOrdID              string
	Symbol               string
	SenderCompID         string
	TargetCompID         string
	Side                 enum.Side
	OrdType              enum.OrdType
	Price                decimal.Decimal
	Quantity             decimal.Decimal
	ExecutedQuantity     decimal.Decimal
	openQuantity         *decimal.Decimal
	AvgPx                decimal.Decimal
	insertTime           time.Time
	LastExecutedQuantity decimal.Decimal
	LastExecutedPrice    decimal.Decimal
}

func (o Order) IsClosed() bool {
	return o.OpenQuantity().Equals(decimal.Zero)
}

func (o Order) OpenQuantity() decimal.Decimal {
	if o.openQuantity == nil {
		return o.Quantity.Sub(o.ExecutedQuantity)
	}

	return *o.openQuantity
}

func (o *Order) Execute(price, quantity decimal.Decimal) {
	o.ExecutedQuantity = o.ExecutedQuantity.Add(quantity)
	o.LastExecutedPrice = price
	o.LastExecutedQuantity = quantity
}

func (o *Order) Cancel() {
	openQuantity := decimal.Zero
	o.openQuantity = &openQuantity
}
