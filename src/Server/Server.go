package main

import (
    "github.com/gorilla/rpc"
    "github.com/gorilla/rpc/json"
    "net/http"
    s "TradingServices/service"
)  

// ----------------------------------------------------------------------------
// Server: JSON-RPC interface for Virtual Stock Trading System
// ----------------------------------------------------------------------------

//This trading engine has JSON-RPC interface for two features(Services) - Buying Stocks and Checking portfolio
//Both services are registered with this server.
 
func main() {
    jsonRPC := rpc.NewServer()
    jsonCodec := json.NewCodec()
    jsonRPC.RegisterCodec(jsonCodec, "application/json")
    jsonRPC.RegisterCodec(jsonCodec, "application/json; charset=UTF-8") // For firefox 11 and other browsers which append the charset=UTF-8
    jsonRPC.RegisterService(new(s.BuyingStocksService), "")
    jsonRPC.RegisterService(new(s.CheckingPortfolioService), "")
    http.Handle("/rpc", jsonRPC)
    http.ListenAndServe(":8080", jsonRPC)    
}

