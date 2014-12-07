package server

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/websocket"

	"warcluster/config"
	"warcluster/entities"
	"warcluster/entities/db"
	"warcluster/leaderboard"
)

var (
	testServer     *Server
	incompleteUser = Request{Command: "login", TwitterID: "some twitter ID"}
	user           = Request{Command: "login", Username: "gophie", TwitterID: "some twitter ID"}
	setupParams    = Request{Command: "setup_parameters", Race: 0, SunTextureId: 0}
	setup          = Request{Command: "setup", Race: 0, SunTextureId: 0}
)

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
	)

	go testServer.Start()
	for !testServer.isRunning {
		time.Sleep(100 * time.Millisecond)
	}
}

type WebSocketTestSuite struct {
	suite.Suite
	conn    redis.Conn
	ws      *websocket.Conn
	message map[string]interface{}
}

func (w *WebSocketTestSuite) Dial() (*websocket.Conn, error) {
	origin := "http://localhost/"
	url := "ws://localhost:7013/universe"
	return websocket.Dial(url, "", origin)
}

func (w *WebSocketTestSuite) SetupTest() {
	var err error

	w.message = make(map[string]interface{})
	w.conn = db.Pool.Get()
	w.conn.Do("FLUSHDB")
	w.ws, err = w.Dial()
	if err != nil {
		log.Fatal(err)
	}

	cfg.Load()
	InitLeaderboard(leaderboard.New())
}

func (w *WebSocketTestSuite) TearDownTest() {
	w.ws.Close()
	w.conn.Close()
}

func (w *WebSocketTestSuite) assertReceive(command string) {
	w.message = make(map[string]interface{})

	receive := func() <-chan struct{} {
		ch := make(chan struct{})
		go func() {
			websocket.JSON.Receive(w.ws, &w.message)
			ch <- struct{}{}
		}()
		return ch
	}

	select {
	case <-time.After(10 * time.Second):
		w.T().Fatalf("Did not receive %s after 10 seconds", command)
	case <-receive():
		assert.Equal(w.T(), command, w.message["Command"])
	}
}

func (w *WebSocketTestSuite) assertSend(request *Request) {
	send := func() <-chan struct{} {
		ch := make(chan struct{})
		go func() {
			websocket.JSON.Send(w.ws, request)
			ch <- struct{}{}
		}()
		return ch
	}

	select {
	case <-time.After(10 * time.Second):
		w.T().Fatalf("Did not send %s after 10 seconds", request.Command)
	case <-send():
	}
}

func NewFakeClient(player *entities.Player) *Client {
	client := NewClient(new(websocket.Conn), player)
	client.codec = new(fakeCodec)
	return client
}

type fakeCodec struct {
	Messages  [][]byte
	Marshal   func(v interface{}) (data []byte, payloadType byte, err error)
	Unmarshal func(data []byte, payloadType byte, v interface{}) (err error)
}

func (c *fakeCodec) Receive(ws *websocket.Conn, v interface{}) error {
	if len(c.Messages) > 0 {
		message := c.Messages[0]
		c.Messages = c.Messages[1:]
		return json.Unmarshal(message, v)
	}
	return errors.New("No message received")
}

func (c *fakeCodec) Send(ws *websocket.Conn, v interface{}) error {
	m, err := json.Marshal(v)
	c.Messages = append(c.Messages, m)
	return err
}
