package entities

import (
	"encoding/json"
	"fmt"
	"github.com/Vladimiroff/vec2d"
	"time"
)

type Mission struct {
	Color       Color
	Source      []int
	Target      []int
	Type        string
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

func (m *Mission) MarshalJSON() ([]byte, error) {
	m.CurrentTime = time.Now().UnixNano() / 1e6
	return json.Marshal((*missionMarshalHook)(m))
}

func (m *Mission) CalculateTravelTime() {
	start_vector := vec2d.New(float64(m.Source[0]), float64(m.Source[1]))
	end_vector := vec2d.New(float64(m.Target[0]), float64(m.Target[1]))
	distance := vec2d.GetDistance(end_vector, start_vector)
	m.TravelTime = int64(distance / float64(m.GetSpeed()) * 100)
}

func EndMission(target *Planet, missionInfo *Mission) (excessShips int) {
	switch missionInfo.Type {
	case "Attack":
		if target.Owner == missionInfo.Player {
			target.SetShipCount(target.ShipCount + missionInfo.ShipCount)
		} else {
			if missionInfo.ShipCount < target.ShipCount {
				target.SetShipCount(target.ShipCount - missionInfo.ShipCount)
			} else {
				if target.IsHome {
					target.SetShipCount(0)
					excessShips = missionInfo.ShipCount - target.ShipCount
				} else {
					target.SetShipCount(missionInfo.ShipCount - target.ShipCount)
					target.Owner = missionInfo.Player
					target.Color = missionInfo.Color
				}
			}
		}
	case "Supply":
		target.SetShipCount(target.ShipCount + missionInfo.ShipCount)
	case "Spy":
	}

	return
}
