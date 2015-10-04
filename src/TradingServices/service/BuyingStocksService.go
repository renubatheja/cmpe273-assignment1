package service

import (
	"net/http"
	"strings"
	"fmt"
	"strconv"
	"math"
	"errors"
)

// ----------------------------------------------------------------------------
// BuyingStock Service - Request and Response
// ----------------------------------------------------------------------------

// BuyingStocksServiceRequest represents a JSON-RPC request params sent by a client.
type BuyingStocksServiceRequest struct {
    Budget float32					 // 20000
    StockSymbolAndPercentage string  //(E.g. “GOOG:50%,YHOO:50%”)
}


// BuyingStocksServiceResponse represents a JSON-RPC response returned by server.
type BuyingStocksServiceResponse struct {
    Stocks string					//(E.g. “GOOG:100:$500.25”, “YHOO:200:$31.40”)
    UnvestedAmount float32
    TradeId int32
}


type BuyingStocksService struct{}

// ----------------------------------------------------------------------------
// Data Structure Definitions - For Storing Responses by Server
// ----------------------------------------------------------------------------

// Data structure to store the Trading response information 
type TradeInfo struct {
	TradeId int32
	StockSymbols TradeStockSymbols
	StockPurchasePrices TradeStockPurchasePrices
	StockPercentages TradeStockPercentages
	StockNumbers TradeStockNumbers
	StockUnvestedAmounts TradeStockUnvestedAmounts
}


type TradeStockSymbols []string
type TradeStockPurchasePrices []float32
type TradeStockPercentages []float32
type TradeStockNumbers []int
type TradeStockUnvestedAmounts []float32

// ----------------------------------------------------------------------------
// Variable Declarations
// ----------------------------------------------------------------------------

//System saves only LATEST 3 responses in the memory
var TradeDetails = make([]TradeInfo, 3)

var CurrentResponseId int32  //This helps for storing currentResponseIds in the memory
var InvalidSymbolsEntered bool


// ----------------------------------------------------------------------------
// Function Definition - For Buying Stocks
// ----------------------------------------------------------------------------

