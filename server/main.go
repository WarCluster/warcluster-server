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

var host       string
var port       int
var is_running bool

var sessions *sockjs.SessionPool = sockjs.NewSessionPool()

func Start(host string, port int) error {
	log.Print("Server is starting...")
	if is_running {
		return errors.New("Server is already started!")
	}
	mux := sockjs.NewServeMux(http.DefaultServeMux)
	conf := sockjs.NewConfig()
	http.HandleFunc("/", indexHandler)
	http.Handle("/static", http.FileServer(http.Dir("./static")))
	mux.Handle("/chat", handler, conf)

	err := http.ListenAndServe(":8081", mux)
	if err == nil {
		host = host
		port = port
		is_running = true
		sessions = sockjs.NewSessionPool()
		log.Println("Server is up and running!")
	} else {
		log.Println(err)
		return err
	}
	return err
}

func Stop() error {
	log.Println("Server is shutting down...")
	if !is_running {
		err := errors.New("Server is already stopped!")
		log.Println(err)
		return err
	}

	is_running = false
	return nil
}

func Restart() {
	Stop()
	Start(host, port)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func login(session sockjs.Session) (*Client, error) {
	nickname, player := authenticate(session)
	// TODO: Assume the login attempt could be wrong

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
// 	for {
// 		select {
// 		case msg := <-msgchan:
// 			log.Printf("New message: %s", msg)
// 			for _, ch := range self.users {
// 				go func(mch chan<- string) { mch <- msg }(ch)
// 			}
// 		case client := <-addchan:
// 			log.Printf("New client: %v\n", client.nickname)
// 			self.users[client.conn] = client.channel
// 		case client := <-rmchan:
// 			log.Printf("Client disconnects: %v\n", client.nickname)
// 			delete(self.users, client.conn)
// 		}
// 	}
// }
