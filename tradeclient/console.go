package main

import (
	"bufio"
	"fmt"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/fix"
	"github.com/quickfixgo/quickfix/fix/enum"
	"github.com/quickfixgo/quickfix/fix/field"
	"github.com/quickfixgo/quickfix/fix40"
	"github.com/quickfixgo/quickfix/fix41"
	"github.com/quickfixgo/quickfix/fix42"
	"github.com/quickfixgo/quickfix/fix43"
	"github.com/quickfixgo/quickfix/fix44"
	"github.com/quickfixgo/quickfix/fix50"
	"github.com/quickfixgo/quickfix/message"
	"os"
	"strconv"
	"strings"
)

func queryAction() (string, error) {
	fmt.Println()
	fmt.Println("1) Enter Order")
	fmt.Print("Action: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text(), scanner.Err()
}

func queryVersion() (field.BeginString, error) {
	fmt.Println()
	fmt.Println("1) FIX.4.0")
	fmt.Println("2) FIX.4.1")
	fmt.Println("3) FIX.4.2")
	fmt.Println("4) FIX.4.3")
	fmt.Println("5) FIX.4.4")
	fmt.Println("6) FIXT.1.1 (FIX.5.0)")
	fmt.Print("BeginString: ")

	var beginString field.BeginString
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return beginString, scanner.Err()
	}

	switch scanner.Text() {
	case "1":
		beginString.Value = fix.BeginString_FIX40
	case "2":
		beginString.Value = fix.BeginString_FIX41
	case "3":
		beginString.Value = fix.BeginString_FIX42
	case "4":
		beginString.Value = fix.BeginString_FIX43
	case "5":
		beginString.Value = fix.BeginString_FIX44
	case "6":
		beginString.Value = fix.BeginString_FIXT11
	case "7":
		beginString.Value = "A"
	default:
		return beginString, fmt.Errorf("unknown BeginString choice: %v", scanner.Text())
	}

	return beginString, nil
}

func queryClOrdID() (field.ClOrdID, error) {
	fmt.Print("ClOrdID: ")
	var clOrdID field.ClOrdID
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	clOrdID.Value = scanner.Text()

	return clOrdID, scanner.Err()
}

func querySymbol() (field.Symbol, error) {
	fmt.Println()
	fmt.Print("Symbol: ")
	scanner := bufio.NewScanner(os.Stdin)

	var symbol field.Symbol
	scanner.Scan()
	symbol.Value = scanner.Text()

	return symbol, scanner.Err()
}

func querySide() (field.Side, error) {

	fmt.Println()
	fmt.Println("1) Buy")
	fmt.Println("2) Sell")
	fmt.Println("3) Sell Short")
	fmt.Println("4) Sell Short Exempt")
	fmt.Println("5) Cross")
	fmt.Println("6) Cross Short")
	fmt.Println("7) Cross Short Exempt")
	fmt.Print("Side: ")

	var side field.Side
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return side, scanner.Err()
	}

	switch scanner.Text() {
	case "1":
		side.Value = enum.Side_BUY
	case "2":
		side.Value = enum.Side_SELL
	case "3":
		side.Value = enum.Side_SELL_SHORT
	case "4":
		side.Value = enum.Side_SELL_SHORT_EXEMPT
	case "5":
		side.Value = enum.Side_CROSS
	case "6":
		side.Value = enum.Side_CROSS_SHORT
	case "7":
		side.Value = "A"
	default:
		return side, fmt.Errorf("unknown side choice: %v", scanner.Text())
	}

	return side, nil
}

func queryOrdType() (field.OrdType, error) {
	fmt.Println()
	fmt.Println("1) Market")
	fmt.Println("2) Limit")
	fmt.Println("3) Stop")
	fmt.Println("4) Stop Limit")
	fmt.Print("OrdType: ")

	var ordType field.OrdType
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return ordType, scanner.Err()
	}

	switch scanner.Text() {
	case "1":
		ordType.Value = enum.OrdType_MARKET
	case "2":
		ordType.Value = enum.OrdType_LIMIT
	case "3":
		ordType.Value = enum.OrdType_STOP
	case "4":
		ordType.Value = enum.OrdType_STOP_LIMIT
	default:
		return ordType, fmt.Errorf("invalid ordtype choice: %v", scanner.Text())
	}

	return ordType, nil
}

