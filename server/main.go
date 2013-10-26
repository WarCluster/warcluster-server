// This is the server package.
// The purpouse of this package is to map a connection to each player(who is online) so we have a comunication chanel.
package server

import (
	"errors"
	"fmt"
	"github.com/fzzy/sockjs-go/sockjs"
	"log"
	"net"
	"net/http"
	"path"
	"runtime"
	"runtime/debug"
	"warcluster/entities"
	"warcluster/entities/db"
	"warcluster/server/response"
)

var listener net.Listener

var sessions *sockjs.SessionPool = sockjs.NewSessionPool() //This is the SockJs sessions pull (a list of all the currently active client's sessions).

// This function goes trough all the procedurs needed for the werver to be initialized.
// Create an empty connections pool and start the listening foe messages loop.
func Start(host string, port int) error {
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

// This function is called from the message handler to parse the first message for every new connection.
// It check for existing user in the DB and logs him if the password is correct.
// If the user is new he is initiated and a new home planet nad solar system are generated.
func login(session sockjs.Session) (*Client, error) {
	nickname, player, err := authenticate(session)
	if err != nil {
		response.Send(response.NewLoginFailed(), session.Send)
		return nil, errors.New("Login failed")
	}

	client := &Client{
		Session:  session,
		Nickname: nickname,
		Player:   player,
	}

	home_planet_entity, _ := db.GetEntity(client.Player.HomePlanet)
	home_planet := home_planet_entity.(*entities.Planet)

	login_success := response.NewLoginSuccess()
	login_success.Username = client.Nickname
	login_success.Position = home_planet.GetCoords()
	response.Send(login_success, session.Send)
	return client, nil
}

// On the first received message from each connection the server will call the handler.
// Add new session to the session pool, call the login func to validate the connection and
// if the connection is valid enters "while true" state and uses ParseRequest to parse the requests.
//
// Shocking right?!?!
func handler(session sockjs.Session) {
	defer func() {
		if panicked := recover(); panicked != nil {
			log.Println(string(debug.Stack()))
			return
		}
	}()
	defer sessions.Remove(session)

	if client, err := login(session); err == nil {
		sessions.Add(session)
		for {
			message := session.Receive()
			if message == nil {
				break
			}

			if request, err := UnmarshalRequest(message, client); err == nil {
				if action, err := ParseRequest(request); err == nil {
					if err := action(request); err != nil {
						log.Println(err)
					}
				} else {
					log.Println("Error in server.main.handler.ParseRequest:", err.Error())
				}
			} else {
				log.Println("Error in server.main.handler.UnmarshalRequest:", err.Error())
			}
		}
	} else {
		session.End()
	}
}

// getStaticDir return an absolute path to the static files
func getStaticDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "../static")
}
