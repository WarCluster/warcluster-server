package server

import (
	"fmt"
	"log"
	"time"

	"warcluster/entities"
	"warcluster/server/response"
)

func FetchMissionTarget(targetKey string) (*entities.Planet, *response.StateChange, error) {
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

// StartMissionary is used when a call to initiate a new mission is rescived.
// 1. When the delay ends the thread ends the mission calling EndMission
// 2. The end of the mission is bradcasted to all clients and the mission entry is erased from the DB.
func StartMissionary(mission *entities.Mission) {
	var (
		timeSlept   time.Duration = 0
		excessShips int32
		target      *entities.Planet
		stateChange *response.StateChange
		err         error
	)
	entities.Save(mission)
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

	var player *entities.Player

	if mission.Target.Owner == "" {
		player = nil
	} else {
		playerEntity, pErr := entities.Get(fmt.Sprintf("player.%s", mission.Target.Owner))
		if pErr != nil {
			log.Println("Error in target planet owner fetch:", pErr.Error())
			return
		}
		player = playerEntity.(*entities.Player)
	}

	switch mission.Type {
	case "Attack":
		target, stateChange, err = FetchMissionTarget(targetKey)
		if err != nil {
			log.Print("Error in target planet fetch:", err.Error())
		}
		ownerBefore := target.Owner
		excessShips = mission.EndAttackMission(target)
		clients.BroadcastToAll(stateChange)
		if ownerBefore != target.Owner {
			fmt.Println("Owner changed")
			go func(owned, owner string) {
				leaderBoard.Transfer(owned, owner)
			}(ownerBefore, target.Owner)

			if player != nil {
				ownerChange := response.NewOwnerChange()
				ownerChange.RawPlanet = map[string]*entities.Planet{
					target.Key(): target,
				}
				clients.Send(player, ownerChange)
			}
		}
	case "Supply":
		target, stateChange, err = FetchMissionTarget(targetKey)
		if err != nil {
			log.Print("Error in target planet fetch:", err.Error())
		}
		excessShips = mission.EndSupplyMission(target)
		if player != nil {
			clients.Send(player, stateChange)
		}
	case "Spy":
		for {
			target, stateChange, err = FetchMissionTarget(targetKey)
			if err != nil {
				log.Print("Error in target planet fetch:", err.Error())
			}
			// All spy pilots die if planet is overtaken (they are killed)
			// Other possible solution is to generate a supply mission back (they flee)
			if target.Owner != mission.Target.Owner {
				fmt.Println("ownera se smeni trugvam si")
				break
			}
			mission.EndSpyMission(target)
			fmt.Println("endSpyMission")
			updateSpyReports(mission, stateChange)
			if mission.ShipCount > 0 {
				time.Sleep(entities.SPY_REPORT_VALIDITY * time.Second)
			} else {
				fmt.Println("Nqma poveche shpioni")
				break
			}
		}
		time.Sleep(entities.SPY_REPORT_VALIDITY * time.Second)
		updateSpyReports(mission, stateChange)
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
