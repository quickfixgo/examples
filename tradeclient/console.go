package main

import (
	"bufio"
	"fmt"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/fix"
	"github.com/quickfixgo/quickfix/fix/enum"
	"github.com/quickfixgo/quickfix/fix/field"

	fix40nos "github.com/quickfixgo/quickfix/fix40/newordersingle"
	fix41nos "github.com/quickfixgo/quickfix/fix41/newordersingle"
	fix42nos "github.com/quickfixgo/quickfix/fix42/newordersingle"
	fix43nos "github.com/quickfixgo/quickfix/fix43/newordersingle"
	fix44nos "github.com/quickfixgo/quickfix/fix44/newordersingle"
	fix50nos "github.com/quickfixgo/quickfix/fix50/newordersingle"

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

func queryVersion() (*field.BeginStringField, error) {
	fmt.Println()
	fmt.Println("1) FIX.4.0")
	fmt.Println("2) FIX.4.1")
	fmt.Println("3) FIX.4.2")
	fmt.Println("4) FIX.4.3")
	fmt.Println("5) FIX.4.4")
	fmt.Println("6) FIXT.1.1 (FIX.5.0)")
	fmt.Print("BeginString: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return nil, scanner.Err()
	}

	switch scanner.Text() {
	case "1":
		return field.NewBeginString(fix.BeginString_FIX40), nil
	case "2":
		return field.NewBeginString(fix.BeginString_FIX41), nil
	case "3":
		return field.NewBeginString(fix.BeginString_FIX42), nil
	case "4":
		return field.NewBeginString(fix.BeginString_FIX43), nil
	case "5":
		return field.NewBeginString(fix.BeginString_FIX44), nil
	case "6":
		return field.NewBeginString(fix.BeginString_FIXT11), nil
	case "7":
		return field.NewBeginString("A"), nil
	}

	return nil, fmt.Errorf("unknown BeginString choice: %v", scanner.Text())
}

func queryClOrdID() (*field.ClOrdIDField, error) {
	fmt.Print("ClOrdID: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return field.NewClOrdID(scanner.Text()), scanner.Err()
}

func querySymbol() (*field.SymbolField, error) {
	fmt.Println()
	fmt.Print("Symbol: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return field.NewSymbol(scanner.Text()), scanner.Err()
}

func querySide() (*field.SideField, error) {

	fmt.Println()
	fmt.Println("1) Buy")
	fmt.Println("2) Sell")
	fmt.Println("3) Sell Short")
	fmt.Println("4) Sell Short Exempt")
	fmt.Println("5) Cross")
	fmt.Println("6) Cross Short")
	fmt.Println("7) Cross Short Exempt")
	fmt.Print("Side: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return nil, scanner.Err()
	}

	switch scanner.Text() {
	case "1":
		return field.NewSide(enum.Side_BUY), nil
	case "2":
		return field.NewSide(enum.Side_SELL), nil
	case "3":
		return field.NewSide(enum.Side_SELL_SHORT), nil
	case "4":
		return field.NewSide(enum.Side_SELL_SHORT_EXEMPT), nil
	case "5":
		return field.NewSide(enum.Side_CROSS), nil
	case "6":
		return field.NewSide(enum.Side_CROSS_SHORT), nil
	case "7":
		return field.NewSide("A"), nil
	}

	return nil, fmt.Errorf("unknown side choice: %v", scanner.Text())
}

func queryOrdType() (*field.OrdTypeField, error) {
	fmt.Println()
	fmt.Println("1) Market")
	fmt.Println("2) Limit")
	fmt.Println("3) Stop")
	fmt.Println("4) Stop Limit")
	fmt.Print("OrdType: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return nil, scanner.Err()
	}

	switch scanner.Text() {
	case "1":
		return field.NewOrdType(enum.OrdType_MARKET), nil
	case "2":
		return field.NewOrdType(enum.OrdType_LIMIT), nil
	case "3":
		return field.NewOrdType(enum.OrdType_STOP), nil
	case "4":
		return field.NewOrdType(enum.OrdType_STOP_LIMIT), nil
	}

	return nil, fmt.Errorf("invalid ordtype choice: %v", scanner.Text())
}

func queryTimeInForce() (*field.TimeInForceField, error) {
	fmt.Println()
	fmt.Println("1) Day")
	fmt.Println("2) IOC")
	fmt.Println("3) OPG")
	fmt.Println("4) GTC")
	fmt.Println("5) GTX")
	fmt.Print("TimeInForce: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return nil, scanner.Err()
	}

	switch scanner.Text() {
	case "1":
		return field.NewTimeInForce(enum.TimeInForce_DAY), nil
	case "2":
		return field.NewTimeInForce(enum.TimeInForce_IMMEDIATE_OR_CANCEL), nil
	case "3":
		return field.NewTimeInForce(enum.TimeInForce_AT_THE_OPENING), nil
	case "4":
		return field.NewTimeInForce(enum.TimeInForce_GOOD_TILL_CANCEL), nil
	case "5":
		return field.NewTimeInForce(enum.TimeInForce_GOOD_TILL_CROSSING), nil
	}

	return nil, fmt.Errorf("invalid choice: %v", scanner.Text())
}

func queryOrderQty() (*field.OrderQtyField, error) {
	fmt.Println()
	fmt.Print("OrderQty: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return nil, scanner.Err()
	}

	val, err := strconv.ParseFloat(scanner.Text(), 64)
	if err != nil {
		return nil, err
	}

	return field.NewOrderQty(val), err
}

func queryPrice() (*field.PriceField, error) {
	fmt.Println()
	fmt.Print("Price: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return nil, scanner.Err()
	}

	val, err := strconv.ParseFloat(scanner.Text(), 64)
	if err != nil {
		return nil, err
	}
	return field.NewPrice(val), nil
}

func queryStopPx() (*field.StopPxField, error) {
	fmt.Println()
	fmt.Print("Stop Price: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return nil, scanner.Err()
	}

	val, err := strconv.ParseFloat(scanner.Text(), 64)
	if err != nil {
		return nil, err
	}
	return field.NewStopPx(val), nil
}

func querySenderCompID() (*field.SenderCompIDField, error) {
	fmt.Println()
	fmt.Print("SenderCompID: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return field.NewSenderCompID(scanner.Text()), scanner.Err()
}

func queryTargetCompID() (*field.TargetCompIDField, error) {
	fmt.Println()
	fmt.Print("TargetCompID: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return field.NewTargetCompID(scanner.Text()), scanner.Err()
}

func queryTargetSubID() (*field.TargetSubIDField, error) {
	fmt.Println()
	fmt.Print("TargetSubID: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return field.NewTargetSubID(scanner.Text()), scanner.Err()
}

func queryConfirm(prompt string) bool {
	fmt.Println()
	fmt.Printf("%v?: ", prompt)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return strings.ToUpper(scanner.Text()) == "Y"
}

func queryHeader(header quickfix.FieldMap) error {
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

func queryNewOrderSingle40() (fix40nos.MessageBuilder, error) {
	var builder fix40nos.MessageBuilder

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

	builder = fix40nos.Builder(clOrdID, field.NewHandlInst("1"), symbol, side, orderQty, ordType)

	timeInForce, err := queryTimeInForce()
	if err != nil {
		return builder, err
	}

	builder.Body().Set(timeInForce)
	if ordType.Value == enum.OrdType_LIMIT || ordType.Value == enum.OrdType_STOP_LIMIT {
		price, err := queryPrice()
		if err != nil {
			return builder, err
		}
		builder.Body().Set(price)
	}

	if ordType.Value == enum.OrdType_STOP || ordType.Value == enum.OrdType_STOP_LIMIT {
		stopPx, err := queryStopPx()
		if err != nil {
			return builder, err
		}
		builder.Body().Set(stopPx)
	}

	queryHeader(builder.Header())

	return builder, nil
}

func queryNewOrderSingle41() (fix41nos.MessageBuilder, error) {
	var builder fix41nos.MessageBuilder

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

	builder = fix41nos.Builder(clOrdID, field.NewHandlInst("1"), symbol, side, ordType)
	orderQty, err := queryOrderQty()
	if err != nil {
		return builder, err
	}
	builder.Body().Set(orderQty)

	timeInForce, err := queryTimeInForce()
	if err != nil {
		return builder, err
	}

	builder.Body().Set(timeInForce)
	if ordType.Value == enum.OrdType_LIMIT || ordType.Value == enum.OrdType_STOP_LIMIT {
		price, err := queryPrice()
		if err != nil {
			return builder, err
		}
		builder.Body().Set(price)
	}

	if ordType.Value == enum.OrdType_STOP || ordType.Value == enum.OrdType_STOP_LIMIT {
		stopPx, err := queryStopPx()
		if err != nil {
			return builder, err
		}
		builder.Body().Set(stopPx)
	}

	queryHeader(builder.Header())

	return builder, nil
}

func queryNewOrderSingle42() (fix42nos.MessageBuilder, error) {
	var builder fix42nos.MessageBuilder

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

	transactTime := &field.TransactTimeField{}

	builder = fix42nos.Builder(clOrdID, field.NewHandlInst("1"), symbol, side, transactTime, ordType)

	orderQty, err := queryOrderQty()
	if err != nil {
		return builder, err
	}
	builder.Body().Set(orderQty)

	timeInForce, err := queryTimeInForce()
	if err != nil {
		return builder, err
	}

	builder.Body().Set(timeInForce)
	if ordType.Value == enum.OrdType_LIMIT || ordType.Value == enum.OrdType_STOP_LIMIT {
		price, err := queryPrice()
		if err != nil {
			return builder, err
		}
		builder.Body().Set(price)
	}

	if ordType.Value == enum.OrdType_STOP || ordType.Value == enum.OrdType_STOP_LIMIT {
		stopPx, err := queryStopPx()
		if err != nil {
			return builder, err
		}
		builder.Body().Set(stopPx)
	}

	queryHeader(builder.Header())

	return builder, nil
}

func queryNewOrderSingle43() (fix43nos.MessageBuilder, error) {
	var builder fix43nos.MessageBuilder

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

	transactTime := &field.TransactTimeField{}

	builder = fix43nos.Builder(clOrdID, field.NewHandlInst("1"), side, transactTime, ordType)

	symbol, err := querySymbol()
	if err != nil {
		return builder, err
	}
	builder.Body().Set(symbol)

	orderQty, err := queryOrderQty()
	if err != nil {
		return builder, err
	}
	builder.Body().Set(orderQty)

	timeInForce, err := queryTimeInForce()
	if err != nil {
		return builder, err
	}

	builder.Body().Set(timeInForce)
	if ordType.Value == enum.OrdType_LIMIT || ordType.Value == enum.OrdType_STOP_LIMIT {
		price, err := queryPrice()
		if err != nil {
			return builder, err
		}
		builder.Body().Set(price)
	}

	if ordType.Value == enum.OrdType_STOP || ordType.Value == enum.OrdType_STOP_LIMIT {
		stopPx, err := queryStopPx()
		if err != nil {
			return builder, err
		}
		builder.Body().Set(stopPx)
	}

	queryHeader(builder.Header())

	return builder, nil
}

func queryNewOrderSingle44() (fix44nos.MessageBuilder, error) {
	var builder fix44nos.MessageBuilder

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

	transactTime := &field.TransactTimeField{}

	builder = fix44nos.Builder(clOrdID, side, transactTime, ordType)

	builder.Body().Set(field.NewHandlInst("1"))
	symbol, err := querySymbol()
	if err != nil {
		return builder, err
	}
	builder.Body().Set(symbol)

	orderQty, err := queryOrderQty()
	if err != nil {
		return builder, err
	}
	builder.Body().Set(orderQty)

	timeInForce, err := queryTimeInForce()
	if err != nil {
		return builder, err
	}

	builder.Body().Set(timeInForce)
	if ordType.Value == enum.OrdType_LIMIT || ordType.Value == enum.OrdType_STOP_LIMIT {
		price, err := queryPrice()
		if err != nil {
			return builder, err
		}
		builder.Body().Set(price)
	}

	if ordType.Value == enum.OrdType_STOP || ordType.Value == enum.OrdType_STOP_LIMIT {
		stopPx, err := queryStopPx()
		if err != nil {
			return builder, err
		}
		builder.Body().Set(stopPx)
	}

	queryHeader(builder.Header())

	return builder, nil
}

func queryNewOrderSingle50() (fix50nos.MessageBuilder, error) {
	var builder fix50nos.MessageBuilder

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

	transactTime := &field.TransactTimeField{}

	builder = fix50nos.Builder(clOrdID, side, transactTime, ordType)

	builder.Body().Set(field.NewHandlInst("1"))
	symbol, err := querySymbol()
	if err != nil {
		return builder, err
	}
	builder.Body().Set(symbol)

	orderQty, err := queryOrderQty()
	if err != nil {
		return builder, err
	}
	builder.Body().Set(orderQty)

	timeInForce, err := queryTimeInForce()
	if err != nil {
		return builder, err
	}

	builder.Body().Set(timeInForce)
	if ordType.Value == enum.OrdType_LIMIT || ordType.Value == enum.OrdType_STOP_LIMIT {
		price, err := queryPrice()
		if err != nil {
			return builder, err
		}
		builder.Body().Set(price)
	}

	if ordType.Value == enum.OrdType_STOP || ordType.Value == enum.OrdType_STOP_LIMIT {
		stopPx, err := queryStopPx()
		if err != nil {
			return builder, err
		}
		builder.Body().Set(stopPx)
	}

	queryHeader(builder.Header())

	return builder, nil
}

func queryEnterOrder() error {
	beginString, err := queryVersion()
	if err != nil {
		return err
	}

	var order *quickfix.MessageBuilder
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
