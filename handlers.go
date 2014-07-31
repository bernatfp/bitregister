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


//Helper functions

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
// GET /orders/<id>/ => returns order <id>
// POST /orders/ => creates a new order
// DELETE /orders/<id>/ => deletes order <id>
//

// Orders HTTP Handler
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
		
		case req.Method == "DELETE":
			data = removeOrder(id)
	}

	w.Write(data)
}

// Debugging
func faviconHandle(w http.ResponseWriter, req *http.Request) {
	w.Write(nil)
}

// Debugging BTC commands
func rootHandle(w http.ResponseWriter, req *http.Request) {
	reply, err := sendCommand()
	if err != nil {
		log.Println("Error sending command: ", err)
	}

	_ = reply
	//data, err := json.Marshal(reply)
	data, err := json.Marshal(bitcoinRates)
	if err != nil {
		log.Println("Can't marshal reply: ", err)
	}

	w.Write(data)

}