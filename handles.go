package main

import (
	"net/http"
	"net/url"
	"log"
	"encoding/json"
	"strings"
)

//debugging
var _ = url.Parse

//retrieve all pending orders
func pendingOrdersHandle(w http.ResponseWriter, req *http.Request) {

	w.Write(nil)
}

//retrieve all completed orders
func completedOrdersHandle(w http.ResponseWriter, req *http.Request) {

	w.Write(nil)
}

//retrieve an order by id or create a new one
func idOrdersHandle(w http.ResponseWriter, req *http.Request) {

	params := strings.Split(req.URL.Path, "/")
	idIndex := len(params) - 1
	
	if params[idIndex] == "" {
		idIndex--
	}
	if idIndex != 3 || params[idIndex] == "" {
		w.Write([]byte("Error, URL must follow this form: /orders/id/<id_value>"))
		return
	}
	
	id := params[idIndex]
	var _ = id

	var data []byte

	switch req.Method {
		case "GET":
			//retrieve an order
			data = []byte("GET id: " + id)

		case "POST":
			//create an order with this id or update it
			data = make([]byte, 0)

		default:
			//unknown method, return error
			data = make([]byte, 0)

	}

	w.Write(data)
}

//retrieve all orders
func ordersHandle(w http.ResponseWriter, req *http.Request) {

	w.Write(nil)
}

func faviconHandle(w http.ResponseWriter, req *http.Request) {
	w.Write(nil)
}

func rootHandle(w http.ResponseWriter, req *http.Request) {

	reply, err:= sendCommand()
	if err != nil {
		log.Println("Error sending command: ", err)
	}

	data, err := json.Marshal(reply)
	if err != nil {
		log.Println("Can't marshal reply: ", err)
	}

	w.Write(data)

}