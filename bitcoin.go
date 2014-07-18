package main

import (
	"net/http"
	"encoding/json"
	"log"
	"io/ioutil"
	"fmt"
)

type Price struct {
	toEUR string `json:"btc_to_eur"`
	toGBP string `json:"btc_to_gbp"`
	toUSD string `json:"btc_to_usd"`
	btc_to_usd string
	test bool
}

func updatePrice() (*Price, error) {
	
	resp, err := http.Get("https://coinbase.com/api/v1/currencies/exchange_rates")
	if err != nil {
		log.Println("Error sending GET: ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var price Price
	err = json.Unmarshal(body, &price)
	if err != nil {
		log.Println("Error decoding JSON: ", err)
	}

	price.test = true;

	fmt.Printf("%s\n", string(body))
	fmt.Printf("%+v\n", price)

	return &price, nil

}