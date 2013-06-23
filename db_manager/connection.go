/*
	DB_manager is the package responsible for database maintenence and managing DB I/O.
*/
package db_manager

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"sync"
)

//connection is the pointer to the db used for db comunication.
var connection redis.Conn

//I think we all know that mutex is like I/O trafic light needed to avoid really bad stuff
var mutex sync.Mutex

/*
This function is called in order to insure propper db acsess.
It creates the DB connection and stores it in the connection variable.
*/
func Connect(network, host string, port int) {
	var err error
	log.Print("Initializing database connection... ")
	if connection, err = redis.Dial(network, fmt.Sprintf("%v:%v", host, port)); err != nil {
		log.Fatal(err)
	}
}

/*
Finalize is called upon the death of the server(intended or not :)
it ensures the propper closing of the DB connection.
*/
func Finalize() {
	log.Print("Closing database connection... ")
	if err := connection.Close(); err != nil {
		log.Fatal(err)
	}
}
