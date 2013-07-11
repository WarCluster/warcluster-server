/*
	This is the server package.
	The purpouse of this package is to map a connection to each player(who is online) so we have a comunication chanel.

*/
package server

import (
	"runtime/debug"
	"fmt"
	"github.com/fzzy/sockjs-go/sockjs"
	"log"
	"net"
	"net/http"
	"path"
	"runtime"
	"warcluster/db_manager"
	"warcluster/entities"
)

var listener net.Listener

var sessions *sockjs.SessionPool = sockjs.NewSessionPool() //This is the SockJs sessions pull (a list of all the currently active client's sessions).

/*This function goes trough all the procedurs needed for the werver to be initialized.
Create an empty connections pool and start the listening foe messages loop.*/
func Start(host string, port int) error {
	log.Print("Server is starting...")
	log.Println("Server is up and running!")

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

//Die biatch and get the fuck out.
func Stop() error {
	log.Println("Server is shutting down...")
	listener.Close()
	log.Println("Server has stopped.")
	return nil
}

//Returns the HTML page needed to display the debug page (server "chat" window).
func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join(getStaticDir(), "/index.html"))
}

/*This function is called from the message handler to parse the first message for every new connection.
It check for existing user in the DB and logs him if the password is correct.
If the user is new he is initiated and a new home planet nad solar system are generated.*/
func login(session sockjs.Session) (*Client, error) {
	nickname, player := authenticate(session)

	client := &Client{
		Session:  session,
		Nickname: nickname,
		Player:   player,
	}

	home_planet_entity, _ := db_manager.GetEntity(client.Player.HomePlanet)
	home_planet := home_planet_entity.(*entities.Planet)
	session.Send([]byte(fmt.Sprintf("{\"Command\": \"login_success\", \"Username\": \"%s\", \"Position\": [%d, %d] }",
		client.Nickname, home_planet.GetCoords()[0], home_planet.GetCoords()[1])))
	return client, nil
}

/*On the first received message from each connection the server will call the handler.
Add new session to the session pool, call the login func to validate the connection and
if the connection is valid enters "while true" state and uses ParseRequest to parse the requests.

Shocking right?!?!*/
func handler(session sockjs.Session) {
	sessions.Add(session)
	defer sessions.Remove(session)
	defer func() {
		if panicked := recover(); panicked != nil {
			log.Println(string(debug.Stack()))
			return
		}
	}()

	if client, err := login(session); err == nil {
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
					log.Println("Error in server.main.handler:", err.Error())
				}
			} else {
				log.Println("Error in server.main.handler:", err.Error())
			}
		}
	} else {
		session.End()
	}
}

func getStaticDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "../static")
}
