package server

import (
	"errors"

	"github.com/Vladimiroff/vec2d"

	"warcluster/entities"
	"warcluster/server/response"
)

// calculateCanvasSize is used to determine where is the viewable by client's area
func calculateCanvasSize(position *vec2d.Vector, resolution []int) (*vec2d.Vector, *vec2d.Vector) {
	topLeft := vec2d.New(
		position.X-float64(resolution[0])/2,
		position.Y+float64(resolution[1])/2,
	)

	bottomRight := vec2d.New(
		position.X+float64(resolution[0])/2,
		position.Y-float64(resolution[1])/2,
	)
	return topLeft, bottomRight
}

// scopeOfView is not finished yet but the purpose of the function is
// to call calculateCanvasSize and give the player the information
// contained in the given borders.
func scopeOfView(request *Request) error {
	res := response.NewScopeOfView()

	populateEntities := func(query string) map[string]entities.Entity {
		result := make(map[string]entities.Entity)
		entities := entities.Find(query)
		for _, entity := range entities {
			result[entity.Key()] = entity
		}
		return result
	}

	res.Missions = populateEntities("mission.*")
	res.Planets = populateEntities("planet.*")
	res.Suns = populateEntities("sun.*")
	request.Client.Player.ScreenPosition = request.Position
	go entities.Save(request.Client.Player)

	return response.Send(res, request.Client.Session.Send)
}

// This function makes all the checks needed for creation of a new mission.
func parseAction(request *Request) error {
	var err error = nil

	defer func() error {
		if panicked := recover(); panicked != nil {
			err = errors.New("Invalid action!")
		}
		return nil
	}()

	source, err := entities.Get(request.StartPlanet)
	if err != nil {
		return errors.New("Start planet does not exist")
	}

	target, err := entities.Get(request.EndPlanet)
	if err != nil {
		return errors.New("End planet does not exist")
	}

	if source.(*entities.Planet).Owner != request.Client.Player.Username {
		return errors.New("This is not your home!")
	}

	if request.Type != "Attack" && request.Type != "Supply" && request.Type != "Spy" {
		return errors.New("Invalid mission type!")
	}

	mission := request.Client.Player.StartMission(
		source.(*entities.Planet),
		target.(*entities.Planet),
		request.Fleet,
		request.Type,
	)

	if mission.ShipCount == 0 {
		missionFailed := response.NewSendMissionFailed()
		missionFailed.Error = "Not enough pilots on source planet!"
		response.Send(missionFailed, request.Client.Session.Send)
	}

	go StartMissionary(mission)
	entities.Save(mission)
	entities.Save(source)

	sendMission := response.NewSendMission()
	sendMission.Mission = mission
	err = response.Send(sendMission, sessions.Broadcast)
	if err != nil {
		return err
	}

	stateChange := response.NewStateChange()
	stateChange.Planets = map[string]entities.Entity{
		source.Key(): source,
	}
	return response.Send(stateChange, sessions.Broadcast)
}
