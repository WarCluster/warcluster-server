package entities

import (
	"encoding/json"
	"fmt"
	"time"
)

type Mission struct {
	Source      []int
	Target      []int
	CurrentTime time.Time
	StartTime   time.Time
	ArrivalTime time.Time
	Player      string
	ShipCount   int
}

func (self *Mission) String() string {
	return self.GetKey()
}

func (self *Mission) GetKey() string {
	return fmt.Sprintf(
		"mission.%d%d_%d_%d",
		self.StartTime.Unix(),
		self.StartTime.Nanosecond()/1e6,
		self.Source[0],
		self.Source[1],
	)
}

func (self *Mission) GetSpeed() int {
	return 5
}

func (self *Mission) Serialize() (string, []byte, error) {
	self.CurrentTime = time.Now()
	result, err := json.Marshal(self)
	if err != nil {
		return self.GetKey(), nil, err
	}
	return self.GetKey(), result, nil
}

func (self *Mission) GetStartTime() time.Time {
	return self.StartTime
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
