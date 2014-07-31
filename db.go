//TO-DO
//
// * Optionally select index at the start 
// * Implement get and set for sets and hashes
//



package main

import (
	"github.com/garyburd/redigo/redis"
	"log"
)


//DB is a small wrapper for Redis that initializes a thread safe pool of connections
//Contains functions for calling needed operations (GET, SET...)
type DB struct {
	pool *redis.Pool
}


//A simple example:
//
// db := new(DB)
// db.init()
// defer db.close()
// 
// result, err := db.get("foo")
// if err != nil {
//   log.Println(err)	
// }
// 
// log.Println(result)
//


//Returns a new connection
//It is called by the pool when we ask for a new connection
func dial() (redis.Conn, error){	
	c, err := redis.Dial("tcp", ":6379")
	
	if err != nil {
		log.Println("Error connecting to Redis: ", err)
	}

	return c, err
}

//Initializes the pool of connections
//Must be called to initialize pool
func (db *DB) init() {	
	db.pool = redis.NewPool(dial, 3)
}

//Closes the pool of connections
//Must be called before finishing
func (db *DB) close() error {
	return db.pool.Close()
}


//Converts arguments for Redis command
func convertArgs(args []interface{}) redis.Args {
	switch v := args[0].(type) {
		case string:
			return redis.Args{}.Add(v)
		
		case []string:
			redisArgs := redis.Args{}
			for _, arg := range v {
				redisArgs = redisArgs.Add(arg)
			}
			return redisArgs

		case *Order:
			redisArgs := redis.Args{}
			return redisArgs.Add(v.Id).AddFlat(v)

		case redis.Args:
			return v
	}

	return nil

}


//Generic function that does the appropiate handling to send commands to Redis
//Returns the response as received and must be manipulated by the caller
func (db *DB) sendCommand(command string, args ...interface{}) (interface{}, error) {
	//Ask for a new connection
	c := db.pool.Get()
	defer c.Close()

	//generate redis compatible arguments
	redisArgs := convertArgs(args)

	//Send command
	res, err := c.Do(command, redisArgs...)
	if err != nil {
		log.Println("Error sending command: ", err)
	}

	return res, err
}

//SET operation
//Receives a key and a value
func (db *DB) set(k string, v string) error {
	_, err := db.sendCommand("SET", k, v)
	if err != nil {
		log.Println("Error sending SET operation: ", err)
	}

	return err
}

//GET operation
//Receives a key and returns a value
func (db *DB) get(k string) (string, error) {	
	res, err := db.sendCommand("GET", k)
	if err != nil {
		//Happens when getted value is empty
		log.Println("Error sending GET operation: ", err)
	}

	return redis.String(res, err)
}

//GET operation
//Receives a key and returns a value
func (db *DB) del(k string) error {
	_, err := db.sendCommand("DEL", k)
	if err != nil {
		//Happens when getted value is empty
		log.Println("Error sending DEL operation: ", err)
	}

	return err
}


//The following are Redis functions adapted for Order
//Useful to code additional behavior between DB and orders layer
//HGETALL for type Order
func (db *DB) retrieveOrder(id string) (*Order, error){
	//res, err := c.Do("HGETALL", id)
	res, err := db.sendCommand("HGETALL", id)
	if err != nil {
		log.Println("Error sending HGETALL operation: ", err)
	}
	
	v, err := redis.Values(res, err)
    if err != nil {
        panic(err)
    }

    order := new(Order)
    if err := redis.ScanStruct(v, order); err != nil {
        panic(err)
    }

    return order, err
}

//HMSET for type Order
func (db *DB) insertOrder(order *Order) error {
	_, err := db.sendCommand("HMSET", order)
	if err != nil {
		log.Println("Error sending HMSET operation: ", err)
	}

	return err
}

//DEL for type Order (not very useful right now, I know)
func (db *DB) removeOrder(id string) error {
	err := db.del(id) 
	if err != nil {
		log.Println("Error deleting order: ", err)
	}

	return err
}

