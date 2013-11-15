// Package db is responsible for database maintenence and managing DB I/O.
package db

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

// Pool maintains a pool of connections to the database
var Pool *redis.Pool

// This function is called in order to insure propper db acsess.
// It creates the DB connection and stores it in the connection variable.
func InitPool(host string, port uint16) {
	log.Print("Initializing database connection... ")
	serverAddr := fmt.Sprintf("%v:%v", host, port)

	Pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", serverAddr)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}
			return conn, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
