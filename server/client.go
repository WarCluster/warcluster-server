package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"../entities"
)

type Client struct {
	conn     net.Conn
	nickname string
	channel  chan string
	player   *entities.Player
}

func (self *Client) ReadLinesInto(ch chan<- string) {
	bufc := bufio.NewReader(self.conn)
	for {
		line, err := bufc.ReadString('\n')
		if err != nil {
			break
		}
		ch <- fmt.Sprintf("%s: %s", self.nickname, line)
	}
}

func (self *Client) WriteLinesFrom(ch <-chan string) {
	for msg := range ch {
		if _, err := io.WriteString(self.conn, msg); err != nil {
			return
		}
	}
}
