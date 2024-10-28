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

package cmd

import (
	"github.com/quickfixgo/examples/cmd/executor"
	"github.com/quickfixgo/examples/cmd/ordermatch"
	"github.com/quickfixgo/examples/cmd/tradeclient"
	"github.com/quickfixgo/examples/version"
	"github.com/spf13/cobra"
)

var (
	// versionF flag prints the version and exits.
	versionF bool
)

// Execute adds all child commands to the root command sets flags appropriately.
func Execute() error {

	c := &cobra.Command{
		Use: "qf",
		RunE: func(cmd *cobra.Command, args []string) error {
			if versionF {
				version.PrintVersion()
				return nil
			}
			return cmd.Usage()
		},
	}

	c.AddCommand(executor.Cmd)
	c.AddCommand(ordermatch.Cmd)
	c.AddCommand(tradeclient.Cmd)
	c.Flags().BoolVarP(&versionF, "version", "v", false, "show the version and exit")
	return c.Execute()
}
