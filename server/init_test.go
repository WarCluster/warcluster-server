package server

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fzzy/sockjs-go/sockjs"
	)

type testSession struct {
	session_id string
	Messages   [][]byte
}

func (s *testSession) Receive() (m []byte) {
	return []byte{}
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