package main

import (
	"log"
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
)

type Order struct {
	Id string `json:"id" redis:"id"`
	Address string `json:"address" redis:"address"`
	Status string `json:"status" redis:"status"`
	//This one is ignored by Redis as declared types are not supported
	Price PriceType `json:"price" redis:"-"`
	//To fix this issue we declare amount and currency and tell json to ignore it
	Amount float64 `json:"-" redis:"amount"`
	Currency string `json:"-" redis:"currency"`
}

type PriceType struct {
	Amount float64 `json:"amount" redis:"amount"`
	Currency string `json:"currency" redis:"currency"`
}

// GET /orders/
func retrieveOrders() []byte {
	return []byte{}
}

// GET /orders/<id>/
func retrieveOrder(id string) []byte {
	var data []byte
	
	db := new(DB)
	db.init()
	defer db.close()

	order, err := db.retrieveOrder(id)
	if err != nil {
		log.Println(fmt.Sprintf("Error retrieve order: ", err))
		data = []byte("Error, this id does not exist")
	} else {
		//OK
		order.populatePriceJson()
		data, err = json.Marshal(order)
		if err != nil {
			log.Println(fmt.Sprintf("Error marshaling: ", err))
		}
	}

	return data
}

// POST /orders/
func createOrder(req *http.Request) []byte {
	var data []byte

	order, err := parseOrder(req)
	order.convertPrice()

	db := new(DB)
	db.init()
	defer db.close()

	err = db.insertOrder(order)
	if err != nil {
		log.Println("Error set: ", err)
		data = []byte("Error SET")	
	} else {
		data = []byte("Order stored")	
	}

	//Have to return payment data (id, address and amount in BTC)
	return data

}


// DELETE /orders/<id>/
func removeOrder(id string) []byte {
	var data []byte

	db := new(DB)
	db.init()
	defer db.close()

	err := db.removeOrder(id)

	//TO-DO return status code and message according to err
	if err != nil {
		log.Println(fmt.Sprintf("Error DEL id: %s\n", id))
		data = []byte("Error deleting id")
	} else {
		data = []byte("OK")
	}

	return data
}

//Helpers
func parseOrder(req *http.Request) (*Order, error) {
	
	//using decoder instead of unmarshal
	dec := json.NewDecoder(req.Body)
	order := new(Order)
	err := dec.Decode(order)
	
	return order, err
}

func (order *Order) convertPrice() {
	switch order.Price.Currency {
		case "EUR":
			order.populatePriceRedis(bitcoinRates.ToEUR)
		case "GBP":
			order.populatePriceRedis(bitcoinRates.ToGBP)
		case "USD":
			order.populatePriceRedis(bitcoinRates.ToUSD)
	}
}

func (order *Order) populatePriceRedis(price string) {
	p, err := strconv.ParseFloat(price, 64)
	if err != nil {
		log.Println(fmt.Sprintf("Error parse float: ", err))
	}
	order.Currency = "BTC"
	order.Amount = order.Price.Amount / p
}

func (order *Order) populatePriceJson() {
	order.Price.Amount = order.Amount
	order.Price.Currency = "BTC"
}


