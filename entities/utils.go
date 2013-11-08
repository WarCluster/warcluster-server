package entities

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"github.com/Vladimiroff/vec2d"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Color struct {
	R int
	G int
	B int
}

func Construct(key string, data []byte) Entity {
	entityType := strings.Split(key, ".")[0]

	switch entityType {
	case "player":
		player := new(Player)
		json.Unmarshal(data, player)
		return player
	case "planet":
		planet := new(Planet)
		json.Unmarshal(data, planet)
		return planet
	case "mission":
		mission := new(Mission)
		json.Unmarshal(data, mission)
		return mission
	case "sun":
		sun := new(Sun)
		json.Unmarshal(data, sun)
		return sun
	}
	return nil
}

func generateHash(username string) string {
	return simplifyHash(usernameHash(username))
}

func usernameHash(username string) []byte {
	hash := sha512.New()
	io.WriteString(hash, username)
	return hash.Sum(nil)
}

func simplifyHash(hash []byte) string {
	result := ""
	for ix := 0; ix < len(hash); ix++ {
		lastDigit := hash[ix] % 10
		result += strconv.Itoa(int(lastDigit))
	}
	return result
}

func getRandomStartPosition(scope int) *vec2d.Vector {
	xSeed := time.Now().UTC().UnixNano()
	ySeed := time.Now().UTC().UnixNano()
	xGenerator := rand.New(rand.NewSource(xSeed))
	yGenerator := rand.New(rand.NewSource(ySeed))
	return vec2d.New(float64(xGenerator.Intn(scope)-scope/2), float64(yGenerator.Intn(scope)-scope/2))
}

// Gets the first three letters from twitter's username
// and returns them in upper-case. If there are no three
// letters there (like in @r2d2) we take as many non-letter
// symbols as we need (@r2d2 should go to RD2)
func extractUsernameInitials(nickname string) string {
	letters := []rune{}
	nonLetters := []rune{}
	for _, s := range nickname {
		symbol := rune(s)
		if unicode.IsLetter(symbol) {
			letters = append(letters, symbol)
		} else {
			nonLetters = append(nonLetters, symbol)
		}
	}
	for len(letters) < 3 {
		letters = append(letters, nonLetters[0])
	}
	return fmt.Sprintf(
		"%c%c%c",
		unicode.ToUpper(letters[0]),
		unicode.ToUpper(letters[1]),
		unicode.ToUpper(letters[2]),
	)
}
