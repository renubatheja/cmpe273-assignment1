// Copyright 2009 The Go Authors. All rights reserved.
// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
    "github.com/gorilla/rpc"
    service "TradingServices/service"	
)

// ----------------------------------------------------------------------------------
// A test file to register services with the JSON-RPC server
// ----------------------------------------------------------------------------------

func TestRegisterService(t *testing.T) {
	var err error
	server := rpc.NewServer()
	service1 := new(service.BuyingStocksService)
	service2 := new(service.CheckingPortfolioService)

	err = server.RegisterService(service1, "")
	if err != nil || !server.HasMethod("BuyingStocksService.BuyStocks") {
		t.Errorf("Expected to be registered: BuyingStocksService.BuyStocks")
	}

	err = server.RegisterService(service2, "")
	if err != nil || !server.HasMethod("CheckingPortfolioService.CheckPortfolio") {
		t.Errorf("Expected to be registered: CheckingPortfolioService.CheckPortfolio")
	}
}

