package main

import (
	"fmt"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/cracker"
	"github.com/quickfixgo/quickfix/errors"
	"github.com/quickfixgo/quickfix/fix"
	"github.com/quickfixgo/quickfix/fix/enum"
	"github.com/quickfixgo/quickfix/fix/field"
	"github.com/quickfixgo/quickfix/fix40"
	"github.com/quickfixgo/quickfix/fix41"
	"github.com/quickfixgo/quickfix/fix42"
	"github.com/quickfixgo/quickfix/fix43"
	"github.com/quickfixgo/quickfix/fix44"
	"github.com/quickfixgo/quickfix/fix50"
	"github.com/quickfixgo/quickfix/log"
	"github.com/quickfixgo/quickfix/message"
	"github.com/quickfixgo/quickfix/settings"
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
	var symbol field.Symbol
	if err = msg.GetSymbol(&symbol); err != nil {
		return
	}

	var side field.Side
	if err = msg.GetSide(&side); err != nil {
		return
	}

	var orderQty field.OrderQty
	if err = msg.GetOrderQty(&orderQty); err != nil {
		return
	}

	var ordType field.OrdType
	if err = msg.GetOrdType(&ordType); err != nil {
		return
	}

	if ordType.Value != enum.OrdType_LIMIT {
		err = errors.ValueIsIncorrect(ordType.Tag())
		return
	}

	var price field.Price
	if err = msg.GetPrice(&price); err != nil {
		return
	}

	var clOrdID field.ClOrdID
	if err = msg.GetClOrdID(&clOrdID); err != nil {
		return
	}

	execReport := fix40.CreateExecutionReportBuilder(
		field.BuildOrderID(e.genOrderID()),
		field.BuildExecID(e.genExecID()),
		field.BuildExecTransType(enum.ExecTransType_NEW),
		field.BuildOrdStatus(enum.OrdStatus_FILLED),
		symbol, side, orderQty, field.BuildLastShares(orderQty.Value), field.BuildLastPx(price.Value), field.BuildCumQty(orderQty.Value), field.BuildAvgPx(price.Value))

	execReport.Body.Set(clOrdID)

	if acct, err := msg.Account(); err != nil {
		execReport.Body.Set(acct)
	}

	quickfix.SendToTarget(execReport.MessageBuilder, sessionID)

	return
}

func (e *Executor) OnFIX41NewOrderSingle(msg fix41.NewOrderSingle, sessionID quickfix.SessionID) (err errors.MessageRejectError) {
	var symbol field.Symbol
	if err = msg.GetSymbol(&symbol); err != nil {
		return
	}

	var side field.Side
	if err = msg.GetSide(&side); err != nil {
		return
	}

	var orderQty field.OrderQty
	if err = msg.GetOrderQty(&orderQty); err != nil {
		return
	}

	var ordType field.OrdType
	if err = msg.GetOrdType(&ordType); err != nil {
		return
	}

	if ordType.Value != enum.OrdType_LIMIT {
		err = errors.ValueIsIncorrect(ordType.Tag())
		return
	}

	var price field.Price
	if err = msg.GetPrice(&price); err != nil {
		return
	}

	var clOrdID field.ClOrdID
	if err = msg.GetClOrdID(&clOrdID); err != nil {
		return
	}

	execReport := fix41.CreateExecutionReportBuilder(
		field.BuildOrderID(e.genOrderID()),
		field.BuildExecID(e.genExecID()),
		field.BuildExecTransType(enum.ExecTransType_NEW),
		field.BuildExecType(enum.ExecType_FILL),
		field.BuildOrdStatus(enum.OrdStatus_FILLED),
		symbol,
		side,
		orderQty,
		field.BuildLastShares(orderQty.Value),
		field.BuildLastPx(price.Value),
		field.BuildLeavesQty(0),
		field.BuildCumQty(orderQty.Value),
		field.BuildAvgPx(price.Value))

	execReport.Body.Set(clOrdID)

	if acct, err := msg.Account(); err != nil {
		execReport.Body.Set(acct)
	}

	quickfix.SendToTarget(execReport.MessageBuilder, sessionID)

	return
}

