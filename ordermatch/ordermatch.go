package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/enum"
	"github.com/quickfixgo/quickfix/fix42/executionreport"
	"github.com/quickfixgo/quickfix/fix42/marketdatarequest"
	"github.com/quickfixgo/quickfix/fix42/newordersingle"
	"github.com/quickfixgo/quickfix/fix42/ordercancelrequest"
	"github.com/quickfixgo/quickfix/tag"
	"os"
	"os/signal"
	"strconv"
)

//Application implements the quickfix.Application interface
type Application struct {
	*quickfix.MessageRouter
	*OrderMatcher
	execID int
}

func newApplication() *Application {
	app := &Application{
		MessageRouter: quickfix.NewMessageRouter(),
		OrderMatcher:  NewOrderMatcher(),
	}
	app.AddRoute(newordersingle.Route(app.onNewOrderSingle))
	app.AddRoute(ordercancelrequest.Route(app.onOrderCancelRequest))
	app.AddRoute(marketdatarequest.Route(app.onMarketDataRequest))

	return app
}

//OnCreate implemented as part of Application interface
func (a Application) OnCreate(sessionID quickfix.SessionID) { return }

//OnLogon implemented as part of Application interface
func (a Application) OnLogon(sessionID quickfix.SessionID) { return }

//OnLogout implemented as part of Application interface
func (a Application) OnLogout(sessionID quickfix.SessionID) { return }

//ToAdmin implemented as part of Application interface
func (a Application) ToAdmin(msg quickfix.Message, sessionID quickfix.SessionID) { return }

//ToApp implemented as part of Application interface
func (a Application) ToApp(msg quickfix.Message, sessionID quickfix.SessionID) error {
	return nil
}

//FromAdmin implemented as part of Application interface
func (a Application) FromAdmin(msg quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

//FromApp implemented as part of Application interface, uses Router on incoming application messages
func (a *Application) FromApp(msg quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return a.Route(msg, sessionID)
}

func (a *Application) onNewOrderSingle(msg newordersingle.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	if msg.Price == nil {
		err = quickfix.ConditionallyRequiredFieldMissing(tag.Price)
		return
	}
	if msg.OrderQty == nil {
		err = quickfix.ConditionallyRequiredFieldMissing(tag.OrderQty)
		return
	}

	order := Order{
		ClOrdID:      msg.ClOrdID,
		Symbol:       msg.Symbol,
		SenderCompID: msg.Header.SenderCompID,
		TargetCompID: msg.Header.TargetCompID,
		Side:         msg.Side,
		OrdType:      msg.OrdType,
		Price:        *msg.Price,
		Quantity:     *msg.OrderQty,
	}

	a.Insert(order)
	a.acceptOrder(order)

	matches := a.Match(order.Symbol)

	for len(matches) > 0 {
		a.fillOrder(matches[0])
		matches = matches[1:]
	}

	return
}

func (a *Application) onOrderCancelRequest(msg ordercancelrequest.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	order := a.Cancel(msg.OrigClOrdID, msg.Symbol, msg.Side)
	if order != nil {
		a.cancelOrder(*order)
	}

	return
}

func (a *Application) onMarketDataRequest(msg marketdatarequest.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	fmt.Printf("%+v\n", msg)
	return
}

func (a *Application) acceptOrder(order Order) {
	a.updateOrder(order, enum.OrdStatus_NEW)
}

func (a *Application) fillOrder(order Order) {
	status := enum.OrdStatus_FILLED
	if !order.IsClosed() {
		status = enum.OrdStatus_PARTIALLY_FILLED
	}
	a.updateOrder(order, status)
}

func (a *Application) cancelOrder(order Order) {
	a.updateOrder(order, enum.OrdStatus_CANCELED)
}

func (a *Application) genExecID() string {
	a.execID++
	return strconv.Itoa(a.execID)
}

func (a *Application) updateOrder(order Order, status string) {
	execReport := executionreport.Message{
		OrderID:       order.ClOrdID,
		ClOrdID:       &order.ClOrdID,
		ExecID:        a.genExecID(),
		ExecTransType: enum.ExecTransType_NEW,
		ExecType:      status,
		OrdStatus:     status,
		Symbol:        order.Symbol,
		Side:          order.Side,
		OrderQty:      &order.Quantity,
		LeavesQty:     order.OpenQuantity(),
		CumQty:        order.ExecutedQuantity,
		AvgPx:         order.AvgPx,
	}

	execReport.Header.TargetCompID = order.SenderCompID
	execReport.Header.SenderCompID = order.TargetCompID

	quickfix.Send(execReport)
}

func main() {
	flag.Parse()

	cfgFileName := "ordermatch.cfg"
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
	app := newApplication()

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
	go func() {
		<-interrupt
		acceptor.Stop()
		os.Exit(0)
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()

		switch value := scanner.Text(); value {
		case "#symbols":
			app.Display()
		default:
			app.DisplayMarket(value)
		}
	}
}
