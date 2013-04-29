/*
	This is the server package.
	The purpouse of this package is to map a connection to each player(who is online) so we have a comunication chanel.

*/
package server

import (
	"../db_manager"
	"../entities"
	"errors"
	"fmt"
	"github.com/fzzy/sockjs-go/sockjs"
	"log"
	"net/http"
)

var HOST string     //Server scope constant that keeps the server host address.
var PORT int        //Server scope constant that keeps the server port number.
var IS_RUNNING bool //Server scope variable that represents the is active flag.

var sessions *sockjs.SessionPool = sockjs.NewSessionPool() //This is the SockJs sessions pull (a list of all the currently active client's sessions). 

/*This function goes trough all the procedurs needed for the werver to be initialized.
1.Create an empty connections pool
2.Starts the listening foe messages loop.*/
func Start(host string, port int) error {
	log.Print("Server is starting...")
	if IS_RUNNING {
		return errors.New("Server is already started!")
	} else {
		HOST = host
		PORT = port
		IS_RUNNING = true
	}
	log.Println("Server is up and running!")
	mux := sockjs.NewServeMux(http.DefaultServeMux)
	conf := sockjs.NewConfig()

	http.HandleFunc("/console", staticHandler)
	http.Handle("/static", http.FileServer(http.Dir("./static")))
	mux.Handle("/universe", handler, conf)

	if err := http.ListenAndServe(fmt.Sprintf("%v:%v", HOST, PORT), mux); err != nil {
		log.Println(err)
		return err
	}
	return Stop()
}


//Die biatch and get the fuck out.
func Stop() error {
	log.Println("Server is shutting down...")
	if !IS_RUNNING {
		err := errors.New("Server is already stopped!")
		log.Println(err)
		return err
	}

	IS_RUNNING = false
	return nil
}


//Stop + Start = Restart
func Restart() {
	Stop()
	Start(HOST, PORT)
}

//Returns the HTML page needed to display the debug page (server "chat" window). 
func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
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
	home_planet := home_planet_entity.(entities.Planet)
	session.Send([]byte(fmt.Sprintf("{username: '%s', position: [%d, %d] }",
		client.Nickname, home_planet.GetCoords()[0], home_planet.GetCoords()[1])))
	return client, nil
}

/*On the first rescived message from each connection the server will call the handler.
So it can complete the following actions:
1.Adding a new session to the session pool.
2.Call the login func to validate the connection
3.If the connection is valid enters "while true" state and uses ParseRequest to parse the requests. Shocking right?!?!*/
func handler(session sockjs.Session) {
	sessions.Add(session)
	defer sessions.Remove(session)

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
