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
	var timeSlept time.Duration = 0

	targetKey := fmt.Sprintf("planet.%s", mission.Target.Name)
	for _, transferPoint := range mission.TransferPoints() {
		timeToSleep := transferPoint.TravelTime - timeSlept
		timeSlept += timeToSleep
		time.Sleep(timeToSleep * time.Millisecond)
		mission.ChangeAreaSet(transferPoint.CoordinateAxis, transferPoint.Direction)

		stateChange := response.NewStateChange()
		stateChange.Missions = map[string]entities.Entity{
			mission.Key(): mission,
		}
		response.Send(stateChange, sessions.Broadcast)
	}

	time.Sleep((mission.TravelTime - timeSlept) * time.Millisecond)

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
		target.Key(): target,
	}
	response.Send(stateChange, sessions.Broadcast)

	entities.RemoveFromArea(mission.Key(), mission.AreaSet())
	entities.Delete(mission.Key())

	if excessShips > 0 {
		startExcessMission(mission, target, excessShips)
	}
}

func startExcessMission(mission *entities.Mission, homePlanet *entities.Planet, ships int32) {
	newTargetKey := fmt.Sprintf("planet.%s", mission.Source.Name)
	newTargetEntity, err := entities.Get(newTargetKey)
	if err != nil {
		log.Print("Error in homePlanet planet fetch: ", err.Error())
		return
	}

	playerEntity, err := entities.Get(fmt.Sprintf("player.%s", mission.Player))
	player := playerEntity.(*entities.Player)

	excessMission := player.StartMission(homePlanet, newTargetEntity.(*entities.Planet), 100, "Attack")
	excessMission.ShipCount = ships
	go StartMissionary(excessMission)
	entities.Save(excessMission)

	sendMission := response.NewSendMission()
	sendMission.Mission = excessMission
	response.Send(sendMission, sessions.Broadcast)
}
