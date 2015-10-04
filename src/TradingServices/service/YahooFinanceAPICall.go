package service

import(
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"strings"
)

// ----------------------------------------------------------------------------
// Data Structure Definition - For Storing YQL (Yahoo Query Lang) Responses
// ----------------------------------------------------------------------------

// Data structures definitions to store the data returned by Yahoo Finance API. 
type MultiStock struct {
	Query struct {
		Created string `json:"Created"`
		Results struct {
			Quote MultiQuotes `json:"quote"`
		}
	}
}

type SingleStock struct {
	Query struct {
		Created string `json:"Created"`
		Results struct {
			Quote SingleQuote `json:"quote"`
		}
	}
}

type SingleQuote QuoteInfo

type MultiQuotes []QuoteInfo

type QuoteInfo struct {
	Symbol string  `json:"Symbol"`
	Ask string	`json:"ask"`
}

// ----------------------------------------------------------------------------------------
// Yahoo Finance API usage - Using YQL to get the current stock price in real-time scenario
// ----------------------------------------------------------------------------------------

func CallYahooAPI(Symbols []string) ([]float32, bool){
	var query string
	var CurrentStockPrices = make([]float32, len(Symbols))
	var symbols string
	if len(Symbols) > 1 {
		symbols = strings.Join(Symbols,"%22%2C%22")
	} else {
		symbols = strings.Join(Symbols,"")
	}

	query = "https://query.yahooapis.com/v1/public/yql?q=select%20symbol%2CAsk%20from%20yahoo.finance.quotes%20where%20symbol%20in%20(%22"+
			symbols +
			"%22)&format=json&diagnostics=true&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback="
	resp, err := http.Get(query)
	
	if err != nil {
		fmt.Println("Encountered an error while running Yahoo Finance API!")
	}
	
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
		
	if len(Symbols) > 1 { 
		var m MultiStock
		err1 := json.Unmarshal(body, &m)
		for index :=0; index < len(Symbols); index++ {
			if m.Query.Results.Quote[index].Ask != "" {
				currPrice,_ := strconv.ParseFloat(m.Query.Results.Quote[index].Ask, 32)
				CurrentStockPrices[index] = float32(currPrice)
			} else {
				InvalidSymbolsEntered = true
				return CurrentStockPrices, InvalidSymbolsEntered
			}	
		} 
		if(err1 != nil) {
			fmt.Println("Error faced while fetching data using Yahoo Finance API : ", err1)
		}
	} else {
		var m SingleStock
		err1 := json.Unmarshal(body, &m)
		if m.Query.Results.Quote.Ask != "" {
			currPrice,_ := strconv.ParseFloat(m.Query.Results.Quote.Ask, 32)
			CurrentStockPrices[0] = float32(currPrice)
		} else {
			InvalidSymbolsEntered = true
			return CurrentStockPrices, InvalidSymbolsEntered
		}
		if(err1 != nil) {
			fmt.Println("Error faced while fetching data using Yahoo Finance API : ", err1)
		}		
	}
	return CurrentStockPrices, InvalidSymbolsEntered
}
