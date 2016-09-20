package internal

import (
	"bufio"
	"fmt"
	"time"

	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/enum"
	"github.com/quickfixgo/quickfix/field"
	"github.com/shopspring/decimal"

	"os"
	"strconv"
	"strings"

	fix40nos "github.com/quickfixgo/quickfix/fix40/newordersingle"
	fix41nos "github.com/quickfixgo/quickfix/fix41/newordersingle"
	fix42nos "github.com/quickfixgo/quickfix/fix42/newordersingle"
	fix43nos "github.com/quickfixgo/quickfix/fix43/newordersingle"
	fix44nos "github.com/quickfixgo/quickfix/fix44/newordersingle"
	fix50nos "github.com/quickfixgo/quickfix/fix50/newordersingle"

	fix40cxl "github.com/quickfixgo/quickfix/fix40/ordercancelrequest"
	fix41cxl "github.com/quickfixgo/quickfix/fix41/ordercancelrequest"
	fix42cxl "github.com/quickfixgo/quickfix/fix42/ordercancelrequest"
	fix43cxl "github.com/quickfixgo/quickfix/fix43/ordercancelrequest"
	fix44cxl "github.com/quickfixgo/quickfix/fix44/ordercancelrequest"
	fix50cxl "github.com/quickfixgo/quickfix/fix50/ordercancelrequest"

	fix42mdr "github.com/quickfixgo/quickfix/fix42/marketdatarequest"
	fix43mdr "github.com/quickfixgo/quickfix/fix43/marketdatarequest"
	fix44mdr "github.com/quickfixgo/quickfix/fix44/marketdatarequest"
	fix50mdr "github.com/quickfixgo/quickfix/fix50/marketdatarequest"
)

func queryString(fieldName string) string {
	fmt.Printf("%v: ", fieldName)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return scanner.Text()
}

func queryDecimal(fieldName string) decimal.Decimal {
	val, err := decimal.NewFromString(queryString(fieldName))
	if err != nil {
		panic(err)
	}

	return val
}

func queryFieldChoices(fieldName string, choices []string, values []string) string {
	for i, choice := range choices {
		fmt.Printf("%v) %v\n", i+1, choice)
	}

	choiceStr := queryString(fieldName)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil || choice < 1 || choice > len(choices) {
		panic(fmt.Errorf("Invalid %v: %v", fieldName, choice))
	}

	if values == nil {
		return choiceStr
	}

	return values[choice-1]
}

func QueryAction() (string, error) {
	fmt.Println()
	fmt.Println("1) Enter Order")
	fmt.Println("2) Cancel Order")
	fmt.Println("3) Request Market Test")
	fmt.Println("4) Quit")
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
	}

	return "", fmt.Errorf("unknown BeginString choice: %v", scanner.Text())
}

func queryClOrdID() field.ClOrdIDField {
	return field.NewClOrdID(queryString("ClOrdID"))
}

func queryOrigClOrdID() field.OrigClOrdIDField {
	return field.NewOrigClOrdID(("OrigClOrdID"))
}

func querySymbol() field.SymbolField {
	return field.NewSymbol(queryString("Symbol"))
}

func querySide() field.SideField {
	choices := []string{
		"Buy",
		"Sell",
		"Sell Short",
		"Sell Short Exempt",
		"Cross",
		"Cross Short",
		"Cross Short Exempt",
	}

	values := []string{
		string(enum.Side_BUY),
		string(enum.Side_SELL),
		string(enum.Side_SELL_SHORT),
		string(enum.Side_SELL_SHORT_EXEMPT),
		string(enum.Side_CROSS),
		string(enum.Side_CROSS_SHORT),
		"A",
	}

	return field.NewSide(enum.Side(queryFieldChoices("Side", choices, values)))
}

func queryOrdType(f *field.OrdTypeField) field.OrdTypeField {
	choices := []string{
		"Market",
		"Limit",
		"Stop",
		"Stop Limit",
	}

	values := []string{
		string(enum.OrdType_MARKET),
		string(enum.OrdType_LIMIT),
		string(enum.OrdType_STOP),
		string(enum.OrdType_STOP_LIMIT),
	}

	f.FIXString = quickfix.FIXString(queryFieldChoices("OrdType", choices, values))
	return *f
}

