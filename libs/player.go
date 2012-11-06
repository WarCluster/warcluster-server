package libs

import (
    "crypto/sha512"
    // "fmt"
    "io"
    // "math"
    "strconv"
)

type Player struct {
    username string
    hash string
    // solar_system []int
}

func (player *Player) String() string {
    return player.hash
}

// func (player *Player) generatePlanets(sun_position []int) []Planet{
//     result := []Planet{}
//     ring_offset := float64(80)
//     planet_radius := 50
// 
//     for ix:=0; ix<9; ix++ {
//         fmt.Println(len(player.hash[4 * ix + 1]))
//         planet_in_creation := Planet{}
//         ring_offset += float64(planet_radius) + float64(player.hash[4 * ix]) - 48
//         planet_in_creation.coords[0] = int(float64(sun_position[0]) + ring_offset * math.Cos(
//             float64(40 * int(player.hash[4 * ix + 1]) - 48)))
//         result = append(result, planet_in_creation)
//     }
//     return result
// 
// }

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
