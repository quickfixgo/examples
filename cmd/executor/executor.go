// Copyright (c) quickfixengine.org  All rights reserved.
//
// This file may be distributed under the terms of the quickfixengine.org
// license as defined by quickfixengine.org and appearing in the file
// LICENSE included in the packaging of this file.
//
// This file is provided AS IS with NO WARRANTY OF ANY KIND, INCLUDING
// THE WARRANTY OF DESIGN, MERCHANTABILITY AND FITNESS FOR A
// PARTICULAR PURPOSE.
//
// See http://www.quickfixengine.org/LICENSE for licensing information.
//
// Contact ask@quickfixengine.org if any conditions of this licensing
// are not clear to you.

package executor

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"syscall"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/examples/cmd/utils"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/tag"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"

	"github.com/quickfixgo/quickfix"

	fix40nos "github.com/quickfixgo/fix40/newordersingle"
	fix41nos "github.com/quickfixgo/fix41/newordersingle"
	fix42nos "github.com/quickfixgo/fix42/newordersingle"
	fix43nos "github.com/quickfixgo/fix43/newordersingle"
	fix44nos "github.com/quickfixgo/fix44/newordersingle"
	fix50nos "github.com/quickfixgo/fix50/newordersingle"

	fix40er "github.com/quickfixgo/fix40/executionreport"
	fix41er "github.com/quickfixgo/fix41/executionreport"
	fix42er "github.com/quickfixgo/fix42/executionreport"
	fix43er "github.com/quickfixgo/fix43/executionreport"
	fix44er "github.com/quickfixgo/fix44/executionreport"
	fix50er "github.com/quickfixgo/fix50/executionreport"

	"os"
	"os/signal"
	"strconv"
)

type executor struct {
	orderID int
	execID  int
	*quickfix.MessageRouter
}

func newExecutor() *executor {
	e := &executor{MessageRouter: quickfix.NewMessageRouter()}
	e.AddRoute(fix40nos.Route(e.OnFIX40NewOrderSingle))
	e.AddRoute(fix41nos.Route(e.OnFIX41NewOrderSingle))
	e.AddRoute(fix42nos.Route(e.OnFIX42NewOrderSingle))
	e.AddRoute(fix43nos.Route(e.OnFIX43NewOrderSingle))
	e.AddRoute(fix44nos.Route(e.OnFIX44NewOrderSingle))
	e.AddRoute(fix50nos.Route(e.OnFIX50NewOrderSingle))

	return e
}

func (e *executor) genOrderID() field.OrderIDField {
	e.orderID++
	return field.NewOrderID(strconv.Itoa(e.orderID))
}

func (e *executor) genExecID() field.ExecIDField {
	e.execID++
	return field.NewExecID(strconv.Itoa(e.execID))
}

// quickfix.Application interface
func (e executor) OnCreate(sessionID quickfix.SessionID)                           {}
func (e executor) OnLogon(sessionID quickfix.SessionID)                            {}
func (e executor) OnLogout(sessionID quickfix.SessionID)                           {}
func (e executor) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID)     {}
func (e executor) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) error { return nil }
func (e executor) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

// Use Message Cracker on Incoming Application Messages
func (e *executor) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return e.Route(msg, sessionID)
}

func (e *executor) OnFIX40NewOrderSingle(msg fix40nos.NewOrderSingle, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	ordType, err := msg.GetOrdType()
	if err != nil {
		return err
	}

	if ordType != enum.OrdType_LIMIT {
		utils.PrintBad("incoming order was not a limit order and was rejected")
		return quickfix.ValueIsIncorrect(tag.OrdType)
	}

	symbol, err := msg.GetSymbol()
	if err != nil {
		return err
	}

	side, err := msg.GetSide()
	if err != nil {
		return err
	}

	orderQty, err := msg.GetOrderQty()
	if err != nil {
		return err
	}

	price, err := msg.GetPrice()
	if err != nil {
		return err
	}

	execReport := fix40er.New(
		e.genOrderID(),
		e.genExecID(),
		field.NewExecTransType(enum.ExecTransType_NEW),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		field.NewSymbol(symbol),
		field.NewSide(side),
		field.NewOrderQty(orderQty, 2),
		field.NewLastShares(orderQty, 2),
		field.NewLastPx(price, 2),
		field.NewCumQty(orderQty, 2),
		field.NewAvgPx(price, 2),
	)

	clOrdID, err := msg.GetClOrdID()
	if err != nil {
		return err
	}
	execReport.SetClOrdID(clOrdID)

	sendErr := quickfix.SendToTarget(execReport, sessionID)
	if sendErr != nil {
		utils.PrintBad(sendErr.Error())
	}

	return nil
}

