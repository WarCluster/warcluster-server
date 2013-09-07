package entities

import (
	"encoding/json"
	"fmt"
	"github.com/Vladimiroff/vec2d"
	"time"
)

type Mission struct {
	Source      []int
	Target      []int
	Type 		string
	CurrentTime int64
	StartTime   int64
	TravelTime  int64
	Player      string
	ShipCount   int
}

func (m *Mission) String() string {
	return m.GetKey()
}

func (m *Mission) GetKey() string {
	return fmt.Sprintf(
		"mission.%d_%d_%d",
		m.StartTime,
		m.Source[0],
		m.Source[1],
	)
}

func (m *Mission) GetSpeed() int {
	return 10
}

func (m *Mission) Serialize() (string, []byte, error) {
	m.CurrentTime = time.Now().UnixNano() / 1e6
	result, err := json.Marshal(m)
	if err != nil {
		return m.GetKey(), nil, err
	}
	return m.GetKey(), result, nil
}

func (m *Mission) CalculateTravelTime() {
	start_vector := vec2d.New(float64(m.Source[0]), float64(m.Source[1]))
	end_vector := vec2d.New(float64(m.Target[0]), float64(m.Target[1]))
	distance := vec2d.GetDistance(end_vector, start_vector)
	m.TravelTime = int64(distance / float64(m.GetSpeed()) * 100)
}

func EndMission(target *Planet, target_owner *Player, missionInfo *Mission) *Planet {
	if target.Owner == missionInfo.Player {
		target.SetShipCount(target.GetShipCount() + missionInfo.ShipCount)
	} else {
		if missionInfo.ShipCount < target.GetShipCount() {
			target.SetShipCount(target.GetShipCount() - missionInfo.ShipCount)
		} else {
			if(target_owner.HomePlanet == target.GetKey()){
				//exess := missionInfo.ShipCount - target.GetShipCount()
				target.SetShipCount(0)
			// We need to create a new mission with the exess ships to sent back to the origin planet
			} else {
				target.SetShipCount(missionInfo.ShipCount - target.GetShipCount())
				target.Owner = missionInfo.Player
			}
		}
	}
	return target
}