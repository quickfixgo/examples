package main

import (
	"fmt"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/fix/enum"
	"github.com/quickfixgo/quickfix/fix/field"

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

func (e *Executor) genExecID() string {
	e.execID++
	return strconv.Itoa(e.execID)
}

//quickfix.Application interface
func (e Executor) OnCreate(sessionID quickfix.SessionID)                                 { return }
func (e Executor) OnLogon(sessionID quickfix.SessionID)                                  { return }
func (e Executor) OnLogout(sessionID quickfix.SessionID)                                 { return }
func (e Executor) ToAdmin(msg quickfix.MessageBuilder, sessionID quickfix.SessionID)     { return }
func (e Executor) ToApp(msg quickfix.MessageBuilder, sessionID quickfix.SessionID) error { return nil }
func (e Executor) FromAdmin(msg quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

//Use Message Cracker on Incoming Application Messages
func (e *Executor) FromApp(msg quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return e.Route(msg, sessionID)
}

func (e *Executor) OnFIX40NewOrderSingle(msg fix40nos.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	symbol, err := msg.Symbol()
	if err != nil {
		return
	}

	side, err := msg.Side()
	if err != nil {
		return
	}

	orderQty, err := msg.OrderQty()
	if err != nil {
		return
	}

	ordType, err := msg.OrdType()
	if err != nil {
		return
	}

	if ordType.Value != enum.OrdType_LIMIT {
		err = quickfix.ValueIsIncorrect(ordType.Tag())
		return
	}

	price, err := msg.Price()
	if err != nil {
		return
	}

	clOrdID, err := msg.ClOrdID()
	if err != nil {
		return
	}

	execReport := fix40er.Builder(
		field.NewOrderID(e.genOrderID()),
		field.NewExecID(e.genExecID()),
		field.NewExecTransType(enum.ExecTransType_NEW),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		symbol, side, orderQty, field.NewLastShares(orderQty.Value), field.NewLastPx(price.Value), field.NewCumQty(orderQty.Value), field.NewAvgPx(price.Value))

	execReport.Body().Set(clOrdID)

	if acct, err := msg.Account(); err != nil {
		execReport.Body().Set(acct)
	}

	quickfix.SendToTarget(execReport.MessageBuilder, sessionID)

	return
}

func (e *Executor) OnFIX41NewOrderSingle(msg fix41nos.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	symbol, err := msg.Symbol()
	if err != nil {
		return
	}

	side, err := msg.Side()
	if err != nil {
		return
	}

	orderQty, err := msg.OrderQty()
	if err != nil {
		return
	}

	ordType, err := msg.OrdType()
	if err != nil {
		return
	}

	if ordType.Value != enum.OrdType_LIMIT {
		err = quickfix.ValueIsIncorrect(ordType.Tag())
		return
	}

	price, err := msg.Price()
	if err != nil {
		return
	}

	clOrdID, err := msg.ClOrdID()
	if err != nil {
		return
	}

	execReport := fix41er.Builder(
		field.NewOrderID(e.genOrderID()),
		field.NewExecID(e.genExecID()),
		field.NewExecTransType(enum.ExecTransType_NEW),
		field.NewExecType(enum.ExecType_FILL),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		symbol,
		side,
		orderQty,
		field.NewLastShares(orderQty.Value),
		field.NewLastPx(price.Value),
		field.NewLeavesQty(0),
		field.NewCumQty(orderQty.Value),
		field.NewAvgPx(price.Value))

	execReport.Body().Set(clOrdID)

	if acct, err := msg.Account(); err != nil {
		execReport.Body().Set(acct)
	}

	quickfix.SendToTarget(execReport.MessageBuilder, sessionID)

	return
}

func (e *Executor) OnFIX42NewOrderSingle(msg fix42nos.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	symbol, err := msg.Symbol()
	if err != nil {
		return
	}

	side, err := msg.Side()
	if err != nil {
		return
	}

	orderQty, err := msg.OrderQty()
	if err != nil {
		return
	}

	ordType, err := msg.OrdType()
	if err != nil {
		return
	}

	if ordType.Value != enum.OrdType_LIMIT {
		err = quickfix.ValueIsIncorrect(ordType.Tag())
		return
	}

	price, err := msg.Price()
	if err != nil {
		return
	}

	clOrdID, err := msg.ClOrdID()
	if err != nil {
		return
	}

	execReport := fix42er.Builder(
		field.NewOrderID(e.genOrderID()),
		field.NewExecID(e.genExecID()),
		field.NewExecTransType(enum.ExecTransType_NEW),
		field.NewExecType(enum.ExecType_FILL),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		symbol,
		side,
		field.NewLeavesQty(0),
		field.NewCumQty(orderQty.Value),
		field.NewAvgPx(price.Value))

	execReport.Body().Set(clOrdID)
	execReport.Body().Set(orderQty)
	execReport.Body().Set(field.NewLastShares(orderQty.Value))
	execReport.Body().Set(field.NewLastPx(price.Value))

	if acct, err := msg.Account(); err != nil {
		execReport.Body().Set(acct)
	}

	quickfix.SendToTarget(execReport.MessageBuilder, sessionID)

	return
}

func (e *Executor) OnFIX43NewOrderSingle(msg fix43nos.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	symbol, err := msg.Symbol()
	if err != nil {
		return
	}

	side, err := msg.Side()
	if err != nil {
		return
	}

	orderQty, err := msg.OrderQty()
	if err != nil {
		return
	}

	ordType, err := msg.OrdType()
	if err != nil {
		return
	}

	if ordType.Value != enum.OrdType_LIMIT {
		err = quickfix.ValueIsIncorrect(ordType.Tag())
		return
	}

	price, err := msg.Price()
	if err != nil {
		return
	}

	clOrdID, err := msg.ClOrdID()
	if err != nil {
		return
	}

	execReport := fix43er.Builder(
		field.NewOrderID(e.genOrderID()),
		field.NewExecID(e.genExecID()),
		field.NewExecType(enum.ExecType_FILL),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		side,
		field.NewLeavesQty(0),
		field.NewCumQty(orderQty.Value),
		field.NewAvgPx(price.Value))

	execReport.Body().Set(clOrdID)
	execReport.Body().Set(symbol)
	execReport.Body().Set(orderQty)
	execReport.Body().Set(field.NewLastShares(orderQty.Value))
	execReport.Body().Set(field.NewLastPx(price.Value))

	if acct, err := msg.Account(); err != nil {
		execReport.Body().Set(acct)
	}

	quickfix.SendToTarget(execReport.MessageBuilder, sessionID)

	return
}

func (e *Executor) OnFIX44NewOrderSingle(msg fix44nos.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	symbol, err := msg.Symbol()
	if err != nil {
		return
	}

	side, err := msg.Side()
	if err != nil {
		return
	}

	orderQty, err := msg.OrderQty()
	if err != nil {
		return
	}

	ordType, err := msg.OrdType()
	if err != nil {
		return
	}

	if ordType.Value != enum.OrdType_LIMIT {
		err = quickfix.ValueIsIncorrect(ordType.Tag())
		return
	}

	price, err := msg.Price()
	if err != nil {
		return
	}

	clOrdID, err := msg.ClOrdID()
	if err != nil {
		return
	}

	execReport := fix44er.Builder(
		field.NewOrderID(e.genOrderID()),
		field.NewExecID(e.genExecID()),
		field.NewExecType(enum.ExecType_FILL),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		side,
		field.NewLeavesQty(0),
		field.NewCumQty(orderQty.Value),
		field.NewAvgPx(price.Value))

	execReport.Body().Set(clOrdID)
	execReport.Body().Set(symbol)
	execReport.Body().Set(orderQty)
	execReport.Body().Set(field.NewLastQty(orderQty.Value))
	execReport.Body().Set(field.NewLastPx(price.Value))

	if acct, err := msg.Account(); err != nil {
		execReport.Body().Set(acct)
	}

	quickfix.SendToTarget(execReport.MessageBuilder, sessionID)

	return
}

func (e *Executor) OnFIX50NewOrderSingle(msg fix50nos.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	symbol, err := msg.Symbol()
	if err != nil {
		return
	}

	side, err := msg.Side()
	if err != nil {
		return
	}

	orderQty, err := msg.OrderQty()
	if err != nil {
		return
	}

	ordType, err := msg.OrdType()
	if err != nil {
		return
	}

	if ordType.Value != enum.OrdType_LIMIT {
		err = quickfix.ValueIsIncorrect(ordType.Tag())
		return
	}

	price, err := msg.Price()
	if err != nil {
		return
	}

	clOrdID, err := msg.ClOrdID()
	if err != nil {
		return
	}

	execReport := fix50er.Builder(
		field.NewOrderID(e.genOrderID()),
		field.NewExecID(e.genExecID()),
		field.NewExecType(enum.ExecType_FILL),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		side,
		field.NewLeavesQty(0),
		field.NewCumQty(orderQty.Value))

	execReport.Body().Set(clOrdID)
	execReport.Body().Set(symbol)
	execReport.Body().Set(orderQty)
	execReport.Body().Set(field.NewLastQty(orderQty.Value))
	execReport.Body().Set(field.NewLastPx(price.Value))
	execReport.Body().Set(field.NewAvgPx(price.Value))

	if acct, err := msg.Account(); err != nil {
		execReport.Body().Set(acct)
	}

	quickfix.SendToTarget(execReport.MessageBuilder, sessionID)

	return
}

func main() {
	cfgFileName := "executor.cfg"
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

	fileLogFactory, err := quickfix.NewFileLogFactory(appSettings)

	if err != nil {
		fmt.Println("Error creating file log factory,", err)
		return
	}

	app := NewExecutor()

	acceptor, err := quickfix.NewAcceptor(app, quickfix.NewMemoryStoreFactory(), appSettings, fileLogFactory)
	if err != nil {
		fmt.Printf("Unable to create Acceptor: %s\n", err)
		return
	}

	acceptor.Start()

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt)
	<-interrupt

	acceptor.Stop()
}
