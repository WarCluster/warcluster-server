package server

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fzzy/sockjs-go/sockjs"

	"warcluster/config"
	"warcluster/entities/db"
	)

type testSession struct {
	session_id string
	Messages   [][]byte
}

func (s *testSession) Receive() (m []byte) {
	defer func() {
		if(len(s.Messages) == 0) {
			s.Messages = s.Messages[1:]
		}
	}()
	return s.Messages[0]
}

func (s *testSession) Send(m []byte) {
	s.Messages = append(s.Messages, m)
	return
}

func (s *testSession) Close(code int, reason string) {
	return
}

func (s *testSession) End() {
	return
}

func (s *testSession) Info() sockjs.RequestInfo {
	return *new(sockjs.RequestInfo)
}

func (s *testSession) Protocol() sockjs.Protocol {
	return *new(sockjs.Protocol)
}

func (s *testSession) String() string {
	rand.Seed(time.Now().Unix())
	return fmt.Sprintf("%d", rand.Int())
	// return session_id
}

func init() {
	var cfg config.Config
	cfg.Load("../config/config.gcfg")
	db.InitPool(cfg.Database.Host, cfg.Database.Port)
	conn := db.Pool.Get()
	defer conn.Close()
	conn.Do("SELECT", 13)
}