func queryTimeInForce() (field.TimeInForce, error) {
	fmt.Println()
	fmt.Println("1) Day")
	fmt.Println("2) IOC")
	fmt.Println("3) OPG")
	fmt.Println("4) GTC")
	fmt.Println("5) GTX")
	fmt.Print("TimeInForce: ")

	var timeInForce field.TimeInForce
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return timeInForce, scanner.Err()
	}

	switch scanner.Text() {
	case "1":
		timeInForce.Value = enum.TimeInForce_DAY
	case "2":
		timeInForce.Value = enum.TimeInForce_IMMEDIATE_OR_CANCEL
	case "3":
		timeInForce.Value = enum.TimeInForce_AT_THE_OPENING
	case "4":
		timeInForce.Value = enum.TimeInForce_GOOD_TILL_CANCEL
	case "5":
		timeInForce.Value = enum.TimeInForce_GOOD_TILL_CROSSING

	default:
		return timeInForce, fmt.Errorf("invalid choice: %v", scanner.Text())
	}

	return timeInForce, nil
}

func queryOrderQty() (field.OrderQty, error) {
	fmt.Println()
	fmt.Print("OrderQty: ")

	var orderQty field.OrderQty
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return orderQty, scanner.Err()
	}

	var err error
	orderQty.Value, err = strconv.ParseFloat(scanner.Text(), 64)

	return orderQty, err
}

func queryPrice() (field.Price, error) {
	fmt.Println()
	fmt.Print("Price: ")

	var price field.Price
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return price, scanner.Err()
	}

	var err error
	price.Value, err = strconv.ParseFloat(scanner.Text(), 64)
	return price, err
}

func queryStopPx() (field.StopPx, error) {
	fmt.Println()
	fmt.Print("Stop Price: ")

	var price field.StopPx
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return price, scanner.Err()
	}

	var err error
	price.Value, err = strconv.ParseFloat(scanner.Text(), 64)
	return price, err
}

