# Example QuickFIX/Go Applications

[![Build Status](https://github.com/quickfixgo/examples/workflows/CI/badge.svg)](https://github.com/quickfixgo/examples/actions) [![GoDoc](https://godoc.org/github.com/quickfixgo/examples?status.png)](https://godoc.org/github.com/quickfixgo/examples) [![Go Report Card](https://goreportcard.com/badge/github.com/quickfixgo/examples)](https://goreportcard.com/report/github.com/quickfixgo/examples)

* [TradeClient](cmd/tradeclient/README.md) is a simple console based trading client
* [Executor](cmd/executor/README.md) is a server that fills every limit order it receives
* [OrderMatch](cmd/ordermatch/README.md) is a primitive matching engine 

All examples have been ported from [QuickFIX](http://quickfixengine.org)

## Installation

### Build From Source
To build and run the examples, you will first need [Go](https://www.golang.org) installed on your machine

Next, clone this repository with `git clone git@github.com:quickfixgo/examples.git`. This project uses go modules, so you just need to type `make build`. This will compile the examples executables in the `./bin` dir in your local copy of the repo. If this exits with exit status 0, then everything is working! You may need to pull the module deps with `go mod download`.

```sh
make build
```

### Running the Examples

Following installation, the examples can be found in `./bin`.  The examples are meant to be run in pairs- the TradeClient as a client of either the Executor or OrderMatch.  By default, the examples will load the default configurations named after the example apps provided in the `config/` root directory.  <i>i.e.</i>, running `./bin/tradeclient` will load the `config/tradeclient.cfg` configuration.  Each example can be run with a custom configuration as a command line argument (`./bin/tradeclient my_trade_client.cfg`).

### Licensing

This software is available under the QuickFIX Software License. Please see the [LICENSE](LICENSE) for the terms specified by the QuickFIX Software License.
