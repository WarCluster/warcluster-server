package server

import (
	"fmt"
	"time"
	"warcluster/db_manager"
	"warcluster/entities"
)

// StartMissionary is used when a call to initiate a new mission is rescived.
// 1. When the delay ends the thread ends the mission calling EndMission
// 2. The end of the mission is bradcasted to all clients and the mission entry is erased from the DB.
func StartMissionary(mission *entities.Mission) {
	target_key := fmt.Sprintf("planet.%d_%d", mission.Target[0], mission.Target[1])
	time.Sleep(time.Duration(mission.TravelTime * 1e6))

	target_entity, err := db_manager.GetEntity(target_key)
	target := target_entity.(*entities.Planet)

	result := entities.EndMission(target, mission)
	key, serialized_planet, err := result.Serialize()
	if err == nil {
		db_manager.SetEntity(result)
		sessions.Broadcast([]byte(fmt.Sprintf("{\"Command\": \"state_change\", \"Planets\": {\"%s\": %s}}", key, serialized_planet)))
	}
	db_manager.DeleteEntity(mission.GetKey())
}
