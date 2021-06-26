# Tradeclient
Tradeclient is a simple FIX initiator console-based trading client

## Features
* Send configurable NewOrderSingle, OrderCancelRequest, and MarketDataRequest messages on a FIX session to a remote server 
* Supports Buy/Sell/Short/Cross/Cross Short order sides 
* Supports Market/Limit/Stop/Stop Limit order types
* Supports Day/IOC/OPG/GTC/GTX time in force

## Usage
A config file similar to the example config [here](../../config/tradeclient.cfg) is required to run the tradeclient.
The cli command usage takes the form of

```sh
qf tradeclient [CONFIG_PATH_FILENAME]
```
where CONFIG_PATH_FILENAME defaults to `config/tradeclient.cfg`


## Example Config Contents
```
[DEFAULT]
SocketConnectHost=127.0.0.1
SocketConnectPort=5001
HeartBtInt=30
SenderCompID=TW
TargetCompID=ISLD
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