func (h *BuyingStocksService) BuyStocks(r *http.Request, requestParams *BuyingStocksServiceRequest, response *BuyingStocksServiceResponse) error {
	fmt.Println("=====================Processing a request for 'Buying Stocks' ===========================")
	InvalidSymbolsEntered = false
	
    //log.Printf(requestParams.StockSymbolAndPercentage)
    
    StockSymbolAndPercentage := requestParams.StockSymbolAndPercentage
    Budget := requestParams.Budget
        
    var SymbolsAndPercentages []string 
    var tempSlices []string
    //Split StockSymbolAndPercentage at comma to get StockSymbol and Percentage separately
    SymbolsAndPercentages = strings.Split(StockSymbolAndPercentage, ",")

    var Symbols = make([]string, len(SymbolsAndPercentages))
    var Percentages = make([]float32, len(SymbolsAndPercentages))
    var CurrentPrices = make([]float32,len(SymbolsAndPercentages))
    var StocksPurchased = make([]int,len(SymbolsAndPercentages))
    var UnvestedAmounts = make([]float32, len(SymbolsAndPercentages))
    var TotalUnverstedAmount float32
    var returnStr string
	var sum float32
	sum = 0.0
    
    for index := 0; index < len(SymbolsAndPercentages); index++ {
 		//Split every stocksymbol and percentage and put them into separate arrays/slices
 		tempSlices = strings.Split(SymbolsAndPercentages[index],":")
 		tempSlices[1] = strings.Trim(tempSlices[1],"%")
 		Symbols[index] = tempSlices[0]
 		percentageForThisStock, _ := strconv.ParseFloat(tempSlices[1], 32)
 		Percentages[index] = float32(percentageForThisStock) //tempSlices[1]
 		sum += Percentages[index]
	}

	//Check if summation of % is = 100
	if(sum < 100.0) {
		fmt.Println("The requested percentages for stocks do not meet the sum total 100%!")
		TotalUnverstedAmount = Budget
    	returnStr = ""
	    response.Stocks = returnStr	
	    response.UnvestedAmount = TotalUnverstedAmount
		response.TradeId = CurrentResponseId	    
		
	} else {
	 	//Get current trading price for every symbol using Yahoo Finance API
		var CurrentStockPrices = make([]float32, len(SymbolsAndPercentages))
	
		CurrentStockPrices, InvalidSymbolsEntered = CallYahooAPI(Symbols)
		if(InvalidSymbolsEntered) {
			fmt.Println("You entered the (one or more) invalid Symbols! Please check them again.")
			TotalUnverstedAmount = Budget
	    	returnStr = ""
		    response.Stocks = returnStr	
		    response.UnvestedAmount = TotalUnverstedAmount
			response.TradeId = CurrentResponseId	    
		} else {
			for index := 0; index < len(SymbolsAndPercentages); index++ {
		 		currentPrice := CurrentStockPrices[index]
		    	CurrentPrices[index] = CurrentStockPrices[index]
		 		
		 		//Calculate % of budget and allocate it to the symbol
		    	//find out how many stocks could be purchased
		    	
		 		budgetForThisStock := (Budget * Percentages[index]) / 100.00
		 		var numOfStocksTBP int 
		 		numOfStocksTBP = int(budgetForThisStock / currentPrice)
		 		unvestedAmtFromThisBudget := math.Mod(float64(budgetForThisStock),float64(currentPrice))
		    	StocksPurchased[index] = numOfStocksTBP
		    	UnvestedAmounts[index] = float32(unvestedAmtFromThisBudget)
			}
			
			//As we are trying to save latest 3 responses only in the memory, once the 4th (onwards) requests come, they should replace
			//the existing responses in the memory
			if(CurrentResponseId == 3) {
				fmt.Println("Going to replace the existing old response, as this system supports storage of only 3 responses in-memory!!")
				CurrentResponseId = 0; //Reset to 0, as save the 4th response now at 0th place.
			}    
			    
			//Store this information in memory
			TradeDetails[CurrentResponseId].TradeId = CurrentResponseId
			TradeDetails[CurrentResponseId].StockSymbols = Symbols
			TradeDetails[CurrentResponseId].StockPurchasePrices = CurrentPrices
			TradeDetails[CurrentResponseId].StockNumbers = StocksPurchased
			TradeDetails[CurrentResponseId].StockPercentages = Percentages
			TradeDetails[CurrentResponseId].StockUnvestedAmounts = UnvestedAmounts
			CurrentResponseId = CurrentResponseId + 1	
					
		    //format reponse.Stocks string E.g. "GOOG:100:$500.25", "YHOO:200:$31.40"
		    var index int
		    TotalUnverstedAmount = 0
		    for index = 0; index < len(SymbolsAndPercentages) - 1; index++ {
		    	s := []string{returnStr, "\"", Symbols[index], ":", strconv.Itoa(StocksPurchased[index]), ":$", strconv.FormatFloat(float64(CurrentPrices[index]), 'f', -1, 32), "\", "};
		    	returnStr = strings.Join(s, "");
		    	TotalUnverstedAmount = TotalUnverstedAmount +UnvestedAmounts[index]
		    }
		    s := []string{returnStr, "\"", Symbols[index], ":", strconv.Itoa(StocksPurchased[index]), ":$", strconv.FormatFloat(float64(CurrentPrices[index]), 'f', -1, 32), "\""};
		    returnStr = strings.Join(s, "");
		    TotalUnverstedAmount = TotalUnverstedAmount +UnvestedAmounts[index]
		    response.Stocks = returnStr	//(E.g. “GOOG:100:$500.25”, “YHOO:200:$31.40”)
		    response.UnvestedAmount = TotalUnverstedAmount
			response.TradeId = CurrentResponseId - 1	    
		    //log.Printf(response.Stocks)
	    }
	}
	fmt.Println("Response is :")
	fmt.Println("TradeID : ",response.TradeId)
	fmt.Println("Stocks : ",response.Stocks)
	fmt.Println("UnvestedAmount : ",response.UnvestedAmount)
	fmt.Println("=========================================================================================")
    return nil
}


//ErrResponseError
var ErrResponseError = errors.New("response error")


func (t *BuyingStocksService) ResponseError(r *http.Request, req *BuyingStocksServiceRequest, res *BuyingStocksServiceResponse) error {
	return ErrResponseError
}
