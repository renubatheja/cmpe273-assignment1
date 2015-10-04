// Copyright 2009 The Go Authors. All rights reserved.
// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	service "TradingServices/service"
)

// ----------------------------------------------------------------------------------
// A test file to send a request to 'BuyingStocks' service and verifying the response
// ----------------------------------------------------------------------------------


type BuyingStocksServiceBadRequest struct {
	M string `json:"method"`
}

func execute(t *testing.T, s *rpc.Server, method string, req, res interface{}) error {
	if !s.HasMethod(method) {
		t.Fatal("Expected to be registered:", method)
	}

	buf, _ := EncodeJSONRPCClientRequest(method, req)
	body := bytes.NewBuffer(buf)
	r, _ := http.NewRequest("POST", "http://localhost:8080/rpc/", body)
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)

	return DecodeJSONRPCClientResponse(w.Body, res)
}


func TestBuyingStocksService(t *testing.T) {
	server := rpc.NewServer()
    jsonCodec := json.NewCodec()
    server.RegisterCodec(jsonCodec, "application/json")
    server.RegisterCodec(jsonCodec, "application/json; charset=UTF-8") // For firefox 11 and other browsers which append the charset=UTF-8
	server.RegisterService(new(service.BuyingStocksService), "")

	var res service.BuyingStocksServiceResponse
	if err := execute(t, server, "BuyingStocksService.BuyStocks", &service.BuyingStocksServiceRequest{100, "GOOG:100%"}, &res); err != nil {
		t.Error("Expected err to be nil, but got:", err)
	}
	if res.Stocks == "" {
		t.Errorf("Wrong response: %v.", res.Stocks)
	}

	if err := execute(t, server, "BuyingStocksService.ResponseError", &service.BuyingStocksServiceRequest{10000, "YHOO:100%"}, &res); err == nil {
		t.Errorf("Expected to get %q, but got nil", service.ErrResponseError)
	} else if err.Error() != service.ErrResponseError.Error() {
		t.Errorf("Expected to get %q, but got %q", service.ErrResponseError, err)
	}

}

