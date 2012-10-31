package main

import (
    "bufio"
    "bytes"
    "fmt"
    "net"
    "strconv"
)

const HOST = "localhost"
const PORT = 3540
const NEW_LINE byte = 10
var clients map[string]*bufio.ReadWriter

func main() {
    clients = make(map[string]*bufio.ReadWriter)
    server, err := net.Listen("tcp", HOST+":"+strconv.Itoa(PORT))
    if server == nil {
        panic("couldn't start listening: " + err.Error())
    }

    fmt.Println("I'm up and running!")
    connections := clientConnections(server)

    for {
        go handleConnection(<-connections)
    }
}

func clientConnections(listener net.Listener) chan net.Conn {
    channel := make(chan net.Conn)
    i := 0
    go func() {
        for {
            client, err := listener.Accept()
            if client == nil {
                fmt.Printf("couldn't accept: " + err.Error())
                continue
            }
            i++
            fmt.Printf("%v <--> %v\n", client.LocalAddr(), client.RemoteAddr())
            channel <- client
        }
    }()
    return channel
}

func handleConnection(client net.Conn) {
    var message bytes.Buffer
    buffer := bufio.NewReadWriter(bufio.NewReader(client), bufio.NewWriter(client))
    clients[client.RemoteAddr().String()] = buffer

    for {
        line, err := buffer.ReadString(NEW_LINE)
        if err != nil {
            fmt.Printf("%v <- $ -> %v\n", client.LocalAddr(), client.RemoteAddr())
            break
        }

        if len(line) > 4 && line[:5] == "/quit" {
            buffer.WriteString("Bye\n")
            buffer.Flush()
            client.Close()
        }

        message.WriteString(client.RemoteAddr().String())
        message.WriteString(" said: ")
        message.WriteString(line)
        go writeToEveryone(message.String())
        message.Reset()
    }
}

func writeToEveryone(message string) {
    for key, _ := range clients {
        clients[key].WriteString(message)
        clients[key].Flush()
    }
}
