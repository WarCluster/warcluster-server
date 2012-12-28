package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const HOST = "localhost"
const PORT = 7000

type Client struct {
	conn     net.Conn
	nickname string
	channel  chan string
}

func main() {
	msgchan := make(chan string)
	addchan := make(chan Client)
	rmchan := make(chan Client)

	ln, err := net.Listen("tcp", fmt.Sprintf("%v:%v", HOST, PORT))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server is up and running!")
	go handleMessages(msgchan, addchan, rmchan)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn, msgchan, addchan, rmchan)
	}
}

func (c Client) ReadLinesInto(ch chan<- string) {
	bufc := bufio.NewReader(c.conn)
	for {
		line, err := bufc.ReadString('\n')
		if err != nil {
			break
		}
		ch <- fmt.Sprintf("%s: %s", c.nickname, line)
	}
}

func (c Client) WriteLinesFrom(ch <-chan string) {
	for msg := range ch {
		if _, err := io.WriteString(c.conn, msg); err != nil {
			return
		}
	}
}

func authenticate(c net.Conn, bufc *bufio.Reader) string {
	io.WriteString(c, "Authenticating...\n")
	io.WriteString(c, "What is your nick?\n> ")
	nick, _, _ := bufc.ReadLine()
	return string(nick)
}

func handleConnection(c net.Conn, msgchan chan<- string, addchan chan<- Client, rmchan chan<- Client) {
	bufc := bufio.NewReader(c)
	client := Client{
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

func handleMessages(msgchan <-chan string, addchan <-chan Client, rmchan <-chan Client) {
	clients := make(map[net.Conn]chan<- string)

	for {
		select {
		case msg := <-msgchan:
			log.Printf("New message: %s", msg)
			for _, ch := range clients {
				go func(mch chan<- string) { mch <- fmt.Sprintf("%v\n%#v\n", msg, clients) }(ch)
			}
		case client := <-addchan:
			log.Printf("New client: %v\n", client.nickname)
			clients[client.conn] = client.channel
		case client := <-rmchan:
			log.Printf("Client disconnects: %v\n", client.nickname)
			delete(clients, client.conn)
		}
	}
}
