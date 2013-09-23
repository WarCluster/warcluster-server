package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"warcluster/db_manager"
	"warcluster/entities"
	"warcluster/server/response"
)

// The three constants bellow are used by calculateCanvasSize to determine
//the size of the area for wich information will be sent to the user.
const (
	BEST_PING = 150
	WORST_PING = 1500
	STEPS = 10
)

// calculateCanvasSize is used to determine how big of an area(information about an area)
// do we need to send to the user to eleminate traces of lag.
func calculateCanvasSize(position []int, resolution []int, lag int) ([]int, []int) {
	step := int(WORST_PING - BEST_PING/STEPS)
	multiply := 1.1 + float32((lag-BEST_PING)/step)*0.1
	end_resolution := []int{
		int(float32(resolution[0]) * multiply),
		int(float32(resolution[1]) * multiply),
	}

	top_left := []int{
		position[0] - int((end_resolution[0]-resolution[0])/2),
		position[1] - int((end_resolution[1]-resolution[1])/2),
	}

	bottom_right := []int{
		position[0] + resolution[0] + int((end_resolution[0]-resolution[0])/2),
		position[1] + resolution[1] + int((end_resolution[1]-resolution[1])/2),
	}
	return top_left, bottom_right
}

// scopeOfView is not finished yet but the purpose of the function is to call calculateCanvasSize
// and give the player the information contained in the given borders.
func scopeOfView(request *Request) error {
	response := new(response.Response)
	response.Command = "scope_of_view_result"

	populate_entities := func(query string) (map[string]*entities.Entity) {
		result := make(map[string]*entities.Entity)
		entities := db_manager.GetEntities(query)
		for _, entity := range entities {
			result[entity.GetKey()] = &entity
		}
		return result
	}

	response.Entities = make(map[string]map[string]*entities.Entity)
	response.Entities["Missions"] = populate_entities("mission.*")
	response.Entities["Planets"]  = populate_entities("planet.*")
	response.Entities["Suns"]     = populate_entities("sun.*")

	if json_response, err := json.Marshal(response); err == nil {
		request.Client.Session.Send(json_response)
	} else {
		return err
	}

	return nil
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

	source, err := db_manager.GetEntity(request.StartPlanet)
	if err != nil {
		return errors.New("Start planet does not exist")
	}

	target, err := db_manager.GetEntity(request.EndPlanet)
	if err != nil {
		return errors.New("End planet does not exist")
	}

	if source.(*entities.Planet).Owner != request.Client.Player.String() {
		err = errors.New("This is not your home!")
	}

	if request.Type != "Attack" || request.Type != "Supply" || request.Type != "Spy" {
		err = errors.New("Invalid mission type!")
	}

	mission := request.Client.Player.StartMission(source.(*entities.Planet), target.(*entities.Planet), request.Fleet, request.Type)
	if _, serialized_mission, err := mission.Serialize(); err == nil {
		go StartMissionary(mission)
		db_manager.SetEntity(mission)
		sessions.Broadcast([]byte(fmt.Sprintf("{ \"Command\": \"send_mission\", \"Mission\": %s}", serialized_mission)))
		if source_key, source_json, source_err := source.Serialize(); source_err == nil {
			sessions.Broadcast([]byte(fmt.Sprintf("{\"Command\": \"state_change\", \"Planets\": {\"%s\": %s}}", source_key, source_json)))
			db_manager.SetEntity(source)
		}
		return nil
	}

	return err
}
