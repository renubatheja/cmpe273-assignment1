# A Virtual Stock Trading System 

An implementation of JSON-RPC over HTTP server/client in Go using Gorilla framework

This system uses real-time pricing via Yahoo finance API(YQL) and has two features:

1. Buying Stocks
2. Checking Portfolio

  
_Note: Currently, this system supports only THREE latest responses stored in-memory._ 

## Build 

ADD current dir to GOPATH using

```
export GOPATH=$GOPATH:$PWD
```

Install Gorilla framework using:

```
go get github.com/gorilla/rpc
```

To build the project run

```
go build
```

## Usage

Following is the list of commands to be used for running server and client from CLI

### For Server

```
cd src\Server
Server.exe
```
OR
```
cd src\Server
go run Server.go
```

### For Client

#### For Buying Stocks :
##### Buying Stock using 'Client.exe'
```
cd src\Client
Client.exe <BUDGET field in Number> <StockSymbol_AND_Percentages field in String>
```
e.g.

```
Client.exe 10000 "GOOG:50%,YHOO:50%"
```

##### Buying Stock using 'go run'

```
cd src\Client
go run Client.go <BUDGET field in Number> <StockSymbol_AND_Percentages field in String>
```

e.g.

```
go run Client.go 10000 "GOOG:50%,YHOO:50%"
```

##### Buying Stock using 'CURL'

Request for buying stock
```
curl -X POST localhost:8080/rpc -H 'Content-Type:application/json' -d '{"Id": 1, "Method": "BuyingStocksService.BuyStocks", "params": [{"Budget":2000,"StockSymbolAndPercentage":"GOOG:100%"}]}'
```

Sample response will look like:
```
{"result":{"Stocks":"\"GOOG:3:$627\"","UnvestedAmount":119,"TradeId":0},"error":null,"id":1}
```

#### For Checking Portfolio :		
##### Checking Portfolio using 'Client.exe'
```	
cd src\Client
Client.exe <TradeId field in Number - O/1/2>
```
e.g.
```
Client.exe 1
```
##### Checking Portfolio using 'go run'
```
cd src\Client
go run Client.go <TradeId field in Number - allowed numbers are O, 1 and 2>
```
e.g.
```
go run Client.go 1
```

##### Checking portfolio using 'CURL'

Requst for checking portfolio
```
curl -X POST localhost:8080/rpc -H 'Content-Type:application/json' -d '{"Id": 1, "Method": "CheckingPortfolioService.CheckPortfolio", "params": [{"TradeId": 1}]}'
```

Sample response will look like:
```
{"result":{"Stocks":"\"GOOG:3:$627\"","CurrentMarketValue":1881,"UnvestedAmount":119},"error":null,"id":1}
```

