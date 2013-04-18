package server

import (
	"../db_manager"
	"../entities"
	"errors"
	"fmt"
	"github.com/fzzy/sockjs-go/sockjs"
	"log"
	"net/http"
	"strings"
)

var host       string
var port       int
var is_running bool

var users *sockjs.SessionPool = sockjs.NewSessionPool()

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
	fmt.Println("wazaaa")
	if err == nil {
		host = host
		port = port
		is_running = true
		users = sockjs.NewSessionPool()
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

func login(session sockjs.Session) bool {
	nickname, player := authenticate(session)
	// TODO: Assume the login attempt could be wrong

	client := &Client{
		session:  session,
		nickname: nickname,
		player:   player,
		channel:  make(chan string),
	}

	home_planet_entity, _ := db_manager.GetEntity(client.player.HomePlanet)
	home_planet := home_planet_entity.(entities.Planet)
	session.Send([]byte(fmt.Sprintf("{username: '%s', position: [%d, %d] }",
		client.nickname, home_planet.GetCoords()[0], home_planet.GetCoords()[1])))
	return true
}

func handler(session sockjs.Session) {
	users.Add(session)
	defer users.Remove(session)

	isLoginAttemptSuccessful := login(session)
	if(isLoginAttemptSuccessful) {
		for {
			message := session.Receive()
			if message == nil {
				break
			}
			fullAddr := session.Info().RemoteAddr
			addr := fullAddr[:strings.LastIndex(fullAddr, ":")]
			message = []byte(fmt.Sprintf("%s: %s", addr, message))
			users.Broadcast(message)
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
