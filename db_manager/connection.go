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

//I think we all know what mutex is but for those that dont(mutex is like I/O trafic light needed to avoid really bad stuff)
var mutex sync.Mutex

//information for the DB connection
const (
	HOSTNAME = "localhost"
	PORT     = 6379
	NETWORK  = "tcp"
)

/*
This function is called on the first package import in order to insure propper db acsess.
*/
func init() {
	var err error
	log.Print("Initializing database connection... ")
	if connection, err = connect(); err != nil {
		log.Fatal(err)
	}
}

/*
connect is called from package init.
The function creates the DB connection and stores it in the connection variable.
*/
func connect() (redis.Conn, error) {
	return redis.Dial("tcp", fmt.Sprintf("%v:%v", HOSTNAME, PORT))
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
