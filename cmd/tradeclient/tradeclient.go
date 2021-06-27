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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/quickfixgo/examples/cmd/tradeclient/internal"
	"github.com/quickfixgo/quickfix"
	"github.com/spf13/cobra"
)

//TradeClient implements the quickfix.Application interface
type TradeClient struct {
}

//OnCreate implemented as part of Application interface
func (e TradeClient) OnCreate(sessionID quickfix.SessionID) {}

//OnLogon implemented as part of Application interface
func (e TradeClient) OnLogon(sessionID quickfix.SessionID) {}

//OnLogout implemented as part of Application interface
func (e TradeClient) OnLogout(sessionID quickfix.SessionID) {}

//FromAdmin implemented as part of Application interface
func (e TradeClient) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return nil
}

//ToAdmin implemented as part of Application interface
func (e TradeClient) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {}

//ToApp implemented as part of Application interface
func (e TradeClient) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) (err error) {
	fmt.Printf("Sending %s\n", msg)
	return
}

//FromApp implemented as part of Application interface. This is the callback for all Application level messages from the counter party.
func (e TradeClient) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	fmt.Printf("FromApp: %s\n", msg.String())
	return
}

const (
	usage = "tradeclient"
	short = "Start a tradeclient"
	long  = "Start a tradeclient."
)

var (
	// Cmd is the quote command.
	Cmd = &cobra.Command{
		Use:     usage,
		Short:   short,
		Long:    long,
		Aliases: []string{"tc"},
		Example: "qf tradeclient config/tradeclient.cfg",
		RunE:    execute,
	}
)

func execute(cmd *cobra.Command, args []string) error {
	var cfgFileName string
	argLen := len(args)
	switch argLen {
	case 0:
		{
			cfgFileName = path.Join("config", "tradeclient.cfg")
		}
	case 1:
		{
			cfgFileName = args[0]
		}
	default:
		{
			return fmt.Errorf("Incorrect argument number")
		}
	}

	cfg, err := os.Open(cfgFileName)
	if err != nil {
		return fmt.Errorf("Error opening %v, %v\n", cfgFileName, err)
	}
	defer cfg.Close()

	stringData, readErr := ioutil.ReadAll(cfg)
	if readErr != nil {
		return fmt.Errorf("Error reading cfg: %s,", readErr)
	}

	appSettings, err := quickfix.ParseSettings(bytes.NewReader(stringData))
	if err != nil {
		return fmt.Errorf("Error reading cfg: %s,", err)
	}

	app := TradeClient{}
	fileLogFactory, err := quickfix.NewFileLogFactory(appSettings)

	if err != nil {
		return fmt.Errorf("Error creating file log factory: %s,", err)
	}

	initiator, err := quickfix.NewInitiator(app, quickfix.NewMemoryStoreFactory(), appSettings, fileLogFactory)
	if err != nil {
		return fmt.Errorf("Unable to create Initiator: %s\n", err)
	}

	err = initiator.Start()
	if err != nil {
		return fmt.Errorf("Unable to start Initiator: %s\n", err)
	}

	printConfig(bytes.NewReader(stringData))

Loop:
	for {
		action, err := internal.QueryAction()
		if err != nil {
			break
		}

		switch action {
		case "1":
			err = internal.QueryEnterOrder()

		case "2":
			err = internal.QueryCancelOrder()

		case "3":
			err = internal.QueryMarketDataRequest()

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
	return nil
}

func printConfig(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	color.Set(color.Bold)
	fmt.Println("Started FIX initiator with config:")
	color.Unset()

	color.Set(color.FgHiMagenta)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	color.Unset()
}
