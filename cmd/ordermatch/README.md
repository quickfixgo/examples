# Ordermatch

## Usage
A config file similar to the example config [here](../../config/ordermatch.cfg) is required to run the ordermatch.
The cli command usage takes the form of

```sh
ordermatch [CONFIG_PATH_FILENAME]
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