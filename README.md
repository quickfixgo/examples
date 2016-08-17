Example QuickFIX/Go Applications
================================

[![Build Status](https://travis-ci.org/quickfixgo/examples.svg?branch=master)](https://travis-ci.org/quickfixgo/examples)

* TradeClient is a simple console based trading client
* Executor is a server that fills every limit order it receives
* OrderMatch is a primitive matching engine 

All examples have been ported from [QuickFIX](http://quickfixengine.org)

Installation
------------

To build and run the examples, you will first need [Go](http://www.golang.org) installed on your machine (version 1.6+ is *required*).

For local dev first make sure Go is properly installed, including setting up a [GOPATH](http://golang.org/doc/code.html#GOPATH).

Next, using [Git](https://git-scm.com/), clone this repository into `$GOPATH/src/github.com/quickfixgo/examples`. All the necessary dependencies are either vendored, so you just need to type `make`. This will compile and install the examples into `$GOPATH/bin`. If this exits with exit status 0, then everything is working!

```sh
$ make
```

Running the Examples
--------------------

Following installation, the examples can be found in `$GOPATH/bin`.  The examples are meant to be run in pairs- the TradeClient as a client of either the Executor or OrderMatch.  By default, the examples will load the default configurations named after the example apps provided in the `config/` root directory.   Eg, running `$GOPATH/bin/tradeclient` will load the `config/tradeclient.cfg` configuration.  Each example can be run with a custom configuration as a command line argument (`$GOPATH/bin/tradeclient my_trade_client.cfg`).

Licensing
---------

This software is available under the QuickFIX Software License. Please see the [LICENSE.txt](https://github.com/quickfixgo/examples/blob/master/LICENSE.txt) for the terms specified by the QuickFIX Software License.
