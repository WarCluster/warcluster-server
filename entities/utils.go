package entities

import (
    "crypto/sha512"
    "encoding/json"
    "io"
    "strconv"
    "strings"
    "fmt"
)

func Construct(key string, data []byte) Entity {
    entity_type := strings.Split(key, ".")[0]
    fmt.Print()

    switch entity_type {
    case "player":
        var player Player
        json.Unmarshal(data, &player)
        player.username = strings.Split(key, "player.")[1]
        return player
    case "planet":
        var planet Planet
        json.Unmarshal(data, &planet)
        planet.coords = ExtractPlanetCoords(key)
        return planet
    }
    return nil
}

func GenerateHash(username string) string {
    return simplifyHash(usernameHash(username))
}

func ExtractPlanetCoords(key string) []int {
    key_coords := strings.Split(key, ".")[1]
    planet_coords := strings.Split(key_coords, "_")
    planet_coords_0, _ := strconv.Atoi(planet_coords[0])
    planet_coords_1, _ := strconv.Atoi(planet_coords[1])
    return []int{planet_coords_0, planet_coords_1}
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

