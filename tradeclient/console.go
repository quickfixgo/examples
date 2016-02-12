package main

import (
	"bufio"
	"fmt"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/enum"
	"github.com/quickfixgo/quickfix/field"

	fix40cxl "github.com/quickfixgo/quickfix/fix40/ordercancelrequest"
	fix41cxl "github.com/quickfixgo/quickfix/fix41/ordercancelrequest"
	fix42cxl "github.com/quickfixgo/quickfix/fix42/ordercancelrequest"
	fix43cxl "github.com/quickfixgo/quickfix/fix43/ordercancelrequest"
	fix44cxl "github.com/quickfixgo/quickfix/fix44/ordercancelrequest"
	fix50cxl "github.com/quickfixgo/quickfix/fix50/ordercancelrequest"

	fix40nos "github.com/quickfixgo/quickfix/fix40/newordersingle"
	fix41nos "github.com/quickfixgo/quickfix/fix41/newordersingle"
	fix42nos "github.com/quickfixgo/quickfix/fix42/newordersingle"
	fix43nos "github.com/quickfixgo/quickfix/fix43/newordersingle"
	fix44nos "github.com/quickfixgo/quickfix/fix44/newordersingle"
	fix50nos "github.com/quickfixgo/quickfix/fix50/newordersingle"

	"os"
	"strconv"
	"strings"
	"time"
)

func queryAction() (string, error) {
	fmt.Println()
	fmt.Println("1) Enter Order")
	fmt.Println("2) Cancel Order")
	fmt.Print("Action: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text(), scanner.Err()
}

func queryVersion() (string, error) {
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
		return "", scanner.Err()
	}

	switch scanner.Text() {
	case "1":
		return enum.BeginStringFIX40, nil
	case "2":
		return enum.BeginStringFIX41, nil
	case "3":
		return enum.BeginStringFIX42, nil
	case "4":
		return enum.BeginStringFIX43, nil
	case "5":
		return enum.BeginStringFIX44, nil
	case "6":
		return enum.BeginStringFIXT11, nil
	case "7":
		return "A", nil
	}

	return "", fmt.Errorf("unknown BeginString choice: %v", scanner.Text())
}

