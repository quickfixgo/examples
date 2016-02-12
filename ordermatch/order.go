package main

import (
	"time"
)

type Order struct {
	ClOrdID          string
	Symbol           string
	SenderCompID     string
	TargetCompID     string
	Side             string
	OrdType          string
	Price            float64
	Quantity         float64
	ExecutedQuantity float64
	openQuantity     *float64
	AvgPx            float64
	insertTime       time.Time
}

func (o Order) IsClosed() bool {
	return o.OpenQuantity() == 0
}

func (o Order) OpenQuantity() float64 {
	if o.openQuantity == nil {
		return o.Quantity - o.ExecutedQuantity
	}

	return *o.openQuantity
}

func (o *Order) Execute(price, quantity float64) {
	o.ExecutedQuantity += quantity
}

func (o *Order) Cancel() {
	openQuantity := float64(0)
	o.openQuantity = &openQuantity
}
