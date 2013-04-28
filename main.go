package main

import (
	"./db_manager"
	"./server"
	"os"
	"os/signal"
)

func main() {
	go final()
	defer func() {
		final()
	}()

	server.Start("0.0.0.0", 7000)
	return
}

func final() {
	sigtermchan := make(chan os.Signal, 1)
	signal.Notify(sigtermchan, os.Interrupt)
	<-sigtermchan

	db_manager.Finalize()
	server.Stop()
	os.Exit(0)
}
