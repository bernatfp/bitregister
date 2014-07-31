package main

import (
	"log"
	"fmt"
	"net/http"
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

	value, err := db.get(id)
	if err != nil {
		log.Println(fmt.Sprintf("Error GET id: %s\n Result: %s", id, value))
		data = []byte("Error, this id does not exist")
	} else {
		data = []byte(value)
	}

	return data
}

// POST /orders/
func createOrder(req *http.Request) []byte {
	var data []byte

	db := new(DB)
	db.init()
	defer db.close()

	err := req.ParseForm()
	if err != nil {
		log.Println("Error ParseForm: ", err)
	}

	key := req.FormValue("id")
	value := req.FormValue("value")

	err = db.set(key, value)
	if err != nil {
		log.Println("Error set: ", err)
		data = []byte("Error SET")	
	} else {
		data = []byte("Order stored")	
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
