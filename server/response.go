package server

import (
	"../db_manager"
	e "../entities"
	"errors"
	"fmt"
)

const BEST_PING = 150
const WORST_PING = 1500
const STEPS = 10

func scopeOfView(position []int, resolution []int, lag int) ([]int, []int) {
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

func actionParser(ch chan<- string, player *e.Player, start_planet_key, end_planet_key string, fleet int) (string, error) {
	var err error = nil
	var result string = ""

	defer func() (string, error) {
		if panicked := recover(); panicked != nil {
			err = errors.New("Invalid action!")
		}
		return result, nil
	}()

	if err != nil {
		return result, errors.New("Player does not exist")
	}

	start_planet, err := db_manager.GetEntity(start_planet_key)
	if err != nil {
		return result, errors.New("Start planet does not exist")
	}

	end_planet, err := db_manager.GetEntity(end_planet_key)
	if err != nil {
		return result, errors.New("End planet does not exist")
	}

	if start_planet.(e.Planet).Owner != player.String() {
		err = errors.New("This is not your home!")
	}

	mission := player.StartMission(start_planet.(e.Planet), end_planet.(e.Planet), fleet)
	if key, serialized_mission, err := mission.Serialize(); err == nil {
		go StartMissionary(ch, mission)
		db_manager.SetEntity(mission)
		return fmt.Sprintf("{%s: %s}", key, serialized_mission), err
	}
	return "", err
}
