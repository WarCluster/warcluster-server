package server

import (
	"errors"
	"fmt"
	"warcluster/db_manager"
	"warcluster/entities"
)

/*
The three constants bellow are used by calculateCanvasSize to determine the size of the area for wich information will be sent to the user.
*/
const BEST_PING = 150
const WORST_PING = 1500
const STEPS = 10

/*
calculateCanvasSize is used to determine how big of an area(information about an area)
 do we need to send to the user to eleminate traces of lag.
*/
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

/*
scopeOfView is not finished yet but the purpose of the function is to call calculateCanvasSize
 and give the player the information contained in the given borders.

 TODO: Make some proper JSON Unmarshaling out here
*/
func scopeOfView(request *Request) error {
	var line string
	entity_list := db_manager.GetEntities("*")
	for _, entity := range entity_list {
		switch t := entity.(type) {
		case entities.Mission, entities.Planet, entities.Player, entities.Sun:
			if key, json, err := t.Serialize(); err == nil {
				line += fmt.Sprintf("\"%v\": %s, ", key, json)
			}
		}
	}
	line = fmt.Sprintf("{\"Command\": \"scope_of_view_result\", \"Planets\": {%v}}", line[:len(line)-2])
	request.Client.Session.Send([]byte(fmt.Sprintf("%v", line)))
	return nil
}

/*
This function makes all the checks needed for creation of a new mission.
*/
func parseAction(request *Request) error {
	var err error = nil

	defer func() error {
		if panicked := recover(); panicked != nil {
			err = errors.New("Invalid action!")
		}
		return nil
	}()

	start_planet, err := db_manager.GetEntity(request.StartPlanet)
	if err != nil {
		return errors.New("Start planet does not exist")
	}

	end_planet, err := db_manager.GetEntity(request.EndPlanet)
	if err != nil {
		return errors.New("End planet does not exist")
	}

	if start_planet.(entities.Planet).Owner != request.Client.Player.String() {
		err = errors.New("This is not your home!")
	}

	mission := request.Client.Player.StartMission(start_planet.(entities.Planet), end_planet.(entities.Planet), request.Fleet)
	if key, serialized_mission, err := mission.Serialize(); err == nil {
		go StartMissionary(mission)
		db_manager.SetEntity(mission)
		sessions.Broadcast([]byte(fmt.Sprintf("{\"%s\": %s}", key, serialized_mission)))
		return nil
	}

	return err
}