func (e *executor) OnFIX41NewOrderSingle(msg fix41nos.NewOrderSingle, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	ordType, err := msg.GetOrdType()
	if err != nil {
		return
	}
	if ordType != enum.OrdType_LIMIT {
		utils.PrintBad("incoming order was not a limit order and was rejected")
		return quickfix.ValueIsIncorrect(tag.OrdType)
	}

	symbol, err := msg.GetSymbol()
	if err != nil {
		return
	}

	side, err := msg.GetSide()
	if err != nil {
		return
	}

	orderQty, err := msg.GetOrderQty()
	if err != nil {
		return
	}

	price, err := msg.GetPrice()
	if err != nil {
		return
	}

	execReport := fix41er.New(
		e.genOrderID(),
		e.genExecID(),
		field.NewExecTransType(enum.ExecTransType_NEW),
		field.NewExecType(enum.ExecType_FILL),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		field.NewSymbol(symbol),
		field.NewSide(side),
		field.NewOrderQty(orderQty, 2),
		field.NewLastShares(orderQty, 2),
		field.NewLastPx(price, 2),
		field.NewLeavesQty(decimal.Zero, 2),
		field.NewCumQty(orderQty, 2),
		field.NewAvgPx(price, 2),
	)

	clOrdID, err := msg.GetClOrdID()
	if err != nil {
		return
	}
	execReport.SetClOrdID(clOrdID)

	sendErr := quickfix.SendToTarget(execReport, sessionID)
	if sendErr != nil {
		utils.PrintBad(sendErr.Error())
	}
	return
}

func (e *executor) OnFIX42NewOrderSingle(msg fix42nos.NewOrderSingle, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	ordType, err := msg.GetOrdType()
	if err != nil {
		return err
	}

	if ordType != enum.OrdType_LIMIT {
		utils.PrintBad("incoming order was not a limit order and was rejected")
		return quickfix.ValueIsIncorrect(tag.OrdType)
	}

	symbol, err := msg.GetSymbol()
	if err != nil {
		return
	}

	side, err := msg.GetSide()
	if err != nil {
		return
	}

	orderQty, err := msg.GetOrderQty()
	if err != nil {
		return
	}

	price, err := msg.GetPrice()
	if err != nil {
		return
	}

	clOrdID, err := msg.GetClOrdID()
	if err != nil {
		return
	}

	execReport := fix42er.New(
		e.genOrderID(),
		e.genExecID(),
		field.NewExecTransType(enum.ExecTransType_NEW),
		field.NewExecType(enum.ExecType_FILL),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		field.NewSymbol(symbol),
		field.NewSide(side),
		field.NewLeavesQty(decimal.Zero, 2),
		field.NewCumQty(orderQty, 2),
		field.NewAvgPx(price, 2),
	)

	execReport.SetClOrdID(clOrdID)
	execReport.SetOrderQty(orderQty, 2)
	execReport.SetLastShares(orderQty, 2)
	execReport.SetLastPx(price, 2)

	if msg.HasAccount() {
		acct, err := msg.GetAccount()
		if err != nil {
			return err
		}
		execReport.SetAccount(acct)
	}

	sendErr := quickfix.SendToTarget(execReport, sessionID)
	if sendErr != nil {
		utils.PrintBad(sendErr.Error())
	}

	return
}

func (e *executor) OnFIX43NewOrderSingle(msg fix43nos.NewOrderSingle, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	ordType, err := msg.GetOrdType()
	if err != nil {
		return err
	}
	if ordType != enum.OrdType_LIMIT {
		utils.PrintBad("incoming order was not a limit order and was rejected")
		return quickfix.ValueIsIncorrect(tag.OrdType)
	}

	symbol, err := msg.GetSymbol()
	if err != nil {
		return
	}

	side, err := msg.GetSide()
	if err != nil {
		return
	}

	orderQty, err := msg.GetOrderQty()
	if err != nil {
		return
	}

	price, err := msg.GetPrice()
	if err != nil {
		return
	}

	clOrdID, err := msg.GetClOrdID()
	if err != nil {
		return
	}

	execReport := fix43er.New(
		e.genOrderID(),
		e.genExecID(),
		field.NewExecType(enum.ExecType_FILL),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		field.NewSide(side),
		field.NewLeavesQty(decimal.Zero, 2),
		field.NewCumQty(orderQty, 2),
		field.NewAvgPx(price, 2),
	)

	execReport.SetClOrdID(clOrdID)
	execReport.SetSymbol(symbol)
	execReport.SetOrderQty(orderQty, 2)
	execReport.SetLastQty(orderQty, 2)
	execReport.SetLastPx(price, 2)

	if msg.HasAccount() {
		acct, err := msg.GetAccount()
		if err != nil {
			return err
		}
		execReport.SetAccount(acct)
	}

	sendErr := quickfix.SendToTarget(execReport, sessionID)
	if sendErr != nil {
		utils.PrintBad(sendErr.Error())
	}

	return
}

