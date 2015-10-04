package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"bytes"
	"net/http"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	s "TradingServices/service"
)

// ----------------------------------------------------------------------------
// Request and Response
// ----------------------------------------------------------------------------

// JSONRPCClientRequest represents a JSON-RPC request sent by a client.
type JSONRPCClientRequest struct {
	// A String containing the name of the method to be invoked.
	Method string `json:"method"`
	// Object to pass as request parameter to the method.
	Params [1]interface{} `json:"params"`
	// The request id. This can be of any type. It is used to match the
	// response with the request that it is replying to.
	Id int32 `json:"id"`
}


// JSONRPCClientResponse represents a JSON-RPC response returned to a client.
type JSONRPCClientResponse struct {
	Result *json.RawMessage `json:"result"`
	Error  interface{}      `json:"error"`
	Id     int32           `json:"id"`
}


// EncodeClientRequest encodes parameters for a JSON-RPC client request.
func EncodeJSONRPCClientRequest(method string, args interface{}) ([]byte, error) {
	c := &JSONRPCClientRequest{
		Method: method,
		Params: [1]interface{}{args},
		Id:     int32(rand.Int31()),
	}
	return json.Marshal(c)
}


// DecodeClientResponse decodes the response body of a client request into
// the interface reply.
func DecodeJSONRPCClientResponse(r io.Reader, reply interface{}) error {
	var c JSONRPCClientResponse
	if err := json.NewDecoder(r).Decode(&c); err != nil {
		if(err.Error() == "EOF") {
			return nil
		} else {
			return err
		}
	}
	if c.Error != nil {
		return fmt.Errorf("%v", c.Error)
	}
	if c.Result == nil {
		return errors.New("result is null")
	}
	return json.Unmarshal(*c.Result, reply)
}


// This function sends the request to the server and decoded the response returned from the server.
func ExecuteJSONRPCRequest(method string, req, res interface{}) error {
	buf, _ := EncodeJSONRPCClientRequest(method, req)
	body := bytes.NewBuffer(buf)
	
	r, _ := http.NewRequest("POST", "http://localhost:8080/rpc/", body)
	r.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, _ := client.Do(r)
	
	defer resp.Body.Close()
	
    body1, _ := ioutil.ReadAll(resp.Body)
    fmt.Printf("Response From Server:%s\n", string(body1));
	return DecodeJSONRPCClientResponse(resp.Body, res)
}


// ----------------------------------------------------------------------------
// Main Function
// ----------------------------------------------------------------------------

//Based on the command line arguments passed, this function decides the service request to be made to the server
//There are two kinds of services registered with the server
//1. Service for buying stocks for user
//2. Service for checking user's portfolio
func main() {
	argsWithProg := os.Args
	if(len(argsWithProg) == 1) {
		fmt.Println("Please enter the command line arguments for below services")
		fmt.Println("To Buy Stocks - Budget and StockSymbol and Percentages")
		fmt.Println("To Check Portfolio - TradeID")
	} else {
		var isInvalidArg = false;
		//Check which request is being made by the client
	    var methodName string
	    if(len(argsWithProg) > 2) {
	    	methodName = "BuyStocks"
	    } else {
	    	methodName = "CheckPortfolio"
	    }	
		
		if(methodName == "BuyStocks") {
			isInvalidArg = ValidateBudget(os.Args[1]);
			if(isInvalidArg) {
				fmt.Println("Please enter a valid NUMBER for 'Budget' field.") 
			} else {
				budget, _ := strconv.ParseFloat(os.Args[1], 32)
				//Budget looks fine, Validate StockSymbol and Percentage field now
				isInvalidArg = ValidateStockSymbolAndPercentage(os.Args[2])
				if(isInvalidArg) {
					fmt.Println("The entered string does not contain Symbols and Percentages in valid format! Please enter it in SYMBOL:PERCENTAGE[,SYMBOL:PERCENTAGE] format.")
				} else {
					isInvalid := ValidateEnteredPercentageFor100OrInvalidSymbols(os.Args[2])
					if(!isInvalid) {
						var response s.BuyingStocksServiceResponse
						if err := ExecuteJSONRPCRequest("BuyingStocksService.BuyStocks", &s.BuyingStocksServiceRequest{float32(budget),os.Args[2]}, &response); err != nil {
							fmt.Println("Error : ", err)
						}
					}
				}
			}
		} else {
			tradeId, errorTradeId := strconv.Atoi(os.Args[1])
			if(errorTradeId != nil) {
				fmt.Println("Entered TradeId is not a valid number! Please enter a valid TradeId.")
			} else if(tradeId > 2){
				fmt.Println("This System allows only latest 3 responses to be checked from Server (Stored In-Memory). The TradeID you requested is the wrong one. Try these - 0 or 1 or 2!")
			} else {
					var response s.CheckingPortfolioServiceResponse
					if err := ExecuteJSONRPCRequest("CheckingPortfolioService.CheckPortfolio", 
						&s.CheckingPortfolioServiceRequest{uint64(tradeId)}, &response); err != nil {
						fmt.Println("Error :",err)
					}
			}	
		}
	}
}


