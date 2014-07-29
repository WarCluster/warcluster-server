package server

import (
	"errors"

	"warcluster/entities"
	"warcluster/server/response"
)

var missionTypes = map[string]struct{}{
	"Attack": {},
	"Supply": {},
	"Spy":    {},
}

// scopeOfView is not finished yet but the purpose of the function is
// to call calculateCanvasSize and give the player the information
// contained in the given borders.
func scopeOfView(request *Request) error {
	response := response.NewScopeOfView(request.Position, request.Resolution)
	request.Client.Player.ScreenPosition = request.Position
	go entities.Save(request.Client.Player)
	clients.Send(request.Client.Player, response)
	return nil
}

// This function makes all the checks needed for creation of a new mission.
func parseAction(request *Request) error {
	var err error

	defer func() error {
		if panicked := recover(); panicked != nil {
			err = errors.New("Invalid action!")
		}
		return err
	}()

	if len(request.StartPlanets) == 0 {
		return errors.New("No start planets provided")
	}

	target, err := entities.Get(request.EndPlanet)
	if err != nil {
		return errors.New("End planet does not exist")
	}
	endPlanet := target.(*entities.Planet)

	if _, isMissionTypeValid := missionTypes[request.Type]; !isMissionTypeValid {
		return errors.New("Invalid mission type!")
	}

	for _, startPlanet := range request.StartPlanets {

		missionErr := prepareMission(startPlanet, endPlanet, request)

		if missionErr != nil {
			clients.Send(request.Client.Player, missionErr)
		}
	}

	return nil
}

func prepareMission(startPlanet string, endPlanet *entities.Planet, request *Request) *response.SendMissionFailed {
	missionFailed := response.NewSendMissionFailed()

	sourceEntity, err := entities.Get(startPlanet)
	if err != nil {
		missionFailed.Error = err.Error()
		return missionFailed
	}

	source := sourceEntity.(*entities.Planet)

	if source.Owner != request.Client.Player.Username {
		missionFailed.Error = "The mission owner does not own the start planet."
		return missionFailed
	}

	if startPlanet == request.EndPlanet {
		missionFailed.Error = "Start and end planet are the same. Mission cancelled."
		return missionFailed
	}

	mission := request.Client.Player.StartMission(
		source,
		endPlanet,
		request.Fleet,
		request.Type,
	)

	if mission.ShipCount == 0 {
		missionFailed.Error = "Not enough pilots on source planet!"
		return missionFailed
	}
	executeMission(mission, source)
	return nil
}

func executeMission(mission *entities.Mission, base *entities.Planet) {
	entities.Save(base)
	go StartMissionary(mission)
	entities.Save(mission)

	sendMission := response.NewSendMission()
	sendMission.Mission = mission
	clients.BroadcastToAll(sendMission)

	stateChange := response.NewStateChange()
	stateChange.RawPlanets = map[string]*entities.Planet{
		base.Key(): base,
	}
	clients.BroadcastToAll(stateChange)
}
