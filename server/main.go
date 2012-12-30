package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
)

type Server struct {
	listener   net.Listener
	clients    map[net.Conn]chan<- string
	is_running bool
}

func (self *Server) Start(host string, port int) error {
	if self.is_running {
		return errors.New("Server is already running!")
	}
	var err error

	sigtermchan := make(chan os.Signal, 1)
	msgchan := make(chan string)
	addchan := make(chan *Client)
	rmchan := make(chan *Client)

	signal.Notify(sigtermchan, os.Interrupt)
	self.listener, err = net.Listen("tcp", fmt.Sprintf("%v:%v", host, port))
	if err == nil {
		self.is_running = true
		self.clients = make(map[net.Conn]chan<- string)
	} else {
		return err
	}

	log.Println("Server is up and running!")
	go self.Stop(sigtermchan)
	go self.handleMessages(msgchan, addchan, rmchan)

	for self.is_running {
		conn, err := self.listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go self.handleConnection(conn, msgchan, addchan, rmchan)
	}
	return nil
}

func (self *Server) Stop(sigtermchan <-chan os.Signal) error {
	if !self.is_running {
		return errors.New("Server is already running!")
	}

	select {
	case <-sigtermchan:
		log.Println("Server is shutting down...")
		for client := range self.clients {
			client.Close()
			delete(self.clients, client)
		}
		self.listener.Close()
		self.is_running = false
	}
	return nil
}

func authenticate(c net.Conn, bufc *bufio.Reader) string {
	io.WriteString(c, "Authenticating...\n")
	io.WriteString(c, "What is your nick?\n> ")
	nick, _, _ := bufc.ReadLine()
	return string(nick)
}

func (self *Server) handleConnection(c net.Conn, msgchan chan<- string, addchan, rmchan chan<- *Client) {
	bufc := bufio.NewReader(c)
	client := &Client{
		conn:     c,
		nickname: authenticate(c, bufc),
		channel:  make(chan string),
	}

	defer func() {
		c.Close()
		msgchan <- fmt.Sprintf("User %s left the chat room.\n", client.nickname)
		log.Printf("Connection from %v closed.\n", c.RemoteAddr())
		rmchan <- client
	}()

	if strings.TrimSpace(client.nickname) == "" {
		io.WriteString(c, "Invalid Username\n")
		return
	}
	addchan <- client
	io.WriteString(c, fmt.Sprintf("Welcome, %s!\n\n", client.nickname))
	msgchan <- fmt.Sprintf("%s has joined.\n", client.nickname)
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
