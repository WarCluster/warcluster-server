package main

import (
	"./db_manager"
	net "./server"
	"os"
	"os/signal"
)

var server *net.Server = &net.Server{}

func main() {
	go final()
	defer func() {
		final()
	}()

	// TODO: Create goroutines for parsing mission
	// and give this instance to the server

	// This goroutine should:
	// * have a sorted list with all the missions
	// *  have direct connection to the db
	// * Add(*Mission) error
	// * delete(mission_key string)
	// * Run()/Loop()/... endless goroutine for 
	//   calculating the mission results
	//
	// TODO: Think for a name of this thing


	server.Start("localhost", 7000)
	return
}

func final() {
	sigtermchan := make(chan os.Signal, 1)
	signal.Notify(sigtermchan, os.Interrupt)
	<-sigtermchan

	server.Stop()
	db_manager.Finalize()
}