func queryTimeInForce() field.TimeInForceField {
	choices := []string{
		"Day",
		"IOC",
		"OPG",
		"GTC",
		"GTX",
	}
	values := []string{
		string(enum.TimeInForce_DAY),
		string(enum.TimeInForce_IMMEDIATE_OR_CANCEL),
		string(enum.TimeInForce_AT_THE_OPENING),
		string(enum.TimeInForce_GOOD_TILL_CANCEL),
		string(enum.TimeInForce_GOOD_TILL_CROSSING),
	}

	return field.NewTimeInForce(enum.TimeInForce(queryFieldChoices("TimeInForce", choices, values)))
}

func queryOrderQty() field.OrderQtyField {
	return field.NewOrderQty(queryDecimal("OrderQty"), 2)
}

func queryPrice() field.PriceField {
	return field.NewPrice(queryDecimal("Price"), 2)
}

func queryStopPx() field.StopPxField {
	return field.NewStopPx(queryDecimal("Stop Price"), 2)
}

func querySenderCompID() field.SenderCompIDField {
	return field.NewSenderCompID(queryString("SenderCompID"))
}

func queryTargetCompID() field.TargetCompIDField {
	return field.NewTargetCompID(queryString("TargetCompID"))
}

func queryTargetSubID() field.TargetSubIDField {
	return field.NewTargetSubID(queryString("TargetSubID"))
}

func queryConfirm(prompt string) bool {
	fmt.Println()
	fmt.Printf("%v?: ", prompt)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return strings.ToUpper(scanner.Text()) == "Y"
}

type header interface {
	Set(f quickfix.FieldWriter) quickfix.FieldMap
}

func queryHeader(h header) {
	h.Set(querySenderCompID())
	h.Set(queryTargetCompID())
	if ok := queryConfirm("Use a TargetSubID"); !ok {
		return
	}

	h.Set(queryTargetSubID())
}

func queryNewOrderSingle40() fix40nos.NewOrderSingle {
	var ordType field.OrdTypeField
	order := fix40nos.New(queryClOrdID(), field.NewHandlInst("1"), querySymbol(), querySide(), queryOrderQty(), queryOrdType(&ordType))

	switch ordType.Value() {
	case enum.OrdType_LIMIT, enum.OrdType_STOP_LIMIT:
		order.Set(queryPrice())
	}

	switch ordType.Value() {
	case enum.OrdType_STOP, enum.OrdType_STOP_LIMIT:
		order.Set(queryStopPx())
	}

	order.Set(queryTimeInForce())
	queryHeader(order.Header.Header)

	return order
}

func queryNewOrderSingle41() (msg quickfix.Message) {
	var ordType field.OrdTypeField
	order := fix41nos.New(queryClOrdID(), field.NewHandlInst("1"), querySymbol(), querySide(), queryOrdType(&ordType))
	order.Set(queryOrderQty())

	switch ordType.Value() {
	case enum.OrdType_LIMIT, enum.OrdType_STOP_LIMIT:
		order.Set(queryPrice())
	}

	switch ordType.Value() {
	case enum.OrdType_STOP, enum.OrdType_STOP_LIMIT:
		order.Set(queryStopPx())
	}

	order.Set(queryTimeInForce())
	msg = order.ToMessage()
	queryHeader(msg.Header)

	return
}

func queryNewOrderSingle42() (msg quickfix.Message) {
	var ordType field.OrdTypeField
	order := fix42nos.New(queryClOrdID(), field.NewHandlInst("1"), querySymbol(), querySide(), field.NewTransactTime(time.Now()), queryOrdType(&ordType))
	order.Set(queryOrderQty())

	switch ordType.Value() {
	case enum.OrdType_LIMIT, enum.OrdType_STOP_LIMIT:
		order.Set(queryPrice())
	}

	switch ordType.Value() {
	case enum.OrdType_STOP, enum.OrdType_STOP_LIMIT:
		order.Set(queryStopPx())
	}

	order.Set(queryTimeInForce())
	msg = order.ToMessage()
	queryHeader(msg.Header)
	return
}

func queryNewOrderSingle43() (msg quickfix.Message) {
	var ordType field.OrdTypeField
	order := fix43nos.New(queryClOrdID(), field.NewHandlInst("1"), querySide(), field.NewTransactTime(time.Now()), queryOrdType(&ordType))
	order.Set(querySymbol())
	order.Set(queryOrderQty())

	switch ordType.Value() {
	case enum.OrdType_LIMIT, enum.OrdType_STOP_LIMIT:
		order.Set(queryPrice())
	}

	switch ordType.Value() {
	case enum.OrdType_STOP, enum.OrdType_STOP_LIMIT:
		order.Set(queryStopPx())
	}

	order.Set(queryTimeInForce())
	msg = order.ToMessage()
	queryHeader(msg.Header)

	return
}

