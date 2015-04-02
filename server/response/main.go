// Package response defines the messages we stream to users
package response

import (
	"encoding/json"
	"time"

	"github.com/pzsz/voronoi"

	"warcluster/entities"
)

var Diagram *voronoi.Diagram

type Responser interface {
	Sanitize(*entities.Player)
}

type Timestamp int64

type baseResponse struct {
	Command   string
	Timestamp Timestamp
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Now().UnixNano() / 1e6)
}

// The sanitizer recieves raw planet and obscures hidden for the player information
func SanitizePlanets(player *entities.Player, planets map[string]*entities.Planet) map[string]*entities.PlanetPacket {
	packets := make(map[string]*entities.PlanetPacket)
	for name, planet := range planets {
		packets[name] = planet.Sanitize(player)
	}
	return packets
}

func Send(r interface{}, player string, sender func(string, []byte)) error {
	serialized, err := json.Marshal(r)
	if err == nil {
		sender(player, serialized)
	}
	return err
}
