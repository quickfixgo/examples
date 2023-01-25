package utils

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/gosuri/uitable"

	"github.com/quickfixgo/quickfix"
)

type screenLog struct {
	prefix string
}

func (l screenLog) OnIncoming(s []byte) {
	table := uitable.New()
	table.MaxColWidth = 150
	table.Wrap = true // wrap columns

	table.AddRow(" |Time:", fmt.Sprintf("%v", time.Now().UTC()))
	table.AddRow(" |Session:", l.prefix)
	table.AddRow(" |Content:", string(s))

	color.Set(color.Bold, color.FgBlue)
	fmt.Println("<=== Incoming FIX Msg: <===")
	fmt.Println(table)
	color.Unset()
}

func (l screenLog) OnOutgoing(s []byte) {
	table := uitable.New()
	table.MaxColWidth = 150
	table.Wrap = true // wrap columns

	table.AddRow(" |Time:", fmt.Sprintf("%v", time.Now().UTC()))
	table.AddRow(" |Session:", l.prefix)
	table.AddRow(" |Content:", string(s))

	color.Set(color.Bold, color.FgMagenta)
	fmt.Println("===> Outgoing FIX Msg: ===>")
	fmt.Println(table)
	color.Unset()
}

func (l screenLog) OnEvent(s string) {

	table := uitable.New()
	table.MaxColWidth = 150
	table.Wrap = true // wrap columns

	table.AddRow(" |Time:", fmt.Sprintf("%v", time.Now().UTC()))
	table.AddRow(" |Session:", l.prefix)
	table.AddRow(" |Content:", s)

	color.Set(color.Bold, color.FgCyan)
	fmt.Println("==== Event: ====")
	fmt.Println(table)
	color.Unset()
}

func (l screenLog) OnEventf(format string, a ...interface{}) {
	l.OnEvent(fmt.Sprintf(format, a...))
}

type screenLogFactory struct{}

func (screenLogFactory) Create() (quickfix.Log, error) {
	log := screenLog{"GLOBAL"}
	return log, nil
}

func (screenLogFactory) CreateSessionLog(sessionID quickfix.SessionID) (quickfix.Log, error) {
	log := screenLog{sessionID.String()}
	return log, nil
}

// NewFancyLog creates an instance of LogFactory that writes messages and events to stdout.
func NewFancyLog() quickfix.LogFactory {
	return screenLogFactory{}
}