func querySenderCompID() (field.SenderCompID, error) {
	fmt.Println()
	fmt.Print("SenderCompID: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	var senderCompID field.SenderCompID
	senderCompID.Value = scanner.Text()
	return senderCompID, scanner.Err()
}

func queryTargetCompID() (field.TargetCompID, error) {
	fmt.Println()
	fmt.Print("TargetCompID: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	var targetCompID field.TargetCompID
	targetCompID.Value = scanner.Text()
	return targetCompID, scanner.Err()
}

func queryTargetSubID() (field.TargetSubID, error) {
	fmt.Println()
	fmt.Print("TargetSubID: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	var targetSubID field.TargetSubID
	targetSubID.Value = scanner.Text()
	return targetSubID, scanner.Err()
}

func queryConfirm(prompt string) bool {
	fmt.Println()
	fmt.Printf("%v?: ", prompt)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return strings.ToUpper(scanner.Text()) == "Y"
}

func queryHeader(header message.Header) error {
	senderCompID, err := querySenderCompID()
	if err != nil {
		return err
	}
	header.Set(senderCompID)

	targetCompID, err := queryTargetCompID()
	if err != nil {
		return err
	}
	header.Set(targetCompID)

	if ok := queryConfirm("Use a TargetSubID"); !ok {
		return nil
	}

	targetSubID, err := queryTargetSubID()
	if err != nil {
		return err
	}
	header.Set(targetSubID)

	return nil
}

func queryNewOrderSingle40() (fix40.NewOrderSingleBuilder, error) {
	var builder fix40.NewOrderSingleBuilder

	clOrdID, err := queryClOrdID()
	if err != nil {
		return builder, err
	}

	symbol, err := querySymbol()
	if err != nil {
		return builder, err
	}

	side, err := querySide()
	if err != nil {
		return builder, err
	}

	orderQty, err := queryOrderQty()
	if err != nil {
		return builder, err
	}

	ordType, err := queryOrdType()
	if err != nil {
		return builder, err
	}

	builder = fix40.CreateNewOrderSingleBuilder(clOrdID, field.BuildHandlInst("1"), symbol, side, orderQty, ordType)

	timeInForce, err := queryTimeInForce()
	if err != nil {
		return builder, err
	}

	builder.Body.Set(timeInForce)
	if ordType.Value == enum.OrdType_LIMIT || ordType.Value == enum.OrdType_STOP_LIMIT {
		price, err := queryPrice()
		if err != nil {
			return builder, err
		}
		builder.Body.Set(price)
	}

	if ordType.Value == enum.OrdType_STOP || ordType.Value == enum.OrdType_STOP_LIMIT {
		stopPx, err := queryStopPx()
		if err != nil {
			return builder, err
		}
		builder.Body.Set(stopPx)
	}

	queryHeader(builder.Header)

	return builder, nil
}

func queryNewOrderSingle41() (fix41.NewOrderSingleBuilder, error) {
	var builder fix41.NewOrderSingleBuilder

	clOrdID, err := queryClOrdID()
	if err != nil {
		return builder, err
	}

	symbol, err := querySymbol()
	if err != nil {
		return builder, err
	}

	side, err := querySide()
	if err != nil {
		return builder, err
	}

	ordType, err := queryOrdType()
	if err != nil {
		return builder, err
	}

	builder = fix41.CreateNewOrderSingleBuilder(clOrdID, field.BuildHandlInst("1"), symbol, side, ordType)
	orderQty, err := queryOrderQty()
	if err != nil {
		return builder, err
	}
	builder.Body.Set(orderQty)

	timeInForce, err := queryTimeInForce()
	if err != nil {
		return builder, err
	}

	builder.Body.Set(timeInForce)
	if ordType.Value == enum.OrdType_LIMIT || ordType.Value == enum.OrdType_STOP_LIMIT {
		price, err := queryPrice()
		if err != nil {
			return builder, err
		}
		builder.Body.Set(price)
	}

	if ordType.Value == enum.OrdType_STOP || ordType.Value == enum.OrdType_STOP_LIMIT {
		stopPx, err := queryStopPx()
		if err != nil {
			return builder, err
		}
		builder.Body.Set(stopPx)
	}

	queryHeader(builder.Header)

	return builder, nil
}

func queryNewOrderSingle42() (fix42.NewOrderSingleBuilder, error) {
	var builder fix42.NewOrderSingleBuilder

	clOrdID, err := queryClOrdID()
	if err != nil {
		return builder, err
	}

	symbol, err := querySymbol()
	if err != nil {
		return builder, err
	}

	side, err := querySide()
	if err != nil {
		return builder, err
	}

	ordType, err := queryOrdType()
	if err != nil {
		return builder, err
	}

	var transactTime field.TransactTime

	builder = fix42.CreateNewOrderSingleBuilder(clOrdID, field.BuildHandlInst("1"), symbol, side, transactTime, ordType)

	orderQty, err := queryOrderQty()
	if err != nil {
		return builder, err
	}
	builder.Body.Set(orderQty)

	timeInForce, err := queryTimeInForce()
	if err != nil {
		return builder, err
	}

	builder.Body.Set(timeInForce)
	if ordType.Value == enum.OrdType_LIMIT || ordType.Value == enum.OrdType_STOP_LIMIT {
		price, err := queryPrice()
		if err != nil {
			return builder, err
		}
		builder.Body.Set(price)
	}

	if ordType.Value == enum.OrdType_STOP || ordType.Value == enum.OrdType_STOP_LIMIT {
		stopPx, err := queryStopPx()
		if err != nil {
			return builder, err
		}
		builder.Body.Set(stopPx)
	}

	queryHeader(builder.Header)

	return builder, nil
}

func queryNewOrderSingle43() (fix43.NewOrderSingleBuilder, error) {
	var builder fix43.NewOrderSingleBuilder

	clOrdID, err := queryClOrdID()
	if err != nil {
		return builder, err
	}

	side, err := querySide()
	if err != nil {
		return builder, err
	}

	ordType, err := queryOrdType()
	if err != nil {
		return builder, err
	}

	var transactTime field.TransactTime

	builder = fix43.CreateNewOrderSingleBuilder(clOrdID, field.BuildHandlInst("1"), side, transactTime, ordType)

	symbol, err := querySymbol()
	if err != nil {
		return builder, err
	}
	builder.Body.Set(symbol)

	orderQty, err := queryOrderQty()
	if err != nil {
		return builder, err
	}
	builder.Body.Set(orderQty)

	timeInForce, err := queryTimeInForce()
	if err != nil {
		return builder, err
	}

	builder.Body.Set(timeInForce)
	if ordType.Value == enum.OrdType_LIMIT || ordType.Value == enum.OrdType_STOP_LIMIT {
		price, err := queryPrice()
		if err != nil {
			return builder, err
		}
		builder.Body.Set(price)
	}

	if ordType.Value == enum.OrdType_STOP || ordType.Value == enum.OrdType_STOP_LIMIT {
		stopPx, err := queryStopPx()
		if err != nil {
			return builder, err
		}
		builder.Body.Set(stopPx)
	}

	queryHeader(builder.Header)

	return builder, nil
}

func queryNewOrderSingle44() (fix44.NewOrderSingleBuilder, error) {
	var builder fix44.NewOrderSingleBuilder

	clOrdID, err := queryClOrdID()
	if err != nil {
		return builder, err
	}

	side, err := querySide()
	if err != nil {
		return builder, err
	}

	ordType, err := queryOrdType()
	if err != nil {
		return builder, err
	}

	var transactTime field.TransactTime

	builder = fix44.CreateNewOrderSingleBuilder(clOrdID, side, transactTime, ordType)

	builder.Body.Set(field.BuildHandlInst("1"))
	symbol, err := querySymbol()
	if err != nil {
		return builder, err
	}
	builder.Body.Set(symbol)

	orderQty, err := queryOrderQty()
	if err != nil {
		return builder, err
	}
	builder.Body.Set(orderQty)

	timeInForce, err := queryTimeInForce()
	if err != nil {
		return builder, err
	}

	builder.Body.Set(timeInForce)
	if ordType.Value == enum.OrdType_LIMIT || ordType.Value == enum.OrdType_STOP_LIMIT {
		price, err := queryPrice()
		if err != nil {
			return builder, err
		}
		builder.Body.Set(price)
	}

	if ordType.Value == enum.OrdType_STOP || ordType.Value == enum.OrdType_STOP_LIMIT {
		stopPx, err := queryStopPx()
		if err != nil {
			return builder, err
		}
		builder.Body.Set(stopPx)
	}

	queryHeader(builder.Header)

	return builder, nil
}

func queryNewOrderSingle50() (fix50.NewOrderSingleBuilder, error) {
	var builder fix50.NewOrderSingleBuilder

	clOrdID, err := queryClOrdID()
	if err != nil {
		return builder, err
	}

	side, err := querySide()
	if err != nil {
		return builder, err
	}

	ordType, err := queryOrdType()
	if err != nil {
		return builder, err
	}

	var transactTime field.TransactTime

	builder = fix50.CreateNewOrderSingleBuilder(clOrdID, side, transactTime, ordType)

	builder.Body.Set(field.BuildHandlInst("1"))
	symbol, err := querySymbol()
	if err != nil {
		return builder, err
	}
	builder.Body.Set(symbol)

	orderQty, err := queryOrderQty()
	if err != nil {
		return builder, err
	}
	builder.Body.Set(orderQty)

	timeInForce, err := queryTimeInForce()
	if err != nil {
		return builder, err
	}

	builder.Body.Set(timeInForce)
	if ordType.Value == enum.OrdType_LIMIT || ordType.Value == enum.OrdType_STOP_LIMIT {
		price, err := queryPrice()
		if err != nil {
			return builder, err
		}
		builder.Body.Set(price)
	}

	if ordType.Value == enum.OrdType_STOP || ordType.Value == enum.OrdType_STOP_LIMIT {
		stopPx, err := queryStopPx()
		if err != nil {
			return builder, err
		}
		builder.Body.Set(stopPx)
	}

	queryHeader(builder.Header)

	return builder, nil
}

func queryEnterOrder() error {
	beginString, err := queryVersion()
	if err != nil {
		return err
	}

	var order *message.MessageBuilder
	switch beginString.Value {
	case fix.BeginString_FIX40:
		fix40Order, err := queryNewOrderSingle40()
		if err != nil {
			return err
		}
		order = &(fix40Order.MessageBuilder)

	case fix.BeginString_FIX41:
		fix41Order, err := queryNewOrderSingle41()
		if err != nil {
			return err
		}
		order = &(fix41Order.MessageBuilder)

	case fix.BeginString_FIX42:
		fix42Order, err := queryNewOrderSingle42()
		if err != nil {
			return err
		}
		order = &(fix42Order.MessageBuilder)

	case fix.BeginString_FIX43:
		fix43Order, err := queryNewOrderSingle43()
		if err != nil {
			return err
		}
		order = &(fix43Order.MessageBuilder)

	case fix.BeginString_FIX44:
		fix44Order, err := queryNewOrderSingle44()
		if err != nil {
			return err
		}
		order = &(fix44Order.MessageBuilder)

	case fix.BeginString_FIXT11:
		fix50Order, err := queryNewOrderSingle50()
		if err != nil {
			return err
		}
		order = &(fix50Order.MessageBuilder)
	}

	return quickfix.Send(*order)
}
