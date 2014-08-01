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
	Filled bool `json:"filled" redis:"filled"`
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
	order, err := parseOrder(req)
	order.convertPrice()

	order.assignAddress()
	order.Filled = false

	db := new(DB)
	db.init()
	defer db.close()

	err = db.insertOrder(order)
	if err != nil {
		log.Println("Error insert order: ", err)
		return []byte("Error insert order")	
	}
	
	order.populatePriceJson()
	data, err := json.Marshal(order)	
	if err != nil {
		log.Println("Error marshal order: ", err)
	}

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

//Redis hashes don't support multiple levels so we must include Price fields as Order fields
func (order *Order) populatePriceRedis(price string) {
	p, err := strconv.ParseFloat(price, 64)
	if err != nil {
		log.Println(fmt.Sprintf("Error parse float: ", err))
	}
	order.Currency = "BTC"
	order.Amount = order.Price.Amount / p
}

//Vice versa
func (order *Order) populatePriceJson() {
	order.Price.Amount = order.Amount
	order.Price.Currency = "BTC"
}

func (order *Order) assignAddress() error {
	addr, err := createAddress(order.Id)
	if err != nil {
		log.Println("Error assigning address to order: ", err)
	}

	order.Address = addr

	return err
}


