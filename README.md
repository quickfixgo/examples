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

Next, using [Git](https://git-scm.com/), clone this repository into `$GOPATH/src/github.com/quickfixgo/examples`. All the necessary dependencies are either vendored, so you just need to type `go build ./cmd/...`. This will compile the example code. If this exits with exit status 0, then everything is working!

Licensing
---------

This software is available under the QuickFIX Software License. Please see the [LICENSE.txt](https://github.com/quickfixgo/examples/blob/master/LICENSE.txt) for the terms specified by the QuickFIX Software License.
