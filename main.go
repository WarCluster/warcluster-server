package main

import (
	"os"
	"os/signal"
	"warcluster/config"
	"warcluster/entities/db"
	"warcluster/server"
)

var cfg config.Config

func main() {
	go final()
	defer final()

	cfg.Load("config/config.gcfg")
	db.Connect(cfg.Database.Network, cfg.Database.Host, cfg.Database.Port)
	server.Start(cfg.Server.Host, cfg.Server.Port)
}

func final() {
	sigtermchan := make(chan os.Signal, 1)
	signal.Notify(sigtermchan, os.Interrupt)
	<-sigtermchan

	db.Finalize()
	server.Stop()
	os.Exit(0)
}
