package main

import (
    "crypto/sha512"
    "fmt"
    "io"
    "strconv"
    "math"
)

func main() {
    player := CreatePlayer("kiril")
    fmt.Println(player.String())
    fmt.Println(generatePlanets(player, []int{1,2}))
}

type Player struct {
    username string
    hash string
    // solar_system []int
}

func (player *Player) String() string {
    return player.hash
}

type Planet struct {
    coords []int
    texture int
    size int
    ship_count int
    max_ship_count int
    owner *Player
}


func CreatePlayer(username string) Player {
    hash := usernameHash(username)
    return Player{username, simplifyHash(hash)}
}

func usernameHash(username string) []byte {
    hash := sha512.New()
    io.WriteString(hash, username)
    return hash.Sum(nil)
}


func simplifyHash(hash []byte) string {
    result := ""
    for ix:=0; ix<len(hash); ix++ {
        last_digit := hash[ix] % 10
        result += strconv.Itoa(int(last_digit))
    }
    return result
}
