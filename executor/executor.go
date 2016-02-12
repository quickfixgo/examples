package main

import (
	"flag"
	"fmt"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/enum"
	"github.com/quickfixgo/quickfix/tag"

	fix40nos "github.com/quickfixgo/quickfix/fix40/newordersingle"
	fix41nos "github.com/quickfixgo/quickfix/fix41/newordersingle"
	fix42nos "github.com/quickfixgo/quickfix/fix42/newordersingle"
	fix43nos "github.com/quickfixgo/quickfix/fix43/newordersingle"
	fix44nos "github.com/quickfixgo/quickfix/fix44/newordersingle"
	fix50nos "github.com/quickfixgo/quickfix/fix50/newordersingle"

	fix40er "github.com/quickfixgo/quickfix/fix40/executionreport"
	fix41er "github.com/quickfixgo/quickfix/fix41/executionreport"
	fix42er "github.com/quickfixgo/quickfix/fix42/executionreport"
	fix43er "github.com/quickfixgo/quickfix/fix43/executionreport"
	fix44er "github.com/quickfixgo/quickfix/fix44/executionreport"
	fix50er "github.com/quickfixgo/quickfix/fix50/executionreport"

	"os"
	"os/signal"
	"strconv"
)

type Executor struct {
	orderID int
	execID  int
	*quickfix.MessageRouter
}

func NewExecutor() *Executor {
	e := &Executor{MessageRouter: quickfix.NewMessageRouter()}
	e.AddRoute(fix40nos.Route(e.OnFIX40NewOrderSingle))
	e.AddRoute(fix41nos.Route(e.OnFIX41NewOrderSingle))
	e.AddRoute(fix42nos.Route(e.OnFIX42NewOrderSingle))
	e.AddRoute(fix43nos.Route(e.OnFIX43NewOrderSingle))
	e.AddRoute(fix44nos.Route(e.OnFIX44NewOrderSingle))
	e.AddRoute(fix50nos.Route(e.OnFIX50NewOrderSingle))

	return e
}

func (e *Executor) genOrderID() string {
	e.orderID++
	return strconv.Itoa(e.orderID)
}

func (e *Executor) genExecID() int {
	e.execID++
	return e.execID
}

