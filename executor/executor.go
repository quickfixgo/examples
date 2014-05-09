package main

import (
	"fmt"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/cracker"
	"github.com/quickfixgo/quickfix/errors"
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
	"os/signal"
	"strconv"
)

type Executor struct {
	cracker.MessageCracker
	orderID int
	execID  int
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
func (e Executor) OnCreate(sessionID quickfix.SessionID)                                { return }
func (e Executor) OnLogon(sessionID quickfix.SessionID)                                 { return }
func (e Executor) OnLogout(sessionID quickfix.SessionID)                                { return }
func (e Executor) ToAdmin(msg message.MessageBuilder, sessionID quickfix.SessionID)     { return }
func (e Executor) ToApp(msg message.MessageBuilder, sessionID quickfix.SessionID) error { return nil }
func (e Executor) FromAdmin(msg message.Message, sessionID quickfix.SessionID) errors.MessageRejectError {
	return nil
}

//Use Message Cracker on Incoming Application Messages
func (e *Executor) FromApp(msg message.Message, sessionID quickfix.SessionID) (reject errors.MessageRejectError) {
	return cracker.Crack(msg, sessionID, e)
}

func (e *Executor) OnFIX40NewOrderSingle(msg fix40.NewOrderSingle, sessionID quickfix.SessionID) (err errors.MessageRejectError) {
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
		err = errors.ValueIsIncorrect(ordType.Tag())
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

	execReport := fix40.CreateExecutionReportBuilder(
		field.NewOrderID(e.genOrderID()),
		field.NewExecID(e.genExecID()),
		field.NewExecTransType(enum.ExecTransType_NEW),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		symbol, side, orderQty, field.NewLastShares(orderQty.Value), field.NewLastPx(price.Value), field.NewCumQty(orderQty.Value), field.NewAvgPx(price.Value))

	execReport.Body.Set(clOrdID)

	if acct, err := msg.Account(); err != nil {
		execReport.Body.Set(acct)
	}

	quickfix.SendToTarget(execReport.MessageBuilder, sessionID)

	return
}

func (e *Executor) OnFIX41NewOrderSingle(msg fix41.NewOrderSingle, sessionID quickfix.SessionID) (err errors.MessageRejectError) {
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
		err = errors.ValueIsIncorrect(ordType.Tag())
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

	execReport := fix41.CreateExecutionReportBuilder(
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

	execReport.Body.Set(clOrdID)

	if acct, err := msg.Account(); err != nil {
		execReport.Body.Set(acct)
	}

	quickfix.SendToTarget(execReport.MessageBuilder, sessionID)

	return
}

func (e *Executor) OnFIX42NewOrderSingle(msg fix42.NewOrderSingle, sessionID quickfix.SessionID) (err errors.MessageRejectError) {
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
		err = errors.ValueIsIncorrect(ordType.Tag())
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

	execReport := fix42.CreateExecutionReportBuilder(
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

	execReport.Body.Set(clOrdID)
	execReport.Body.Set(orderQty)
	execReport.Body.Set(field.NewLastShares(orderQty.Value))
	execReport.Body.Set(field.NewLastPx(price.Value))

	if acct, err := msg.Account(); err != nil {
		execReport.Body.Set(acct)
	}

	quickfix.SendToTarget(execReport.MessageBuilder, sessionID)

	return
}

func (e *Executor) OnFIX43NewOrderSingle(msg fix43.NewOrderSingle, sessionID quickfix.SessionID) (err errors.MessageRejectError) {
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
		err = errors.ValueIsIncorrect(ordType.Tag())
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

	execReport := fix43.CreateExecutionReportBuilder(
		field.NewOrderID(e.genOrderID()),
		field.NewExecID(e.genExecID()),
		field.NewExecType(enum.ExecType_FILL),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		side,
		field.NewLeavesQty(0),
		field.NewCumQty(orderQty.Value),
		field.NewAvgPx(price.Value))

	execReport.Body.Set(clOrdID)
	execReport.Body.Set(symbol)
	execReport.Body.Set(orderQty)
	execReport.Body.Set(field.NewLastShares(orderQty.Value))
	execReport.Body.Set(field.NewLastPx(price.Value))

	if acct, err := msg.Account(); err != nil {
		execReport.Body.Set(acct)
	}

	quickfix.SendToTarget(execReport.MessageBuilder, sessionID)

	return
}

func (e *Executor) OnFIX44NewOrderSingle(msg fix44.NewOrderSingle, sessionID quickfix.SessionID) (err errors.MessageRejectError) {
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
		err = errors.ValueIsIncorrect(ordType.Tag())
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

	execReport := fix44.CreateExecutionReportBuilder(
		field.NewOrderID(e.genOrderID()),
		field.NewExecID(e.genExecID()),
		field.NewExecType(enum.ExecType_FILL),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		side,
		field.NewLeavesQty(0),
		field.NewCumQty(orderQty.Value),
		field.NewAvgPx(price.Value))

	execReport.Body.Set(clOrdID)
	execReport.Body.Set(symbol)
	execReport.Body.Set(orderQty)
	execReport.Body.Set(field.NewLastQty(orderQty.Value))
	execReport.Body.Set(field.NewLastPx(price.Value))

	if acct, err := msg.Account(); err != nil {
		execReport.Body.Set(acct)
	}

	quickfix.SendToTarget(execReport.MessageBuilder, sessionID)

	return
}

func (e *Executor) OnFIX50NewOrderSingle(msg fix50.NewOrderSingle, sessionID quickfix.SessionID) (err errors.MessageRejectError) {
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
		err = errors.ValueIsIncorrect(ordType.Tag())
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

	execReport := fix50.CreateExecutionReportBuilder(
		field.NewOrderID(e.genOrderID()),
		field.NewExecID(e.genExecID()),
		field.NewExecType(enum.ExecType_FILL),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		side,
		field.NewLeavesQty(0),
		field.NewCumQty(orderQty.Value))

	execReport.Body.Set(clOrdID)
	execReport.Body.Set(symbol)
	execReport.Body.Set(orderQty)
	execReport.Body.Set(field.NewLastQty(orderQty.Value))
	execReport.Body.Set(field.NewLastPx(price.Value))
	execReport.Body.Set(field.NewAvgPx(price.Value))

	if acct, err := msg.Account(); err != nil {
		execReport.Body.Set(acct)
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

	app := &Executor{}

	acceptor, err := quickfix.NewAcceptor(app, appSettings, fileLogFactory)
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
