package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServerTest struct {
	suite.Suite
	server *Server
}

func (st *ServerTest) SetupTest() {
	st.server = NewServer(
		cfg.Server.Host,
		7014,
	)
	go st.server.Start()
	for !st.server.isRunning {
		time.Sleep(100 * time.Millisecond)
	}
}

func (st *ServerTest) TearDownTest() {
	st.server.Stop()
}

func (st *ServerTest) TestStopping() {
	assert.True(st.T(), st.server.isRunning)
	st.server.Stop()
	assert.False(st.T(), st.server.isRunning)
}

func (st *ServerTest) TestConsolePermissions() {
	defer func(consoleStatus bool) {
		cfg.Server.Console = consoleStatus
		ExportConfig(cfg)
	}(cfg.Server.Console)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/console", st.server.Addr), nil)
	if err != nil {
		st.Error(err)
	}

	cfg.Server.Console = true
	consoleHandler(w, req)
	assert.Equal(st.T(), 200, w.Code)

	cfg.Server.Console = false
	w = httptest.NewRecorder()
	consoleHandler(w, req)
	assert.Equal(st.T(), 404, w.Code)
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTest))
}
