package main

import (
	"net/http"
	"encoding/json"
	"log"
	"io/ioutil"
)

type Price struct {
	ToEUR string `json:"btc_to_eur"`
	ToGBP string `json:"btc_to_gbp"`
	ToUSD string `json:"btc_to_usd"`
}

//Returns current Coinbase rates for EUR, GBP and USD
func updatePrice() (*Price, error) {
	price := &Price{}
	
	resp, err := http.Get("https://coinbase.com/api/v1/currencies/exchange_rates")
	if err != nil {
		log.Println("Error sending GET: ", err)
		return price, nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, price)
	if err != nil {
		log.Println("Error decoding JSON: ", err)
		return price, err
	}

	return price, nil
}