func queryNewOrderSingle44() (msg quickfix.Message) {
	var ordType field.OrdTypeField
	order := fix44nos.New(queryClOrdID(), querySide(), field.NewTransactTime(time.Now()), queryOrdType(&ordType))
	order.SetHandlInst("1")
	order.Set(querySymbol())
	order.Set(queryOrderQty())

	switch ordType.Value() {
	case enum.OrdType_LIMIT, enum.OrdType_STOP_LIMIT:
		order.Set(queryPrice())
	}

	switch ordType.Value() {
	case enum.OrdType_STOP, enum.OrdType_STOP_LIMIT:
		order.Set(queryStopPx())
	}

	order.Set(queryTimeInForce())
	msg = order.ToMessage()
	queryHeader(msg.Header)

	return
}

func queryNewOrderSingle50() (msg quickfix.Message) {
	var ordType field.OrdTypeField
	order := fix50nos.New(queryClOrdID(), querySide(), field.NewTransactTime(time.Now()), queryOrdType(&ordType))
	order.SetHandlInst("1")
	order.Set(querySymbol())
	order.Set(queryOrderQty())
	order.Set(queryTimeInForce())

	switch ordType.Value() {
	case enum.OrdType_LIMIT, enum.OrdType_STOP_LIMIT:
		order.Set(queryPrice())
	}

	switch ordType.Value() {
	case enum.OrdType_STOP, enum.OrdType_STOP_LIMIT:
		order.Set(queryStopPx())
	}

	msg = order.ToMessage()
	queryHeader(msg.Header)

	return
}

func queryOrderCancelRequest40() (msg quickfix.Message) {
	cancel := fix40cxl.New(queryOrigClOrdID(), queryClOrdID(), field.NewCxlType("F"), querySymbol(), querySide(), queryOrderQty())
	msg = cancel.ToMessage()
	queryHeader(msg.Header)
	return
}

func queryOrderCancelRequest41() (msg quickfix.Message) {
	cancel := fix41cxl.New(queryOrigClOrdID(), queryClOrdID(), querySymbol(), querySide())
	cancel.Set(queryOrderQty())
	msg = cancel.ToMessage()
	queryHeader(msg.Header)
	return
}

func queryOrderCancelRequest42() (msg quickfix.Message) {
	cancel := fix42cxl.New(queryOrigClOrdID(), queryClOrdID(), querySymbol(), querySide(), field.NewTransactTime(time.Now()))
	cancel.Set(queryOrderQty())
	msg = cancel.ToMessage()
	queryHeader(msg.Header)
	return
}

func queryOrderCancelRequest43() (msg quickfix.Message) {
	cancel := fix43cxl.New(queryOrigClOrdID(), queryClOrdID(), querySide(), field.NewTransactTime(time.Now()))
	cancel.Set(querySymbol())
	cancel.Set(queryOrderQty())
	msg = cancel.ToMessage()
	queryHeader(msg.Header)
	return
}

func queryOrderCancelRequest44() (msg quickfix.Message) {
	cancel := fix44cxl.New(queryOrigClOrdID(), queryClOrdID(), querySide(), field.NewTransactTime(time.Now()))
	cancel.Set(querySymbol())
	cancel.Set(queryOrderQty())

	msg = cancel.ToMessage()
	queryHeader(msg.Header)
	return
}

func queryOrderCancelRequest50() (msg quickfix.Message) {
	cancel := fix50cxl.New(queryOrigClOrdID(), queryClOrdID(), querySide(), field.NewTransactTime(time.Now()))
	cancel.Set(querySymbol())
	cancel.Set(queryOrderQty())
	msg = cancel.ToMessage()
	queryHeader(msg.Header)
	return
}

func queryMarketDataRequest42() fix42mdr.MarketDataRequest {
	request := fix42mdr.New(field.NewMDReqID("MARKETDATAID"),
		field.NewSubscriptionRequestType(enum.SubscriptionRequestType_SNAPSHOT),
		field.NewMarketDepth(0),
	)

	entryTypes := fix42mdr.NewNoMDEntryTypesRepeatingGroup()
	entryTypes.Add().SetMDEntryType(enum.MDEntryType_BID)
	request.SetNoMDEntryTypes(entryTypes)

	relatedSym := fix42mdr.NewNoRelatedSymRepeatingGroup()
	relatedSym.Add().SetSymbol("LNUX")
	request.SetNoRelatedSym(relatedSym)

	queryHeader(request.Header)
	return request
}

