package server

import (
	"time"

	"warcluster/config"
	"warcluster/entities/db"

	"code.google.com/p/go.net/websocket"
)

var testServer *Server

func init() {
	var cfg config.Config
	cfg.Load()
	db.InitPool(cfg.Database.Host, cfg.Database.Port, 13)
	conn := db.Pool.Get()
	defer conn.Close()

	conn.Do("FLUSHDB")
	testServer = NewServer(
		cfg.Server.Host,
		7013,
		Handle,
	)

	go testServer.Start()
	for !testServer.isRunning {
		time.Sleep(100 * time.Millisecond)
	}
}

func testHandler(*websocket.Conn) {
}
