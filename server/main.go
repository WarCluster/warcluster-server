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

	"github.com/fzzy/sockjs-go/sockjs"
)

var (
	listener net.Listener
	clients  *ClientPool = NewClientPool()
)

// This function goes trough all the procedurs needed for the werver to be initialized.
// Create an empty connections pool and start the listening foe messages loop.
func Start(host string, port uint16) error {
	log.Print(fmt.Sprintf("Server is running at http://%v:%v/", host, port))
	log.Print("Quit the server with Ctrl-C.")

	mux := sockjs.NewServeMux(http.DefaultServeMux)
	conf := sockjs.NewConfig()

	http.HandleFunc("/console", staticHandler)
	http.Handle("/static", http.FileServer(http.Dir(getStaticDir())))
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
func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join(getStaticDir(), "/index.html"))
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

	client, response, err := login(session)
	if err != nil {
		log.Print("Error in server.main.handler.login:", err.Error())
	}
	clients.Add(client)
	defer clients.Remove(client)

	clients.Send(client.Player, response)
	client.Player.UpdateSpyReports()
	for {
		message := session.Receive()
		if message == nil {
			break
		}

		request, err := UnmarshalRequest(message, client)
		if err != nil {
			log.Println("Error in server.main.handler.UnmarshalRequest:", err.Error())
		}

		action, err := ParseRequest(request)
		if err != nil {
			log.Println("Error in server.main.handler.ParseRequest:", err.Error())
		}

		if err := action(request); err != nil {
			log.Println(err)
		}
	}
}

// getStaticDir return an absolute path to the static files
func getStaticDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "../static")
}
