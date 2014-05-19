package main

import (
	"fmt"
	"github.com/quickfixgo/quickfix"
	"os"
)

type TradeClient struct {
}

func (e TradeClient) OnCreate(sessionID quickfix.SessionID) {
	return
}

func (e TradeClient) OnLogon(sessionID quickfix.SessionID) {
	return
}

func (e TradeClient) OnLogout(sessionID quickfix.SessionID) {
	return
}

func (e TradeClient) FromAdmin(msg quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return
}

func (e TradeClient) ToAdmin(msg quickfix.MessageBuilder, sessionID quickfix.SessionID) {
	return
}

func (e TradeClient) ToApp(msg quickfix.MessageBuilder, sessionID quickfix.SessionID) (err error) {
	return
}

func (e TradeClient) FromApp(msg quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	fmt.Println("FromApp: ", msg)
	return
}

func main() {
	cfgFileName := "tradeclient.cfg"
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

	initiator, err := quickfix.NewInitiator(app, appSettings, fileLogFactory)
	if err != nil {
		fmt.Printf("Unable to create Initiator: %s\n", err)
		return
	}

	initiator.Start()

	for {
		action, err := queryAction()
		if err != nil {
			break
		}

		switch action {
		case "1":
			err = queryEnterOrder()

		default:
			err = fmt.Errorf("unknown action: '%v'", action)
		}

		if err != nil {
			fmt.Printf("%v\n", err)
			break
		}
	}

	initiator.Stop()
}
