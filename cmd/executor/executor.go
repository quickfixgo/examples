package main

import (
	"flag"
	"fmt"
	"path"

	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/enum"
	"github.com/quickfixgo/quickfix/field"
	"github.com/quickfixgo/quickfix/tag"
	"github.com/shopspring/decimal"

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
func (e executor) OnCreate(sessionID quickfix.SessionID)                          { return }
func (e executor) OnLogon(sessionID quickfix.SessionID)                           { return }
func (e executor) OnLogout(sessionID quickfix.SessionID)                          { return }
func (e executor) ToAdmin(msg quickfix.Message, sessionID quickfix.SessionID)     { return }
func (e executor) ToApp(msg quickfix.Message, sessionID quickfix.SessionID) error { return nil }
func (e executor) FromAdmin(msg quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

//Use Message Cracker on Incoming Application Messages
func (e *executor) FromApp(msg quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return e.Route(msg, sessionID)
}

func (e *executor) OnFIX40NewOrderSingle(msg fix40nos.NewOrderSingle, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	var ordType field.OrdTypeField
	if ordType, err = msg.GetOrdType(); err != nil {
		return err
	}
	if ordType.String() != enum.OrdType_LIMIT {
		return quickfix.ValueIsIncorrect(tag.OrdType)
	}

	var symbol field.SymbolField
	if symbol, err = msg.GetSymbol(); err != nil {
		return
	}

	var side field.SideField
	if side, err = msg.GetSide(); err != nil {
		return
	}

	var orderQty field.OrderQtyField
	if orderQty, err = msg.GetOrderQty(); err != nil {
		return
	}

	var price field.PriceField
	if price, err = msg.GetPrice(); err != nil {
		return
	}

	execReport := fix40er.New(
		e.genOrderID(),
		e.genExecID(),
		field.NewExecTransType(enum.ExecTransType_NEW),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		symbol,
		side,
		orderQty,
		field.NewLastShares(orderQty.Decimal, 2),
		field.NewLastPx(price.Decimal, 2),
		field.NewCumQty(orderQty.Decimal, 2),
		field.NewAvgPx(price.Decimal, 2),
	)

	var clOrdID field.ClOrdIDField
	if clOrdID, err = msg.GetClOrdID(); err != nil {
		return
	}
	execReport.Set(clOrdID)

	quickfix.SendToTarget(execReport, sessionID)

	return
}

func (e *executor) OnFIX41NewOrderSingle(msg fix41nos.NewOrderSingle, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	var ordType field.OrdTypeField
	if ordType, err = msg.GetOrdType(); err != nil {
		return err
	}
	if ordType.String() != enum.OrdType_LIMIT {
		return quickfix.ValueIsIncorrect(tag.OrdType)
	}

	var symbol field.SymbolField
	if symbol, err = msg.GetSymbol(); err != nil {
		return
	}

	var side field.SideField
	if side, err = msg.GetSide(); err != nil {
		return
	}

	var orderQty field.OrderQtyField
	if orderQty, err = msg.GetOrderQty(); err != nil {
		return
	}

	var price field.PriceField
	if price, err = msg.GetPrice(); err != nil {
		return
	}

	execReport := fix41er.New(
		e.genOrderID(),
		e.genExecID(),
		field.NewExecTransType(enum.ExecTransType_NEW),
		field.NewExecType(enum.ExecType_FILL),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		symbol,
		side,
		orderQty,
		field.NewLastShares(orderQty.Decimal, 2),
		field.NewLastPx(price.Decimal, 2),
		field.NewLeavesQty(decimal.Zero, 2),
		field.NewCumQty(orderQty.Decimal, 2),
		field.NewAvgPx(price.Decimal, 2),
	)

	var clOrdID field.ClOrdIDField
	if clOrdID, err = msg.GetClOrdID(); err != nil {
		return
	}
	execReport.Set(clOrdID)

	quickfix.SendToTarget(execReport, sessionID)

	return
}

func (e *executor) OnFIX42NewOrderSingle(msg fix42nos.NewOrderSingle, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	var ordType field.OrdTypeField
	if ordType, err = msg.GetOrdType(); err != nil {
		return err
	}
	if ordType.String() != enum.OrdType_LIMIT {
		return quickfix.ValueIsIncorrect(tag.OrdType)
	}

	var symbol field.SymbolField
	if symbol, err = msg.GetSymbol(); err != nil {
		return
	}

	var side field.SideField
	if side, err = msg.GetSide(); err != nil {
		return
	}

	var orderQty field.OrderQtyField
	if orderQty, err = msg.GetOrderQty(); err != nil {
		return
	}

	var price field.PriceField
	if price, err = msg.GetPrice(); err != nil {
		return
	}

	var clOrdID field.ClOrdIDField
	if clOrdID, err = msg.GetClOrdID(); err != nil {
		return
	}

	execReport := fix42er.New(
		e.genOrderID(),
		e.genExecID(),
		field.NewExecTransType(enum.ExecTransType_NEW),
		field.NewExecType(enum.ExecType_FILL),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		symbol,
		side,
		field.NewLeavesQty(decimal.Zero, 2),
		field.NewCumQty(orderQty.Decimal, 2),
		field.NewAvgPx(price.Decimal, 2),
	)

	execReport.Set(clOrdID)
	execReport.Set(orderQty)
	execReport.SetLastShares(orderQty.Decimal, 2)
	execReport.SetLastPx(price.Decimal, 2)

	if msg.HasAccount() {
		var acct field.AccountField
		if acct, err = msg.GetAccount(); err != nil {
			return err
		}
		execReport.Set(acct)
	}

	quickfix.SendToTarget(execReport, sessionID)

	return
}

func (e *executor) OnFIX43NewOrderSingle(msg fix43nos.NewOrderSingle, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	var ordType field.OrdTypeField
	if ordType, err = msg.GetOrdType(); err != nil {
		return err
	}
	if ordType.String() != enum.OrdType_LIMIT {
		return quickfix.ValueIsIncorrect(tag.OrdType)
	}

	var symbol field.SymbolField
	if symbol, err = msg.GetSymbol(); err != nil {
		return
	}

	var side field.SideField
	if side, err = msg.GetSide(); err != nil {
		return
	}

	var orderQty field.OrderQtyField
	if orderQty, err = msg.GetOrderQty(); err != nil {
		return
	}

	var price field.PriceField
	if price, err = msg.GetPrice(); err != nil {
		return
	}

	var clOrdID field.ClOrdIDField
	if clOrdID, err = msg.GetClOrdID(); err != nil {
		return
	}

	execReport := fix43er.New(
		e.genOrderID(),
		e.genExecID(),
		field.NewExecType(enum.ExecType_FILL),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		side,
		field.NewLeavesQty(decimal.Zero, 2),
		field.NewCumQty(orderQty.Decimal, 2),
		field.NewAvgPx(price.Decimal, 2),
	)

	execReport.Set(clOrdID)
	execReport.Set(symbol)
	execReport.Set(orderQty)
	execReport.SetLastQty(orderQty.Decimal, 2)
	execReport.SetLastPx(price.Decimal, 2)

	if msg.HasAccount() {
		var acct field.AccountField
		if acct, err = msg.GetAccount(); err != nil {
			return err
		}
		execReport.Set(acct)
	}

	quickfix.SendToTarget(execReport, sessionID)

	return
}

func (e *executor) OnFIX44NewOrderSingle(msg fix44nos.NewOrderSingle, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	var ordType field.OrdTypeField
	if ordType, err = msg.GetOrdType(); err != nil {
		return err
	}

	if ordType.String() != enum.OrdType_LIMIT {
		return quickfix.ValueIsIncorrect(tag.OrdType)
	}

	var symbol field.SymbolField
	if symbol, err = msg.GetSymbol(); err != nil {
		return
	}

	var side field.SideField
	if side, err = msg.GetSide(); err != nil {
		return
	}

	var orderQty field.OrderQtyField
	if orderQty, err = msg.GetOrderQty(); err != nil {
		return
	}

	var price field.PriceField
	if price, err = msg.GetPrice(); err != nil {
		return
	}

	var clOrdID field.ClOrdIDField
	if clOrdID, err = msg.GetClOrdID(); err != nil {
		return
	}

	execReport := fix44er.New(
		e.genOrderID(),
		e.genExecID(),
		field.NewExecType(enum.ExecType_FILL),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		side,
		field.NewLeavesQty(decimal.Zero, 2),
		field.NewCumQty(orderQty.Decimal, 2),
		field.NewAvgPx(price.Decimal, 2),
	)

	execReport.Set(clOrdID)
	execReport.Set(symbol)
	execReport.Set(orderQty)
	execReport.SetLastQty(orderQty.Decimal, 2)
	execReport.SetLastPx(price.Decimal, 2)

	if msg.HasAccount() {
		var acct field.AccountField
		if acct, err = msg.GetAccount(); err != nil {
			return err
		}
		execReport.Set(acct)
	}

	quickfix.SendToTarget(execReport, sessionID)

	return
}

func (e *executor) OnFIX50NewOrderSingle(msg fix50nos.NewOrderSingle, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	var ordType field.OrdTypeField
	if ordType, err = msg.GetOrdType(); err != nil {
		return err
	}

	if ordType.String() != enum.OrdType_LIMIT {
		return quickfix.ValueIsIncorrect(tag.OrdType)
	}

	var symbol field.SymbolField
	if symbol, err = msg.GetSymbol(); err != nil {
		return
	}

	var side field.SideField
	if side, err = msg.GetSide(); err != nil {
		return
	}

	var orderQty field.OrderQtyField
	if orderQty, err = msg.GetOrderQty(); err != nil {
		return
	}

	var price field.PriceField
	if price, err = msg.GetPrice(); err != nil {
		return
	}

	var clOrdID field.ClOrdIDField
	if clOrdID, err = msg.GetClOrdID(); err != nil {
		return
	}

	execReport := fix50er.New(
		e.genOrderID(),
		e.genExecID(),
		field.NewExecType(enum.ExecType_FILL),
		field.NewOrdStatus(enum.OrdStatus_FILLED),
		side,
		field.NewLeavesQty(decimal.Zero, 2),
		field.NewCumQty(orderQty.Decimal, 2),
	)

	execReport.Set(clOrdID)
	execReport.Set(symbol)
	execReport.Set(orderQty)
	execReport.SetLastQty(orderQty.Decimal, 2)
	execReport.SetLastPx(price.Decimal, 2)
	execReport.SetAvgPx(price.Decimal, 2)

	if msg.HasAccount() {
		var acct field.AccountField
		if acct, err = msg.GetAccount(); err != nil {
			return err
		}
		execReport.Set(acct)
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
