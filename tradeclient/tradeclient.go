package main

import (
	"fmt"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/errors"
	"github.com/quickfixgo/quickfix/fix"
	"github.com/quickfixgo/quickfix/fix/enum"
	"github.com/quickfixgo/quickfix/log"
	"github.com/quickfixgo/quickfix/message"
	"github.com/quickfixgo/quickfix/settings"
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

func (e TradeClient) FromAdmin(msg message.Message, sessionID quickfix.SessionID) (reject errors.MessageRejectError) {
	return
}

func (e TradeClient) ToAdmin(msg message.MessageBuilder, sessionID quickfix.SessionID) {
	return
}

func (e TradeClient) ToApp(msg message.MessageBuilder, sessionID quickfix.SessionID) (err error) {
	return
}

func (e TradeClient) FromApp(msg message.Message, sessionID quickfix.SessionID) (reject errors.MessageRejectError) {
	return
}

func main() {
	globalSettings := settings.NewDictionary()
	globalSettings.SetString(settings.SocketConnectHost, "127.0.0.1")
	globalSettings.SetInt(settings.SocketConnectPort, 5001)
	globalSettings.SetInt(settings.HeartBtInt, 30)
	globalSettings.SetString(settings.SenderCompID, "TW")
	globalSettings.SetString(settings.TargetCompID, "ISLD")
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

	app := new(TradeClient)

	initiator, err := quickfix.NewInitiator(app, appSettings, log.ScreenLogFactory{})
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
