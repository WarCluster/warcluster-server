package entities

import (
	"warcluster/vector"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
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
	case "mission":
		var mission Mission
		json.Unmarshal(data, &mission)
		mission.start_planet, mission.start_time = ExtractMissionsKey(key)
		return mission
	case "sun":
		var sun Sun
		json.Unmarshal(data, &sun)
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

func ExtractSunKey(key string) *vector.Vector {
	params_raw := strings.Split(key, ".")[1]
	params := strings.Split(params_raw, "_")
	sun_coords_0, _ := strconv.ParseFloat(params[0], 64)
	sun_coords_1, _ := strconv.ParseFloat(params[1], 64)
	coords := vector.New(sun_coords_0, sun_coords_1)
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
