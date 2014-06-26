// Package server maps a connection to each player(who is online) so we have a comunication chanel.
package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"path"
	"runtime"
	"runtime/debug"

	"warcluster/config"
	"warcluster/server/response"

	"github.com/fzzy/sockjs-go/sockjs"

	"warcluster/leaderboard"
)

var (
	cfg         config.Config
	clients     *ClientPool = NewClientPool()
	leaderBoard *leaderboard.Leaderboard
	listener    net.Listener
)

// This function goes trough all the procedurs needed for the werver to be initialized.
// Create an empty connections pool and start the listening foe messages loop.
func Start() error {
	cfg.Load("../config/config.gcfg")
	host := cfg.Server.Host
	port := cfg.Server.Port

	log.Print(fmt.Sprintf("Server is running at http://%v:%v/", host, port))
	log.Print("Quit the server with Ctrl-C.")

	mux := sockjs.NewServeMux(http.DefaultServeMux)
	conf := sockjs.NewConfig()

	http.HandleFunc("/console", consoleHandler)
	http.HandleFunc("/leaderboard/players/", leaderboardPlayersHandler)
	http.HandleFunc("/leaderboard/teams/", leaderboardTeamsHandler)
	http.HandleFunc("/search/", searchHandler)
	mux.Handle("/universe", handler, conf)

	if err := ListenAndServe(fmt.Sprintf("%v:%v", host, port), mux); err != nil {
		log.Println(err)
		return err
	}

	return Stop()
}

// ListenAndServe listens on the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections.  If
// srv.Addr is blank, ":http" is used.
func ListenAndServe(address string, mux *sockjs.ServeMux) error {
	var err error

	server := &http.Server{Addr: address, Handler: mux}
	addr := server.Addr
	if addr == "" {
		addr = ":http"
	}
	listener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return server.Serve(listener)
}

// Die biatch and get the fuck out.
func Stop() error {
	log.Println("Server is shutting down...")
	listener.Close()
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
		http.ServeFile(response, request, path.Join(getStaticDir(), "/troll_face.png"))
	}
}

// On the first received message from each connection the server will call the handler.
// Add new session to the session pool, call the login func to validate the connection and
// if the connection is valid enters "while true" state and uses ParseRequest to parse the requests.
//
// Shocking right?!?!
func handler(session sockjs.Session) {
	defer func() {
		if panicked := recover(); panicked != nil {
			log.Println(fmt.Sprintf("%s\n\nStacktrace:\n\n%s", panicked, debug.Stack()))
			return
		}
	}()
	defer session.End()

	client, logResponse, err := login(session)
	if err != nil {
		log.Print("Error in server.main.handler.login:", err.Error())
		log.Println(err.Error())
		message, err := json.Marshal(logResponse)
		if err != nil {
			log.Println(err.Error())
		}
		session.Send(message)
		log.Printf("<- %s", message)

		return
	}
	clients.Add(client)
	defer clients.Remove(client)

	clients.Send(client.Player, logResponse)

	client.Player.UpdateSpyReports()
	for {
		message := session.Receive()
		if message == nil {
			break
		}

		request, err := UnmarshalRequest(message, client)
		if err != nil {
			log.Println("Error in server.main.handler.UnmarshalRequest:", err.Error())
			clients.Send(client.Player, response.NewComsError("Unable to unmarshal request"))
			continue
		}

		action, err := ParseRequest(request)
		if err != nil {
			log.Println("Error in server.main.handler.ParseRequest:", err.Error())
			clients.Send(client.Player, response.NewComsError("Unable to parse request"))
			continue
		}

		if err := action(request); err != nil {
			log.Println(err)
			continue
		}
		log.Printf("-> %s", message)
	}
}

// getStaticDir return an absolute path to the static files
func getStaticDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "../static")
}
