# Example QuickFIX/Go Applications

[![Build Status](https://github.com/quickfixgo/examples/workflows/CI/badge.svg)](https://github.com/quickfixgo/examples/actions) [![GoDoc](https://godoc.org/github.com/quickfixgo/examples?status.png)](https://godoc.org/github.com/quickfixgo/examples) [![Go Report Card](https://goreportcard.com/badge/github.com/quickfixgo/examples)](https://goreportcard.com/report/github.com/quickfixgo/examples)

* [TradeClient](cmd/tradeclient/README.md) is a simple FIX initiator console-based trading client
* [Executor](cmd/executor/README.md) is a FIX acceptor service that fills every limit order it receives
* [OrderMatch](cmd/ordermatch/README.md) is a primitive matching engine and FIX acceptor serice

All examples have been ported from [QuickFIX](http://quickfixengine.org)

## Usage
This project builds a cli tool `qf` with 3 commands corresponding to each example.
The generalized usage is of the form:
```sh
qf [GLOBAL FLAGS] [COMMAND] [COMMAND FLAGS] [ARGS]
```

The examples are meant to be run in pairs- the TradeClient as a client of either the Executor or OrderMatcher. By default, the examples will load the default configurations named after the example apps provided in the `config/` root directory.  <i>i.e.</i>, running `qf tradeclient` will load the `config/tradeclient.cfg` configuration.  Each example can be run with a custom configuration as a command line argument (`qf tradeclient my_trade_client.cfg`).


## Installation
In order to use this awesome tool, you'll need to get it on your machine!

### From Homebrew
If you're on macOS, the easiest way to get the examples is through the homebrew tap.
```sh
brew tap quickfixgo/qf
brew install qf
```
Run the command `qf help` in your shell for the list of possible example subcommands.

### From Release
1. Head over to the official [releases page](https://github.com/quickfixgo/examples/releases)
2. Determine the appropriate distribution for your operating system (mac | windows | linux)
3. Download and untar the distribution. Shortcut for macs:
```sh
curl -sL https://github.com/quickfixgo/examples/releases/download/v{VERSION}/qf_{VERSION}_Darwin_x86_64.tar.gz | tar zx
```
4. Move the binary into your local `$PATH`.
5. Run the command `qf help` in your shell for the list of possible example subcommands.

### From Source
To build and run the examples, you will first need [Go](https://www.golang.org) installed on your machine

Next, clone this repository with `git clone git@github.com:quickfixgo/examples.git`. This project uses go modules, so you just need to type `make build`. This will compile the examples executable in the `./bin` dir in your local copy of the repo. If this exits with exit status 0, then everything is working! You may need to pull the module deps with `go mod download`.
```sh
make build
```
Run the command `./bin/qf help` in your shell for the list of possible example subcommands.

### From Snapcraft
Linux OS users can install the examples through the snap store.
```sh
sudo snap install quickfixgo-qf
```
Run the command `qf help` in your shell for the list of possible example subcommands.

### From Scoop
Windows users can install the examples via the Scoop package manager.
```sh
scoop bucket add auth0 https://github.com/auth0/scoop-auth0-cli.git
scoop install auth0
```
Run the command `qf help` in your shell for the list of possible example subcommands.

### Docker Image
The quickfix examples are also available as a docker image [here](https://hub.docker.com/r/quickfixgo/qf). To pull and run the `latest`, use the following command:
```sh
docker run -it quickfixgo/qf
```
To run a specific example, you can do something like this:
```sh
docker run -it -p 5001:5001 quickfixgo/qf ordermatch
```
Note: The docker image comes pre-loaded with the default configs. If you want to supply your own, you can specify a volume binding to your local directory in the run command.

## Contributing
The preferred way to build the examples for development is using `make`. Run `make build` to compile and test any new fixes or features.

1. Fork the repo and clone your fork
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request :)

## Licensing
This software is available under the QuickFIX Software License. Please see the [LICENSE](LICENSE) for the terms specified by the QuickFIX Software License.
