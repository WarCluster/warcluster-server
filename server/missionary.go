package server

import (
	"fmt"
	"log"
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
	if err != nil {
		log.Print("Error in target planet fetch: ", err.Error())
	}
	target := target_entity.(*entities.Planet)

	var target_owner *entities.Player

	if target.HasOwner() {
		owner_id := fmt.Sprintf("player.%s", target.Owner)
		target_owner_entity, err := db_manager.GetEntity(owner_id)
		if err != nil {
			log.Print("Error in target planet owner fetch: ", err.Error(), " Searched for: ", owner_id)
		}
		if target_owner_entity != nil {
			target_owner = target_owner_entity.(*entities.Player)
		} else {
			log.Print("Error in target planet owner cast. Owner is nil!")
		}
	} else {
		target_owner = nil
	}

	target.UpdateShipCount()

	result := entities.EndMission(target, target_owner, mission)
	key, serialized_planet, err := result.Serialize()
	if err == nil {
		db_manager.SetEntity(result)
		sessions.Broadcast([]byte(fmt.Sprintf("{\"Command\": \"state_change\", \"Planets\": {\"%s\": %s}}", key, serialized_planet)))
	}
	db_manager.DeleteEntity(mission.GetKey())
}
