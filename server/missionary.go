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
	targetKey := fmt.Sprintf("planet.%d_%d", mission.Target[0], mission.Target[1])
	time.Sleep(time.Duration(mission.TravelTime * 1e6))

	targetEntity, err := entities.Get(targetKey)
	if err != nil {
		log.Print("Error in target planet fetch: ", err.Error())
		return
	}
	target := targetEntity.(*entities.Planet)
	target.UpdateShipCount()

	excessShips := entities.EndMission(target, mission)
	entities.Save(target)

	stateChange := response.NewStateChange()
	stateChange.Planets = map[string]entities.Entity{
		target.GetKey(): target,
	}
	response.Send(stateChange, sessions.Broadcast)

	entities.Delete(mission.GetKey())

	if excessShips > 0 {
		sourceKey := fmt.Sprintf("planet.%d_%d", mission.Source[0], mission.Source[1])
		sourceEntity, err := entities.Get(sourceKey)
		if err != nil {
			log.Print("Error in target planet fetch: ", err.Error())
			return
		}

		playerEntity, err := entities.Get(mission.Player)
		player := playerEntity.(*entities.Player)

		excessMission := player.StartMission(target, sourceEntity.(*entities.Planet), 100, "Attack")
		excessMission.ShipCount = excessShips
		go StartMissionary(excessMission)
		entities.Save(excessMission)

		sendMission := response.NewSendMission()
		sendMission.Mission = excessMission
		response.Send(sendMission, sessions.Broadcast)
	}
}
