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
