// Package server maps a connection to each player(who is online) so we have a comunication chanel.
package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"path"
	"runtime"
	"runtime/debug"

	"golang.org/x/net/websocket"

	"warcluster/config"
	"warcluster/leaderboard"
	"warcluster/server/response"
)

type Server struct {
	http.Server
	listener  net.Listener
	isRunning bool
}

var (
	cfg         config.Config
	clients     *ClientPool
	leaderBoard *leaderboard.Leaderboard
	listener    net.Listener
	empty       = struct{}{}
)

// Exports to given loaded config file into server.cfg
func ExportConfig(loadedCfg config.Config) {
	cfg = loadedCfg
}

// Create new server and setup its routes.
func NewServer(host string, port uint16) *Server {
	s := new(Server)
	s.Addr = fmt.Sprintf("%v:%v", host, port)
	return s
}

// This function goes trough all the procedurs needed for the werver to be initialized.
// Create an empty connections pool and start the listening foe messages loop.
func (s *Server) Start() error {
	clients = NewClientPool(13)

	log.Print(fmt.Sprintf("Server is running at http://%s/", s.Addr))
	log.Print("Quit the server with Ctrl-C.")

	if err := s.ListenAndServe(); err != nil {
		log.Println(err)
		return err
	}

	return s.Stop()
}

// ListenAndServe listens on the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections.o
// It's re-defined here in order to have the listener which allows us
// to stop listening and clean up before exit.
func (s *Server) ListenAndServe() error {
	var err error

	s.listener, err = net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	s.isRunning = true
	return s.Serve(s.listener)
}

// Stops the server.
func (s *Server) Stop() error {
	log.Println("Server is shutting down...")
	s.isRunning = false
	s.listener.Close()
	log.Println("Server has stopped.")
	return nil
}

// Returns the HTML page needed to display the debug page (server "chat" window).
func consoleHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
	response.Header().Add("Pragma", "no-cache")                                   // HTTP 1.0.
	response.Header().Add("Expires", "0")                                         // Proxies
	if cfg.Server.Console {
		http.ServeFile(response, request, path.Join(getStaticDir(), "/index.html"))
	} else {
		http.NotFound(response, request)
	}
}

// On the first received message from each connection the server will call the handler.
// Add new session to the session pool, call the login func to validate the connection and
// if the connection is valid enters "while true" state and uses ParseRequest to parse the requests.
func Handle(ws *websocket.Conn) {
	var request Request
	defer func() {
		if panicked := recover(); panicked != nil {
			log.Println(fmt.Sprintf("%s\n\nStacktrace:\n\n%s", panicked, debug.Stack()))
			return
		}
	}()
	defer ws.Close()

	client, logResponse, err := login(ws)
	if err != nil {
		log.Print("Error in server.main.handler.login:", err.Error())
		websocket.JSON.Send(ws, &logResponse)
		return
	}
	clients.Add(client)
	defer clients.Remove(client)

	clients.Send(client.Player, logResponse)

	client.Player.UpdateSpyReports()
	for {
		err := websocket.JSON.Receive(client.Conn, &request)
		if err != nil {
			clients.Send(client.Player, response.NewError(err.Error()))
			continue
		}
		request.Client = client
		action, err := ParseRequest(&request)
		if err != nil {
			log.Println("Error in server.main.handler.ParseRequest:", err.Error())
			clients.Send(client.Player, response.NewError(err.Error()))
			continue
		}

		if err := action(&request); err != nil {
			log.Println(err)
			continue
		}
	}
}

// getStaticDir return an absolute path to the static files
func getStaticDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "../static")
}
