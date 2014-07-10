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

//Generic function that does the appropiate handling to send commands to Redis
//Returns the response as received and must be manipulated by the caller
func (db *DB) sendCommand(command string, strArgs ...string) (interface{}, error) {
	//Ask for a new connection
	c := db.pool.Get()
	defer c.Close()

	//Convert command arguments because they must be of type []interface{}
	args := redis.Args{}
	for _, arg := range strArgs {
		args = args.Add(arg)
	}

	//Send command
	res, err := c.Do(command, args...)
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

