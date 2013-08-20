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
	CurrentTime int64
	StartTime   int64
	ArrivalTime int64
	Player      string
	ShipCount   int
}

func (self *Mission) String() string {
	return self.GetKey()
}

func (self *Mission) GetKey() string {
	return fmt.Sprintf(
		"mission.%d_%d_%d",
		self.StartTime,
		self.Source[0],
		self.Source[1],
	)
}

func (self *Mission) GetSpeed() int {
	return 5
}

func (self *Mission) Serialize() (string, []byte, error) {
	self.CurrentTime = time.Now().UnixNano() / 1e6
	result, err := json.Marshal(self)
	if err != nil {
		return self.GetKey(), nil, err
	}
	return self.GetKey(), result, nil
}

// The CalculateArrivalTime is used to calculate the mission duration.
func (self *Mission) CalculateArrivalTime() {
	start_vector := vec2d.New(float64(self.Source[0]), float64(self.Source[1]))
	end_vector := vec2d.New(float64(self.Target[0]), float64(self.Target[1]))
	distance := vec2d.GetDistance(end_vector, start_vector)
	self.ArrivalTime += int64(distance/float64(self.GetSpeed()) * 100)
}

func EndMission(endPlanet *Planet, missionInfo *Mission) *Planet {
	if endPlanet.Owner == missionInfo.Player {
		endPlanet.SetShipCount(endPlanet.GetShipCount() + missionInfo.ShipCount)
	} else {
		if missionInfo.ShipCount < endPlanet.GetShipCount() {
			endPlanet.SetShipCount(endPlanet.GetShipCount() - missionInfo.ShipCount)
		} else {
			endPlanet.SetShipCount(missionInfo.ShipCount - endPlanet.GetShipCount())
			endPlanet.Owner = missionInfo.Player
		}
	}
	return endPlanet
}