func queryMarketDataRequest43() fix43mdr.MarketDataRequest {
	request := fix43mdr.New(field.NewMDReqID("MARKETDATAID"),
		field.NewSubscriptionRequestType(enum.SubscriptionRequestType_SNAPSHOT),
		field.NewMarketDepth(0),
	)

	entryTypes := fix43mdr.NewNoMDEntryTypesRepeatingGroup()
	entryTypes.Add().SetMDEntryType(enum.MDEntryType_BID)
	request.SetNoMDEntryTypes(entryTypes)

	relatedSym := fix43mdr.NewNoRelatedSymRepeatingGroup()
	relatedSym.Add().SetSymbol("LNUX")
	request.SetNoRelatedSym(relatedSym)

	queryHeader(request.Header)
	return request
}

func queryMarketDataRequest44() fix44mdr.MarketDataRequest {
	request := fix44mdr.New(field.NewMDReqID("MARKETDATAID"),
		field.NewSubscriptionRequestType(enum.SubscriptionRequestType_SNAPSHOT),
		field.NewMarketDepth(0),
	)

	entryTypes := fix44mdr.NewNoMDEntryTypesRepeatingGroup()
	entryTypes.Add().SetMDEntryType(enum.MDEntryType_BID)
	request.SetNoMDEntryTypes(entryTypes)

	relatedSym := fix44mdr.NewNoRelatedSymRepeatingGroup()
	relatedSym.Add().SetSymbol("LNUX")
	request.SetNoRelatedSym(relatedSym)

	queryHeader(request.Header)
	return request
}

func queryMarketDataRequest50() fix50mdr.MarketDataRequest {
	request := fix50mdr.New(field.NewMDReqID("MARKETDATAID"),
		field.NewSubscriptionRequestType(enum.SubscriptionRequestType_SNAPSHOT),
		field.NewMarketDepth(0),
	)

	entryTypes := fix50mdr.NewNoMDEntryTypesRepeatingGroup()
	entryTypes.Add().SetMDEntryType(enum.MDEntryType_BID)
	request.SetNoMDEntryTypes(entryTypes)

	relatedSym := fix50mdr.NewNoRelatedSymRepeatingGroup()
	relatedSym.Add().SetSymbol("LNUX")
	request.SetNoRelatedSym(relatedSym)

	queryHeader(request.Header)
	return request
}

func QueryEnterOrder() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	var beginString string
	beginString, err = queryVersion()
	if err != nil {
		return err
	}

	var order quickfix.Messagable
	switch beginString {
	case enum.BeginStringFIX40:
		order = queryNewOrderSingle40()

	case enum.BeginStringFIX41:
		order = queryNewOrderSingle41()

	case enum.BeginStringFIX42:
		order = queryNewOrderSingle42()

	case enum.BeginStringFIX43:
		order = queryNewOrderSingle43()

	case enum.BeginStringFIX44:
		order = queryNewOrderSingle44()

	case enum.BeginStringFIXT11:
		order = queryNewOrderSingle50()
	}

	return quickfix.Send(order)
}

func QueryCancelOrder() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	var beginString string
	beginString, err = queryVersion()
	if err != nil {
		return err
	}

	var cxl quickfix.Message
	switch beginString {
	case enum.BeginStringFIX40:
		cxl = queryOrderCancelRequest40()

	case enum.BeginStringFIX41:
		cxl = queryOrderCancelRequest41()

	case enum.BeginStringFIX42:
		cxl = queryOrderCancelRequest42()

	case enum.BeginStringFIX43:
		cxl = queryOrderCancelRequest43()

	case enum.BeginStringFIX44:
		cxl = queryOrderCancelRequest44()

	case enum.BeginStringFIXT11:
		cxl = queryOrderCancelRequest50()
	}

	if queryConfirm("Send Cancel") {
		return quickfix.Send(cxl)
	}

	return
}

func QueryMarketDataRequest() error {
	beginString, err := queryVersion()
	if err != nil {
		return err
	}

	var req quickfix.Messagable
	switch beginString {
	case enum.BeginStringFIX42:
		req = queryMarketDataRequest42()

	case enum.BeginStringFIX43:
		req = queryMarketDataRequest43()

	case enum.BeginStringFIX44:
		req = queryMarketDataRequest44()

	case enum.BeginStringFIXT11:
		req = queryMarketDataRequest50()

	default:
		return fmt.Errorf("No test for version %v", beginString)
	}

	if queryConfirm("Send MarketDataRequest") {
		return quickfix.Send(req)
	}

	return nil
}
