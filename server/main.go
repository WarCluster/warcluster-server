package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"../db_manager"
	"../entities"
)

type Server struct {
	host       string
	port       int
	listener   net.Listener
	clients    map[net.Conn]chan<- string
	is_running bool
}

func (self *Server) Start(host string, port int) error {
	log.Print("Server is starting...")
	if self.is_running {
		return errors.New("Server is already started!")
	}
	var err error

	msgchan := make(chan string)
	addchan := make(chan *Client)
	rmchan := make(chan *Client)

	self.listener, err = net.Listen("tcp", fmt.Sprintf("%v:%v", host, port))
	if err == nil {
		self.host = host
		self.port = port
		self.is_running = true
		self.clients = make(map[net.Conn]chan<- string)
		log.Println("Server is up and running!")
	} else {
		log.Println(err)
		return err
	}

	go self.handleMessages(msgchan, addchan, rmchan)

	for self.is_running {
		conn, err := self.listener.Accept()
		if err != nil {
			if self.is_running {
				log.Println(err)
				continue
			} else {
				break
			}
		}
		go self.handleConnection(conn, msgchan, addchan, rmchan)
	}
	return nil
}

func (self *Server) Stop() error {
	log.Println("Server is shutting down...")
	if !self.is_running {
		err := errors.New("Server is already stopped!")
		log.Println(err)
		return err
	}

	for client := range self.clients {
		client.Close()
		delete(self.clients, client)
	}
	self.listener.Close()
	self.is_running = false
	return nil
}

func (self *Server) Restart() {
	self.Stop()
	self.Start(self.host, self.port)
}

func (self *Server) handleConnection(c net.Conn, msgchan chan<- string, addchan, rmchan chan<- *Client) {
	bufc := bufio.NewReader(c)
	nickname, player := authenticate(c, bufc)

	client := &Client{
		conn:     c,
		nickname: nickname,
		player:   player,
		channel:  make(chan string),
	}

	defer func() {
		c.Close()
		log.Printf("Connection from %v closed.\n", c.RemoteAddr())
		rmchan <- client
	}()

	addchan <- client
	home_planet_entity, _ := db_manager.GetEntity(client.player.HomePlanet)
	home_planet := home_planet_entity.(entities.Planet)
	io.WriteString(c, fmt.Sprintf("{username: '%s', position: [%d, %d] }",
		client.nickname, home_planet.GetCoords()[0], home_planet.GetCoords()[1]))
	go client.ReadLinesInto(msgchan)
	client.WriteLinesFrom(client.channel)
}

func (self *Server) handleMessages(msgchan <-chan string, addchan, rmchan <-chan *Client) {
	for {
		select {
		case msg := <-msgchan:
			log.Printf("New message: %s", msg)
			for _, ch := range self.clients {
				go func(mch chan<- string) { mch <- msg }(ch)
			}
		case client := <-addchan:
			log.Printf("New client: %v\n", client.nickname)
			self.clients[client.conn] = client.channel
		case client := <-rmchan:
			log.Printf("Client disconnects: %v\n", client.nickname)
			delete(self.clients, client.conn)
		}
	}
}
