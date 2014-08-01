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
		resp, err := http.Get("https://coinbase.com/api/v1/currencies/exchange_rates")
		if err != nil {
			log.Fatal("Error sending GET: ", err)
			break
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(body, rates)
		if err != nil {
			log.Fatal("Error decoding JSON: ", err)
			break
		}

		log.Println("Updated Bitcoin rates: ", *rates)

		time.Sleep(time.Minute)

	}
	
}


//Sends commands to bitcoind
func sendCommand(command string, argsStr ...string) (btcjson.Reply, error) {

	args := make([]interface{}, len(argsStr))
	for i, v := range argsStr {
		args[i] = v
	}

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




