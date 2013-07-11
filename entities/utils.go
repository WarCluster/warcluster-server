package entities

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"github.com/Vladimiroff/vec2d"
	"math/rand"
	"io"
	"strconv"
	"strings"
	"time"
)

func Construct(key string, data []byte) Entity {
	entity_type := strings.Split(key, ".")[0]

	switch entity_type {
	case "player":
		player := new(Player)
		json.Unmarshal(data, player)
		player.username = strings.Split(key, "player.")[1]
		return player
	case "planet":
		planet := new(Planet)
		json.Unmarshal(data, planet)
		planet.coords = ExtractPlanetCoords(key)
		return planet
	case "mission":
		mission := new(Mission)
		json.Unmarshal(data, mission)
		mission.start_planet, mission.start_time = ExtractMissionsKey(key)
		return mission
	case "sun":
		sun := new(Sun)
		json.Unmarshal(data, sun)
		sun.position = ExtractSunKey(key)
		return sun
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

func ExtractMissionsKey(key string) (string, time.Time) {
	params_raw := strings.Split(key, ".")[1]
	params := strings.Split(params_raw, "_")
	parsed_time, _ := strconv.ParseInt(params[0], 10, 64)
	start_time := time.Unix(parsed_time, 0)
	start_planet := fmt.Sprintf("planet.%s_%s", params[1], params[2])
	return start_planet, start_time
}

func ExtractSunKey(key string) *vec2d.Vector {
	params_raw := strings.Split(key, ".")[1]
	params := strings.Split(params_raw, "_")
	sun_coords_0, _ := strconv.ParseFloat(params[0], 64)
	sun_coords_1, _ := strconv.ParseFloat(params[1], 64)
	coords := vec2d.New(sun_coords_0, sun_coords_1)
	return coords
}

func usernameHash(username string) []byte {
	hash := sha512.New()
	io.WriteString(hash, username)
	return hash.Sum(nil)
}

func simplifyHash(hash []byte) string {
	result := ""
	for ix := 0; ix < len(hash); ix++ {
		last_digit := hash[ix] % 10
		result += strconv.Itoa(int(last_digit))
	}
	return result
}

func getRandomStartPosition(scope int) *vec2d.Vector {
	x_seed := time.Now().UTC().UnixNano()
	y_seed := time.Now().UTC().UnixNano()
	x_generator := rand.New(rand.NewSource(x_seed))
	y_generator := rand.New(rand.NewSource(y_seed))
	return vec2d.New(float64(x_generator.Intn(scope) - scope/2), float64(y_generator.Intn(scope) - scope/2))
}
