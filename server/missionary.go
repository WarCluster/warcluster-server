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
	var (
		timeSlept   time.Duration = 0
		excessShips int32
	)

	targetKey := fmt.Sprintf("planet.%s", mission.Target.Name)
	for _, transferPoint := range mission.TransferPoints() {
		timeToSleep := transferPoint.TravelTime - timeSlept
		timeSlept += timeToSleep
		time.Sleep(timeToSleep * time.Millisecond)
		mission.ChangeAreaSet(transferPoint.CoordinateAxis, transferPoint.Direction)

		stateChange := response.NewStateChange()
		stateChange.Missions = map[string]*entities.Mission{
			mission.Key(): mission,
		}
		clients.BroadcastToAll(stateChange)
	}

	time.Sleep((mission.TravelTime - timeSlept) * time.Millisecond)

	targetEntity, err := entities.Get(targetKey)
	if err != nil {
		log.Print("Error in target planet fetch:", err.Error())
		return
	}
	target := targetEntity.(*entities.Planet)
	target.UpdateShipCount()

	playerEntity, err := entities.Get(fmt.Sprintf("player.%s", mission.Target.Owner))
	if err != nil {
		log.Println("Error in target planet owner fetch:", err.Error())
		return
	}
	player := playerEntity.(*entities.Player)

	stateChange := response.NewStateChange()
	stateChange.RawPlanets = map[string]*entities.Planet{
		target.Key(): target,
	}

	switch mission.Type {
	case "Attack":
		excessShips = mission.EndAttackMission(target)
		clients.BroadcastToAll(stateChange)
	case "Supply":
		excessShips = mission.EndSupplyMission(target)
		clients.Send(player, stateChange)
	case "Spy":
		for {
			mission.EndSpyMission(target)
			updateSpyReports(player, mission, stateChange)
			if mission.ShipCount > 0 {
				time.Sleep(entities.SPY_REPORT_VALIDITY * time.Second)
			} else {
				break
			}
		}
	}

	entities.RemoveFromArea(mission.Key(), mission.AreaSet())
	entities.Delete(mission.Key())
	entities.Save(target)

	if excessShips > 0 {
		startExcessMission(mission, target, excessShips)
	}
}

func startExcessMission(mission *entities.Mission, homePlanet *entities.Planet, ships int32) {
	newTargetKey := fmt.Sprintf("planet.%s", mission.Source.Name)
	newTargetEntity, err := entities.Get(newTargetKey)
	if err != nil {
		log.Print("Error in newTarget planet fetch: ", err.Error())
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
	clients.BroadcastToAll(sendMission)
}

func updateSpyReports(player *entities.Player, mission *entities.Mission, state *response.StateChange) {
	for e := clients.pool[player].Front(); e != nil; e = e.Next() {
		client := e.Value.(*Client)
		if client.Player.Username == mission.Player {
			client.Player.UpdateSpyReports()
			clients.Send(client.Player, state)
		}
	}
}
