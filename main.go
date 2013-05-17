package main

import (
	"os"
	"os/signal"
	"warcluster/config"
	"warcluster/db_manager"
	"warcluster/server"
)

var cfg config.Config

func main() {
	go final()
	defer func() {
		final()
	}()

	cfg.Load("config/config.gcfg")
	db_manager.Connect(cfg.Database.Network, cfg.Database.Host, cfg.Database.Port)
	server.Start(cfg.Server.Host, cfg.Server.Port)
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
