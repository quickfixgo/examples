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

package tradeclient

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"

	"github.com/quickfixgo/examples/cmd/tradeclient/internal"
	"github.com/quickfixgo/examples/cmd/tradeclient/loadtest"
	"github.com/quickfixgo/examples/cmd/utils"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/fix44/logon"
	"github.com/quickfixgo/quickfix"
	"github.com/spf13/cobra"
)

// TradeClient implements the quickfix.Application interface
type TradeClient struct {
}

// OnCreate implemented as part of Application interface
func (e TradeClient) OnCreate(sessionID quickfix.SessionID) {
	fmt.Printf("initiator session Id: %s\n", sessionID)
}

// OnLogon implemented as part of Application interface
func (e TradeClient) OnLogon(sessionID quickfix.SessionID) {}

// OnLogout implemented as part of Application interface
func (e TradeClient) OnLogout(sessionID quickfix.SessionID) {
	fmt.Printf("OnLogout: %s\n", sessionID)
}

// FromAdmin implemented as part of Application interface
func (e TradeClient) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	utils.PrintInfo(fmt.Sprintf("FromAdmin: %s\n", msg.String()))
	return nil
}

const (
	FIX_SEP = "\u0001"

	Publickey = "c"

	Privatekey = "b"

	APIKey = "a"
)

// ToAdmin implemented as part of Application interface
func (e TradeClient) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
	msgType, err := msg.MsgType()
	if err != nil {
		println("wrong message type")
	}

	if msgType == "A" {
		msg.Body.Set(field.NewPassword(APIKey))
		signature, err := e.sign(msg)
		if err != nil {
			println("error in signing the message")
		}
		msg.Body.Set(field.NewRawData(signature))
	}

	utils.PrintInfo(fmt.Sprintf("ToAdmin: %s", msg.String()))
}

func (e TradeClient) sign(logonmsg *quickfix.Message) (string, error) {
	msg := logon.FromMessage(logonmsg)

	sendingTime, err := msg.GetSendingTime()
	if err != nil {
		println("error in getting SendingTime from the logon")
		return "", &quickfix.RejectLogon{Text: "invalid SendingTime"}
	}

	seqNum, err := msg.GetMsgSeqNum()
	if err != nil {
		println("error in getting MsgSeqNum from the logon")
		return "", &quickfix.RejectLogon{Text: "invalid MsgSeqNum"}
	}

	senderCompID, err := msg.GetSenderCompID()
	if err != nil {
		println("error in getting SenderCompID from the logon")
		return "", &quickfix.RejectLogon{Text: "invalid SenderCompID"}
	}

	targetCompID, err := msg.GetTargetCompID()
	if err != nil {
		println("error in getting TargetCompID from the logon")
		return "", &quickfix.RejectLogon{Text: "invalid TargetCompID"}
	}

	msgToSign := sendingTime.Format("20060102-15:04:05.000") + FIX_SEP +
		strconv.Itoa(seqNum) + FIX_SEP +
		senderCompID + FIX_SEP +
		targetCompID

	privateKeyBytes, _ := hex.DecodeString(Privatekey)
	ed25519PrivateKey := ed25519.PrivateKey(privateKeyBytes)
	signature := ed25519.Sign(ed25519PrivateKey, []byte(msgToSign))
	return hex.EncodeToString(signature), nil
}

// ToApp implemented as part of Application interface
func (e TradeClient) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) (err error) {
	utils.PrintInfo(fmt.Sprintf("Sending: %s", msg.String()))
	return
}

// FromApp implemented as part of Application interface. This is the callback for all Application level messages from the counter party.
func (e TradeClient) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	utils.PrintInfo(fmt.Sprintf("FromApp: %s", msg.String()))
	return
}

const (
	usage = "tradeclient"
	short = "Start a tradeclient (FIX initiator) cli trading agent"
	long  = "Start a tradeclient (FIX initiator) cli trading agent."
)

var (
	// Cmd is the quote command.
	Cmd = &cobra.Command{
		Use:     usage,
		Short:   short,
		Long:    long,
		Aliases: []string{"tc"},
		Example: "qf tradeclient [YOUR_FIX_CONFIG_FILE_HERE.cfg] (default is ./config/tradeclient.cfg)",
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
			utils.PrintInfo("attempting to use default location './config/tradeclient.cfg' ...")
			cfgFileName = path.Join("config", "tradeclient.cfg")
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

	app := TradeClient{}
	fileLogFactory, err := quickfix.NewFileLogFactory(appSettings)

	if err != nil {
		return fmt.Errorf("error creating file log factory: %s,", err)
	}

	initiator, err := quickfix.NewInitiator(app, quickfix.NewMemoryStoreFactory(), appSettings, fileLogFactory)
	if err != nil {
		return fmt.Errorf("unable to create initiator: %s", err)
	}

	err = initiator.Start()
	if err != nil {
		return fmt.Errorf("unable to start initiator: %s", err)
	}

	utils.PrintConfig("initiator", bytes.NewReader(stringData))

Loop:
	for {
		action, err := internal.QueryAction()
		if err != nil {
			break
		}

		switch action {
		case "1":
			err = internal.QueryEnterOrder("CUST2_Order", "ANCHORAGE")

		case "2":
			err = internal.QueryCancelOrder()

		case "3":
			err = internal.QueryMarketDataRequest("CUST2_Marketdata", "ANCHORAGE")

		case "4":
			var ordersPerSecond int
			var totalOrders int

			// Prompt the user for orders per second
			fmt.Print("Enter orders per second: ")
			_, err := fmt.Scanf("%d", &ordersPerSecond)
			if err != nil {
				utils.PrintBad("Invalid input for orders per second")
				break
			}

			// Prompt the user for total number of orders
			fmt.Print("Enter total number of orders: ")
			_, err = fmt.Scanf("%d", &totalOrders)
			if err != nil {
				utils.PrintBad("Invalid input for total orders")
				break
			}

			// Create load test configuration
			loadTestConfig := loadtest.LoadTestConfig{
				OrdersPerSecond: ordersPerSecond,
				TotalOrders:     totalOrders,
				SenderCompID:    "CUST2_Order",
				TargetCompID:    "ANCHORAGE",
			}

			// Run the load test
			loadtest.RunLoadTest(loadTestConfig)

		case "5":
			//quit
			break Loop

		default:
			err = fmt.Errorf("unknown action: '%v'", action)
		}

		if err != nil {
			utils.PrintBad(err.Error())
		}
	}

	utils.PrintInfo("stopping FIX initiator ..")
	initiator.Stop()
	utils.PrintInfo("stopped")
	return nil
}