func (e *Executor) OnFIX42NewOrderSingle(msg fix42.NewOrderSingle, sessionID quickfix.SessionID) (err errors.MessageRejectError) {
	var symbol field.Symbol
	if err = msg.GetSymbol(&symbol); err != nil {
		return
	}

	var side field.Side
	if err = msg.GetSide(&side); err != nil {
		return
	}

	var orderQty field.OrderQty
	if err = msg.GetOrderQty(&orderQty); err != nil {
		return
	}

	var ordType field.OrdType
	if err = msg.GetOrdType(&ordType); err != nil {
		return
	}

	if ordType.Value != enum.OrdType_LIMIT {
		err = errors.ValueIsIncorrect(ordType.Tag())
		return
	}

	var price field.Price
	if err = msg.GetPrice(&price); err != nil {
		return
	}

	var clOrdID field.ClOrdID
	if err = msg.GetClOrdID(&clOrdID); err != nil {
		return
	}

	execReport := fix42.CreateExecutionReportBuilder(
		field.BuildOrderID(e.genOrderID()),
		field.BuildExecID(e.genExecID()),
		field.BuildExecTransType(enum.ExecTransType_NEW),
		field.BuildExecType(enum.ExecType_FILL),
		field.BuildOrdStatus(enum.OrdStatus_FILLED),
		symbol,
		side,
		field.BuildLeavesQty(0),
		field.BuildCumQty(orderQty.Value),
		field.BuildAvgPx(price.Value))

	execReport.Body.Set(clOrdID)
	execReport.Body.Set(orderQty)
	execReport.Body.Set(field.BuildLastShares(orderQty.Value))
	execReport.Body.Set(field.BuildLastPx(price.Value))

	if acct, err := msg.Account(); err != nil {
		execReport.Body.Set(acct)
	}

	quickfix.SendToTarget(execReport.MessageBuilder, sessionID)

	return
}

func (e *Executor) OnFIX43NewOrderSingle(msg fix43.NewOrderSingle, sessionID quickfix.SessionID) (err errors.MessageRejectError) {
	var symbol field.Symbol
	if err = msg.GetSymbol(&symbol); err != nil {
		return
	}

	var side field.Side
	if err = msg.GetSide(&side); err != nil {
		return
	}

	var orderQty field.OrderQty
	if err = msg.GetOrderQty(&orderQty); err != nil {
		return
	}

	var ordType field.OrdType
	if err = msg.GetOrdType(&ordType); err != nil {
		return
	}

	if ordType.Value != enum.OrdType_LIMIT {
		err = errors.ValueIsIncorrect(ordType.Tag())
		return
	}

	var price field.Price
	if err = msg.GetPrice(&price); err != nil {
		return
	}

	var clOrdID field.ClOrdID
	if err = msg.GetClOrdID(&clOrdID); err != nil {
		return
	}

	execReport := fix43.CreateExecutionReportBuilder(
		field.BuildOrderID(e.genOrderID()),
		field.BuildExecID(e.genExecID()),
		field.BuildExecType(enum.ExecType_FILL),
		field.BuildOrdStatus(enum.OrdStatus_FILLED),
		side,
		field.BuildLeavesQty(0),
		field.BuildCumQty(orderQty.Value),
		field.BuildAvgPx(price.Value))

	execReport.Body.Set(clOrdID)
	execReport.Body.Set(symbol)
	execReport.Body.Set(orderQty)
	execReport.Body.Set(field.BuildLastShares(orderQty.Value))
	execReport.Body.Set(field.BuildLastPx(price.Value))

	if acct, err := msg.Account(); err != nil {
		execReport.Body.Set(acct)
	}

	quickfix.SendToTarget(execReport.MessageBuilder, sessionID)

	return
}

