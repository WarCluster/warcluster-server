// db is the package responsible for database maintenence and managing DB I/O.
package db

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"sync"
	"time"
)

// connection is the pointer to the db used for db comunication.
var connection redis.Conn

// Pool maintains a pool of connections to the database
var pool redis.Pool

// This function is called in order to insure propper db acsess.
// It creates the DB connection and stores it in the connection variable.
func NewPool(host string, port int) {
	var err error
	log.Print("Initializing database connection... ")
	serverAddr := fmt.Sprintf("%v:%v", host, port)

	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			var err error
			connection, err = redis.Dial("tcp", serverAddr)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}
			return connection, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

}

// Finalize is called upon the death of the server(intended or not :)
// it ensures the propper closing of the DB connection.
func Finalize() {
	log.Print("Closing database connection... ")
	if err := connection.Close(); err != nil {
		log.Fatal(err)
	}
}
