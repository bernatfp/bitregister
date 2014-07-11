package main

import (
	"net/http"
	"net/url"
	"log"
	"encoding/json"
	"strings"
	"fmt"
	"errors"
)


//debugging
var _ = url.Parse
var _ = errors.New
var _ = strings.Join


func retrieveOrders() []byte {
	return []byte{}
}

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

func updateOrder(req *http.Request, id string) []byte {
	return []byte{}
}

func removeOrder(id string) []byte {
	return []byte{}
}

//Parses order ID from the request if included
func parseId(req *http.Request) (string, error) {
	var id string
	var err error

	params := strings.Split(req.URL.Path, "/")
	idIndex := len(params) - 1
	
	//if url is /orders/ or /orders/<id>/ (instead of /orders or /orders/<id>)
	if params[idIndex] == "" {
		idIndex--
	}

	switch idIndex {
		// url is /orders
		case 1:
			id = ""
		// url is /orders/<id>
		case 2:
			id = params[idIndex]
		default:
			id = ""
			err = errors.New("Orders URL is malformed")
	}

	return id, err
}

// Resource: orders
//
// Methods accepted: GET, POST, DELETE
//
// GET /orders/ => lists all orders (accepts filters as part of the querystring)
// GET /orders/<id> => returns order <id>
// POST /orders/ => creates a new order
// POST /orders/<id> => updates order <id>
// DELETE /orders/<id> => deletes order <id>
//
// This is the orders handler function
func ordersHandle(w http.ResponseWriter, req *http.Request) {	
	var data []byte

	id, err := parseId(req)

	if err != nil {
		w.Write([]byte(fmt.Sprintf("%v",err)))
		return
	}

	switch {
		case req.Method == "GET" && id == "":
			data = retrieveOrders()
		
		case req.Method == "GET":
			data = retrieveOrder(id)
		
		case req.Method == "POST" && id == "":
			data = createOrder(req)
		
		case req.Method == "POST":
			data = updateOrder(req, id)
		
		case req.Method == "DELETE":
			data = removeOrder(id)
	}

	w.Write(data)
}

//Debugging
func faviconHandle(w http.ResponseWriter, req *http.Request) {
	w.Write(nil)
}

//Debugging BTC commands
func rootHandle(w http.ResponseWriter, req *http.Request) {

	reply, err := sendCommand()
	if err != nil {
		log.Println("Error sending command: ", err)
	}

	data, err := json.Marshal(reply)
	if err != nil {
		log.Println("Can't marshal reply: ", err)
	}

	w.Write(data)

}