func (e *Executor) OnFIX44NewOrderSingle(msg fix44.NewOrderSingle, sessionID quickfix.SessionID) (err errors.MessageRejectError) {
	var symbol field.Symbol
	if err = msg.GetSymbol(&symbol); err != nil {
		return
	}

	var side field.Side
	if err = msg.GetSide(&side); err != nil {
		return
	}

	var orderQty field.OrderQty
	if err = msg.GetOrderQty(&orderQty); err != nil {
		return
	}

	var ordType field.OrdType
	if err = msg.GetOrdType(&ordType); err != nil {
		return
	}

	if ordType.Value != enum.OrdType_LIMIT {
		err = errors.ValueIsIncorrect(ordType.Tag())
		return
	}

	var price field.Price
	if err = msg.GetPrice(&price); err != nil {
		return
	}

	var clOrdID field.ClOrdID
	if err = msg.GetClOrdID(&clOrdID); err != nil {
		return
	}

	execReport := fix44.CreateExecutionReportBuilder(
		field.BuildOrderID(e.genOrderID()),
		field.BuildExecID(e.genExecID()),
		field.BuildExecType(enum.ExecType_FILL),
		field.BuildOrdStatus(enum.OrdStatus_FILLED),
		side,
		field.BuildLeavesQty(0),
		field.BuildCumQty(orderQty.Value),
		field.BuildAvgPx(price.Value))

	execReport.Body.Set(clOrdID)
	execReport.Body.Set(symbol)
	execReport.Body.Set(orderQty)
	execReport.Body.Set(field.BuildLastQty(orderQty.Value))
	execReport.Body.Set(field.BuildLastPx(price.Value))

	if acct, err := msg.Account(); err != nil {
		execReport.Body.Set(acct)
	}

	quickfix.SendToTarget(execReport.MessageBuilder, sessionID)

	return
}

func (e *Executor) OnFIX50NewOrderSingle(msg fix50.NewOrderSingle, sessionID quickfix.SessionID) (err errors.MessageRejectError) {
	var symbol field.Symbol
	if err = msg.GetSymbol(&symbol); err != nil {
		return
	}

	var side field.Side
	if err = msg.GetSide(&side); err != nil {
		return
	}

	var orderQty field.OrderQty
	if err = msg.GetOrderQty(&orderQty); err != nil {
		return
	}

	var ordType field.OrdType
	if err = msg.GetOrdType(&ordType); err != nil {
		return
	}

	if ordType.Value != enum.OrdType_LIMIT {
		err = errors.ValueIsIncorrect(ordType.Tag())
		return
	}

	var price field.Price
	if err = msg.GetPrice(&price); err != nil {
		return
	}

	var clOrdID field.ClOrdID
	if err = msg.GetClOrdID(&clOrdID); err != nil {
		return
	}

	execReport := fix50.CreateExecutionReportBuilder(
		field.BuildOrderID(e.genOrderID()),
		field.BuildExecID(e.genExecID()),
		field.BuildExecType(enum.ExecType_FILL),
		field.BuildOrdStatus(enum.OrdStatus_FILLED),
		side,
		field.BuildLeavesQty(0),
		field.BuildCumQty(orderQty.Value))

	execReport.Body.Set(clOrdID)
	execReport.Body.Set(symbol)
	execReport.Body.Set(orderQty)
	execReport.Body.Set(field.BuildLastQty(orderQty.Value))
	execReport.Body.Set(field.BuildLastPx(price.Value))
	execReport.Body.Set(field.BuildAvgPx(price.Value))

	if acct, err := msg.Account(); err != nil {
		execReport.Body.Set(acct)
	}

	quickfix.SendToTarget(execReport.MessageBuilder, sessionID)

	return
}

func main() {
	globalSettings := settings.NewDictionary()
	globalSettings.SetInt(settings.SocketAcceptPort, 5001)
	globalSettings.SetString(settings.SenderCompID, "ISLD")
	globalSettings.SetString(settings.TargetCompID, "TW")
	globalSettings.SetBool(settings.ResetOnLogon, true)

	appSettings := settings.NewApplicationSettings(globalSettings)

	appSettings.AddSession("FIX40", settings.NewDictionary().
		SetString(settings.BeginString, fix.BeginString_FIX40))

	appSettings.AddSession("FIX41", settings.NewDictionary().
		SetString(settings.BeginString, fix.BeginString_FIX41))

	appSettings.AddSession("FIX42", settings.NewDictionary().
		SetString(settings.BeginString, fix.BeginString_FIX42))

	appSettings.AddSession("FIX43", settings.NewDictionary().
		SetString(settings.BeginString, fix.BeginString_FIX43))

	appSettings.AddSession("FIX44", settings.NewDictionary().
		SetString(settings.BeginString, fix.BeginString_FIX44))

	appSettings.AddSession("FIX50", settings.NewDictionary().
		SetString(settings.BeginString, fix.BeginString_FIXT11).
		SetString(settings.DefaultApplVerID, enum.ApplVerID_FIX50))

	app := new(Executor)

	acceptor, err := quickfix.NewAcceptor(app, appSettings, log.ScreenLogFactory{})
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
