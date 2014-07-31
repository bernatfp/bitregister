package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"github.com/garyburd/redigo/redis"
	"time"
	"bytes"
	"io/ioutil"
)

var bitcoinRates *Rates


//debugging, remove afterwards
var _ = json.Marshal
var _ = fmt.Println
var _ = strings.Split
var _ = url.Parse
var _ = redis.Dial
var _ = ioutil.ReadAll
var _ = bytes.NewReader


func main() {

	go testbox()

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


//Debugging requests here
func testbox() {
	time.Sleep(time.Second * 5)

	log.Println("Sending commands...")

	//test create order
	log.Println("Create order...")
	o := &Order{Id: "111", Price: PriceType{Amount: 1000, Currency: "USD"}}
	buf, err := json.Marshal(o)
	if err != nil {
		log.Println("Marshal: ", err)
	}

	log.Println(string(buf))

	r := bytes.NewReader(buf)
	resp, err := http.Post("http://localhost:12345/orders/", "application/json", r)
	if err != nil {
		log.Println("Post: ", err)
	}
	defer resp.Body.Close()

	log.Println("Post sent...")


	//test retrieve order
	log.Println("Retrieve order...")
	res, err := http.Get("http://localhost:12345/orders/111")
	if err != nil {
		log.Println(err)
	}

	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%s\n", content)

	//test delete order
	log.Println("Delete order...")
	
	c := &http.Client{}
	//req := &http.Request{Method: "DELETE", URL: &url.URL{Scheme: "http", Host: "localhost:12345", Path: "orders/111"}}
	req, err := http.NewRequest("DELETE", "http://localhost:12345/orders/111", nil)
	respo, err := c.Do(req)
	log.Println("Delete order SENT...")
	cont, err := ioutil.ReadAll(respo.Body)
	respo.Body.Close()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%s", cont)

}