package server

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fzzy/sockjs-go/sockjs"

	"warcluster/config"
	"warcluster/entities/db"
)

func init() {
	var cfg config.Config
	cfg.Load("../config/config.gcfg")
	db.InitPool(cfg.Database.Host, cfg.Database.Port, 13)
	conn := db.Pool.Get()
	defer conn.Close()

	conn.Do("FLUSHDB")
}

type testSession struct {
	session_id string
	Messages   [][]byte
}

func (s *testSession) Receive() (m []byte) {
	result := s.Messages[0]
	if len(s.Messages) > 0 {
		s.Messages = s.Messages[1:]
	}

	return result
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
}
