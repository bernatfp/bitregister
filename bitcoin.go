//Operations to implement
//
// * Move between accounts (from each order to master account)
//   move <from_account> <to_account> <amount>  
//
// * Send transaction
//   sendtoaddress <bitcoinaddress> <amount> [comment] [comment-to]





package main

import (
	"net/http"
	"encoding/json"
	"log"
	"io/ioutil"
	"time"
	"github.com/conformal/btcjson"
	"errors"
)


//Temporary, must be loaded from config file
const (
	user     = "bitcoinrpc"
	password = "HNZX6G4hvd1aUUNqVPadV1KdWmfeySCa3gtmh2tmZ5hX"
	server   = "127.0.0.1:8332"
)

type Rates struct {
	ToEUR string `json:"btc_to_eur"`
	ToGBP string `json:"btc_to_gbp"`
	ToUSD string `json:"btc_to_usd"`
}

//Returns current Coinbase rates for EUR, GBP and USD
func updateRates(rates *Rates) {
	for {
		//HTTP request
		resp, err := http.Get("https://coinbase.com/api/v1/currencies/exchange_rates")
		if err != nil {
			log.Fatal("Error sending GET: ", err)
			break
		}
		defer resp.Body.Close()

		//Read body and decode JSON
		body, err := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, rates)
		if err != nil {
			log.Fatal("Error decoding JSON: ", err)
			break
		}

		log.Println("Updated Bitcoin rates: ", *rates)

		//60-sec update period to avoid being banned
		time.Sleep(time.Minute)

	}
	
}

//Sends commands to bitcoind
func sendCommand(command string, argsStr ...string) (btcjson.Reply, error) {
	//Type conversion
	//args can't be of type []string, it must be []interface{}
	args := make([]interface{}, len(argsStr))
	for i, v := range argsStr {
		args[i] = v
	}

	//Craft message and send to bitcoind
	msg, err := btcjson.CreateMessage(command, args...)
	reply, err := btcjson.RpcCommand(user, password, server, msg)

	return reply, err
}

//Creates an address for the id provided if it does not exist, otherwise returns existing address
func createAddress(id string) (string, error) {
	reply, err := sendCommand("getaccountaddress", id)
	if err != nil {
		log.Println("Error with getaccountaddress command: ", err)
	}

	v, ok := reply.Result.(string)
	if ok != true {
		log.Println("Error converting address to string")
		err = errors.New("Error converting address to string")
	}

	return v, err
}

//Obtains the account's current balance
func getBalance(id string) (float64, error) {
	reply, err := sendCommand("getbalance", id)
	if err != nil {
		log.Println("Error with getbalance command: ", err)
	}

	v, ok := reply.Result.(float64)
	if ok != true {
		log.Println("Error converting address to float64")
		err = errors.New("Error converting address to float64")
	}

	return v, err
}

//Obtain total balance of the wallet
func getTotalBalance() (float64, error) {
	//No account provided returns total balance
	return getBalance("")
}


