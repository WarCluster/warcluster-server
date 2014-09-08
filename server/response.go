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
	request.Client.MoveToAreas(response.Areas())

	return nil
}

// This function makes all the checks needed for creation of a new mission.
func parseAction(request *Request) error {
	var err error

	sendMissionMessage := response.NewSendMissions()

	defer func() error {
		if panicked := recover(); panicked != nil {
			err = errors.New("Parse action panic!")
		}
		return err
	}()

	if len(request.StartPlanets) == 0 {
		errorMessage := "No start planets provided"
		sendMissionMessage.FailedMissions["Global"] = errorMessage
		return errors.New(errorMessage)
	}

	target, err := entities.Get(request.EndPlanet)
	if err != nil {
		errorMessage := "End planet does not exist"
		sendMissionMessage.FailedMissions["Global"] = errorMessage
		return errors.New(errorMessage)
	}
	endPlanet := target.(*entities.Planet)

	if _, isMissionTypeValid := missionTypes[request.Type]; !isMissionTypeValid {
		errorMessage := "Invalid mission type!"
		sendMissionMessage.FailedMissions["Global"] = errorMessage
		return errors.New(errorMessage)
	}

	for _, startPlanet := range request.StartPlanets {

		mission, err := prepareMission(startPlanet, endPlanet, request)

		if err == nil {
			sendMissionMessage.Missions[mission.Key()] = mission
		} else {
			sendMissionMessage.FailedMissions[startPlanet] = err.Error()
		}
	}

	request.Client.Send(sendMissionMessage)

	return nil
}

func prepareMission(startPlanet string, endPlanet *entities.Planet, request *Request) (*entities.Mission, error) {
	sourceEntity, err := entities.Get(startPlanet)
	if err != nil {
		return nil, err
	}

	source := sourceEntity.(*entities.Planet)

	if source.Owner != request.Client.Player.Username {

		return nil, errors.New("The mission owner does not own the start planet.")
	}

	if startPlanet == request.EndPlanet {
		return nil, errors.New("Start and end planet are the same.")
	}

	mission := request.Client.Player.StartMission(
		source,
		endPlanet,
		request.Fleet,
		request.Type,
	)

	if mission.ShipCount == 0 {
		return nil, errors.New("Not enough pilots on source planet!")
	}

	entities.Save(source)
	go StartMissionary(mission)
	entities.Save(mission)
	clients.Broadcast(source)

	return mission, nil
}
