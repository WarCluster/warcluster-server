package entities

import (
    "crypto/sha512"
    "encoding/json"
    "io"
    "strconv"
    "strings"
)

func Construct(key string, data []byte) Entity {
    entity_type := strings.Split(key, ".")[0]

    switch entity_type {
    case "player":
        var player Player
        json.Unmarshal(data, &player)
        player.username = strings.Split(key, "player.")[1]
        return player
    }
    return nil
}

func GenerateHash(username string) string {
    return simplifyHash(usernameHash(username))
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
