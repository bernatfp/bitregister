package main

import (
	"log"
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
)

//Every order must contain this information
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

//Struct returned when retrieving all orders
type Orders struct {
	Filled []string `json:"filled"`
	Pending []string `json:"pending"`
}

// GET /orders/
func retrieveOrders() []byte {
	
	//Start DB connection
	db := new(DB)
	db.init()
	defer db.close()

	orders := new(Orders)

	//Keys "bitregister:orders:*"

	//Loop
		//goroutine calls HGET of id to retrieve field filled

	//Wait until all goroutines are finished and return marshaled result

	data, err := json.Marshal(orders)
	if err != nil {
		log.Println(fmt.Sprintf("Error marshaling: ", err))
	}

	return data
}

// GET /orders/<id>/
func retrieveOrder(id string) []byte {
	var data []byte
	
	//Start DB connection
	db := new(DB)
	db.init()
	defer db.close()

	//Get order info
	order, err := db.retrieveOrder(id)
	if err != nil {
		log.Println(fmt.Sprintf("Error retrieve order: ", err))
		return []byte("Error, this id does not exist")
	} 
	
	//Check the amount of money this order has received
	amount, err := getBalance(order.Id)
	if err != nil {
		log.Println("Error getBalance of order: ", err)
	}

	//Check if amount received is enough
	if amount >= order.Amount {
		order.Filled = true
		//Update info on background
		go db.insertOrder(order)
	}

	//Format order ot be returned as JSON
	order.populatePriceJson()
	data, err = json.Marshal(order)
	if err != nil {
		log.Println(fmt.Sprintf("Error marshaling: ", err))
	}

	return data
}

// POST /orders/
func createOrder(req *http.Request) []byte {
	//Parse order from request
	order, err := parseOrder(req)

	//Start DB
	db := new(DB)
	db.init()
	defer db.close()

	//If order already exists, return stored order
	if ok, err := db.existsOrder(order.Id); err != nil {
		log.Println("Error check order exists: ", err)
		return []byte("Error check order exists")
	} else if ok {
		return retrieveOrder(order.Id)
	}

	//Complete order information
	order.convertPrice()
	order.assignAddress()
	order.Filled = false

	//Store order
	err = db.insertOrder(order)
	if err != nil {
		log.Println("Error insert order: ", err)
		return []byte("Error insert order")	
	}
	
	//Update bitcoin price on order JSON
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

	//Start DB
	db := new(DB)
	db.init()
	defer db.close()

	err := db.removeOrder(id)
	if err != nil {
		log.Println(fmt.Sprintf("Error DEL id: %s\n", id))
		return []byte("Error deleting id")
	} else {
		data = []byte("Order deleted")
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

//Creates an address and assigns it to the order
func (order *Order) assignAddress() error {
	addr, err := createAddress(order.Id)
	if err != nil {
		log.Println("Error assigning address to order: ", err)
	}

	order.Address = addr

	return err
}


