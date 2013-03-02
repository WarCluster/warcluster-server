package db_manager

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"sync"
)

var connection redis.Conn
var mutex sync.Mutex

const (
	HOSTNAME = "localhost"
	PORT     = 6379
	NETWORK  = "tcp"
)

func init() {
	var err error
	log.Print("Initializing database connection... ")
	if connection, err = connect(); err != nil {
		log.Fatal(err)
	}
}

func connect() (redis.Conn, error) {
	return redis.Dial("tcp", fmt.Sprintf("%v:%v", HOSTNAME, PORT))
}

func Finalize() {
	log.Print("Closing database connection... ")
	if err := connection.Close(); err != nil {
		log.Fatal(err)
	}
}
