// Real-time massively multiplayer online space strategy arcade browser game!
package main

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"warcluster/config"
	"warcluster/entities/db"
	"warcluster/leaderboard"
	"warcluster/server"
)

var cfg config.Config

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	cfg.Load()
	db.InitPool(cfg.Database.Host, cfg.Database.Port, 8)
	server.ExportConfig(cfg)
	server.InitLeaderboard(leaderboard.New())
	server.SpawnDbMissions()

	s := server.NewServer(
		cfg.Server.Host,
		cfg.Server.Port,
		server.Handle,
	)
	go final(s)

	s.Start()
}

func final(s *server.Server) {
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT)
	signal.Notify(exitChan, syscall.SIGKILL)
	signal.Notify(exitChan, syscall.SIGTERM)
	<-exitChan

	s.Stop()
	os.Exit(0)
}