func (e *executor) OnFIX44NewOrderSingle(msg fix44nos.NewOrderSingle, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	ordType, err := msg.GetOrdType()
	if err != nil {
		return err
	}

	if ordType != enum.OrdType_LIMIT {
		utils.PrintBad("incoming order was not a limit order and was rejected")
		return quickfix.ValueIsIncorrect(tag.OrdType)
	}

	symbol, err := msg.GetSymbol()
	if err != nil {
		return
	}

	side, err := msg.GetSide()
	if err != nil {
		return
	}

	orderQty, err := msg.GetOrderQty()
	if err != nil {
		return
	}

	price, err := msg.GetPrice()
	if err != nil {
		return
	}

	clOrdID, err := msg.GetClOrdID()
	if err != nil {
		return
	}

	execReport := fix44er.New(
		e.genOrderID(),
		e.genExecID(),
		field.NewExecType(enum.ExecType_FILL),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		field.NewSide(side),
		field.NewLeavesQty(decimal.Zero, 2),
		field.NewCumQty(orderQty, 2),
		field.NewAvgPx(price, 2),
	)

	execReport.SetClOrdID(clOrdID)
	execReport.SetSymbol(symbol)
	execReport.SetOrderQty(orderQty, 2)
	execReport.SetLastQty(orderQty, 2)
	execReport.SetLastPx(price, 2)

	if msg.HasAccount() {
		acct, err := msg.GetAccount()
		if err != nil {
			return err
		}
		execReport.SetAccount(acct)
	}

	sendErr := quickfix.SendToTarget(execReport, sessionID)
	if sendErr != nil {
		utils.PrintBad(sendErr.Error())
	}

	return
}

func (e *executor) OnFIX50NewOrderSingle(msg fix50nos.NewOrderSingle, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	ordType, err := msg.GetOrdType()
	if err != nil {
		return err
	}

	if ordType != enum.OrdType_LIMIT {
		utils.PrintBad("incoming order was not a limit order and was rejected")
		return quickfix.ValueIsIncorrect(tag.OrdType)
	}

	symbol, err := msg.GetSymbol()
	if err != nil {
		return
	}

	side, err := msg.GetSide()
	if err != nil {
		return
	}

	orderQty, err := msg.GetOrderQty()
	if err != nil {
		return
	}

	price, err := msg.GetPrice()
	if err != nil {
		return
	}

	clOrdID, err := msg.GetClOrdID()
	if err != nil {
		return
	}

	execReport := fix50er.New(
		e.genOrderID(),
		e.genExecID(),
		field.NewExecType(enum.ExecType_FILL),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		field.NewSide(side),
		field.NewLeavesQty(decimal.Zero, 2),
		field.NewCumQty(orderQty, 2),
	)

	execReport.SetClOrdID(clOrdID)
	execReport.SetSymbol(symbol)
	execReport.SetOrderQty(orderQty, 2)
	execReport.SetLastQty(orderQty, 2)
	execReport.SetLastPx(price, 2)
	execReport.SetAvgPx(price, 2)

	if msg.HasAccount() {
		acct, err := msg.GetAccount()
		if err != nil {
			return err
		}
		execReport.SetAccount(acct)
	}

	sendErr := quickfix.SendToTarget(execReport, sessionID)
	if sendErr != nil {
		utils.PrintBad(sendErr.Error())
	}

	return
}

const (
	usage = "executor"
	short = "Start an order execution (FIX acceptor) service"
	long  = "Start an order execution (FIX acceptor) service."
)

var (
	// Cmd is the executor command.
	Cmd = &cobra.Command{
		Use:     usage,
		Short:   short,
		Long:    long,
		Aliases: []string{"x"},
		Example: "qf executor [YOUR_FIX_CONFIG_FILE_HERE.cfg] (default is ./config/executor.cfg)",
		RunE:    execute,
	}
)

func execute(cmd *cobra.Command, args []string) error {
	var cfgFileName string
	argLen := len(args)
	switch argLen {
	case 0:
		{
			utils.PrintInfo("FIX config file not provided...")
			utils.PrintInfo("attempting to use default location './config/executor.cfg' ...")
			cfgFileName = path.Join("config", "executor.cfg")
		}
	case 1:
		{
			cfgFileName = args[0]
		}
	default:
		{
			return fmt.Errorf("incorrect argument number")
		}
	}
	cfg, err := os.Open(cfgFileName)
	if err != nil {
		return fmt.Errorf("error opening %v, %v", cfgFileName, err)
	}
	defer cfg.Close()
	stringData, readErr := io.ReadAll(cfg)
	if readErr != nil {
		return fmt.Errorf("error reading cfg: %s,", readErr)
	}

	appSettings, err := quickfix.ParseSettings(bytes.NewReader(stringData))
	if err != nil {
		return fmt.Errorf("error reading cfg: %s,", err)
	}

	logger := utils.NewFancyLog()
	app := newExecutor()

	utils.PrintConfig("acceptor", bytes.NewReader(stringData))
	acceptor, err := quickfix.NewAcceptor(app, quickfix.NewMemoryStoreFactory(), appSettings, logger)
	if err != nil {
		return fmt.Errorf("unable to create acceptor: %s", err)
	}

	err = acceptor.Start()
	if err != nil {
		return fmt.Errorf("unable to start FIX acceptor: %s", err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt

	utils.PrintInfo("stopping FIX acceptor service..")
	acceptor.Stop()
	utils.PrintInfo("stopped")

	return nil
}
