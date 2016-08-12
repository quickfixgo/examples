package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/quickfixgo/examples/cmd/tradeclient/internal"
	"github.com/quickfixgo/quickfix"
)

//TradeClient implements the quickfix.Application interface
type TradeClient struct {
}

//OnCreate implemented as part of Application interface
func (e TradeClient) OnCreate(sessionID quickfix.SessionID) {
	return
}

//OnLogon implemented as part of Application interface
func (e TradeClient) OnLogon(sessionID quickfix.SessionID) {
	return
}

//OnLogout implemented as part of Application interface
func (e TradeClient) OnLogout(sessionID quickfix.SessionID) {
	return
}

//FromAdmin implemented as part of Application interface
func (e TradeClient) FromAdmin(msg quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return
}

//ToAdmin implemented as part of Application interface
func (e TradeClient) ToAdmin(msg quickfix.Message, sessionID quickfix.SessionID) {
	return
}

//ToApp implemented as part of Application interface
func (e TradeClient) ToApp(msg quickfix.Message, sessionID quickfix.SessionID) (err error) {
	msg.Build()
	fmt.Printf("Sending %s\n", &msg)
	return
}

//FromApp implemented as part of Application interface. This is the callback for all Application level messages from the counter party.
func (e TradeClient) FromApp(msg quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	fmt.Printf("FromApp: %s\n", msg.String())
	return
}

func main() {
	flag.Parse()

	cfgFileName := "tradeclient.cfg"
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

	app := TradeClient{}
	fileLogFactory, err := quickfix.NewFileLogFactory(appSettings)

	if err != nil {
		fmt.Println("Error creating file log factory,", err)
		return
	}

	initiator, err := quickfix.NewInitiator(app, quickfix.NewMemoryStoreFactory(), appSettings, fileLogFactory)
	if err != nil {
		fmt.Printf("Unable to create Initiator: %s\n", err)
		return
	}

	initiator.Start()

Loop:
	for {
		action, err := internal.QueryAction()
		if err != nil {
			break
		}

		switch action {
		case "1":
			err = internal.QueryEnterOrder()

			/*		case "2":
						err = queryCancelOrder()

					case "3":
						err = queryMarketDataRequest()
			*/

		case "4":
			//quit
			break Loop

		default:
			err = fmt.Errorf("unknown action: '%v'", action)
		}

		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}

	initiator.Stop()
}
