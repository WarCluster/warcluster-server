// $ 6g echo.go && 6l -o echo echo.6
// $ ./echo
//
// ~ in another terminal ~
//
// $ nc localhost 3540

package main

import (
    "net"
    "bufio"
    "strconv"
    "fmt"
)

const HOST = "localhost"
const PORT = 3540

func main() {
    server, err := net.Listen("tcp", HOST + ":" + strconv.Itoa(PORT))
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
    b := bufio.NewReader(client)
    for {
        line, err := b.ReadBytes('\n')
        if err != nil { // EOF, or worse
            fmt.Printf("%v <- $ -> %v\n", client.LocalAddr(), client.RemoteAddr())
            break
        }
        client.Write([]byte(client.RemoteAddr().String()))
        client.Write([]byte(" said: "))
        client.Write(line)
        client.Write([]byte(">>> "))
    }
}