//quickfix.Application interface
func (e Executor) OnCreate(sessionID quickfix.SessionID)                          { return }
func (e Executor) OnLogon(sessionID quickfix.SessionID)                           { return }
func (e Executor) OnLogout(sessionID quickfix.SessionID)                          { return }
func (e Executor) ToAdmin(msg quickfix.Message, sessionID quickfix.SessionID)     { return }
func (e Executor) ToApp(msg quickfix.Message, sessionID quickfix.SessionID) error { return nil }
func (e Executor) FromAdmin(msg quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

//Use Message Cracker on Incoming Application Messages
func (e *Executor) FromApp(msg quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return e.Route(msg, sessionID)
}

func (e *Executor) OnFIX40NewOrderSingle(msg fix40nos.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	if msg.OrdType != enum.OrdType_LIMIT {
		err = quickfix.ValueIsIncorrect(tag.OrdType)
		return
	}

	if msg.Price == nil {
		err = quickfix.ConditionallyRequiredFieldMissing(tag.Price)
		return
	}

	execReport := fix40er.Message{
		ClOrdID:       &msg.ClOrdID,
		Account:       msg.Account,
		OrderID:       e.genOrderID(),
		ExecID:        e.genExecID(),
		ExecTransType: enum.ExecTransType_NEW,
		OrdStatus:     enum.OrdStatus_FILLED,
		Symbol:        msg.Symbol,
		Side:          msg.Side,
		OrderQty:      msg.OrderQty,
		LastShares:    msg.OrderQty,
		LastPx:        *msg.Price,
		CumQty:        msg.OrderQty,
		AvgPx:         *msg.Price,
	}

	quickfix.SendToTarget(execReport, sessionID)

	return
}

func (e *Executor) OnFIX41NewOrderSingle(msg fix41nos.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	if msg.OrdType != enum.OrdType_LIMIT {
		err = quickfix.ValueIsIncorrect(tag.OrdType)
		return
	}

	if msg.Price == nil {
		err = quickfix.ConditionallyRequiredFieldMissing(tag.Price)
		return
	}

	if msg.OrderQty == nil {
		err = quickfix.ConditionallyRequiredFieldMissing(tag.OrderQty)
		return
	}

	execReport := fix41er.Message{
		ClOrdID:       &msg.ClOrdID,
		Account:       msg.Account,
		OrderID:       e.genOrderID(),
		ExecID:        strconv.Itoa(e.genExecID()),
		ExecTransType: enum.ExecTransType_NEW,
		ExecType:      enum.ExecType_FILL,
		OrdStatus:     enum.OrdStatus_FILLED,
		Symbol:        msg.Symbol,
		Side:          msg.Side,
		OrderQty:      *msg.OrderQty,
		LastShares:    *msg.OrderQty,
		LastPx:        *msg.Price,
		LeavesQty:     0,
		CumQty:        *msg.OrderQty,
		AvgPx:         *msg.Price,
	}

	quickfix.SendToTarget(execReport, sessionID)

	return
}

func (e *Executor) OnFIX42NewOrderSingle(msg fix42nos.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	if msg.OrdType != enum.OrdType_LIMIT {
		err = quickfix.ValueIsIncorrect(tag.OrdType)
		return
	}

	if msg.Price == nil {
		err = quickfix.ConditionallyRequiredFieldMissing(tag.Price)
		return
	}

	if msg.OrderQty == nil {
		err = quickfix.ConditionallyRequiredFieldMissing(tag.OrderQty)
		return
	}

	execReport := fix42er.Message{
		ClOrdID:       &msg.ClOrdID,
		Account:       msg.Account,
		OrderID:       e.genOrderID(),
		ExecID:        strconv.Itoa(e.genExecID()),
		ExecTransType: enum.ExecTransType_NEW,
		ExecType:      enum.ExecType_FILL,
		OrdStatus:     enum.OrdStatus_FILLED,
		Symbol:        msg.Symbol,
		Side:          msg.Side,
		OrderQty:      msg.OrderQty,
		LeavesQty:     0,
		LastShares:    msg.OrderQty,
		CumQty:        *msg.OrderQty,
		AvgPx:         *msg.Price,
		LastPx:        msg.Price,
	}

	quickfix.SendToTarget(execReport, sessionID)

	return
}

func (e *Executor) OnFIX43NewOrderSingle(msg fix43nos.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	if msg.OrdType != enum.OrdType_LIMIT {
		err = quickfix.ValueIsIncorrect(tag.OrdType)
		return
	}

	if msg.Price == nil {
		err = quickfix.ConditionallyRequiredFieldMissing(tag.Price)
		return
	}

	if msg.OrderQtyData.OrderQty == nil {
		err = quickfix.ConditionallyRequiredFieldMissing(tag.OrderQty)
		return
	}

	execReport := fix43er.Message{
		ClOrdID:      &msg.ClOrdID,
		Account:      msg.Account,
		OrderID:      e.genOrderID(),
		ExecID:       strconv.Itoa(e.genExecID()),
		ExecType:     enum.ExecType_FILL,
		OrdStatus:    enum.OrdStatus_FILLED,
		Side:         msg.Side,
		Instrument:   msg.Instrument,
		OrderQtyData: msg.OrderQtyData,
		LeavesQty:    0,
		LastQty:      msg.OrderQtyData.OrderQty,
		CumQty:       *msg.OrderQtyData.OrderQty,
		AvgPx:        *msg.Price,
		LastPx:       msg.Price,
	}
	quickfix.SendToTarget(execReport, sessionID)

	return
}

func (e *Executor) OnFIX44NewOrderSingle(msg fix44nos.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	if msg.OrdType != enum.OrdType_LIMIT {
		err = quickfix.ValueIsIncorrect(tag.OrdType)
		return
	}

	if msg.Price == nil {
		err = quickfix.ConditionallyRequiredFieldMissing(tag.Price)
		return
	}

	if msg.OrderQtyData.OrderQty == nil {
		err = quickfix.ConditionallyRequiredFieldMissing(tag.OrderQty)
		return
	}

	execReport := fix44er.Message{
		ClOrdID:      &msg.ClOrdID,
		OrderID:      e.genOrderID(),
		Account:      msg.Account,
		ExecID:       strconv.Itoa(e.genExecID()),
		ExecType:     enum.ExecType_FILL,
		OrdStatus:    enum.OrdStatus_FILLED,
		Side:         msg.Side,
		Instrument:   msg.Instrument,
		LeavesQty:    0,
		CumQty:       *msg.OrderQtyData.OrderQty,
		LastQty:      msg.OrderQtyData.OrderQty,
		OrderQtyData: msg.OrderQtyData,
		AvgPx:        *msg.Price,
		LastPx:       msg.Price,
	}

	quickfix.SendToTarget(execReport, sessionID)

	return
}

func (e *Executor) OnFIX50NewOrderSingle(msg fix50nos.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	if msg.OrdType != enum.OrdType_LIMIT {
		err = quickfix.ValueIsIncorrect(tag.OrdType)
		return
	}

	if msg.Price == nil {
		err = quickfix.ConditionallyRequiredFieldMissing(tag.Price)
		return
	}

	if msg.OrderQtyData.OrderQty == nil {
		err = quickfix.ConditionallyRequiredFieldMissing(tag.OrderQty)
		return
	}

	execReport := fix50er.Message{
		ClOrdID:      &msg.ClOrdID,
		Instrument:   msg.Instrument,
		Account:      msg.Account,
		OrderID:      e.genOrderID(),
		ExecID:       strconv.Itoa(e.genExecID()),
		ExecType:     enum.ExecType_FILL,
		OrdStatus:    enum.OrdStatus_FILLED,
		Side:         msg.Side,
		OrderQtyData: msg.OrderQtyData,
		LastQty:      msg.OrderQtyData.OrderQty,
		LeavesQty:    0,
		CumQty:       *msg.OrderQtyData.OrderQty,
		LastPx:       msg.Price,
		AvgPx:        msg.Price,
	}

	quickfix.SendToTarget(execReport, sessionID)

	return
}

func main() {
	flag.Parse()

	cfgFileName := "executor.cfg"
	if flag.NArg() > 0 {
		cfgFileName = flag.Arg(0)
	}

	cfg, err := os.Open(cfgFileName)
	if err != nil {
		fmt.Printf("Error opening %v, %v\n", cfgFileName, err)
		return
	}

	appSettings, err := quickfix.ParseSettings(cfg)
	if err != nil {
		fmt.Println("Error reading cfg,", err)
		return
	}

	logFactory := quickfix.NewScreenLogFactory()
	app := NewExecutor()

	acceptor, err := quickfix.NewAcceptor(app, quickfix.NewMemoryStoreFactory(), appSettings, logFactory)
	if err != nil {
		fmt.Printf("Unable to create Acceptor: %s\n", err)
		return
	}

	err = acceptor.Start()
	if err != nil {
		fmt.Printf("Unable to start Acceptor: %s\n", err)
		return
	}

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt)
	<-interrupt

	acceptor.Stop()
}
