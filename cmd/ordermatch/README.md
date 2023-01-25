# Ordermatch
Ordermatch is a simple matching engine (A set of orderbooks) with a FIX acceptor service as the point of ingress. 

## Features
* Accept any canonical `NewOrderSingle` message for an instrument, with the instrument symbol consisting of an arbitrary string
* Accept any canonical `OrderCancelRequest` message for any order resting in the book
* Accept any canonical `MarketDataRequest` message for any book 
* Sends `ExecutionReport` messages when orders are matched, either partially or in full
* Reads text from `stdin`, either `#symbols` to display the active market symbols, or your symbol, <i>i.e.</i> `AAPL` and will display the state of the book for that symbol 


## Usage
A config file similar to the example config [here](../../config/ordermatch.cfg) is required to run the ordermatch.
The cli command usage takes the form of

```sh
qf ordermatch [CONFIG_PATH_FILENAME]
```
where CONFIG_PATH_FILENAME defaults to `config/ordermatch.cfg`

## Example Config Contents
```
[DEFAULT]
SocketAcceptPort=5001
SenderCompID=ISLD
TargetCompID=TW
ResetOnLogon=Y
FileLogPath=tmp

[SESSION]
BeginString=FIX.4.0

[SESSION]
BeginString=FIX.4.1

[SESSION]
BeginString=FIX.4.2

[SESSION]
BeginString=FIX.4.3

[SESSION]
BeginString=FIX.4.4

[SESSION]
BeginString=FIXT.1.1
DefaultApplVerID=7
```