package main

import (
	"os"
	"os/signal"
	"syscall"
	"warcluster/config"
	"warcluster/entities/db"
	"warcluster/server"
)

var cfg config.Config

func main() {
	go final()

	cfg.Load("config/config.gcfg")
	db.InitPool(cfg.Database.Host, cfg.Database.Port)
	server.Start(cfg.Server.Host, cfg.Server.Port)
}

func final() {
	exit_chan := make(chan os.Signal, 1)
	signal.Notify(exit_chan, syscall.SIGINT)
	signal.Notify(exit_chan, syscall.SIGKILL)
	signal.Notify(exit_chan, syscall.SIGTERM)
	<-exit_chan

	server.Stop()
	os.Exit(0)
}
