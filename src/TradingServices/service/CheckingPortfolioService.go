package service

import (
	"net/http"
	"fmt"
	"strconv"
	"strings"
	"errors"
)

// ----------------------------------------------------------------------------
// CheckingPortfolio Service - Request and Response
// ----------------------------------------------------------------------------

// CheckingPortfolioServiceRequest represents a JSON-RPC request params sent by a client.
type CheckingPortfolioServiceRequest struct {
    TradeId uint64
}


// CheckingPortfolioServiceRequest represents a JSON-RPC response sent by a server.
type CheckingPortfolioServiceResponse struct {
    Stocks string  //(E.g. “GOOG:100:+$520.25”, “YHOO:200:-$30.40”)
    CurrentMarketValue float32
    UnvestedAmount float32
}


type CheckingPortfolioService struct{}

// ----------------------------------------------------------------------------
// Function Definition - For Checking Portfolio
// ----------------------------------------------------------------------------

func (h *CheckingPortfolioService) CheckPortfolio(r *http.Request, requestParams *CheckingPortfolioServiceRequest, response *CheckingPortfolioServiceResponse) error {
	fmt.Println("=====================Processing a request for 'Checking Portfolio' ======================")
    
    InvalidSymbolsEntered = false
        
    TradeId := requestParams.TradeId
         
    //Format above information in required response params.
	if(len(TradeDetails[TradeId].StockSymbols) == 0) {
		fmt.Println("No earlier response for this TradeId is present in the system! Please try again later after buying stocks!")
		return errors.New("No earlier response for this TradeId is present in the system! Please try again later after buying stocks!")
	} else {
	    //Fetch currentPrices again using Yahoo Finance API, to determine losses/gains by the client
	    var CurrentStockPrices = make([]float32, len(TradeDetails[TradeId].StockSymbols))
	    
		CurrentStockPrices,InvalidSymbolsEntered = CallYahooAPI(TradeDetails[TradeId].StockSymbols)
	    
	    //variable to save if-there-is-any-profit-boolean per stock
	    var IsAnyProfit = make([]string, len(TradeDetails[TradeId].StockSymbols))
	
		var StockString string
		var index = 0
		var profit string
		for index = 0; index < len(TradeDetails[TradeId].StockSymbols) - 1; index++ {		
			//Determine if value has increased or decresed for the StockPrice
			profit = "false"
			if(TradeDetails[TradeId].StockPurchasePrices[index] < CurrentStockPrices[index]) {
				profit = "true"
			} else if(TradeDetails[TradeId].StockPurchasePrices[index] > CurrentStockPrices[index]) {
				profit = "false"
			} else if(TradeDetails[TradeId].StockPurchasePrices[index] == CurrentStockPrices[index]) {
				profit = "nil"
			}
			
			//Store profit-boolean
			IsAnyProfit[index] = profit
			
			tempPrice := strconv.FormatFloat(float64(CurrentStockPrices[index]), 'f', -1, 32)
			if(profit == "true") {
				s := []string{StockString, "\"", TradeDetails[TradeId].StockSymbols[index], ":", 
	    			strconv.Itoa(TradeDetails[TradeId].StockNumbers[index]), ":$", "+", tempPrice, "\", "};
	    		StockString = strings.Join(s, "");
	    	} else if(profit == "false") {
				s := []string{StockString, "\"", TradeDetails[TradeId].StockSymbols[index], ":", 
	    			strconv.Itoa(TradeDetails[TradeId].StockNumbers[index]), ":$", "-", tempPrice, "\", "};
	    		StockString = strings.Join(s, "");
	    	} else if(profit == "nil") {
				s := []string{StockString, "\"", TradeDetails[TradeId].StockSymbols[index], ":", 
	    			strconv.Itoa(TradeDetails[TradeId].StockNumbers[index]), ":$", tempPrice, "\", "};
	    		StockString = strings.Join(s, "");
	    	}	
		}
		profit = "false"
		if(TradeDetails[TradeId].StockPurchasePrices[index] < CurrentStockPrices[index]) {
			profit = "true"
		} else if(TradeDetails[TradeId].StockPurchasePrices[index] > CurrentStockPrices[index]) {
			profit = "false"
		} else if(TradeDetails[TradeId].StockPurchasePrices[index] == CurrentStockPrices[index]) {
			profit = "nil"
		}
		//Store profit-boolean
		IsAnyProfit[index] = profit
		
		tempPrice := strconv.FormatFloat(float64(CurrentStockPrices[index]), 'f', -1, 32)
	    if(profit == "true" ) {
	    	s := []string{StockString, "\"", TradeDetails[TradeId].StockSymbols[index], ":", 
	    		strconv.Itoa(TradeDetails[TradeId].StockNumbers[index]), ":$","+",tempPrice, "\""};
	    		StockString = strings.Join(s, "") //(E.g. “GOOG:100:+$520.25”, “YHOO:200:-$30.40”)
	    } else if(profit == "false") {
	    	s := []string{StockString, "\"", TradeDetails[TradeId].StockSymbols[index], ":", 
	    		strconv.Itoa(TradeDetails[TradeId].StockNumbers[index]), ":$","-",tempPrice, "\""};
	    		StockString = strings.Join(s, "") //(E.g. “GOOG:100:+$520.25”, “YHOO:200:-$30.40”)
	    } else if(profit == "nil") {
	    	s := []string{StockString, "\"", TradeDetails[TradeId].StockSymbols[index], ":", 
	    		strconv.Itoa(TradeDetails[TradeId].StockNumbers[index]), ":$",tempPrice, "\""};
	    		StockString = strings.Join(s, "") //(E.g. “GOOG:100:+$520.25”, “YHOO:200:-$30.40”)
	    }
	    
	    //Calculate CurrentMarketValue
	    var CurrentMarketValue,UnvestedAmount float32
	    UnvestedAmount = 0 
	    CurrentMarketValue = 0
	    
	    for index = 0; index < len(TradeDetails[TradeId].StockSymbols); index++ {
	    	CurrentMarketValue = CurrentMarketValue + (float32(TradeDetails[TradeId].StockNumbers[index]) * CurrentStockPrices[index])
	    	UnvestedAmount =  UnvestedAmount + TradeDetails[TradeId].StockUnvestedAmounts[index]	
	    }
	
		//Set the response params
	    response.Stocks = StockString  //(E.g. “GOOG:100:+$520.25”, “YHOO:200:-$30.40”)
	    response.CurrentMarketValue = CurrentMarketValue
	    response.UnvestedAmount = UnvestedAmount
		fmt.Println("Response is :")
		fmt.Println("Stocks : ",response.Stocks)
		fmt.Println("CurrentMarketValue : ",response.CurrentMarketValue)
		fmt.Println("UnvestedAmount : ",response.UnvestedAmount)
		fmt.Println("=========================================================================================")	    
	    return nil
	}	  
	return nil
}

