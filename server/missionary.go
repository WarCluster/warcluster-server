package server

import (
	"fmt"
	"log"
	"time"
	"warcluster/entities"
	"warcluster/server/response"
)

// StartMissionary is used when a call to initiate a new mission is rescived.
// 1. When the delay ends the thread ends the mission calling EndMission
// 2. The end of the mission is bradcasted to all clients and the mission entry is erased from the DB.
func StartMissionary(mission *entities.Mission) {
	target_key := fmt.Sprintf("planet.%d_%d", mission.Target[0], mission.Target[1])
	time.Sleep(time.Duration(mission.TravelTime * 1e6))

	target_entity, err := entities.Get(target_key)
	if err != nil {
		log.Print("Error in target planet fetch: ", err.Error())
		return
	}
	target := target_entity.(*entities.Planet)

	target.UpdateShipCount()

	result := entities.EndMission(target, mission)
	state_change := response.NewStateChange()
	state_change.Planets = map[string]entities.Entity{
		result.GetKey(): result,
	}
	response.Send(state_change, sessions.Broadcast)
	entities.Delete(mission.GetKey())
}
