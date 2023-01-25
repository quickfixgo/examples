# Executor
Executor is a FIX acceptor service that fills every limit order it receives. 

(Note: it will reject any non-limit order type)

## Features
* Accept any canonical `NewOrderSingle` message for an instrument, with the instrument symbol consisting of an arbitrary string
* Sends `ExecutionReport` messages as responses indicating order fills

## Usage
A config file similar to the example config [here](../../config/executor.cfg) is required to run the executor.
The cli command usage takes the form of

```sh
qf executor [CONFIG_PATH_FILENAME]
```
where CONFIG_PATH_FILENAME defaults to `config/executor.cfg`

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