// ----------------------------------------------------------------------------
// Validations for Command Line Arguments
// ----------------------------------------------------------------------------

//Function to validate if client has passed the valid command line argument for 'Budget' field
func ValidateBudget(arg string) bool {
	var isInvalid bool
	budget, errorFloat := strconv.ParseFloat(os.Args[1], 32)
	if(errorFloat != nil) {
		fmt.Println("The entered budget is not in NUMBER format! ")
		isInvalid = true;
	} else if(budget == 0){
		fmt.Println("The entered ZERO as your 'Budget'! ")
		isInvalid = true;
	} else if(budget < 0){
		fmt.Println("The entered 'Budget' is less than 0! ")
		isInvalid = true;
	}		
	return isInvalid
}


//Function to validate if client has passed the valid command line argument for 'StockSymbol and Percentage' field
func ValidateStockSymbolAndPercentage(arg string) bool {
	var isInvalid bool
	isInvalid = false
	if(strings.TrimSpace(arg) == "") {
		fmt.Println("You entered an empty string!")
		isInvalid = true
	} else if !strings.Contains(arg, ":") || !strings.Contains(arg,"%") {
		fmt.Println("Invalid String!")
		isInvalid = true
	}
	return isInvalid
}


//Function to Validate if client has entered total sum of 100% for all the stock symbols requested
//OR entered any of the Invalid Stock Symbols
func ValidateEnteredPercentageFor100OrInvalidSymbols(arg string) bool {
    var symbolsAndpercentages []string 
    var tempslices []string
    symbolsAndpercentages = strings.Split(arg, ",")
	var sum float32
	sum = 0.0
	
	var symbols = make([]string, len(symbolsAndpercentages))
    var percentages = make([]float32, len(symbolsAndpercentages))
    for index := 0; index < len(symbolsAndpercentages); index++ {
 		//Split every stocksymbol and percentage and put them into separate arrays/slices
 		tempslices = strings.Split(symbolsAndpercentages[index],":")
 		tempslices[1] = strings.Trim(tempslices[1],"%")
 		percentageForThisStock, _ := strconv.ParseFloat(tempslices[1], 32)
 		symbols[index] = tempslices[0]
 		percentages[index] = float32(percentageForThisStock) //tempSlices[1]
		sum += percentages[index]
	}
	if(sum < 100.0 || sum > 100.0) {
		fmt.Println("The entered breakage of percentages does not meet the summation of 100%! Please enter the correct percentages.")
		return true		
	} else {
		_, invalidSymbols := s.CallYahooAPI(symbols)
		if(invalidSymbols) {
			fmt.Println("You entered (one or more) invalid Symbols! Please check them again.")
			return true
		}
	}	
	return false
}

