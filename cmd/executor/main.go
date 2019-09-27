package main

import (
	"flag"
	"fmt"
	"path"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
	"github.com/shopspring/decimal"

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

//quickfix.Application interface
func (e executor) OnCreate(sessionID quickfix.SessionID)                           { return }
func (e executor) OnLogon(sessionID quickfix.SessionID)                            { return }
func (e executor) OnLogout(sessionID quickfix.SessionID)                           { return }
func (e executor) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID)     { return }
func (e executor) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) error { return nil }
func (e executor) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

//Use Message Cracker on Incoming Application Messages
func (e *executor) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return e.Route(msg, sessionID)
}

func (e *executor) OnFIX40NewOrderSingle(msg fix40nos.NewOrderSingle, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	ordType, err := msg.GetOrdType()
	if err != nil {
		return err
	}

	if ordType != enum.OrdType_LIMIT {
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

	quickfix.SendToTarget(execReport, sessionID)

	return nil
}

func (e *executor) OnFIX41NewOrderSingle(msg fix41nos.NewOrderSingle, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	ordType, err := msg.GetOrdType()
	if err != nil {
		return
	}
	if ordType != enum.OrdType_LIMIT {
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

	quickfix.SendToTarget(execReport, sessionID)

	return
}

func (e *executor) OnFIX42NewOrderSingle(msg fix42nos.NewOrderSingle, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	ordType, err := msg.GetOrdType()
	if err != nil {
		return err
	}

	if ordType != enum.OrdType_LIMIT {
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

	quickfix.SendToTarget(execReport, sessionID)

	return
}

func (e *executor) OnFIX43NewOrderSingle(msg fix43nos.NewOrderSingle, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	ordType, err := msg.GetOrdType()
	if err != nil {
		return err
	}
	if ordType != enum.OrdType_LIMIT {
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

	quickfix.SendToTarget(execReport, sessionID)

	return
}

func (e *executor) OnFIX44NewOrderSingle(msg fix44nos.NewOrderSingle, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	ordType, err := msg.GetOrdType()
	if err != nil {
		return err
	}

	if ordType != enum.OrdType_LIMIT {
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

	quickfix.SendToTarget(execReport, sessionID)

	return
}

func (e *executor) OnFIX50NewOrderSingle(msg fix50nos.NewOrderSingle, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	ordType, err := msg.GetOrdType()
	if err != nil {
		return err
	}

	if ordType != enum.OrdType_LIMIT {
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

	quickfix.SendToTarget(execReport, sessionID)

	return
}

func main() {
	flag.Parse()

	cfgFileName := path.Join("config", "executor.cfg")
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
	app := newExecutor()

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
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	<-interrupt

	acceptor.Stop()
}
