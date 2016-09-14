package internal

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ClOrdID              string
	Symbol               string
	SenderCompID         string
	TargetCompID         string
	Side                 string
	OrdType              string
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
