package server

import (
	"errors"

	"warcluster/entities"
	"warcluster/server/response"
)

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
// TODO: Do not stream on spy mission
func parseAction(request *Request) error {
	var err error

	defer func() error {
		if panicked := recover(); panicked != nil {
			err = errors.New("Invalid action!")
		}
		return nil
	}()

	sourceEntity, err := entities.Get(request.StartPlanet)
	if err != nil {
		return errors.New("Start planet does not exist")
	}
	source := sourceEntity.(*entities.Planet)

	target, err := entities.Get(request.EndPlanet)
	if err != nil {
		return errors.New("End planet does not exist")
	}

	if source.Owner != request.Client.Player.Username {
		return errors.New("The mission owner does not own the start planet.")
	}

	// FIXME: Why not a simple list with the possible attacks?
	if request.Type != "Attack" && request.Type != "Supply" && request.Type != "Spy" {
		return errors.New("Invalid mission type!")
	}

	if request.StartPlanet == request.EndPlanet {
		clients.Send(request.Client.Player, response.NewComsError("Unable to target the same planet."))
		return errors.New("Start and end planet are the same. Mission cancelled.")
	}

	mission := request.Client.Player.StartMission(
		source,
		target.(*entities.Planet),
		request.Fleet,
		request.Type,
	)

	if mission.ShipCount == 0 {
		missionFailed := response.NewSendMissionFailed()
		missionFailed.Error = "Not enough pilots on source planet!"
		clients.Send(request.Client.Player, missionFailed)
	}

	entities.Save(source)
	go StartMissionary(mission, 0)
	entities.Save(mission)

	sendMission := response.NewSendMission()
	sendMission.Mission = mission
	clients.BroadcastToAll(sendMission)

	stateChange := response.NewStateChange()
	stateChange.RawPlanets = map[string]*entities.Planet{
		source.Key(): source,
	}
	clients.BroadcastToAll(stateChange)
	return nil
}