func queryClOrdID() (string, error) {
	fmt.Print("ClOrdID: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text(), scanner.Err()
}

func queryOrigClOrdID() (string, error) {
	fmt.Print("OrigClOrdID: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text(), scanner.Err()
}

func querySymbol() (string, error) {
	fmt.Println()
	fmt.Print("Symbol: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return scanner.Text(), scanner.Err()
}

func querySide() (string, error) {

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
		return "", scanner.Err()
	}

	switch scanner.Text() {
	case "1":
		return enum.Side_BUY, nil
	case "2":
		return enum.Side_SELL, nil
	case "3":
		return enum.Side_SELL_SHORT, nil
	case "4":
		return enum.Side_SELL_SHORT_EXEMPT, nil
	case "5":
		return enum.Side_CROSS, nil
	case "6":
		return enum.Side_CROSS_SHORT, nil
	case "7":
		return "A", nil
	}

	return "", fmt.Errorf("unknown side choice: %v", scanner.Text())
}

func queryOrdType() (string, error) {
	fmt.Println()
	fmt.Println("1) Market")
	fmt.Println("2) Limit")
	fmt.Println("3) Stop")
	fmt.Println("4) Stop Limit")
	fmt.Print("OrdType: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return "", scanner.Err()
	}

	switch scanner.Text() {
	case "1":
		return enum.OrdType_MARKET, nil
	case "2":
		return enum.OrdType_LIMIT, nil
	case "3":
		return enum.OrdType_STOP, nil
	case "4":
		return enum.OrdType_STOP_LIMIT, nil
	}

	return "", fmt.Errorf("invalid ordtype choice: %v", scanner.Text())
}

func queryTimeInForce() (string, error) {
	fmt.Println()
	fmt.Println("1) Day")
	fmt.Println("2) IOC")
	fmt.Println("3) OPG")
	fmt.Println("4) GTC")
	fmt.Println("5) GTX")
	fmt.Print("TimeInForce: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return "", scanner.Err()
	}

	switch scanner.Text() {
	case "1":
		return enum.TimeInForce_DAY, nil
	case "2":
		return enum.TimeInForce_IMMEDIATE_OR_CANCEL, nil
	case "3":
		return enum.TimeInForce_AT_THE_OPENING, nil
	case "4":
		return enum.TimeInForce_GOOD_TILL_CANCEL, nil
	case "5":
		return enum.TimeInForce_GOOD_TILL_CROSSING, nil
	}

	return "", fmt.Errorf("invalid choice: %v", scanner.Text())
}

func queryOrderQty() (float64, error) {
	fmt.Println()
	fmt.Print("OrderQty: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return 0, scanner.Err()
	}

	val, err := strconv.ParseFloat(scanner.Text(), 64)
	if err != nil {
		return 0, err
	}

	return val, err
}

func queryPrice() (float64, error) {
	fmt.Println()
	fmt.Print("Price: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return 0, scanner.Err()
	}

	val, err := strconv.ParseFloat(scanner.Text(), 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func queryStopPx() (float64, error) {
	fmt.Println()
	fmt.Print("Stop Price: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return 0, scanner.Err()
	}

	val, err := strconv.ParseFloat(scanner.Text(), 64)
	if err != nil {
		return 0, err
	}
	return val, nil
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

func queryNewOrderSingle40() (msg quickfix.Message, err error) {
	order := fix40nos.Message{
		HandlInst: "1",
	}

	if order.ClOrdID, err = queryClOrdID(); err != nil {
		return
	}

	if order.Symbol, err = querySymbol(); err != nil {
		return
	}

	if order.Side, err = querySide(); err != nil {
		return
	}

	if order.OrdType, err = queryOrdType(); err != nil {
		return
	}

	if fVal, err := queryOrderQty(); err != nil {
		return msg, err
	} else {
		order.OrderQty = int(fVal)
	}

	switch order.OrdType {
	case enum.OrdType_LIMIT, enum.OrdType_STOP_LIMIT:
		if fVal, err := queryPrice(); err != nil {
			return msg, err
		} else {
			order.Price = &fVal
		}
	}

	switch order.OrdType {
	case enum.OrdType_STOP, enum.OrdType_STOP_LIMIT:
		if fVal, err := queryStopPx(); err != nil {
			return msg, err
		} else {
			order.StopPx = &fVal
		}
	}

	if tif, err := queryTimeInForce(); err != nil {
		return msg, err
	} else {
		order.TimeInForce = &tif
	}

	msg = quickfix.Marshal(order)
	queryHeader(msg.Header)

	return
}

func queryNewOrderSingle41() (msg quickfix.Message, err error) {
	order := fix41nos.Message{
		HandlInst: "1",
	}

	if order.ClOrdID, err = queryClOrdID(); err != nil {
		return
	}

	if order.Symbol, err = querySymbol(); err != nil {
		return
	}

	if order.Side, err = querySide(); err != nil {
		return
	}

	if order.OrdType, err = queryOrdType(); err != nil {
		return
	}

	if fVal, err := queryOrderQty(); err != nil {
		return msg, err
	} else {
		iVal := int(fVal)
		order.OrderQty = &iVal
	}

	switch order.OrdType {
	case enum.OrdType_LIMIT, enum.OrdType_STOP_LIMIT:
		if fVal, err := queryPrice(); err != nil {
			return msg, err
		} else {
			order.Price = &fVal
		}
	}

	switch order.OrdType {
	case enum.OrdType_STOP, enum.OrdType_STOP_LIMIT:
		if fVal, err := queryStopPx(); err != nil {
			return msg, err
		} else {
			order.StopPx = &fVal
		}
	}

	if tif, err := queryTimeInForce(); err != nil {
		return msg, err
	} else {
		order.TimeInForce = &tif
	}

	msg = quickfix.Marshal(order)
	queryHeader(msg.Header)

	return
}

func queryNewOrderSingle42() (msg quickfix.Message, err error) {
	order := fix42nos.Message{
		HandlInst:    "1",
		TransactTime: time.Now(),
	}

	if order.ClOrdID, err = queryClOrdID(); err != nil {
		return
	}

	if order.Symbol, err = querySymbol(); err != nil {
		return
	}

	if order.Side, err = querySide(); err != nil {
		return
	}

	if order.OrdType, err = queryOrdType(); err != nil {
		return
	}

	if fVal, err := queryOrderQty(); err != nil {
		return msg, err
	} else {
		order.OrderQty = &fVal
	}

	switch order.OrdType {
	case enum.OrdType_LIMIT, enum.OrdType_STOP_LIMIT:
		if fVal, err := queryPrice(); err != nil {
			return msg, err
		} else {
			order.Price = &fVal
		}
	}

	switch order.OrdType {
	case enum.OrdType_STOP, enum.OrdType_STOP_LIMIT:
		if fVal, err := queryStopPx(); err != nil {
			return msg, err
		} else {
			order.StopPx = &fVal
		}
	}

	if tif, err := queryTimeInForce(); err != nil {
		return msg, err
	} else {
		order.TimeInForce = &tif
	}

	msg = quickfix.Marshal(order)
	queryHeader(msg.Header)
	return
}

func queryNewOrderSingle43() (msg quickfix.Message, err error) {
	order := fix43nos.Message{
		HandlInst:    "1",
		TransactTime: time.Now(),
	}
	if order.ClOrdID, err = queryClOrdID(); err != nil {
		return
	}

	if sym, err := querySymbol(); err != nil {
		return msg, err
	} else {
		order.Instrument.Symbol = &sym
	}

	if order.Side, err = querySide(); err != nil {
		return
	}

	if order.OrdType, err = queryOrdType(); err != nil {
		return
	}

	if fVal, err := queryOrderQty(); err != nil {
		return msg, err
	} else {
		order.OrderQtyData.OrderQty = &fVal
	}

	switch order.OrdType {
	case enum.OrdType_LIMIT, enum.OrdType_STOP_LIMIT:
		if fVal, err := queryPrice(); err != nil {
			return msg, err
		} else {
			order.Price = &fVal
		}
	}

	switch order.OrdType {
	case enum.OrdType_STOP, enum.OrdType_STOP_LIMIT:
		if fVal, err := queryStopPx(); err != nil {
			return msg, err
		} else {
			order.StopPx = &fVal
		}
	}

	if tif, err := queryTimeInForce(); err != nil {
		return msg, err
	} else {
		order.TimeInForce = &tif
	}

	msg = quickfix.Marshal(order)
	queryHeader(msg.Header)

	return
}

func queryNewOrderSingle44() (msg quickfix.Message, err error) {
	handleInst := "1"
	order := fix44nos.Message{
		HandlInst:    &handleInst,
		TransactTime: time.Now(),
	}

	if order.ClOrdID, err = queryClOrdID(); err != nil {
		return
	}

	if sym, err := querySymbol(); err != nil {
		return msg, err
	} else {
		order.Instrument.Symbol = &sym
	}

	if order.Side, err = querySide(); err != nil {
		return
	}

	if order.OrdType, err = queryOrdType(); err != nil {
		return
	}

	if fVal, err := queryOrderQty(); err != nil {
		return msg, err
	} else {
		order.OrderQtyData.OrderQty = &fVal
	}

	switch order.OrdType {
	case enum.OrdType_LIMIT, enum.OrdType_STOP_LIMIT:
		if fVal, err := queryPrice(); err != nil {
			return msg, err
		} else {
			order.Price = &fVal
		}
	}

	switch order.OrdType {
	case enum.OrdType_STOP, enum.OrdType_STOP_LIMIT:
		if fVal, err := queryStopPx(); err != nil {
			return msg, err
		} else {
			order.StopPx = &fVal
		}
	}

	if tif, err := queryTimeInForce(); err != nil {
		return msg, err
	} else {
		order.TimeInForce = &tif
	}

	msg = quickfix.Marshal(order)
	queryHeader(msg.Header)

	return
}

func queryNewOrderSingle50() (msg quickfix.Message, err error) {
	handleInst := "1"
	order := fix50nos.Message{
		HandlInst:    &handleInst,
		TransactTime: time.Now(),
	}

	if order.ClOrdID, err = queryClOrdID(); err != nil {
		return
	}

	if sym, err := querySymbol(); err != nil {
		return msg, err
	} else {
		order.Instrument.Symbol = &sym
	}

	if order.Side, err = querySide(); err != nil {
		return
	}

	if order.OrdType, err = queryOrdType(); err != nil {
		return
	}

	if fVal, err := queryOrderQty(); err != nil {
		return msg, err
	} else {
		order.OrderQtyData.OrderQty = &fVal
	}

	if tif, err := queryTimeInForce(); err != nil {
		return msg, err
	} else {
		order.TimeInForce = &tif
	}

	switch order.OrdType {
	case enum.OrdType_LIMIT, enum.OrdType_STOP_LIMIT:
		if fVal, err := queryPrice(); err != nil {
			return msg, err
		} else {
			order.Price = &fVal
		}
	}

	switch order.OrdType {
	case enum.OrdType_STOP, enum.OrdType_STOP_LIMIT:
		if fVal, err := queryStopPx(); err != nil {
			return msg, err
		} else {
			order.StopPx = &fVal
		}
	}

	msg = quickfix.Marshal(order)
	queryHeader(msg.Header)

	return
}

func queryOrderCancelRequest40() (msg quickfix.Message, err error) {
	cancel := fix40cxl.Message{
		CxlType: "F",
	}

	if cancel.OrigClOrdID, err = queryOrigClOrdID(); err != nil {
		return
	}

	if cancel.ClOrdID, err = queryClOrdID(); err != nil {
		return
	}

	if cancel.Symbol, err = querySymbol(); err != nil {
		return
	}

	if cancel.Side, err = querySide(); err != nil {
		return
	}

	if fVal, err := queryOrderQty(); err != nil {
		return msg, err
	} else {
		cancel.OrderQty = int(fVal)
	}

	msg = cancel.Marshal()
	queryHeader(msg.Header)
	return
}

func queryOrderCancelRequest41() (msg quickfix.Message, err error) {
	var cancel fix41cxl.Message

	if cancel.OrigClOrdID, err = queryOrigClOrdID(); err != nil {
		return
	}

	if cancel.ClOrdID, err = queryClOrdID(); err != nil {
		return
	}

	if cancel.Symbol, err = querySymbol(); err != nil {
		return
	}

	if cancel.Side, err = querySide(); err != nil {
		return
	}

	if fVal, err := queryOrderQty(); err != nil {
		return msg, err
	} else {
		iVal := int(fVal)
		cancel.OrderQty = &iVal
	}

	msg = cancel.Marshal()
	queryHeader(msg.Header)
	return
}

func queryOrderCancelRequest42() (msg quickfix.Message, err error) {
	cancel := fix42cxl.Message{
		TransactTime: time.Now(),
	}

	if cancel.OrigClOrdID, err = queryOrigClOrdID(); err != nil {
		return
	}

	if cancel.ClOrdID, err = queryClOrdID(); err != nil {
		return
	}

	if cancel.Symbol, err = querySymbol(); err != nil {
		return
	}

	if cancel.Side, err = querySide(); err != nil {
		return
	}

	if fVal, err := queryOrderQty(); err != nil {
		return msg, err
	} else {
		cancel.OrderQty = &fVal
	}

	msg = cancel.Marshal()
	queryHeader(msg.Header)
	return
}

func queryOrderCancelRequest43() (msg quickfix.Message, err error) {
	cancel := fix43cxl.Message{
		TransactTime: time.Now(),
	}

	if cancel.OrigClOrdID, err = queryOrigClOrdID(); err != nil {
		return
	}

	if cancel.ClOrdID, err = queryClOrdID(); err != nil {
		return
	}

	if sym, err := querySymbol(); err != nil {
		return msg, err
	} else {
		cancel.Instrument.Symbol = &sym
	}

	if cancel.Side, err = querySide(); err != nil {
		return
	}

	if fVal, err := queryOrderQty(); err != nil {
		return msg, err
	} else {
		cancel.OrderQtyData.OrderQty = &fVal
	}

	msg = cancel.Marshal()
	queryHeader(msg.Header)
	return
}

func queryOrderCancelRequest44() (msg quickfix.Message, err error) {
	cancel := fix44cxl.Message{
		TransactTime: time.Now(),
	}

	if cancel.OrigClOrdID, err = queryOrigClOrdID(); err != nil {
		return
	}

	if cancel.ClOrdID, err = queryClOrdID(); err != nil {
		return
	}

	if sym, err := querySymbol(); err != nil {
		return msg, err
	} else {
		cancel.Instrument.Symbol = &sym
	}

	if cancel.Side, err = querySide(); err != nil {
		return
	}

	if fVal, err := queryOrderQty(); err != nil {
		return msg, err
	} else {
		cancel.OrderQtyData.OrderQty = &fVal
	}

	msg = cancel.Marshal()
	queryHeader(msg.Header)
	return
}

func queryOrderCancelRequest50() (msg quickfix.Message, err error) {
	cancel := fix50cxl.Message{
		TransactTime: time.Now(),
	}

	if cancel.OrigClOrdID, err = queryOrigClOrdID(); err != nil {
		return
	}

	if cancel.ClOrdID, err = queryClOrdID(); err != nil {
		return
	}

	if sym, err := querySymbol(); err != nil {
		return msg, err
	} else {
		cancel.Instrument.Symbol = &sym
	}

	if cancel.Side, err = querySide(); err != nil {
		return
	}

	if fVal, err := queryOrderQty(); err != nil {
		return msg, err
	} else {
		cancel.OrderQtyData.OrderQty = &fVal
	}

	msg = cancel.Marshal()
	queryHeader(msg.Header)
	return
}

func queryEnterOrder() error {
	beginString, err := queryVersion()
	if err != nil {
		return err
	}

	var order quickfix.Message
	switch beginString {
	case enum.BeginStringFIX40:
		order, err = queryNewOrderSingle40()

	case enum.BeginStringFIX41:
		order, err = queryNewOrderSingle41()

	case enum.BeginStringFIX42:
		order, err = queryNewOrderSingle42()

	case enum.BeginStringFIX43:
		order, err = queryNewOrderSingle43()

	case enum.BeginStringFIX44:
		order, err = queryNewOrderSingle44()

	case enum.BeginStringFIXT11:
		order, err = queryNewOrderSingle50()
	}

	if err != nil {
		return err
	}
	bytes, _ := order.Build()

	fmt.Println("Sending ", string(bytes))

	return quickfix.Send(order)
}

func queryCancelOrder() error {
	beginString, err := queryVersion()
	if err != nil {
		return err
	}

	var cxl quickfix.Message
	switch beginString {
	case enum.BeginStringFIX40:
		cxl, err = queryOrderCancelRequest40()

	case enum.BeginStringFIX41:
		cxl, err = queryOrderCancelRequest41()

	case enum.BeginStringFIX42:
		cxl, err = queryOrderCancelRequest42()

	case enum.BeginStringFIX43:
		cxl, err = queryOrderCancelRequest43()

	case enum.BeginStringFIX44:
		cxl, err = queryOrderCancelRequest44()

	case enum.BeginStringFIXT11:
		cxl, err = queryOrderCancelRequest50()
	}

	if err != nil {
		return err
	}

	if queryConfirm("Send Cancel") {
		return quickfix.Send(cxl)
	}

	return nil
}
