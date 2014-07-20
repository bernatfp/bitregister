package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"github.com/garyburd/redigo/redis"
)

var bitcoinRates *Rates


//debugging, remove afterwards
var _ = json.Marshal
var _ = fmt.Println
var _ = strings.Split
var _ = url.Parse
var _ = redis.Dial


func main() {

	//Corroutine updates Bitcoin rates every minute
	bitcoinRates = &Rates{}
	go updateRates(bitcoinRates)

	//Register HTTP server handlers
	http.HandleFunc("/", rootHandle) //root handle isn't expected to do anything, at the moment it's used for debugging
	http.HandleFunc("/favicon.ico", faviconHandle) //temporary
	http.HandleFunc("/orders/", ordersHandle) //Orders resource
	
	log.Println("Starting server on port 12345...")
	
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}


}
