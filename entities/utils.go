package entities

import (
    "crypto/sha512"
    "io"
    "strconv"
)

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
