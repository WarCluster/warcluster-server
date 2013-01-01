package server

import (
	"../entities"
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
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

		if strings.HasPrefix(line, "sm;") {
			params := strings.Split(line, ";")
			if len(params) != 4 {
				continue
			}
			fleet, _ := strconv.Atoi(params[3])
			if err := actionParser(self.nickname, params[1], params[2], fleet); err == nil {
				ch <- fmt.Sprintf("%s: %s", self.nickname, line)
			}
		} else if strings.HasPrefix(line, "scope:") {
			ch <- fmt.Sprintf("%s: %s", self.nickname, line)
		}
	}
}

func (self *Client) WriteLinesFrom(ch <-chan string) {
	for msg := range ch {
		if _, err := io.WriteString(self.conn, msg); err != nil {
			return
		}
	}
}
