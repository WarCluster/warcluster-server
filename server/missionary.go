package server

import (
	"fmt"
	"log"
	"time"

	"warcluster/entities"
	"warcluster/server/response"
)

// Spawns missionary for all mission records found
// in the database when the server is started
func SpawnDbMissions() {
	for _, entity := range entities.Find("mission.*") {
		mission, ok := entity.(*entities.Mission)
		if !ok {
			log.Printf("Record %s does not seem to be a mission!?\n", mission.Key())
		}

		sourceKey := fmt.Sprintf("planet.%s", mission.Source.Name)
		sourceEntity, err := entities.Get(sourceKey)
		if err != nil {
			log.Printf("Can't find planet %s for mission %s!?\n", sourceKey, mission.Key())
		}
		source := sourceEntity.(*entities.Planet)

		mission.SetAreaSet(source.AreaSet())
		log.Printf(
			"Spawning %s's mission from %s to %s...\n",
			mission.Player,
			mission.Source.Name,
			mission.Target.Name,
		)
		go StartMissionary(mission)
	}
}

// StartMissionary is used when a call to initiate a new mission is rescived.
// 1. When the delay ends the thread ends the mission calling EndMission
// 2. The end of the mission is bradcasted to all clients and the mission entry is erased from the DB.
func StartMissionary(mission *entities.Mission) {
	var (
		err             error
		excessShips     int32
		ownerHasChanged bool
		foundStartPoint bool
		player          *entities.Player
		stateChange     *response.StateChange
		target          *entities.Planet
		timeSlept       time.Duration
	)

	initialTimeSlept := time.Duration(time.Now().UnixNano()/1e6 - mission.StartTime)
	if initialTimeSlept > 0 {
		timeSlept = initialTimeSlept
	} else {
		foundStartPoint = true
	}

	entities.Save(mission)
	targetKey := fmt.Sprintf("planet.%s", mission.Target.Name)
	for _, transferPoint := range mission.TransferPoints() {
		if !foundStartPoint {
			if initialTimeSlept > transferPoint.TravelTime {
				initialTimeSlept -= transferPoint.TravelTime
				mission.ChangeAreaSet(transferPoint.CoordinateAxis, transferPoint.Direction)
				continue
			} else {
				foundStartPoint = true
			}
		}

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
	target, stateChange, err = fetchMissionTarget(targetKey)
	ownerBeforeMission := target.Owner

	if ownerBeforeMission == "" {
		player = nil
	} else {
		playerEntity, pErr := entities.Get(fmt.Sprintf("player.%s", ownerBeforeMission))
		if pErr != nil {
			log.Println("Error in target planet owner fetch:", pErr.Error())
			return
		}
		player = playerEntity.(*entities.Player)
	}

	switch mission.Type {
	case "Attack":
		if err != nil {
			log.Print("Error in target planet fetch:", err.Error())
		}
		excessShips, ownerHasChanged = mission.EndAttackMission(target)
		clients.BroadcastToAll(stateChange)
	case "Supply":
		if err != nil {
			log.Print("Error in target planet fetch:", err.Error())
		}
		excessShips, ownerHasChanged = mission.EndSupplyMission(target)
		if player != nil {
			clients.Send(player, stateChange)
		}
	case "Spy":
		for {
			if err != nil {
				log.Print("Error in target planet fetch:", err.Error())
			}
			// All spy pilots die if planet is overtaken (they are killed)
			// Other possible solution is to generate a supply mission back (they flee)
			if target.Owner != mission.Target.Owner {
				break
			}
			mission.EndSpyMission(target)
			updateSpyReports(mission, stateChange)
			if mission.ShipCount > 0 {
				time.Sleep(entities.Settings.SpyReportValidity * time.Second)
			} else {
				break
			}
			target, stateChange, err = fetchMissionTarget(targetKey)
		}
		time.Sleep(entities.Settings.SpyReportValidity * time.Second)
		updateSpyReports(mission, stateChange)
	}

	entities.RemoveFromArea(mission.Key(), mission.AreaSet())
	entities.Delete(mission.Key())
	entities.Save(target)

	if ownerHasChanged {
		go func(owned, owner string) {
			leaderBoard.Channel <- [2]string{owned, owner}
		}(ownerBeforeMission, target.Owner)

		if player != nil {
			ownerChange := response.NewOwnerChange()
			ownerChange.RawPlanet = map[string]*entities.Planet{
				target.Key(): target,
			}
			clients.Send(player, ownerChange)
		}
	}

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

func updateSpyReports(mission *entities.Mission, state *response.StateChange) {
	var (
		player *entities.Player
		err    error
	)

	if mission.Player == "" {
		log.Print("Error! Found mission with empty owner.")
		return
	}

	player, err = clients.Player(mission.Player)
	if err != nil {
		return
	}

	for element := clients.pool[mission.Player].Front(); element != nil; element = element.Next() {
		client := element.Value.(*Client)
		client.Player.UpdateSpyReports()
	}
	clients.Send(player, state)
}

func fetchMissionTarget(targetKey string) (*entities.Planet, *response.StateChange, error) {
	targetEntity, err := entities.Get(targetKey)
	if err != nil {
		return nil, nil, err
	}
	target := targetEntity.(*entities.Planet)
	target.UpdateShipCount()

	stateChange := response.NewStateChange()
	stateChange.RawPlanets = map[string]*entities.Planet{
		target.Key(): target,
	}

	return target, stateChange, nil
}
