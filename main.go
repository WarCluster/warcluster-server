// Real-time massively multiplayer online space strategy arcade browser game!
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
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT)
	signal.Notify(exitChan, syscall.SIGKILL)
	signal.Notify(exitChan, syscall.SIGTERM)
	<-exitChan

	server.Stop()
	os.Exit(0)
}
