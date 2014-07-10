package main

import (
	"encoding/json"
	"fmt"
	"github.com/conformal/btcjson"
	"log"
	"net/http"
	"net/url"
	"strings"
	"github.com/garyburd/redigo/redis"
)

const (
	user     = "bitcoinrpc"
	password = "HNZX6G4hvd1aUUNqVPadV1KdWmfeySCa3gtmh2tmZ5hX"
	server   = "127.0.0.1:8332"
)

/*
This server uses Redis for data persistence at this point

Other options to consider:

- Use BerkeleyDB (already installed with any bitcoind) for data storage
	Libraries:
	- https://bitbucket.org/chust/goprotodb
	  The problem with this one is that any user will need to install protobuffers

	- https://github.com/pebbe/dbxml
	  The problem with this one is that it's not clear whether BDB XML comes with any BDB install or it is totally separated

- Alternatively, we can use an embedded Go data store:
	- https://github.com/peterbourgon/diskv
	- https://github.com/steveyen/gkvlite
	- https://github.com/HouzuoGuo/tiedot
	- https://github.com/boltdb/bolt
*/


type Order struct {
	id string `json:"id"`
	amount float64 `json:"amount"`
	currency string `json:"currency"`
}

func sendCommand() (btcjson.Reply, error) {
	msg, err := btcjson.CreateMessage("getinfo")
	reply, err := btcjson.RpcCommand(user, password, server, msg)

	return reply, err
}


//debugging, remove afterwards
var _ = json.Marshal
var _ = fmt.Println
var _ = strings.Split
var _ = url.Parse
var _ = redis.Dial


func main() {

	log.Println("Starting server on port 12345...")

	//Register HTTP server handlers
	http.HandleFunc("/", rootHandle)
	http.HandleFunc("/favicon.ico", faviconHandle)
	http.HandleFunc("/orders/", ordersHandle)
	http.HandleFunc("/orders/pending/", pendingOrdersHandle)
	http.HandleFunc("/orders/completed/", completedOrdersHandle)
	http.HandleFunc("/orders/id/", idOrdersHandle)	

	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
